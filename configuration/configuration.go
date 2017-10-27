package configuration

import (
	"errors"
	log "github.com/macarrie/flemzerd/logging"
	"github.com/spf13/viper"
)

var customConfigFilePath string
var customConfigFile bool

type Configuration struct {
	Providers   map[string]map[string]string
	Indexers    map[string][]map[string]string
	Notifiers   map[string]map[string]string
	Downloaders map[string]map[string]string
	Shows       []string
}

func UseFile(filePath string) {
	log.WithFields(log.Fields{
		"file": filePath,
	}).Info("Using specified configuration file")
	customConfigFile = true
	customConfigFilePath = filePath
}

func Check(config Configuration) error {
	if len(config.Shows) == 0 {
		return errors.New("No Shows defined")
	}

	if len(config.Providers) == 0 {
		return errors.New("No Providers defined")
	}

	_, tvdb := config.Providers["tvdb"]
	if tvdb {
		_, tvdbapikey := config.Providers["tvdb"]["apikey"]

		if !tvdbapikey {
			return errors.New("Missing key(s)for tvdb provider (apikey required)")
		}
	}

	if len(config.Notifiers) == 0 {
		return errors.New("No Notifiers defined")
	}

	_, pushbullet := config.Notifiers["pushbullet"]
	if pushbullet {
		_, pushbulletaccesstoken := config.Notifiers["pushbullet"]["accesstoken"]

		if !pushbulletaccesstoken {
			return errors.New("Missing key for pushbullet notifier (accessToken required)")
		}
	}

	return nil
}

func Load() (Configuration, error) {
	viper.SetConfigType("yaml")

	if customConfigFile {
		viper.SetConfigFile(customConfigFilePath)
	} else {
		viper.SetConfigName("flemzerd")
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
