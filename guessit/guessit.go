package guessit

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
)

const baseURL = "https://api.guessit.io/"

func GetInfo(name string) (MediaInfo, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(HTTP_TIMEOUT * time.Second),
	}

	urlObject, _ := url.ParseRequestURI(baseURL)

	var request *http.Request
	var err error

	params := url.Values{}
	params.Add("filename", name)
	urlObject.RawQuery = params.Encode()

	request, err = http.NewRequest("GET", urlObject.String(), nil)
	if err != nil {
		log.Error("[Guessit] Could not create request:", err)
		return MediaInfo{}, nil
	}

	request.Close = true

	response, err := httpClient.Do(request)
	if err != nil {
		log.Error(err)
		log.Error("[Guessit] Could not perform request:", err)
		return MediaInfo{}, nil
	}
	defer response.Body.Close()

	body, readError := ioutil.ReadAll(response.Body)
	if readError != nil {
		log.Error("[Guessit] Could not read request result:", err)
		return MediaInfo{}, readError
	}

	var mediaInfo MediaInfo
	parseErr := json.Unmarshal(body, &mediaInfo)
	if parseErr != nil {
		log.Error("[Guessit] Could not parse request result:", parseErr)
	}

	return mediaInfo, nil
}
