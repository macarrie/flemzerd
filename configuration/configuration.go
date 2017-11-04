package configuration

import (
	"errors"

	log "github.com/macarrie/flemzerd/logging"
	"github.com/spf13/viper"
)

var customConfigFilePath string
var customConfigFile bool

var Config Configuration

// YAML tags do not work in viper, mapstructure tags have to be used instead
// https://github.com/spf13/viper/issues/385
type Configuration struct {
	Providers     map[string]map[string]string
	Indexers      map[string][]map[string]string
	Notifiers     map[string]map[string]string
	Downloaders   map[string]map[string]string
	Notifications struct {
		Enabled                bool `mapstructure:"enabled"`
		NotifyNewEpisode       bool `mapstructure:"notify_new_episode"`
		NotifyDownloadComplete bool `mapstructure:"notify_download_complete"`
		NotifyFailure          bool `mapstructure:"notify_failure"`
	}
	System struct {
		EpisodeCheckInterval int `mapstructure:"episode_check_interval"`
	}
	Shows []string
}

func UseFile(filePath string) {
	log.WithFields(log.Fields{
		"file": filePath,
	}).Info("Using specified configuration file")
	customConfigFile = true
	customConfigFilePath = filePath
}

func Check() error {
	if len(Config.Shows) == 0 {
		return errors.New("No Shows defined")
	}

	if len(Config.Providers) == 0 {
		return errors.New("No Providers defined")
	}

	_, tvdb := Config.Providers["tvdb"]
	if tvdb {
		_, tvdbapikey := Config.Providers["tvdb"]["apikey"]

		if !tvdbapikey {
			return errors.New("Missing key(s)for tvdb provider (apikey required)")
		}
	}

	if len(Config.Notifiers) == 0 {
		return errors.New("No Notifiers defined")
	}

	_, pushbullet := Config.Notifiers["pushbullet"]
	if pushbullet {
		_, pushbulletaccesstoken := Config.Notifiers["pushbullet"]["accesstoken"]

		if !pushbulletaccesstoken {
			return errors.New("Missing key for pushbullet notifier (accessToken required)")
		}
	}

	return nil
}

func Load() error {
	viper.SetConfigType("yaml")

	if customConfigFile {
		viper.SetConfigFile(customConfigFilePath)
	} else {
		viper.SetConfigName("flemzerd")
		viper.AddConfigPath("$HOME/.config/flemzerd/")
		viper.AddConfigPath("/etc/flemzerd/")
	}

	readErr := viper.ReadInConfig()
	if readErr != nil {
		return readErr
	}

	viper.SetDefault("notifications.enabled", true)
	viper.SetDefault("notifications.notify_new_episode", true)
	viper.SetDefault("notifications.notify_download_complete", true)
	viper.SetDefault("notifications.notify_failure", true)

	viper.SetDefault("system.episode_check_interval", 15)

	var conf Configuration
	unmarshalError := viper.Unmarshal(&conf)
	if unmarshalError != nil {
		return unmarshalError
	}

	log.WithFields(log.Fields{
		"file": viper.ConfigFileUsed(),
	}).Info("Configuration file loaded")

	Config = conf
	return nil
}
