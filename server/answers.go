package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var flagMap map[string]interface{}

func configChange(e fsnotify.Event) {
	flagMap = viper.GetStringMap("flags")
}

func checkFlag(submission, page string) bool {
	flag, ok := flagMap[page]
	if !ok {
		return false
	}

	if submission != flag {
		return false
	}

	return true
}

func parseTrainingFlags() {
	viper.SetConfigName("flags")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")
	viper.OnConfigChange(configChange)

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error in flags config file: %s\n", err))
	}

	flagMap = viper.GetStringMap("flags")

	viper.WatchConfig()
}
