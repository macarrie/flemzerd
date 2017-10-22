package transmission

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	log "flemzerd/logging"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type TransmissionDownloader struct {
	Address   string
	Port      int
	User      string
	Password  string
	SessionId string
}

func New(address string, port int, user string, password string) *TransmissionDownloader {
	return &TransmissionDownloader{
		Address:  address,
		Port:     port,
		User:     user,
		Password: password,
	}
}

func (downloader *TransmissionDownloader) Init() error {
	params := map[string]string{"pouet": "pouet"}
	response, _ := performAPIRequest("torrent-get", downloader.Address, downloader.Port, "", params)
	if response.StatusCode == 409 {
		downloader.SessionId = response.Header["X-Transmission-Session-Id"][0]
	}

	return nil
}

func (downloader TransmissionDownloader) AddTorrent(url string) error {
	log.Debug("Add torrent: ", url)
	params := map[string]string{"pouet": "lol"}
	response, content := performAPIRequest("torrent-get", downloader.Address, downloader.Port, downloader.SessionId, params)
	log.Debug(fmt.Sprintf("%+v\n", response))
	log.Debug(string(content))
	return nil
}

func performAPIRequest(method string, host string, port int, sessionId string, paramsMap map[string]string) (http.Response, []byte) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}

	urlObject := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", host, port),
		Path:   "/transmission/rpc",
	}
	log.Debug("URL String: ", urlObject.String())

	var request *http.Request
	var err error

	jsonParams, _ := json.Marshal(paramsMap)

	request, err = http.NewRequest("POST", urlObject.String(), bytes.NewReader(jsonParams))
	if err != nil {
		log.Fatal("NewRequest: ", err)
	}

	request.Header.Set("X-Transmission-Session-Id", sessionId)

	response, err := httpClient.Do(request)
	if err != nil {
		log.Fatal("API Request: ", err)
	}

	body, readError := ioutil.ReadAll(response.Body)
	if readError != nil {
		log.Fatal("API Response read: ", readError)
	}

	return *response, body
}
