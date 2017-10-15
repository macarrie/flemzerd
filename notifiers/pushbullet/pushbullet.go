package pushbullet

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	log "flemzerd/logging"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	//"flemzerd/notifier"
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

func performAPIRequest(method string, path string, paramsMap map[string]string) (int, []byte) {
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
			log.Fatal("NewRequest: ", err)
		}
	} else if method == "POST" {
		jsonParams, _ := json.Marshal(paramsMap)

		request, err = http.NewRequest(method, urlObject.String(), bytes.NewReader(jsonParams))
		if err != nil {
			log.Fatal("NewRequest: ", err)
		}
	} else {
		log.WithFields(log.Fields{
			"method": method,
		}).Fatal("HTTP method unknown")
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Access-Token", fmt.Sprintf("%s", APIToken.Token))

	response, err := httpClient.Do(request)
	if err != nil {
		log.Fatal("API Request: ", err)
	}

	body, readError := ioutil.ReadAll(response.Body)
	if readError != nil {
		log.Fatal("API Response read: ", readError)
	}

	return response.StatusCode, body
}

func (pushbulletNotifier *PushbulletNotifier) Setup(credentials map[string]string) {
	APIToken.Token = credentials["AccessToken"]
	pushbulletNotifier.Name = "Pushbullet"
}

func (notifier *PushbulletNotifier) Init() bool {
	log.Debug("Checking Pushbullet user token validity")

	_, body := performAPIRequest("GET", "v2/users/me", nil)
	//fmt.Println("Status: ", status)
	//fmt.Println("Content: ", string(body[:]))

	var authResult User
	err := json.Unmarshal(body, &authResult)
	if err != nil {
		log.Fatal(err)
		return false
	} else {
		log.Debug("Pushbullet token valid")

		return true
	}
}

func (notifier *PushbulletNotifier) Send(title, content string) error {
	log.WithFields(log.Fields{
		"title":   title,
		"content": content,
	}).Info("Sending Pushbullet notification")

	params := map[string]string{"type": "note", "title": title, "body": content}
	log.Debug("Pushbullet Mock: Notification sent", params)
	//status, _ := performAPIRequest("POST", "v2/pushes", params)

	//if status != 200 {
	//log.WithFields(log.Fields{
	//"http_error": status,
	//}).Warning("Unable to send notification")
	//}

	return nil
}

func New(credentials map[string]string) *PushbulletNotifier {
	var notifier PushbulletNotifier
	notifier.Setup(credentials)

	return &notifier
}
