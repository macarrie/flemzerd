package jellyfin

import (
	"crypto/tls"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/macarrie/flemzerd/configuration"
	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
)

type JellyfinMediaCenter struct {
	Address string
	Port    int
	Apikey  string
}

var module Module

func New() (k *JellyfinMediaCenter, err error) {
	module = Module{
		Name: "jellyfin",
		Type: "mediacenter",
		Status: ModuleStatus{
			Alive: false,
		},
	}
	port, err := strconv.Atoi(configuration.Config.MediaCenters["jellyfin"]["port"])
	if err != nil {
		log.WithFields(log.Fields{
			"error":            err,
			"mediacenter":      "jellyfin",
			"port_from_config": configuration.Config.MediaCenters["jellyfin"]["port"],
		}).Warning("Impossible to use port specified in configuration. Using port 8096 instead")
		port = 8096
	}

	k = &JellyfinMediaCenter{
		Address: configuration.Config.MediaCenters["jellyfin"]["address"],
		Port:    port,
		Apikey:  configuration.Config.MediaCenters["jellyfin"]["apikey"],
	}

	module.Status.Alive = true

	return k, nil
}

func (j *JellyfinMediaCenter) Status() (Module, error) {
	log.Debug("Checking jellyfin mediacenter status")

	resp, err := apiRequest("GET", j.Address, j.Port, j.Apikey, "/System/Info")
	if err != nil {
		module.Status.Alive = false
		module.Status.Message = err.Error()
		return module, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			module.Status.Alive = false
			module.Status.Message = err.Error()

			return module, errors.Wrap(err, "could not read response content")
		}

		module.Status.Alive = false
		module.Status.Message = string(body)

		return module, err
	}

	module.Status.Alive = true
	module.Status.Message = ""

	return module, nil
}

func (j *JellyfinMediaCenter) GetName() string {
	return "jellyfin"
}

func (j *JellyfinMediaCenter) RefreshLibrary() error {
	log.Debug("Refreshing jellyfin library")

	resp, err := apiRequest("POST", j.Address, j.Port, j.Apikey, "/emby/Library/Refresh")
	if err != nil {
		return errors.Wrap(err, "could not perform jellyfin API request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "could not read response content")
		}

		return fmt.Errorf("could not refresh jellyfin library (http %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func apiRequest(method string, address string, port int, apikey string, endpoint string) (*http.Response, error) {
	baseURL := fmt.Sprintf("%s:%d/%s", address, port, endpoint)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(HTTP_TIMEOUT * time.Second),
	}

	var request *http.Request

	request, err := http.NewRequest(method, baseURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error while constructing HTTP request")
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("X-Emby-Token", apikey)
	request.Close = true

	return httpClient.Do(request)
}
