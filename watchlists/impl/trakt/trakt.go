package trakt

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
	"github.com/macarrie/flemzerd/retention"
)

const (
	TRAKT_CLIENT_ID     = "4d5d84f3a95adc528d08b1e2f6df6c8197548d128e72e68b4c33fc0dfa634926"
	TRAKT_CLIENT_SECRET = "bb0816bb9f73e6a4fbc837f583ce8a0fb8160648a9a9d60f71e55fa47cca4b42"
	TRAKT_API_URL       = "https://api.trakt.tv"
)

type TraktWatchlist struct {
	DeviceCode TraktDeviceCode
	Token      TraktToken
}

type TraktDeviceCode struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationUrl string `json:"verification_url"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

type TraktToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	CreatedAt    int    `json:"created_at"`
}

type TraktShow struct {
	Title string `json:"title"`
	Year  int    `json:"year"`
	Ids   struct {
		Trakt  int    `json:"trakt"`
		Slug   string `json:"slug"`
		Tvdb   int    `json:"tvdb"`
		Imdb   string `json:"imdb"`
		Tmdb   int    `json:"tmdb"`
		TvRage int    `json:"tvrage,omitempty"`
	} `json:"ids"`
}

type traktWatchlistRequestResult struct {
	Rank     int       `json:"rank"`
	ListedAt string    `json:"listed_at"`
	Type     string    `json:"type"`
	Show     TraktShow `json:"show"`
}

var module Module

func (t *TraktWatchlist) performAPIRequest(method string, path string, paramsMap map[string]string) (http.Response, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(HTTP_TIMEOUT * time.Second),
	}

	urlObject, _ := url.ParseRequestURI(TRAKT_API_URL)
	urlObject.Path = path
	if paramsMap == nil {
		paramsMap = map[string]string{}
	}

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
	request.Header.Set("trakt-api-key", TRAKT_CLIENT_ID)
	request.Header.Set("trakt-api-version", "2")

	if t.Token.AccessToken != "" {
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.Token.AccessToken))
	}
	request.Close = true

	response, err := httpClient.Do(request)
	if err != nil {
		return http.Response{}, err
	}

	return *response, nil
}

func New() (t *TraktWatchlist, err error) {
	t = &TraktWatchlist{}

	token := retention.LoadTraktToken()
	if token != "" {
		t.Token.AccessToken = token
	} else {
		log.Warning("No Trakt token found. User will need to authorize access to Trakt")
	}

	module = Module{
		Name: "Trakt",
		Type: "watchlist",
		Status: ModuleStatus{
			Alive:   true,
			Message: "",
		},
	}

	return t, nil
}

func (t *TraktWatchlist) RequestDeviceCode() error {
	paramsMap := map[string]string{"client_id": TRAKT_CLIENT_ID}
	response, err := t.performAPIRequest("POST", "/oauth/device/code", paramsMap)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	var deviceCode TraktDeviceCode
	if response.StatusCode == http.StatusOK {
		parseErr := json.Unmarshal(body, &deviceCode)
		if parseErr != nil {
			log.Error(err)
			return err
		}
	}

	t.DeviceCode = deviceCode
	return nil
}

func (t *TraktWatchlist) RequestToken() (http.Response, error) {
	params := map[string]string{
		"client_id":     TRAKT_CLIENT_ID,
		"client_secret": TRAKT_CLIENT_SECRET,
		"code":          t.DeviceCode.DeviceCode,
	}
	response, err := t.performAPIRequest("POST", "/oauth/device/token", params)
	if err != nil {
		return http.Response{}, err
	}

	return response, nil
}

func (t *TraktWatchlist) Auth() error {
	if t.DeviceCode.DeviceCode != "" {
		log.Debug("Trakt auth process already in progress. Skipping")
		return nil
	}
	if t.Token.AccessToken != "" {
		log.Debug("Trakt app already authorized")
		return nil
	}

	err := t.RequestDeviceCode()
	fmt.Printf("%+v\n", t)
	if err != nil {
		log.Error("RequestDeviceError: ", err)
		return err
	}

	errorsChannel := make(chan error, 1)
	doneChannel := make(chan bool, 1)

	pollTicker := time.NewTicker(time.Duration(t.DeviceCode.Interval) * time.Second)
	go func() {
		for {
			response, err := t.RequestToken()
			if err != nil {
				log.Error(err)
			}
			defer response.Body.Close()

			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Error(err)
			}
			log.Warning(response.Status)
			log.Warning(string(body))

			switch response.StatusCode {
			case http.StatusOK:
				t.DeviceCode = TraktDeviceCode{}

				var token TraktToken
				parseErr := json.Unmarshal(body, &token)
				if parseErr != nil {
					errorsChannel <- err
					doneChannel <- true
				}

				t.Token = token
				doneChannel <- true

			case http.StatusBadRequest:
				// Waiting for the user to authorize
				log.Debug("Trakt auth: Waiting for user authorization")

			case http.StatusNotFound:
				t.DeviceCode = TraktDeviceCode{}
				errorsChannel <- errors.New("[Trakt auth] Invalid device code")
				doneChannel <- true

			case http.StatusConflict:
				// Code already approved, no error returned
				t.DeviceCode = TraktDeviceCode{}
				doneChannel <- true

			case http.StatusGone:
				t.DeviceCode = TraktDeviceCode{}
				t.Token = TraktToken{}
				errorsChannel <- errors.New("[Trakt auth] Token expired")
				doneChannel <- true

			case http.StatusTeapot:
				t.DeviceCode = TraktDeviceCode{}
				errorsChannel <- errors.New("[Trakt auth] Authorization refused by used")
				doneChannel <- true

			case http.StatusTooManyRequests:
				log.Warning("Trakt auth: Polling too quickly")

			default:
				log.Error("Trakt auth: ", "Unknown status code ", response.StatusCode)
				doneChannel <- true
			}

			<-pollTicker.C
		}
	}()

	go func() {
		time.Sleep(time.Duration(t.DeviceCode.ExpiresIn) * time.Second)
		//time.Sleep(8 * time.Second)
		errorsChannel <- errors.New("[Trakt auth] Device code expired")
		t.DeviceCode = TraktDeviceCode{}
		doneChannel <- true
	}()

	<-doneChannel
	pollTicker.Stop()

	select {
	case err := <-errorsChannel:
		log.Error("TRAKT ERROR: ", err)
		return err
	default:
		log.Info("Trakt auth ok")
		t.DeviceCode = TraktDeviceCode{}
		retention.SaveTraktToken(t.Token.AccessToken)
		return nil
	}
}

func (t *TraktWatchlist) Status() (Module, error) {
	log.Warning("Checking Trakt watchlist status")

	response, err := t.performAPIRequest("GET", "/users/settings", nil)
	log.Warning(response.StatusCode)
	if err != nil {
		module.Status.Alive = false
		module.Status.Message = err.Error()
	} else {
		if response.StatusCode != http.StatusOK {
			log.Error("Not alive")
			t.DeviceCode = TraktDeviceCode{}
			t.Token = TraktToken{}

			module.Status.Alive = false
			module.Status.Message = "Not authenticated into Trakt"
		} else {
			module.Status.Alive = true
		}
	}

	return module, err
}

func (t *TraktWatchlist) GetTvShows() ([]string, error) {
	log.WithFields(log.Fields{
		"watchlist": "trakt",
	}).Debug("Getting TV shows from watchlist")

	if t.Token.AccessToken == "" {
		return []string{}, errors.New("Not authenticated into Trakt")
	}

	response, err := t.performAPIRequest("GET", "/users/me/watchlist/shows", nil)
	if err != nil {
		return []string{}, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []string{}, err
	}

	var results []traktWatchlistRequestResult
	if response.StatusCode == http.StatusOK {
		parseErr := json.Unmarshal(body, &results)
		if parseErr != nil {
			log.Error("JSON Unmarshal error: ", err)
			return []string{}, err
		}

		var toReturn []string
		for _, val := range results {
			toReturn = append(toReturn, val.Show.Title)
		}

		toReturn = append(toReturn, retention.GetShowsFromWatchlist("trakt")...)

		return toReturn, nil
	}
	return []string{}, errors.New("Unknown HTTP return code from trakt watchlist call")
}
func (t *TraktWatchlist) GetMovies() ([]string, error) {
	log.WithFields(log.Fields{
		"watchlist": "trakt",
	}).Debug("Getting movies from watchlist")

	return []string{}, nil
}
