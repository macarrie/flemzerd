package pushbullet

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	log "github.com/macarrie/flemzerd/logging"
	//"github.com/macarrie/flemzerd/notifier"
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

type Notification struct {
	Title   string
	Content string
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
	httpClient := &http.Client{Transport: tr}

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

	response, err := httpClient.Do(request)
	if err != nil {
		return http.Response{}, err
	}

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

func (notifier *PushbulletNotifier) IsAlive() error {
	log.Debug("Checking Pushbullet notifier status")
	response, err := performAPIRequest("GET", "v2/users/me", nil)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Pushbullet request return %d status code", response.StatusCode))
	}

	return nil
}

func (notifier *PushbulletNotifier) Send(title, content string) error {
	log.WithFields(log.Fields{
		"title":   title,
		"content": content,
	}).Debug("Sending Pushbullet notification")

	params := map[string]string{"type": "note", "title": title, "body": content}
	response, err := performAPIRequest("POST", "v2/pushes", params)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"http_error": response.StatusCode,
		}).Warning("Unable to send notification")
	}

	return nil
}

func New(credentials map[string]string) *PushbulletNotifier {
	var notifier PushbulletNotifier
	notifier.Setup(credentials)

	return &notifier
}
