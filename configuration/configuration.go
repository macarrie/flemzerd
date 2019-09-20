package configuration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-multierror"
	log "github.com/macarrie/flemzerd/logging"
	"github.com/macarrie/flemzerd/vidocq"
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

type StatusCode int

const (
	OK StatusCode = iota
	WARNING
	CRITICAL
	UNKNOWN
)

func (s StatusCode) String() string {
	display := []string{
		"OK",
		"WARNING",
		"CRITICAL",
		"UNKNOWN",
	}

	if s > CRITICAL || s < OK {
		return display[UNKNOWN]
	}

	return display[s]
}

type ConfigurationError struct {
	Status  StatusCode
	Key     string
	Value   string
	Message string
}

func (e ConfigurationError) Error() string {
	if e.Key == "" {
		return fmt.Sprintf("[%s] %s", e.Status, e.Message)
	}

	return fmt.Sprintf("[%s] %s (%s: %s)", e.Status, e.Message, e.Key, e.Value)
}

type ConfigurationFileError struct {
	File string
	Err  error
}

func (e ConfigurationFileError) Error() string {
	return fmt.Sprintf("%s (%s)", e.Err.Error(), e.File)
}

// YAML tags do not work in viper, mapstructure tags have to be used instead
// https://github.com/spf13/viper/issues/385
type Configuration struct {
	Interface struct {
		Enabled bool
		Port    int
		Auth    struct {
			Username string
			Password string
		}
		UseSSL       bool   `mapstructure:"use_ssl"`
		SSLCert      string `mapstructure:"ssl_cert"`
		SSLServerKey string `mapstructure:"ssl_server_key"`
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
		HealthcheckInterval          int    `mapstructure:"healthcheck_interval"`
		TorrentDownloadAttemptsLimit int    `mapstructure:"torrent_download_attempts_limit"`
		TrackShows                   bool   `mapstructure:"track_shows"`
		TrackMovies                  bool   `mapstructure:"track_movies"`
		AutomaticShowDownload        bool   `mapstructure:"automatic_show_download"`
		AutomaticMovieDownload       bool   `mapstructure:"automatic_movie_download"`
		ShowDownloadDelay            int    `mapstructure:"show_download_delay"`
		MovieDownloadDelay           int    `mapstructure:"movie_download_delay"`
		PreferredMediaQuality        string `mapstructure:"preferred_media_quality"`
		ExcludedReleaseTypes         string `mapstructure:"excluded_release_types"`
		StrictTorrentCheck           bool   `mapstructure:"strict_torrent_check"`
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
	viper.SetDefault("interface.auth.username", "admin")
	viper.SetDefault("interface.auth.password", "flemzerd")
	viper.SetDefault("interface.use_ssl", true)
	viper.SetDefault("interface.ssl_cert", "/var/lib/flemzerd/certs/cert.pem")
	viper.SetDefault("interface.ssl_server_key", "/var/lib/flemzerd/certs/server.key")

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
	viper.SetDefault("system.healthcheck_interval", 5)
	viper.SetDefault("system.torrent_download_attempts_limit", 20)
	viper.SetDefault("system.track_shows", true)
	viper.SetDefault("system.track_movies", true)
	viper.SetDefault("system.automatic_show_download", true)
	viper.SetDefault("system.automatic_movie_download", false)
	viper.SetDefault("system.show_download_delay", 1)
	viper.SetDefault("system.movie_download_delay", 7)
	viper.SetDefault("system.preferred_media_quality", "720p")
	viper.SetDefault("system.excluded_release_types", "cam,screener,telesync,telecine")
	viper.SetDefault("system.strict_torrent_check", true)
}

func UseFile(filePath string) {
	log.WithFields(log.Fields{
		"file": filePath,
	}).Info("Using specified configuration file")
	customConfigFile = true
	customConfigFilePath = filePath
}

func Check() error {
	var errorList *multierror.Error
	var configError ConfigurationError

	if len(Config.Watchlists) == 0 {
		configError = ConfigurationError{
			Status:  CRITICAL,
			Message: "No watchlist configured",
		}
		log.WithFields(log.Fields{
			"error": configError,
		}).Error("Configuration error")
		errorList = multierror.Append(errorList, configError)
	}
	_, traktWatchlist := Config.Watchlists["trakt"]
	if traktWatchlist {
		if TRAKT_CLIENT_SECRET == "" {
			configError = ConfigurationError{
				Status:  CRITICAL,
				Message: "No client secret defined for trakt watchlist. Set FLZ_TRAKT_CLIENT_SECRET env var",
			}
			log.WithFields(log.Fields{
				"error": configError,
			}).Error("Configuration error")
			errorList = multierror.Append(errorList, configError)
		}
	}

	if len(Config.Providers) == 0 {
		configError = ConfigurationError{
			Status:  CRITICAL,
			Message: "No provider configured",
		}
		log.WithFields(log.Fields{
			"error": configError,
		}).Error("Configuration error")
		errorList = multierror.Append(errorList, configError)
	}
	_, tmdbProvider := Config.Providers["tmdb"]
	if tmdbProvider {
		if TMDB_API_KEY == "" {
			configError = ConfigurationError{
				Status:  CRITICAL,
				Message: "No API key defined for TMDB provider. Set FLZ_TMDB_API_KEY env var",
			}
			log.WithFields(log.Fields{
				"error": configError,
			}).Error("Configuration error")
			errorList = multierror.Append(errorList, configError)
		}
	}
	_, tvdbProvider := Config.Providers["tvdb"]
	if tvdbProvider {
		if TVDB_API_KEY == "" {
			configError = ConfigurationError{
				Status:  CRITICAL,
				Message: "No API key defined for TVDB provider. Set FLZ_TVDB_API_KEY env var",
			}
			log.WithFields(log.Fields{
				"error": configError,
			}).Error("Configuration error")
			errorList = multierror.Append(errorList, configError)
		}
	}

	if len(Config.Notifiers) == 0 && Config.Notifications.Enabled {
		configError = ConfigurationError{
			Status:  WARNING,
			Message: "Notifications are enabled but no notifiers are defined. No notifications will be sent.",
		}
		log.WithFields(log.Fields{
			"error": configError,
		}).Warning("Configuration warning")
		errorList = multierror.Append(errorList, configError)
	}

	_, pushbullet := Config.Notifiers["pushbullet"]
	if pushbullet {
		_, pushbulletaccesstoken := Config.Notifiers["pushbullet"]["accesstoken"]

		if !pushbulletaccesstoken {
			configError = ConfigurationError{
				Status:  WARNING,
				Message: "No access token defined for pushbullet notifier",
				Key:     "notifiers.pushbullet.accesstoken",
				Value:   "",
			}
			log.WithFields(log.Fields{
				"error": configError,
			}).Warning("Configuration warning")
			errorList = multierror.Append(errorList, configError)
		}
	}

	_, telegram := Config.Notifiers["telegram"]
	if telegram {
		if TELEGRAM_BOT_TOKEN == "" {
			configError = ConfigurationError{
				Status:  WARNING,
				Message: "No bot token defined for telegram notifier",
			}
			log.WithFields(log.Fields{
				"error": configError,
			}).Warning("Configuration warning")
			errorList = multierror.Append(errorList, configError)
		}
	}

	if !filepath.IsAbs(Config.Library.ShowPath) {
		configError = ConfigurationError{
			Status:  CRITICAL,
			Message: "Library show path must be an absolute path",
			Key:     "library.show_path",
			Value:   Config.Library.ShowPath,
		}
		log.WithFields(log.Fields{
			"error": configError,
		}).Error("Configuration error")
		errorList = multierror.Append(errorList, configError)
	}

	err := unix.Access(Config.Library.ShowPath, unix.W_OK)
	if err != nil {
		configError = ConfigurationError{
			Status:  CRITICAL,
			Message: "Cannot write into library show path. Downloaded show episodes will not be able to be moved in library folder and will stay in temporary folder",
			Key:     "library.show_path",
			Value:   Config.Library.ShowPath,
		}
		log.WithFields(log.Fields{
			"access_error": err,
			"error":        configError,
		}).Error("Configuration error")
		errorList = multierror.Append(errorList, configError)
	}

	err = unix.Access(Config.Library.MoviePath, unix.W_OK)
	if err != nil {
		configError = ConfigurationError{
			Status:  CRITICAL,
			Message: "Cannot write into library movie path. Downloaded show episodes will not be able to be moved in library folder and will stay in temporary folder",
			Key:     "library.movie_path",
			Value:   Config.Library.MoviePath,
		}
		log.WithFields(log.Fields{
			"access_error": err,
			"error":        configError,
		}).Error("Configuration error")
		errorList = multierror.Append(errorList, configError)
	}

	err = unix.Access(Config.Library.CustomTmpPath, unix.W_OK)
	if err != nil {
		configError = ConfigurationError{
			Status:  CRITICAL,
			Message: "Cannot write into tmp path. Media will not be able to be downloaded.",
			Key:     "library.custom_tmp_path",
			Value:   Config.Library.CustomTmpPath,
		}
		log.WithFields(log.Fields{
			"access_error": err,
			"error":        configError,
		}).Error("Configuration error")
		errorList = multierror.Append(errorList, configError)
	}

	_, kodi := Config.MediaCenters["kodi"]
	if kodi {
		_, kodiAddress := Config.MediaCenters["kodi"]["address"]
		if !kodiAddress {
			configError = ConfigurationError{
				Status:  WARNING,
				Message: "Kodi mediacenter address not defined. Using 'localhost'",
				Key:     "mediacenters.kodi.address",
				Value:   "",
			}
			log.WithFields(log.Fields{
				"access_error": err,
				"error":        configError,
			}).Warning("Configuration warning")
			errorList = multierror.Append(errorList, configError)
			Config.MediaCenters["kodi"]["address"] = "localhost"
		}
		_, kodiPort := Config.MediaCenters["kodi"]["port"]
		if !kodiPort {
			configError = ConfigurationError{
				Status:  WARNING,
				Message: "Kodi mediacenter port not defined. Using '9090'",
				Key:     "mediacenters.kodi.port",
				Value:   "",
			}
			log.WithFields(log.Fields{
				"access_error": err,
				"error":        configError,
			}).Warning("Configuration warning")
			errorList = multierror.Append(errorList, configError)
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
			configError = ConfigurationError{
				Status:  WARNING,
				Message: "Invalid media quality preference parameter. No filter will be done on quality",
				Key:     "system.preferred_media_quality",
				Value:   Config.System.PreferredMediaQuality,
			}
			log.WithFields(log.Fields{
				"error":                   configError,
				"invalid_quality_setting": filter,
			}).Warning("Configuration warning")
			errorList = multierror.Append(errorList, configError)

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
			configError = ConfigurationError{
				Status:  WARNING,
				Message: "Invalid release type preference parameter. No filter will be done on release types",
				Key:     "system.excluded_release_types",
				Value:   Config.System.ExcludedReleaseTypes,
			}
			log.WithFields(log.Fields{
				"error":                        configError,
				"invalid_release_type_setting": releaseType,
			}).Warning("Configuration warning")
			errorList = multierror.Append(errorList, configError)

			Config.System.ExcludedReleaseTypes = ""
		}
	}

	vidocq.CheckVidocq()
	if !vidocq.LocalVidocqAvailable {
		configError = ConfigurationError{
			Status:  CRITICAL,
			Message: "Vidocq has not been found on the server or could not be correctly executed. Torrent search will encounter issues, especially if strict torrent checking is enabled",
		}
		log.WithFields(log.Fields{
			"error": configError,
		}).Error("Configuration error")
		errorList = multierror.Append(errorList, configError)
	}

	if Config.Interface.Enabled && Config.Interface.Auth.Username == "admin" && Config.Interface.Auth.Password == "flemzerd" {
		configError = ConfigurationError{
			Status:  WARNING,
			Message: "Auth credentials defined in configuration are the default ones. Please change them for a better security",
		}
		log.WithFields(log.Fields{
			"error": configError,
		}).Warning("Configuration warning")
		errorList = multierror.Append(errorList, configError)
	}

	return errorList.ErrorOrNil()
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
		return ConfigurationFileError{
			File: viper.ConfigFileUsed(),
			Err:  errors.Wrap(readErr, "cannot read configuration file"),
		}
	}

	setDefaultValues()

	var conf Configuration
	unmarshalError := viper.Unmarshal(&conf)
	if unmarshalError != nil {
		return ConfigurationFileError{
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
