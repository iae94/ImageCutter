package lru

import (
	"ImageCutter/pkg/models"
	logging "ImageCutter/pkg/logger"
	cfg "ImageCutter/pkg/config"
	"os"
	"path"
	"sync"
	"testing"
)

func TestCache_Add(t *testing.T) {
	type fields struct {
		Cache	*Cache
	}
	type args struct {
		img *models.Image
	}

	// Logger
	logger, err := logging.CreateLogger(&cfg.Logger{Level: "info", Encoding: "console", OutputPaths:[]string{"stdout"}, ErrorOutputPaths:[]string{"stderr"}})
	if err != nil {
		t.Errorf("Add() create logger give error: %v", err)
	}

	// Creating temp folder for cache
	cacheFolder := "test_images"
	if _, err := os.Stat(cacheFolder); os.IsExist(err) {
		err = os.RemoveAll(cacheFolder)
		if err != nil {
			t.Errorf("Add() Cannot delete cache folder at %v\n", cacheFolder)
		}
	}
	err = os.MkdirAll(cacheFolder, os.ModePerm)
	if err != nil {
		t.Errorf("Add() Cannot create cache folder at %v\n", cacheFolder)
	}


	// Emtpty cache for first test
	emptyCache := &Cache{
		MaxSize:       10 * 1024 * 1024,
		CleanInterval: 5,
		Folder:        cacheFolder,
		Storage:       make([]*models.Image, 0),
		Logger:        logger,
		lock:          &sync.RWMutex{},
	}
	if err != nil{
		t.Errorf("Add () Cannot create emptyCache instance:%v", err)
	}
	// Full cache for second test
	fullCache := &Cache{
		MaxSize:       2 * 1024 * 1024,
		CleanInterval: 5,
		Folder:        cacheFolder,
		Storage:       make([]*models.Image, 0),
		Logger:        logger,
		lock:          &sync.RWMutex{},
	}
	if err != nil{
		t.Errorf("Add () Cannot create fullCache instance:%v", err)
	}

	// Img for first test
	simpleImage := &models.Image{
		Name:       "simpleImage.jpg",
		MimeType:   "image/jpeg",
		Url:        "someurl",
		Size:       1 * 1024 * 1024,
		Headers:    nil,
		FetchCount: 0,
	}
	// Img for second test
	veryBigImage := &models.Image{
		Name:       "veryBigImage.jpg",
		MimeType:   "image/jpeg",
		Url:        "someurl",
		Size:       emptyCache.MaxSize + 1024,
		Headers:    nil,
		FetchCount: 0,
	}
	// Img for third test
	bigImage := &models.Image{
		Name:       "bigImage.jpg",
		MimeType:   "image/jpeg",
		Url:        "someurl",
		Size:       1.5 * 1024 * 1024,
		Headers:    nil,
		FetchCount: 0,
	}
	err = fullCache.Add(bigImage)
	if err != nil{
		t.Errorf("Add() Cannot add bigImage to fullCache:%v", err)
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "Add image", fields: fields{Cache: emptyCache}, args: args{img: simpleImage}, wantErr: false},
		{name: "Add image which size more than cache size", fields: fields{Cache: emptyCache}, args: args{img: veryBigImage}, wantErr: true},
		{name: "Add image in full cache", fields: fields{Cache: fullCache}, args: args{img: simpleImage}, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc := tt.fields.Cache
			if err := cc.Add(tt.args.img); (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCache_Delete(t *testing.T) {
	type fields struct {
		Cache	*Cache
	}
	type args struct {
		image *models.Image
	}

	// Logger
	logger, err := logging.CreateLogger(&cfg.Logger{Level: "info", Encoding: "console", OutputPaths:[]string{"stdout"}, ErrorOutputPaths:[]string{"stderr"}})
	if err != nil {
		t.Errorf("Delete() create logger give error: %v", err)
	}

	// Creating temp folder for cache
	cacheFolder := "test_images"
	if _, err := os.Stat(cacheFolder); os.IsExist(err) {
		err = os.RemoveAll(cacheFolder)
		if err != nil {
			t.Errorf("Add() Cannot delete cache folder at %v\n", cacheFolder)
		}
	}
	err = os.MkdirAll(cacheFolder, os.ModePerm)
	if err != nil {
		t.Errorf("Add() Cannot create cache folder at %v\n", cacheFolder)
	}


	// Emtpty cache
	emptyCache := &Cache{
		MaxSize:       10 * 1024 * 1024,
		CleanInterval: 5,
		Folder:        cacheFolder,
		Storage:       make([]*models.Image, 0),
		Logger:        logger,
		lock:          &sync.RWMutex{},
	}
	if err != nil{
		t.Errorf("Delete() Cannot create emptyCache instance:%v", err)
	}


	notExistedImage := &models.Image{
		Name:       "nonexist.jpg",
		MimeType:   "image/jpeg",
		Url:        "url not exist",
		Size:       1 * 1024 * 1024,
		Headers:    nil,
		FetchCount: 0,
	}
	onlyOnDiskImage := &models.Image{
		Name:       "onlyondisk.jpg",
		MimeType:   "image/jpeg",
		Url:        "url on disk",
		Size:       1 * 1024 * 1024,
		Headers:    nil,
		FetchCount: 0,
	}
	onlyInCacheImage := &models.Image{
		Name:       "onlyincache.jpg",
		MimeType:   "image/jpeg",
		Url:        "url in cache",
		Size:       1 * 1024 * 1024,
		Headers:    nil,
		FetchCount: 0,
	}
	existedImage := &models.Image{
		Name:       "exist.jpg",
		MimeType:   "image/jpeg",
		Url:        "url exist",
		Size:       1 * 1024 * 1024,
		Headers:    nil,
		FetchCount: 0,
	}

	existingImageFile, err := os.Create(path.Join(cacheFolder, existedImage.Name))
	if err != nil{
		t.Errorf("Delete() Cannot create existingImageFile:%v", err)
	}
	existingImageFile.Close()

	onlyOnDiskImageFile, err := os.Create(path.Join(cacheFolder, onlyOnDiskImage.Name))
	if err != nil{
		t.Errorf("Delete() Cannot create onlyOnDiskImageFile:%v", err)
	}
	onlyOnDiskImageFile.Close()


	err = emptyCache.Add(existedImage)
	if err != nil{
		t.Errorf("Delete() Cannot add existedImage to cache:%v", err)
	}
	err = emptyCache.Add(onlyOnDiskImage)
	if err != nil{
		t.Errorf("Delete() Cannot add onlyOnDiskImage to cache:%v", err)
	}
	err = emptyCache.Add(onlyInCacheImage)
	if err != nil{
		t.Errorf("Delete() Cannot add onlyInCacheImage to cache:%v", err)
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "Image in cache and on disk", fields: fields{Cache: emptyCache}, args: args{image: existedImage}, wantErr: false},
		{name: "Image only in cache", fields: fields{Cache: emptyCache}, args: args{image: onlyInCacheImage}, wantErr: false},
		{name: "Image only on disk", fields: fields{Cache: emptyCache}, args: args{image: onlyOnDiskImage}, wantErr: false},
		{name: "Image not exist", fields: fields{Cache: emptyCache}, args: args{image: notExistedImage}, wantErr: false},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc := tt.fields.Cache
			if err := cc.Delete(tt.args.image); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCache_GetImageByUrl(t *testing.T) {
	type fields struct {
		Cache	*Cache
	}
	type args struct {
		url string
	}

	// Logger
	logger, err := logging.CreateLogger(&cfg.Logger{Level: "info", Encoding: "console", OutputPaths:[]string{"stdout"}, ErrorOutputPaths:[]string{"stderr"}})
	if err != nil {
		t.Errorf("GetImageByUrl() create logger give error: %v", err)
	}

	// Creating temp folder for cache
	cacheFolder := "test_images"
	if _, err := os.Stat(cacheFolder); os.IsExist(err) {
		err = os.RemoveAll(cacheFolder)
		if err != nil {
			t.Errorf("Add() Cannot delete cache folder at %v\n", cacheFolder)
		}
	}
	err = os.MkdirAll(cacheFolder, os.ModePerm)
	if err != nil {
		t.Errorf("Add() Cannot create cache folder at %v\n", cacheFolder)
	}


	// Emtpty cache
	emptyCache := &Cache{
		MaxSize:       10 * 1024 * 1024,
		CleanInterval: 5,
		Folder:        cacheFolder,
		Storage:       make([]*models.Image, 0),
		Logger:        logger,
		lock:          &sync.RWMutex{},
	}
	if err != nil{
		t.Errorf("GetImageByUrl() Cannot create emptyCache instance:%v", err)
	}

	// Img for first test
	simpleImage := &models.Image{
		Name:       "simpleImage.jpg",
		MimeType:   "image/jpeg",
		Url:        "url1",
		Size:       1 * 1024 * 1024,
		Headers:    nil,
		FetchCount: 0,
	}
	// Img for first test
	simpleImageFile, err := os.Create(path.Join(cacheFolder, simpleImage.Name))
	if err != nil{
		t.Errorf("GetImageByUrl() Cannot create simpleImageFile:%v", err)
	}
	simpleImageFile.Close()

	// Img for second test
	imageWithoutFile := &models.Image{
		Name:       "imageWithoutFile.jpg",
		MimeType:   "image/jpeg",
		Url:        "url2",
		Size:       1 * 1024 * 1024,
		Headers:    nil,
		FetchCount: 0,
	}
	// Img for third test
	imageNotInCache := &models.Image{
		Name:       "imageNotInCache.jpg",
		MimeType:   "image/jpeg",
		Url:        "url3",
		Size:       1 * 1024 * 1024,
		Headers:    nil,
		FetchCount: 0,
	}

	err = emptyCache.Add(simpleImage)
	if err != nil{
		t.Errorf("GetImageByUrl() Cannot add simpleImage to cache:%v", err)
	}
	err = emptyCache.Add(imageWithoutFile)
	if err != nil{
		t.Errorf("GetImageByUrl() Cannot add imageWithoutFile to cache:%v", err)
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "Image not in cache", fields: fields{Cache: emptyCache}, args: args{url: imageNotInCache.Url}, wantErr: true},
		{name: "Image in cache", fields: fields{Cache: emptyCache}, args: args{url: simpleImage.Url}, wantErr: false},
		{name: "Image in cache but not found on disk", fields: fields{Cache: emptyCache}, args: args{url: imageWithoutFile.Url}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc := tt.fields.Cache
			if _, err := cc.GetImageByUrl(tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("GetImageByUrl() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCache_GetImageIndex(t *testing.T) {
	type fields struct {
		Cache	*Cache
	}
	type args struct {
		image *models.Image
	}

	// Logger
	logger, err := logging.CreateLogger(&cfg.Logger{Level: "info", Encoding: "console", OutputPaths:[]string{"stdout"}, ErrorOutputPaths:[]string{"stderr"}})
	if err != nil {
		t.Errorf("GetImageIndex() create logger give error: %v", err)
	}

	// Creating temp folder for cache
	cacheFolder := "test_images"
	if _, err := os.Stat(cacheFolder); os.IsExist(err) {
		err = os.RemoveAll(cacheFolder)
		if err != nil {
			t.Errorf("Add() Cannot delete cache folder at %v\n", cacheFolder)
		}
	}
	err = os.MkdirAll(cacheFolder, os.ModePerm)
	if err != nil {
		t.Errorf("Add() Cannot create cache folder at %v\n", cacheFolder)
	}


	// Emtpty cache
	emptyCache := &Cache{
		MaxSize:       10 * 1024 * 1024,
		CleanInterval: 5,
		Folder:        cacheFolder,
		Storage:       make([]*models.Image, 0),
		Logger:        logger,
		lock:          &sync.RWMutex{},
	}
	if err != nil{
		t.Errorf("GetImageIndex() Cannot create emptyCache instance:%v", err)
	}

	// Img for first test
	existingImage := &models.Image{
		Name:       "existingImage.jpg",
		MimeType:   "image/jpeg",
		Url:        "url_exist",
		Size:       1 * 1024 * 1024,
		Headers:    nil,
		FetchCount: 0,
	}

	// Img for second test
	notExistingImage := &models.Image{
		Name:       "notExistingImage.jpg",
		MimeType:   "image/jpeg",
		Url:        "url_not_exist",
		Size:       1 * 1024 * 1024,
		Headers:    nil,
		FetchCount: 0,
	}

	err = emptyCache.Add(existingImage)
	if err != nil{
		t.Errorf("GetImageIndex() Cannot add simpleImage to cache:%v", err)
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "Image exist", fields: fields{Cache: emptyCache}, args: args{image: existingImage}, want:0, wantErr: false},
		{name: "Image not exist", fields: fields{Cache: emptyCache}, args: args{image: notExistingImage}, want:-1, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc := tt.fields.Cache
			got, err := cc.GetImageIndex(tt.args.image)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetImageIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetImageIndex() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_RemoveOldest(t *testing.T) {
	type fields struct {
		Cache	*Cache
	}
	// Logger
	logger, err := logging.CreateLogger(&cfg.Logger{Level: "info", Encoding: "console", OutputPaths:[]string{"stdout"}, ErrorOutputPaths:[]string{"stderr"}})
	if err != nil {
		t.Errorf("RemoveOldest() create logger give error: %v", err)
	}

	// Creating temp folder for cache
	cacheFolder := "test_images"
	if _, err := os.Stat(cacheFolder); os.IsExist(err) {
		err = os.RemoveAll(cacheFolder)
		if err != nil {
			t.Errorf("Add() Cannot delete cache folder at %v\n", cacheFolder)
		}
	}
	err = os.MkdirAll(cacheFolder, os.ModePerm)
	if err != nil {
		t.Errorf("Add() Cannot create cache folder at %v\n", cacheFolder)
	}


	// Emtpty cache
	emptyCache := &Cache{
		MaxSize:       10 * 1024 * 1024,
		CleanInterval: 5,
		Folder:        cacheFolder,
		Storage:       make([]*models.Image, 0),
		Logger:        logger,
		lock:          &sync.RWMutex{},
	}
	if err != nil{
		t.Errorf("RemoveOldest() Cannot create emptyCache instance:%v", err)
	}
	// Emtpty cache
	oneElemCache := &Cache{
		MaxSize:       10 * 1024 * 1024,
		CleanInterval: 5,
		Folder:        cacheFolder,
		Storage:       make([]*models.Image, 0),
		Logger:        logger,
		lock:          &sync.RWMutex{},
	}
	if err != nil{
		t.Errorf("RemoveOldest() Cannot create oneElemCache instance:%v", err)
	}
	// Emtpty cache
	twoElemCache := &Cache{
		MaxSize:       10 * 1024 * 1024,
		CleanInterval: 5,
		Folder:        cacheFolder,
		Storage:       make([]*models.Image, 0),
		Logger:        logger,
		lock:          &sync.RWMutex{},
	}
	if err != nil{
		t.Errorf("RemoveOldest() Cannot create twoElemCache instance:%v", err)
	}

	// Img for first test
	firstImage := &models.Image{
		Name:       "firstImage.jpg",
		MimeType:   "image/jpeg",
		Url:        "url_first",
		Size:       1 * 1024 * 1024,
		Headers:    nil,
		FetchCount: 0,
	}

	// Img for second test
	secondImage := &models.Image{
		Name:       "secondImage.jpg",
		MimeType:   "image/jpeg",
		Url:        "url_second",
		Size:       1 * 1024 * 1024,
		Headers:    nil,
		FetchCount: 0,
	}

	err = oneElemCache.Add(firstImage)
	if err != nil{
		t.Errorf("RemoveOldest() Cannot add firstImage to oneElemCache cache:%v", err)
	}
	err = twoElemCache.Add(firstImage)
	if err != nil{
		t.Errorf("RemoveOldest() Cannot add firstImage to twoElemCache cache:%v", err)
	}
	err = twoElemCache.Add(secondImage)
	if err != nil{
		t.Errorf("RemoveOldest() Cannot add secondImage to twoElemCache cache:%v", err)
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
		// Empty cache
		// Cache with 1 elem
		// Cache with 1+ different elems
		{name: "Remove oldest in EMPTY cache", fields: fields{Cache: emptyCache}, wantErr: false},
		{name: "Remove oldest in one elem cache", fields: fields{Cache: oneElemCache}, wantErr: false},
		{name: "Remove oldest in many elem cache", fields: fields{Cache: twoElemCache}, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc := tt.fields.Cache
			if err := cc.RemoveOldest(); (err != nil) != tt.wantErr {
				t.Errorf("RemoveOldest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}