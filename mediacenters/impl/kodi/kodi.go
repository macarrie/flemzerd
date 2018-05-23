package kodi

import (
	"errors"
	"fmt"
	"time"

	"github.com/macarrie/flemzerd/configuration"
	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"

	"github.com/pdf/kodirpc"
)

type KodiMediaCenter struct {
	Client *kodirpc.Client
}

var module Module

func New() (k *KodiMediaCenter, err error) {
	k = &KodiMediaCenter{}

	address := fmt.Sprintf("%s:%s", configuration.Config.MediaCenters["kodi"]["address"], configuration.Config.MediaCenters["kodi"]["port"])
	client, err := kodirpc.NewClient(address, &kodirpc.Config{
		ReadTimeout:         time.Duration(HTTP_TIMEOUT * time.Second),
		ConnectTimeout:      time.Duration(HTTP_TIMEOUT * time.Second),
		Reconnect:           false,
		ConnectBackoffScale: 1,
	})
	if err != nil {
		k.Client = nil
		return k, fmt.Errorf("Cannot connect to kodi mediacenter: %s", err.Error())
	}
	k.Client = client

	module = Module{
		Name: "kodi",
		Type: "mediacenter",
		Status: ModuleStatus{
			Alive:   true,
			Message: "",
		},
	}

	return k, nil
}

func (k *KodiMediaCenter) Status() (Module, error) {
	log.Debug("Checking kodi mediacenter status")

	if k.Client == nil {
		module.Status.Alive = false
		module.Status.Message = "Could not connect to kodi: no client"
		return module, errors.New(module.Status.Message)
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

func (k *KodiMediaCenter) GetName() string {
	return "kodi"
}

func (k *KodiMediaCenter) RefreshLibrary() error {
	log.Debug("Refreshing kodi library")

	if k.Client == nil {
		return errors.New("Could not connect to kodi: no client")
	}

	err := k.Client.Notify("VideoLibrary.Scan", nil)
	if err != nil {
		log.WithFields(log.Fields{
			"mediacenter": k.GetName(),
			"error":       err,
		}).Warning("Could not refresh media center library")
	}
	return nil
}
