package kodi_notifier

import (
	"fmt"

	"github.com/macarrie/flemzerd/configuration"
	kodi_helper "github.com/macarrie/flemzerd/helpers/kodi"
	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"

	"github.com/pdf/kodirpc"
	"github.com/pkg/errors"
)

type KodiNotifier struct {
	Client *kodirpc.Client
}

var module Module

func New() (k *KodiNotifier, err error) {
	module = Module{
		Name: "kodi",
		Type: "notifier",
		Status: ModuleStatus{
			Alive: false,
		},
	}

	k = &KodiNotifier{}

	client, err := kodi_helper.CreateKodiClient(configuration.Config.MediaCenters["kodi"]["address"], configuration.Config.MediaCenters["kodi"]["port"])
	if err != nil {
		k.Client = nil
		msg := fmt.Sprintf("Cannot connect to kodi mediacenter: %s", err.Error())
		module.Status.Message = msg
		return k, errors.Wrap(err, "cannot connect to kodi mediacenter")
	}

	k.Client = client
	module.Status.Alive = true

	return k, nil
}

func (k *KodiNotifier) Status() (Module, error) {
	log.Debug("Checking kodi notifier status")

	if k.Client == nil {
		client, err := kodi_helper.CreateKodiClient(configuration.Config.MediaCenters["kodi"]["address"], configuration.Config.MediaCenters["kodi"]["port"])
		if err != nil || client == nil {
			module.Status.Alive = false
			module.Status.Message = "Could not connect to kodi: no client"
			return module, errors.New(module.Status.Message)
		}
		k.Client = client
	}

	if _, err := k.Client.Call("JSONRPC.Ping", nil); err != nil {
		module.Status.Alive = false
		module.Status.Message = err.Error()

		return module, errors.Wrap(err, "cannot ping kodi mediacenter")
	}

	module.Status.Alive = true
	module.Status.Message = ""

	return module, nil
}

func (k *KodiNotifier) GetName() string {
	return "kodi"
}

func (k *KodiNotifier) Send(notif Notification) error {
	log.Debug("Sending Kodi notification")

	if k.Client == nil {
		return errors.New("Could not contact kodi server to send notification")
	}

	title := ""
	content := ""

	switch notif.Type {
	case NOTIFICATION_NEW_EPISODE:
		title = fmt.Sprintf("%v S%03dE%03d ", notif.Episode.TvShow.OriginalName, notif.Episode.Season, notif.Episode.Number)
		content = "Episode aired"

	case NOTIFICATION_NEW_MOVIE:
		title = fmt.Sprintf("%s", notif.Movie.OriginalTitle)
		content = "Movie found in watchlists"

	case NOTIFICATION_DOWNLOAD_START:
		if notif.Episode.ID != 0 {
			title = fmt.Sprintf("%v S%03dE%03d", notif.Episode.TvShow.OriginalName, notif.Episode.Season, notif.Episode.Number)
		}
		if notif.Movie.ID != 0 {
			title = fmt.Sprintf("%v", notif.Movie.OriginalTitle)
		}
		content = "Starting download"

	case NOTIFICATION_DOWNLOAD_SUCCESS:
		if notif.Episode.ID != 0 {
			title = fmt.Sprintf("%v S%03dE%03d", notif.Episode.TvShow.OriginalName, notif.Episode.Season, notif.Episode.Number)
			content = "Episode downloaded"
		}
		if notif.Movie.ID != 0 {
			title = fmt.Sprintf("%v", notif.Movie.OriginalTitle)
			content = "Movie downloaded"
		}

	case NOTIFICATION_DOWNLOAD_FAILURE:
		if notif.Episode.ID != 0 {
			title = fmt.Sprintf("%v S%03dE%03d", notif.Episode.TvShow.OriginalName, notif.Episode.Season, notif.Episode.Number)
		}
		if notif.Movie.ID != 0 {
			title = fmt.Sprintf("%v", notif.Movie.OriginalTitle)
		}
		content = "Download failed"

	case NOTIFICATION_NO_TORRENTS:
		if notif.Episode.ID != 0 {
			title = fmt.Sprintf("%v S%03dE%03d", notif.Episode.TvShow.OriginalName, notif.Episode.Season, notif.Episode.Number)
		}
		if notif.Movie.ID != 0 {
			title = fmt.Sprintf("%v", notif.Movie.OriginalTitle)
		}
		content = "No torrents found"

	default:
		return fmt.Errorf("Unable to send notification: Unknown notification type (%d)", notif.Type)
	}

	params := map[string]interface{}{
		"title":       title,
		"message":     content,
		"displaytime": 5000,
	}

	if err := k.Client.Notify("GUI.ShowNotification", params); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warning("Unable to send notification")
		return errors.Wrap(err, "unable to send kodi notification")
	}

	return nil
}
