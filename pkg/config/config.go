package config

import (
	"github.com/spf13/viper"
	"log"
)
type Logger struct {
	Level            string   `mapstructure:"level"`
	Encoding         string   `mapstructure:"encoding"`
	OutputPaths      []string `mapstructure:"outputPaths"`
	ErrorOutputPaths []string `mapstructure:"errorOutputPaths"`
}
type Cache struct {
	Size            int64   `mapstructure:"size"`
	Folder         string   `mapstructure:"folder"`
	CleanInterval int `mapstructure:"cleantime"`
}

type CutterConfig struct {
	Cutter struct {
		Port   int `mapstructure:"Port"`
		Cache Cache `mapstructure:"Cache"`
		Logger Logger `mapstructure:"Logger"`
	} `mapstructure:"Cutter"`
}

func ReadConfig() (conf *CutterConfig, err error) {
	viper.SetConfigName("cutter")
	viper.AddConfigPath("configs")
	viper.AddConfigPath("../../configs")
	viper.AddConfigPath("../configs")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		log.Printf("Reading cutter config error: %v \n", err)
		return nil, err
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		log.Printf("Unmarshaling cutter config error: %v \n", err)
		return nil, err
	}

	return conf, nil
}