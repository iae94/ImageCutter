package cutter

import (
	cfg "ImageCutter/pkg/config"
	"ImageCutter/pkg/cropper"
	"ImageCutter/pkg/lru"
	"ImageCutter/pkg/models"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io"
	"net/http"
	urllib "net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)


type CutterService struct {
	Logger *zap.Logger
	Config *cfg.CutterConfig
	Cropper *cropper.Cropper
	Cache *lru.Cache
}

func NewCutterService(logger *zap.Logger, config *cfg.CutterConfig) (*CutterService, error) {
	cp := cropper.NewCropper(logger, config)

	logger.Sugar().Infof("Init Cache instance with parameters:\nCACHESIZE=%v\nCACHECLEAN=%v\nCACHEFOLDER=%v\n", config.Cutter.Cache.Size, config.Cutter.Cache.CleanInterval, config.Cutter.Cache.Folder)

	cache, err := lru.NewCache(logger, config.Cutter.Cache.Size, config.Cutter.Cache.Folder, config.Cutter.Cache.CleanInterval)
	if err != nil {
		logger.Sugar().Errorf("Creating instance of Cache give error: %v", err)
		return nil, err
	}

	return &CutterService{
		Logger: logger,
		Config: config,
		Cropper: cp,
		Cache: cache,
	}, nil
}



func (cs *CutterService) Start (){
	router := mux.NewRouter()

	router.HandleFunc("/crop/{width}/{height}/{url:(?:.+)}", cs.Crop)
	router.HandleFunc("/cache/{url:(?:.+)}", cs.CheckCache)


	http.Handle("/", router)

	address := fmt.Sprintf(":%v", cs.Config.Cutter.Port)
	cs.Logger.Sugar().Infof("Start cutter service at address: %v", address)
	err := http.ListenAndServe(address, nil)
	cs.Logger.Sugar().Fatalf("HTTP Listener give error: %v", err)
}

func (cs *CutterService) CheckCache(w http.ResponseWriter, r *http.Request) {

	cs.Logger.Info("Try check image in cache...")
	args := mux.Vars(r)

	url := args["url"]
	u, _ := urllib.Parse(url)
	if u.Scheme == ""{
		mess := fmt.Sprintf("Remote image url is incorrect: %v. Protocol is missed(require http:// or https://)", url)
		cs.Logger.Error(mess)
		http.Error(w, mess, 400)
		return
	}
	// Url scheme must contains two slash not one
	if u.Host == ""{
		url = fmt.Sprintf("%v:/%v", u.Scheme, u.Path)
	}

	// Try get from cache
	_, err := cs.Cache.GetImageByUrl(url)
	if err != nil{
		cs.Logger.Sugar().Infof("Image with url: %v not in cache", url)
		http.Error(w, fmt.Sprintf("Image with url: %v not in cache :(", url), 404)
	} else {
		cs.Logger.Sugar().Infof("Image with url: %v in cache", url)
		http.Error(w, fmt.Sprintf("Image with url: %v in cache :)", url), 200)
	}

}

func (cs *CutterService) Crop(w http.ResponseWriter, r *http.Request) {
	cs.Logger.Info("Try crop image...")
	args := mux.Vars(r)

	if args["url"] == ""{
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, "Image url is required", 400)
		return
	}
	if args["width"] == ""{
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, "Crop width is required", 400)
		return
	}
	if args["height"] == ""{
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, "Crop height is required", 400)
		return
	}

	url := args["url"]
	u, _ := urllib.Parse(url)
	if u.Scheme == ""{
		mess := fmt.Sprintf("Remote image url is incorrect: %v. Protocol is missed(require http:// or https://)", url)
		cs.Logger.Error(mess)
		http.Error(w, mess, 400)
		return
	}
	// Url scheme must contains two slash not one
	if u.Host == ""{
		url = fmt.Sprintf("%v:/%v", u.Scheme, u.Path)
	}

	width, err := strconv.Atoi(args["width"])
	if err != nil {
		mess := fmt.Sprintf("Cannot convert width to int from string: %v", err)
		cs.Logger.Error(mess)
		http.Error(w, mess, 500)
		return
	}
	height, err := strconv.Atoi(args["height"])
	if err != nil {
		mess := fmt.Sprintf("Cannot convert height to int from string: %v", err)
		cs.Logger.Error(mess)
		http.Error(w, mess, 500)
		return
	}
	if width == 0 && height == 0{
		mess := fmt.Sprintf("Both width and height are zero!")
		cs.Logger.Error(mess)
		http.Error(w, mess, 400)
		return
	}

	// Try get from cache
	cacheImage, err := cs.Cache.GetImageByUrl(url)

	// If image not in cache
	if err != nil {
		// Get image from remote server
		var code int
		cacheImage, code, err = cs.FetchImage(url)
		if err != nil {
			mess := fmt.Sprintf("Fetching url: %v give error: %v", url, err)
			cs.Logger.Error(mess)
			http.Error(w, mess, code)
			return
		}

		cs.Logger.Sugar().Infof("Successfully fetched new image: %v", cacheImage.Name)

		// Add new image to cache
		err := cs.Cache.Add(cacheImage)
		if err != nil {
			cs.Logger.Sugar().Warnf("Cannot add image: %v to cache. Reason: %v",cacheImage.Name, err)
		} else {
			cs.Logger.Sugar().Infof("Image %v now in cache!", cacheImage.Url)
		}

	} else {
		cs.Logger.Sugar().Infof("Take image %v from cache", cacheImage.Url)
	}

	cacheImage.FetchCount += 1 // Increment fetch count

	croppedImage, err := cs.Cropper.Crop(width, height, cacheImage)
	if err != nil {
		mess := fmt.Sprintf("Cropping image give error: %v", err)
		cs.Logger.Error(mess)
		http.Error(w, mess, 500)
		return
	}

	if _, err := w.Write(croppedImage); err != nil {
		cs.Logger.Sugar().Errorf("Unable to write cropped image to writer: %v", err)
	}
}


