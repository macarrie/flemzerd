package configuration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/macarrie/flemzerd/logging"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"golang.org/x/sys/unix"
)

var Version string = "dev"

var customConfigFilePath string
var customConfigFile bool

var Config Configuration

var ConfigurationFatal func(code int) = os.Exit

// SECRET KEYS FROM ENV
var TRAKT_CLIENT_SECRET string
var TELEGRAM_BOT_TOKEN string
var TMDB_API_KEY string
var TVDB_API_KEY string

type ConfigurationError struct {
	File string
	Err  error
}

func (e ConfigurationError) Error() string {
	return fmt.Sprintf("%s (%s)", e.Err.Error(), e.File)
}

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
	MediaCenters  map[string]map[string]string
	Notifications struct {
		Enabled                bool `mapstructure:"enabled"`
		NotifyNewEpisode       bool `mapstructure:"notify_new_episode"`
		NotifyNewMovie         bool `mapstructure:"notify_new_movie"`
		NotifyDownloadStart    bool `mapstructure:"notify_download_start"`
		NotifyDownloadComplete bool `mapstructure:"notify_download_complete"`
		NotifyFailure          bool `mapstructure:"notify_failure"`
	}
	System struct {
		CheckInterval                int    `mapstructure:"check_interval"`
		TorrentDownloadAttemptsLimit int    `mapstructure:"torrent_download_attempts_limit"`
		TrackShows                   bool   `mapstructure:"track_shows"`
		TrackMovies                  bool   `mapstructure:"track_movies"`
		AutomaticShowDownload        bool   `mapstructure:"automatic_show_download"`
		AutomaticMovieDownload       bool   `mapstructure:"automatic_movie_download"`
		ShowDownloadDelay            int    `mapstructure:"show_download_delay"`
		MovieDownloadDelay           int    `mapstructure:"movie_download_delay"`
		PreferredMediaQuality        string `mapstructure:"preferred_media_quality"`
		ExcludedReleaseTypes         string `mapstructure:"excluded_release_types"`
	}
	Library struct {
		ShowPath      string `mapstructure:"show_path"`
		MoviePath     string `mapstructure:"movie_path"`
		CustomTmpPath string `mapstructure:"custom_tmp_path"`
	}
	Version string
}

