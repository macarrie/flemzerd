package guessit

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
	"time"

	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
)

const baseURL = "https://api.guessit.io/"

var localGuessitAvailable bool

func init() {
	_, err := exec.LookPath("guessit")
	if err != nil {
		log.Debug("Local guessit executable not found. Requests will be performed using external API")
		localGuessitAvailable = false
	} else {
		localGuessitAvailable = true
	}
}

func makeAPIRequest(name string) (http.Response, error) {
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
		return http.Response{}, err
	}

	request.Close = true

	response, err := httpClient.Do(request)
	if err != nil {
		return http.Response{}, err
	}

	return *response, nil
}

func makeLocalRequest(name string) ([]byte, error) {
	if !localGuessitAvailable {
		return []byte{}, fmt.Errorf("guessit not available locally")
	}

	cmd := exec.Command("guessit", "--json", name)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		log.Debug("[Guessit] Could not get execute guessit: ", err)
		return []byte{}, err
	}

	return stdout.Bytes(), nil
}

func performGuessitRequest(name string) ([]byte, error) {
	var response []byte

	response, err := makeLocalRequest(name)
	if err != nil {
		httpResponse, httpErr := makeAPIRequest(name)
		if httpErr != nil {
			return []byte{}, httpErr
		}
		defer httpResponse.Body.Close()
		body, readError := ioutil.ReadAll(httpResponse.Body)
		if readError != nil {
			return []byte{}, readError
		}

		return body, nil
	}

	return response, nil
}

func GetInfo(name string) (MediaInfo, error) {
	response, err := performGuessitRequest(name)
	if err != nil {
		log.Error("[Guessit] Error while performing request: ", err)
		return MediaInfo{}, err
	}

	var mediaInfo MediaInfo
	parseErr := json.Unmarshal(response, &mediaInfo)
	if parseErr != nil {
		log.Error("[Guessit] Could not parse request result:", parseErr)
	}

	return mediaInfo, nil
}

func GetEpisodeInfo(name string, season int, episode int) (EpisodeTorrentInfo, error) {
	response, err := performGuessitRequest(name)
	if err != nil {
		log.Error("[Guessit] Error while performing request: ", err)
		return EpisodeTorrentInfo{}, err
	}

	var episodeInfo EpisodeTorrentInfo
	parseErr := json.Unmarshal(response, &episodeInfo)
	if parseErr != nil {
		log.Error("[Guessit] Could not parse request result:", parseErr)
		return EpisodeTorrentInfo{}, parseErr
	}

	return episodeInfo, nil
}