func (cs *CutterService) FetchImage(url string) (*models.Image, int, error) {

	resp, err := http.Get(url)

	// If server does not exist
	if err != nil {
		cs.Logger.Sugar().Errorf("Fetching url: %v give error: %v", url, err)
		return nil, 503, err
	}
	// If server return 500 code
	if resp.StatusCode == 500{
		mess := fmt.Sprintf("Remote server error return 500 code for url: %v", url)
		cs.Logger.Info(mess)
		return nil, 500, errors.New(mess)
	}

	// If file does not exists
	if resp.StatusCode == 404{
		mess := fmt.Sprintf("File not found on url: %v", url)
		cs.Logger.Info(mess)
		return nil, 404, errors.New(mess)
	}

	// If file mime type not image
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "image"){
		mess := fmt.Sprintf("Fetching file is not image: %v", contentType)
		cs.Logger.Warn(mess)
		return nil, 422, errors.New(mess)
	}

	// Generate file name as hash of url
	extension := strings.ReplaceAll(contentType, "image/", "")
	imageName := fmt.Sprintf("%x.%v", md5.Sum([]byte(url)), extension)

	imagePath := filepath.Join(cs.Config.Cutter.Cache.Folder, imageName)

	imageFile, err := os.Create(imagePath)
	if err != nil {
		cs.Logger.Sugar().Errorf("Creating file for image give error: %v", err)
		return nil, 500, err
	}

	// Close resp and file
	defer func(){
		err := resp.Body.Close()
		if err != nil {
			cs.Logger.Sugar().Errorf("Response body closing give error: %v", err)
		}
		err = imageFile.Close()
		if err != nil {
			cs.Logger.Sugar().Errorf("Image file closing give error: %v", err)
		}
	}()

	_, err = io.Copy(imageFile, resp.Body)
	if err != nil {
		cs.Logger.Sugar().Errorf("Copying image to file give error: %v", err)
		return nil, 500, err
	}
	imageStat, err := imageFile.Stat()
	if err != nil {
		cs.Logger.Sugar().Errorf("Cannot get size of image, error: %v", err)
		return nil, 500, err
	}

	headers := make(map[string]string)
	for key, value := range resp.Header{
		headers[key] = strings.Join(value, ";")
	}

	return &models.Image{
		Name: imageName,
		Url:url,
		Headers: headers,
		FetchCount: 0,
		Size: imageStat.Size(),
		MimeType: contentType,
	}, 200, nil

}