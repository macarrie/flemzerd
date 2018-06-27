package kodi_notifier

import (
	"errors"
	"fmt"

	"github.com/macarrie/flemzerd/configuration"
	kodi_helper "github.com/macarrie/flemzerd/helpers/kodi"
	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"

	"github.com/pdf/kodirpc"
)

type KodiNotifier struct {
	Client *kodirpc.Client
}

var module Module

func New() (k *KodiNotifier, err error) {
	k = &KodiNotifier{}

	client, err := kodi_helper.CreateKodiClient(configuration.Config.MediaCenters["kodi"]["address"], configuration.Config.MediaCenters["kodi"]["port"])
	if err != nil {
		k.Client = nil
		return k, fmt.Errorf("Cannot connect to kodi mediacenter: %s", err.Error())
	}
	k.Client = client

	module = Module{
		Name: "kodi",
		Type: "notifier",
		Status: ModuleStatus{
			Alive:   true,
			Message: "",
		},
	}

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

	_, err := k.Client.Call("JSONRPC.Ping", nil)
	if err != nil {
		module.Status.Alive = false
		module.Status.Message = err.Error()
	}

	module.Status.Alive = true
	module.Status.Message = ""

	return module, nil
}

func (k *KodiNotifier) GetName() string {
	return "kodi"
}

func (k *KodiNotifier) Send(title, content string) error {
	log.WithFields(log.Fields{
		"title":   title,
		"content": content,
	}).Debug("Sending Kodi notification")

	if k.Client == nil {
		return errors.New("Could not contact kodi server to send notification")
	}

	params := map[string]interface{}{
		"title":       title,
		"message":     content,
		"displaytime": 5000,
	}

	err := k.Client.Notify("GUI.ShowNotification", params)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warning("Unable to send notification")
		return err
	}

	return nil
}
