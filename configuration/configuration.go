package configuration

import (
	log "flemzerd/logging"
	"github.com/spf13/viper"
)

var customConfigFilePath string
var customConfigFile bool

type Configuration struct {
	Providers map[string]map[string]string
	Notifiers map[string]map[string]string
	Shows     []string
}

func UseFile(filePath string) {
	log.WithFields(log.Fields{
		"file": filePath,
	}).Info("Using specified configuration file")
	customConfigFile = true
	customConfigFilePath = filePath
}

func Load() (Configuration, error) {
	viper.SetConfigType("yaml")

	if customConfigFile {
		viper.SetConfigFile(customConfigFilePath)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
	}

	readErr := viper.ReadInConfig()
	if readErr != nil {
		return Configuration{}, readErr
	}

	var config Configuration
	unmarshalError := viper.Unmarshal(&config)
	if unmarshalError != nil {
		return Configuration{}, unmarshalError
	}

	log.WithFields(log.Fields{
		"file": viper.ConfigFileUsed(),
	}).Info("Configuration file loaded")

	return config, nil
}
