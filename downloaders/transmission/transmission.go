package transmission

import (
	"errors"
	"fmt"
	//log "flemzerd/logging"
	"bytes"
	"crypto/tls"
	"encoding/json"
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

type TransmissionRequestParams struct {
	Method    string            `json:"method"`
	Arguments map[string]string `json:"arguments"`
}

type TransmissionRequestResponse struct {
	Result string `json:"result"`
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
	params := map[string]string{"filename": url}
	response, err := performAPIRequest("torrent-add", downloader.Address, downloader.Port, downloader.SessionId, params)
	if err != nil {
		return err
	}

	body, readError := ioutil.ReadAll(response.Body)
	if readError != nil {
		return errors.New(fmt.Sprintf("API Response read: %v", readError))
	}

	var result TransmissionRequestResponse
	parseErr := json.Unmarshal(body, &result)
	if parseErr != nil {
		return parseErr
	}
	if result.Result != "success" {
		return errors.New(result.Result)
	}

	return nil
}

func performAPIRequest(method string, host string, port int, sessionId string, paramsMap map[string]string) (http.Response, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}

	urlObject := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", host, port),
		Path:   "/transmission/rpc",
	}

	var request *http.Request
	var err error

	requestContent := TransmissionRequestParams{
		Method:    method,
		Arguments: paramsMap,
	}
	jsonParams, _ := json.Marshal(requestContent)

	request, err = http.NewRequest("POST", urlObject.String(), bytes.NewReader(jsonParams))
	if err != nil {
		return http.Response{}, err
	}

	request.Header.Set("X-Transmission-Session-Id", sessionId)

	response, err := httpClient.Do(request)
	if err != nil {
		return http.Response{}, err
	}

	return *response, nil
}
