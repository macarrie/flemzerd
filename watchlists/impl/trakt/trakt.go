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

	"github.com/macarrie/flemzerd/configuration"
	"github.com/macarrie/flemzerd/db"
	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
)

const (
	TRAKT_API_URL = "https://api.trakt.tv"
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

type TraktMovie struct {
	Title string `json:"title"`
	Year  int    `json:"year"`
	Ids   struct {
		Trakt int    `json:"trakt"`
		Slug  string `json:"slug"`
		Imdb  string `json:"imdb"`
		Tmdb  int    `json:"tmdb"`
	} `json:"ids"`
}

type traktTVWatchlistRequestResult struct {
	Rank     int       `json:"rank"`
	ListedAt string    `json:"listed_at"`
	Type     string    `json:"type"`
	Show     TraktShow `json:"show"`
}

type traktMoviesWatchlistRequestResult struct {
	Rank     int        `json:"rank"`
	ListedAt string     `json:"listed_at"`
	Type     string     `json:"type"`
	Movie    TraktMovie `json:"movie"`
}

var module Module
var authErrors []error

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

	token := db.Session.TraktToken
	if token != "" {
		t.Token.AccessToken = token
	} else {
		log.Warning("No Trakt token found. User will need to authorize access to Trakt")
	}

	module = Module{
		Name: t.GetName(),
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
		"client_secret": configuration.TRAKT_CLIENT_SECRET,
		"code":          t.DeviceCode.DeviceCode,
	}
	response, err := t.performAPIRequest("POST", "/oauth/device/token", params)
	if err != nil {
		return http.Response{}, err
	}

	return response, nil
}

func (t *TraktWatchlist) Auth() {
	authErrors = []error{}

	if t.DeviceCode.DeviceCode != "" {
		log.Debug("Trakt auth process already in progress. Skipping")
		return
	}
	if t.Token.AccessToken != "" {
		log.Debug("Trakt app already authorized")
		return
	}

	err := t.RequestDeviceCode()
	if err != nil {
		log.Error("RequestDeviceError: ", err)
		return
	}

	errorsChannel := make(chan error, 1)
	doneChannel := make(chan bool, 1)

	pollTicker := time.NewTicker(time.Duration(t.DeviceCode.Interval) * time.Second)
	go func() {
		for {
			response, err := t.RequestToken()
			if err != nil {
				log.Error(err)
				continue
			}
			defer response.Body.Close()

			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Error(err)
			}

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
				log.Debug("[Trakt auth] Waiting for user authorization")

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
				log.Warning("[Trakt auth] Polling too quickly")

			default:
				log.Error("[Trakt auth] ", "Unknown status code ", response.StatusCode)
				doneChannel <- true
			}

			<-pollTicker.C
		}
	}()

	go func() {
		time.Sleep(time.Duration(t.DeviceCode.ExpiresIn) * time.Second)
		errorsChannel <- errors.New("[Trakt auth] Device code expired")
		t.DeviceCode = TraktDeviceCode{}
		doneChannel <- true
	}()

	<-doneChannel
	pollTicker.Stop()

	select {
	case err := <-errorsChannel:
		authErrors = append(authErrors, err)

		return
	default:
		t.DeviceCode = TraktDeviceCode{}
		db.SaveTraktToken(t.Token.AccessToken)
		authErrors = []error{}

		return
	}
}

func (t *TraktWatchlist) GetAuthErrors() []error {
	return authErrors
}

func (t *TraktWatchlist) IsAuthenticated() error {
	if t.Token.AccessToken == "" {
		return errors.New("No trakt token found")
	}

	response, err := t.performAPIRequest("GET", "/users/settings", nil)
	if err != nil {
		return err
	} else {
		if response.StatusCode != http.StatusOK {
			t.DeviceCode = TraktDeviceCode{}
			t.Token = TraktToken{}

			return errors.New("Not authenticated into Trakt")
		} else {
			return nil
		}
	}

}

func (t *TraktWatchlist) Status() (Module, error) {
	log.Warning("Checking Trakt watchlist status")

	err := t.IsAuthenticated()
	if err != nil {
		module.Status.Alive = false
		module.Status.Message = err.Error()
	} else {
		module.Status.Alive = true
		module.Status.Message = ""
	}

	return module, err
}

func (t *TraktWatchlist) GetName() string {
	return "trakt"
}

func (t *TraktWatchlist) GetTvShows() ([]MediaIds, error) {
	log.WithFields(log.Fields{
		"watchlist": t.GetName(),
	}).Debug("Getting TV shows from watchlist")

	if t.Token.AccessToken == "" {
		return []MediaIds{}, errors.New("Not authenticated into Trakt")
	}

	response, err := t.performAPIRequest("GET", "/users/me/watchlist/shows", nil)
	if err != nil {
		return []MediaIds{}, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []MediaIds{}, err
	}

	var results []traktTVWatchlistRequestResult
	if response.StatusCode == http.StatusOK {
		parseErr := json.Unmarshal(body, &results)
		if parseErr != nil {
			log.Error("JSON Unmarshal error: ", err)
			return []MediaIds{}, err
		}

		var toReturn []MediaIds
		for _, val := range results {
			toReturn = append(toReturn, MediaIds{
				Name:  val.Show.Title,
				Trakt: val.Show.Ids.Trakt,
				Tmdb:  val.Show.Ids.Tmdb,
				Imdb:  val.Show.Ids.Imdb,
				Tvdb:  val.Show.Ids.Tvdb,
			})
		}

		return toReturn, nil
	}
	return []MediaIds{}, errors.New("Unknown HTTP return code from trakt watchlist call")
}

func (t *TraktWatchlist) GetMovies() ([]MediaIds, error) {
	log.WithFields(log.Fields{
		"watchlist": t.GetName(),
	}).Debug("Getting movies from watchlist")

	if t.Token.AccessToken == "" {
		return []MediaIds{}, errors.New("Not authenticated into Trakt")
	}

	response, err := t.performAPIRequest("GET", "/users/me/watchlist/movies", nil)
	if err != nil {
		return []MediaIds{}, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []MediaIds{}, err
	}

	var results []traktMoviesWatchlistRequestResult
	if response.StatusCode == http.StatusOK {
		parseErr := json.Unmarshal(body, &results)
		if parseErr != nil {
			log.Error("JSON Unmarshal error: ", err)
			return []MediaIds{}, err
		}

		var toReturn []MediaIds
		for _, val := range results {
			toReturn = append(toReturn, MediaIds{
				Name:  val.Movie.Title,
				Trakt: val.Movie.Ids.Trakt,
				Tmdb:  val.Movie.Ids.Tmdb,
				Imdb:  val.Movie.Ids.Imdb,
			})
		}

		return toReturn, nil
	}
	return []MediaIds{}, errors.New("Unknown HTTP return code from trakt watchlist call")
}
