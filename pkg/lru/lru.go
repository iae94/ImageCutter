package lru

import (
	"ImageCutter/pkg/models"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Cache struct {
	CurrentSize int64
	MaxSize int64
	CleanInterval int
	Folder string
	Storage []*models.Image
	Logger *zap.Logger
	lock *sync.RWMutex
}

func NewCache (logger *zap.Logger, size int64, folder string, cleanInterval int) (*Cache, error) {

	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err = os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			mess := fmt.Sprintf("Cannot create cache folder at %v\n", folder)
			logger.Error(mess)
			return nil, err
		}
		logger.Sugar().Infof("Cache folder: '%v' was created", folder)
	} else {
		logger.Sugar().Infof("Cache folder: '%v' is exist", folder)
	}
	cache := &Cache{
		MaxSize: size * 1024 * 1024,
		Folder: folder,
		CleanInterval: cleanInterval,
		Storage: make([]*models.Image, 0),
		Logger: logger,
		lock: &sync.RWMutex{},
	}
	logger.Info("Start cache cleaner goroutine")
	go cache.Cleaner() // Cache cleaner

	return cache, nil
}


func (cc *Cache) Add(img *models.Image) error{
	// if image size too big - not put it in cache
	if img.Size > cc.MaxSize {
		mess := fmt.Sprintf("Image size is higher than maximum cache size! %v Kb vs %v Kb. This image will not be caching!", img.Size / 1024, cc.MaxSize / 1024)
		cc.Logger.Info(mess)
		return errors.New(mess)
	}

	// will remove cache images until get enough cache space for incoming image
	tries := 1
	for cc.CurrentSize + img.Size > cc.MaxSize {
		if len(cc.Storage) == 1{
			cc.lock.Lock() // Lock for safety
			err := cc.Delete(cc.Storage[0])
			cc.lock.Unlock()
			if err != nil {
				cc.Logger.Sugar().Errorf("Deleting cache image give error: %v", err)
				return err
			}
			break
		}
		cc.Logger.Sugar().Infof("Free cache space is not enough for incoming image with size: %v Kb. Try remove oldest images from cache (%v try)", img.Size / 1024, tries)
		err := cc.RemoveOldest()
		if err != nil {
			cc.Logger.Sugar().Errorf("RemoveOldest give error: %v", err)
			return err
		}
		tries += 1
	}
	cc.lock.Lock()
	cc.Storage = append(cc.Storage, img)
	cc.CurrentSize += img.Size
	cc.Logger.Sugar().Infof("Cache size increased from %v/%v KB to %v/%v KB", cc.CurrentSize-img.Size, cc.MaxSize, cc.CurrentSize, cc.MaxSize)
	cc.lock.Unlock()

	return nil
}

func (cc *Cache) Delete(image *models.Image) error{
	imagePath := filepath.Join(cc.Folder, image.Name)

	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		cc.Logger.Sugar().Errorf("Image %v is not found on disk!", imagePath)
	} else {
		err := os.Remove(imagePath)
		if err != nil{
			cc.Logger.Sugar().Errorf("Removing image: %v from disk give error: %v", imagePath, err)
			return err
		}
	}
	ind, err := cc.GetImageIndex(image)
	if err != nil {
		cc.Logger.Sugar().Errorf("Image %v is not found in cache storage!", imagePath)
	} else {
		// Delete elem from slice
		if ind < len(cc.Storage) - 1 {
			copy(cc.Storage[ind:], cc.Storage[ind+1:])
		}
		cc.Storage[len(cc.Storage)-1] = nil
		cc.Storage = cc.Storage[:len(cc.Storage)-1]
		cc.CurrentSize -= image.Size // Decrease current cache size
	}

	return nil
}


func (cc *Cache) GetImageByUrl(url string) (*models.Image, error) {
	for _, img := range cc.Storage{
		if img.Url == url {

			imagePath := filepath.Join(cc.Folder, img.Name)
			if _, err := os.Stat(imagePath); os.IsNotExist(err) {
				mess := fmt.Sprintf("Already cached image %v is not found on disk!", img.Url)
				cc.Logger.Warn(mess)
				return nil, errors.New(mess)
			}
			return img, nil
		}
	}
	mess := fmt.Sprintf("Image with url: %v not in cache", url)
	cc.Logger.Info(mess)

	return nil, errors.New(mess)
}

func (cc *Cache) GetImageIndex(image *models.Image) (int, error) {
	for ind, img := range cc.Storage {
		if img.Name == image.Name && img.Url == image.Url {
			return ind, nil
		}
	}
	return -1, fmt.Errorf("Image with name: %v not found in cache storage", image.Name)
}

func (cc *Cache) RemoveOldest() error{
	if len(cc.Storage) <= 1{
		if len(cc.Storage) == 0{
			cc.Logger.Sugar().Infof("Cache is empty!")
		}
		if len(cc.Storage) == 1{
			cc.Logger.Sugar().Infof("Only 1 image in cache. Cache should keep at least 1 image")
		}
		return nil
	}
	cc.Logger.Sugar().Infof("Cache size before clean: %v/%v KB", cc.CurrentSize / 1024, cc.MaxSize / 1024)

	minFetch := cc.Storage[0].FetchCount

	for _, img := range cc.Storage{
		if img.FetchCount < minFetch {
			minFetch = img.FetchCount
		}
	}
	deletedImages := make([]*models.Image, 0)
	for _, img := range cc.Storage{
		if img.FetchCount == minFetch {
			deletedImages = append(deletedImages, img)
		}
	}
	// If all elems have equal FetchCount -> we should keep at least one elem in cache
	leftOne := false
	if len(deletedImages) == len(cc.Storage){
		cc.Logger.Sugar().Infof("%v images have equal FetchCount -> At least 1 image will be kept in cache", len(deletedImages))
		leftOne = true
	}
	deleted := 0
	for _, img := range deletedImages{
		if leftOne { // Keep first image
			leftOne = false
			continue
		}
		cc.lock.Lock() // Lock for safety
		err := cc.Delete(img)
		cc.lock.Unlock()
		if err != nil {
			cc.Logger.Sugar().Errorf("Deleting cache image give error: %v", err)
			return err
		}
		deleted += 1
	}
	cc.Logger.Sugar().Infof("%v oldest images was deleted from cache", deleted)
	cc.Logger.Sugar().Infof("Cache size after clean: %v/%v KB", cc.CurrentSize / 1024, cc.MaxSize / 1024)
	return nil

}


func (cc *Cache) Cleaner() {
	sleepTime := time.Minute * time.Duration(cc.CleanInterval)
	for {
		cc.Logger.Info("Cache cleaner try remove oldest instances from cache...")
		_ = cc.RemoveOldest()
		time.Sleep(sleepTime)
	}
}