package pushbullet

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	notifier_helper "github.com/macarrie/flemzerd/helpers/notifiers"
	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"

	"github.com/pkg/errors"
)

///////////////
// Constants //
///////////////

const baseURL = "https://api.pushbullet.com/"

/////////////////////////////////
// API Structures declarations //
/////////////////////////////////

type Token struct {
	Token string `json:"token"`
}

type PushbulletNotifier struct {
	Name string
}

type User struct {
	Active bool   `json:"active"`
	Iden   string `json:"iden"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

//////////////////
// Package vars //
//////////////////

var APIToken Token

func performAPIRequest(method string, path string, paramsMap map[string]string) (http.Response, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(HTTP_TIMEOUT * time.Second),
	}

	urlObject, _ := url.ParseRequestURI(baseURL)
	urlObject.Path = path

	var request *http.Request
	var err error

	if method == "GET" {
		params := url.Values{}
		for key, value := range paramsMap {
			params.Add(key, value)
		}
		urlObject.RawQuery = params.Encode()

		request, err = http.NewRequest(method, urlObject.String(), nil)
		if err != nil {
			return http.Response{}, err
		}
	} else if method == "POST" {
		jsonParams, _ := json.Marshal(paramsMap)

		request, err = http.NewRequest(method, urlObject.String(), bytes.NewReader(jsonParams))
		if err != nil {
			return http.Response{}, err
		}
	} else {
		log.WithFields(log.Fields{
			"method": method,
		}).Error("HTTP method unknown")

		return http.Response{}, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Access-Token", fmt.Sprintf("%s", APIToken.Token))
	request.Close = true

	response, err := httpClient.Do(request)
	if err != nil {
		return http.Response{}, err
	}
	defer response.Body.Close()

	return *response, nil
}

func (pushbulletNotifier *PushbulletNotifier) Setup(credentials map[string]string) {
	APIToken.Token = credentials["AccessToken"]
	pushbulletNotifier.Name = "Pushbullet"

	pushbulletNotifier.Init()
}

func (notifier *PushbulletNotifier) Init() error {
	return nil
	//log.Debug("Checking Pushbullet user token validity")

	//response, err := performAPIRequest("GET", "v2/users/me", nil)
	//if err != nil {
	//return err
	//}

	//body, err := ioutil.ReadAll(response.Body)
	//if err != nil {
	//return err
	//}

	//var authResult User
	//parseErr := json.Unmarshal(body, &authResult)
	//if parseErr != nil {
	//return parseErr
	//} else {
	//log.Debug("Pushbullet token valid")

	//return nil
	//}
}

func (notifier *PushbulletNotifier) Status() (Module, error) {
	returnStruct := Module{
		Name: "Pushbullet",
		Type: "notifier",
		Status: ModuleStatus{
			Alive: false,
		},
	}

	log.Debug("Checking Pushbullet notifier status")
	response, err := performAPIRequest("GET", "v2/users/me", nil)
	if err != nil {
		returnStruct.Status.Message = err.Error()
		return returnStruct, errors.Wrap(err, "error during Pushbullet API request")
	}

	if response.StatusCode != 200 {
		msg := fmt.Sprintf("Pushbullet notifier request returned %d status code", response.StatusCode)
		returnStruct.Status.Message = msg
		return returnStruct, errors.New(msg)
	}

	returnStruct.Status.Alive = true
	return returnStruct, nil
}

func (n *PushbulletNotifier) GetName() string {
	return "Pushbullet"
}

func (notifier *PushbulletNotifier) Send(notif Notification) error {
	log.Debug("Sending Pushbullet notification")

	title, content, err := notifier_helper.GetNotificationText(notif)
	if err != nil {
		return err
	}

	params := map[string]string{"type": "note", "title": title, "body": content}
	response, err := performAPIRequest("POST", "v2/pushes", params)
	if err != nil {
		return errors.Wrap(err, "error during Pushbullet API request")
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Unable to send notification: HTTP status code %d", response.StatusCode)
	}

	return nil
}

func New(credentials map[string]string) *PushbulletNotifier {
	var notifier PushbulletNotifier
	notifier.Setup(credentials)

	return &notifier
}
