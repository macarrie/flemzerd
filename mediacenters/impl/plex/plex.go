package plex

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	plex "github.com/macarrie/go-plex-client"

	"github.com/macarrie/flemzerd/configuration"
	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
)

type PlexMediaCenter struct {
	Client *plex.Plex
}

var module Module

func New() (p *PlexMediaCenter, err error) {
	module = Module{
		Name: "plex",
		Type: "mediacenter",
		Status: ModuleStatus{
			Alive: false,
		},
	}

	p = &PlexMediaCenter{}
	client, err := plex.New(fmt.Sprintf("%s:%s", configuration.Config.MediaCenters["plex"]["address"], configuration.Config.MediaCenters["plex"]["port"]), configuration.Config.MediaCenters["plex"]["token"])
	if err != nil {
		p.Client = nil
		msg := fmt.Sprintf("Cannot connect to plex mediacenter: %s", err.Error())
		module.Status.Message = msg
		return p, errors.New(msg)
	}

	p.Client = client
	module.Status.Alive = true

	return p, nil
}

func (p *PlexMediaCenter) Status() (Module, error) {
	log.Debug("Checking plex mediacenter status")

	if p.Client == nil {
		client, err := plex.New(fmt.Sprintf("%s:%s", configuration.Config.MediaCenters["plex"]["address"], configuration.Config.MediaCenters["plex"]["port"]), configuration.Config.MediaCenters["plex"]["token"])
		if err != nil || client == nil {
			module.Status.Alive = false
			module.Status.Message = "Could not connect to plex: no client"
			return module, errors.New(module.Status.Message)
		}
		p.Client = client
	}

	success, err := p.Client.Test()
	if err != nil {
		module.Status.Alive = false
		module.Status.Message = err.Error()
		return module, err
	}
	if !success {
		module.Status.Alive = false
		module.Status.Message = "Could not connect to Plex media server"
		return module, errors.New(module.Status.Message)
	}

	module.Status.Alive = true
	module.Status.Message = ""

	return module, nil
}

func (p *PlexMediaCenter) GetName() string {
	return "plex"
}

func (p *PlexMediaCenter) RefreshLibrary() error {
	log.Debug("Refreshing plex library")

	libs, err := p.GetLibraries()
	if err != nil {
		return errors.Wrap(err, "could not list plex libraries to refresh")
	}

	var errorList *multierror.Error
	for _, lib := range libs {
		if err := p.Client.ScanLibrary(lib.Key); err != nil {
			errorList = multierror.Append(errorList, errors.Wrap(err, fmt.Sprintf("could not refresh library '%s' (id %s)", lib.Title, lib.Key)))
		}
	}

	return errorList.ErrorOrNil()
}

func (p *PlexMediaCenter) GetLibraries() ([]plex.Directory, error) {
	if p.Client == nil {
		return []plex.Directory{}, errors.New("nil plex client")
	}

	libs, err := p.Client.GetLibraries()
	if err != nil {
		return []plex.Directory{}, errors.Wrap(err, "cannot get library list from plex server")
	}

	return libs.MediaContainer.Directory, nil
}