func setDefaultValues() {
	viper.SetDefault("interface.enabled", true)
	viper.SetDefault("interface.port", 8080)

	viper.SetDefault("notifications.enabled", true)
	viper.SetDefault("notifications.notify_new_episode", true)
	viper.SetDefault("notifications.notify_new_movie", true)
	viper.SetDefault("notifications.notify_download_start", true)
	viper.SetDefault("notifications.notify_download_complete", true)
	viper.SetDefault("notifications.notify_failure", true)

	viper.SetDefault("library.show_path", "/var/lib/flemzerd/library/shows")
	viper.SetDefault("library.movie_path", "/var/lib/flemzerd/library/movies")
	viper.SetDefault("library.custom_tmp_path", "/var/lib/flemzerd/tmp")

	viper.SetDefault("system.check_interval", 15)
	viper.SetDefault("system.torrent_download_attempts_limit", 20)
	viper.SetDefault("system.track_shows", true)
	viper.SetDefault("system.track_movies", true)
	viper.SetDefault("system.automatic_show_download", true)
	viper.SetDefault("system.automatic_movie_download", false)
	viper.SetDefault("system.show_download_delay", 1)
	viper.SetDefault("system.movie_download_delay", 7)
	viper.SetDefault("system.preferred_media_quality", "720p")
	viper.SetDefault("system.excluded_release_types", "cam,screener,telesync,telecine")
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
		}).Error("Configuration error")
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

	err := unix.Access(Config.Library.ShowPath, unix.W_OK)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"path":  Config.Library.ShowPath,
		}).Error("Cannot write into library show path. Downloaded show episodes will not be able to be moved in library folder and will stay in temporary folder")
	}
	err = unix.Access(Config.Library.MoviePath, unix.W_OK)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"path":  Config.Library.MoviePath,
		}).Error("Cannot write into library movie path. Downloaded movies will not be able to be moved in library folder and will stay in temporary folder")
	}
	err = unix.Access(Config.Library.CustomTmpPath, unix.W_OK)
	if err != nil {
		log.WithFields(log.Fields{
			"path":  Config.Library.CustomTmpPath,
			"error": err,
		}).Error("Cannot write into tmp path. Media will not be able to be downloaded.")
	}

	_, kodi := Config.MediaCenters["kodi"]
	if kodi {
		_, kodiAddress := Config.MediaCenters["kodi"]["address"]
		if !kodiAddress {
			log.Warning("Kodi mediacenter address not defined. Using 'localhost'")
			Config.MediaCenters["kodi"]["address"] = "localhost"
		}
		_, kodiPort := Config.MediaCenters["kodi"]["port"]
		if !kodiPort {
			log.Warning("Kodi mediacenter port not defined. Using '9090'")
			Config.MediaCenters["kodi"]["port"] = "9090"
		}
	}

	qualityFilters := strings.Split(Config.System.PreferredMediaQuality, ",")
	for i := range qualityFilters {
		qualityFilters[i] = strings.TrimSpace(qualityFilters[i])
	}
	for _, filter := range qualityFilters {
		switch filter {
		case "480p", "576p", "720p", "900p", "1080p", "1440p", "2160p", "5k", "8k", "16k":
		default:
			log.WithFields(log.Fields{
				"preferred_media_quality": Config.System.PreferredMediaQuality,
				"invalid_quality_setting": filter,
			}).Error("Invalid media quality preference parameter. No filter will be done on quality")
			Config.System.PreferredMediaQuality = ""
		}
	}

	releaseTypeFilters := strings.Split(Config.System.ExcludedReleaseTypes, ",")
	for i := range qualityFilters {
		releaseTypeFilters[i] = strings.TrimSpace(releaseTypeFilters[i])
	}
	for _, releaseType := range releaseTypeFilters {
		switch releaseType {
		case "cam", "telesync", "telecine", "screener", "dvdrip", "hdtv", "webdl", "blurayrip":
		default:
			log.WithFields(log.Fields{
				"excluded_release_types":       Config.System.ExcludedReleaseTypes,
				"invalid_release_type_setting": releaseType,
			}).Error("Invalid release type preference parameter. No filter will be done on release types")
			Config.System.ExcludedReleaseTypes = ""
		}
	}
}

func Load() error {
	viper.SetConfigType("toml")

	if customConfigFile {
		viper.SetConfigFile(customConfigFilePath)
	} else {
		viper.SetConfigName("flemzerd")
		viper.AddConfigPath("$HOME/.config/flemzerd/")
		viper.AddConfigPath("/etc/flemzerd/")
	}

	viper.SetEnvPrefix("FLZ")
	viper.BindEnv("TRAKT_CLIENT_SECRET")
	viper.BindEnv("TELEGRAM_BOT_TOKEN")
	viper.BindEnv("TMDB_API_KEY")
	viper.BindEnv("TVDB_API_KEY")

	if key := viper.GetString("TRAKT_CLIENT_SECRET"); key != "" {
		TRAKT_CLIENT_SECRET = key
	}
	if key := viper.GetString("TELEGRAM_BOT_TOKEN"); key != "" {
		TELEGRAM_BOT_TOKEN = viper.GetString("TELEGRAM_BOT_TOKEN")
	}
	if key := viper.GetString("TMDB_API_KEY"); key != "" {
		TMDB_API_KEY = viper.GetString("TMDB_API_KEY")
	}
	if key := viper.GetString("TVDB_API_KEY"); key != "" {
		TVDB_API_KEY = viper.GetString("TVDB_API_KEY")
	}

	readErr := viper.ReadInConfig()
	if readErr != nil {
		return ConfigurationError{
			File: viper.ConfigFileUsed(),
			Err:  errors.Wrap(readErr, "cannot read configuration file"),
		}
	}

	setDefaultValues()

	var conf Configuration
	unmarshalError := viper.Unmarshal(&conf)
	if unmarshalError != nil {
		return ConfigurationError{
			File: viper.ConfigFileUsed(),
			Err:  errors.Wrap(unmarshalError, "cannot parse configuration file"),
		}
	}

	conf.Version = Version

	log.WithFields(log.Fields{
		"file": viper.ConfigFileUsed(),
	}).Info("Configuration file loaded")

	Config = conf
	return nil
}
