package main

import (
	cfg "ImageCutter/pkg/config"
	logging "ImageCutter/pkg/logger"
	"ImageCutter/pkg/services/cutter"
	"log"
	"os"
	"strconv"
)

func main() {
	path, _ := os.Getwd()
	_ = path

	// Read config
	config, err := cfg.ReadConfig()
	if err != nil {
		log.Fatalf("Reading cutter config give error: %v\n", err)
	}
	envPort := os.Getenv("PORT")
	envCacheSize := os.Getenv("CACHESIZE")
	envCacheClean := os.Getenv("CACHECLEAN")
	envCacheFolder := os.Getenv("CACHEFOLDER")

	// Replace config settings by env settings if they are not nil
	if envPort != "" {
		port, err := strconv.Atoi(envPort)
		if err != nil {
			log.Fatalf("Cannot convert env var PORT: %v to int, err: %v", envPort, err)
		}
		config.Cutter.Port = port
	}
	if envCacheSize != "" {
		cacheSize, err := strconv.ParseInt(envCacheSize, 10, 64)
		if err != nil {
			log.Fatalf("Cannot convert env var CACHESIZE: %v to int, err: %v", cacheSize, err)
		}
		config.Cutter.Cache.Size = cacheSize
	}
	if envCacheClean != "" {
		cacheClean, err := strconv.Atoi(envCacheClean)
		if err != nil {
			log.Fatalf("Cannot convert env var CACHECLEAN: %v to int, err: %v", envCacheClean, err)
		}
		config.Cutter.Cache.CleanInterval = cacheClean
	}
	if envCacheFolder != ""{
		config.Cutter.Cache.Folder = envCacheFolder
	}


	// Create logger
	logger, err := logging.CreateLogger(&config.Cutter.Logger)
	if err != nil {
		log.Fatalf("Creating cutter logger give error: %v\n", err)
	}

	// Create cutter instance
	cutterService, err := cutter.NewCutterService(logger, config)
	if err != nil {
		log.Fatalf("Creating cutter service instance give error: %v\n", err)
	}

	// Start cutter listener
	cutterService.Start()
}