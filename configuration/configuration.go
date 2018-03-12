package configuration

import (
	"path/filepath"

	log "github.com/macarrie/flemzerd/logging"
	"github.com/spf13/viper"
)

var customConfigFilePath string
var customConfigFile bool

var Config Configuration

// YAML tags do not work in viper, mapstructure tags have to be used instead
// https://github.com/spf13/viper/issues/385
type Configuration struct {
	Interface struct {
		Enabled bool
		Port    int
	}
	Providers     map[string]map[string]string
	Indexers      map[string][]map[string]string
	Notifiers     map[string]map[string]string
	Downloaders   map[string]map[string]string
	Watchlists    map[string]interface{}
	Notifications struct {
		Enabled                bool `mapstructure:"enabled"`
		NotifyNewEpisode       bool `mapstructure:"notify_new_episode"`
		NotifyDownloadComplete bool `mapstructure:"notify_download_complete"`
		NotifyFailure          bool `mapstructure:"notify_failure"`
	}
	System struct {
		EpisodeCheckInterval         int `mapstructure:"episode_check_interval"`
		TorrentDownloadAttemptsLimit int `mapstructure:"torrent_download_attempts_limit"`
	}
	Library struct {
		ShowPath  string `mapstructure:"show_path"`
		MoviePath string `mapstructure:"movie_path"`
	}
}

func setDefaultValues() {
	viper.SetDefault("interface.enabled", true)
	viper.SetDefault("interface.port", 8080)

	viper.SetDefault("notifications.enabled", true)
	viper.SetDefault("notifications.notify_new_episode", true)
	viper.SetDefault("notifications.notify_download_complete", true)
	viper.SetDefault("notifications.notify_failure", true)

	viper.SetDefault("library.show_path", "/var/lib/flemzerd/library/shows")
	viper.SetDefault("library.movie_path", "/var/lib/flemzerd/library/movies")

	viper.SetDefault("system.episode_check_interval", 15)
	viper.SetDefault("torrent_download_attempts_limit", 20)
}

func UseFile(filePath string) {
	log.WithFields(log.Fields{
		"file": filePath,
	}).Info("Using specified configuration file")
	customConfigFile = true
	customConfigFilePath = filePath
}

func Check() {
	if len(Config.Watchlists) == 0 {
		log.WithFields(log.Fields{
			"error": "No watchlists defined",
		}).Error("Configuration error")
	}

	if len(Config.Providers) == 0 {
		log.WithFields(log.Fields{
			"error": "No providers defined",
		}).Fatal("Configuration error")
	}

	_, tvdb := Config.Providers["tvdb"]
	if tvdb {
		_, tvdbapikey := Config.Providers["tvdb"]["apikey"]

		if !tvdbapikey {
			log.WithFields(log.Fields{
				"error": "Missing key(s)for tvdb provider (apikey required)",
			}).Error("Configuration error")
		}
	}

	if len(Config.Notifiers) == 0 && Config.Notifications.Enabled {
		log.WithFields(log.Fields{
			"error": "Notifications are enabled but no notifiers are defined. No notifications will be sent.",
		}).Warning("Configuration warning")
	}

	_, pushbullet := Config.Notifiers["pushbullet"]
	if pushbullet {
		_, pushbulletaccesstoken := Config.Notifiers["pushbullet"]["accesstoken"]

		if !pushbulletaccesstoken {
			log.WithFields(log.Fields{
				"error": "Missing key for pushbullet notifier (accessToken required)",
			}).Error("Configuration error")
		}
	}

	if !filepath.IsAbs(Config.Library.ShowPath) {
		log.WithFields(log.Fields{
			"error": "Library show path must be an absolute path",
		}).Error("Configuration error")
	}
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

	setDefaultValues()

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
