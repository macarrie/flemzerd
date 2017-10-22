package tvdb

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/macarrie/flemzerd/logging"
	provider "github.com/macarrie/flemzerd/providers"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
)

///////////////
// Constants //
///////////////

const baseURL = "https://api.thetvdb.com/"

/////////////////////////////////
// API Structures declarations //
/////////////////////////////////

type TVDBProvider struct {
	ApiKey   string
	Username string
	UserKey  string
	Token    string
}

type Token struct {
	Token string `json:"token"`
}

type TVDBShow struct {
	Aliases    []string `json:"aliases"`
	Banner     string   `json:"banner"`
	FirstAired string   `json:"firstAired"`
	Id         int      `json:"id"`
	Network    string   `json:"network"`
	Overview   string   `json:"overview"`
	SeriesName string   `json:"seriesName"`
	Status     string   `json:"status"`
}

type TVDBEpisode struct {
	AbsoluteNumber     int `json:"absoluteNumber"`
	AiredEpisodeNumber int `json:"airedEpisodeNumber"`
	AiredSeason        int `json:"airedSeason"`
	//DvdEpisodeNumber    int         `json:"dvdEpisodeNumber"`
	//DvdSeason           int         `json:"dvdSeason"`
	EpisodeName string `json:"episodeName"`
	FirstAired  string `json:"firstAired"`
	Id          int    `json:"id"`
	LastUpdated int    `json:"lastUpdated"`
	Overview    string `json:"overview"`
}

type ShowsSearchResults struct {
	Results []TVDBShow `json:"data"`
}

type ShowSearchResult struct {
	//data interface{} `json:"data"`
	Result TVDBShow `json:"data"`
}

type EpisodesSearchResults struct {
	Results []TVDBEpisode `json:"data"`
	Pages   struct {
		First int `json:"first"`
		Last  int `json:"last"`
	} `json:"links"`
}

//////////////////
// Package vars //
//////////////////

var APIToken Token

func convertShow(s TVDBShow) provider.Show {
	return provider.Show{
		Aliases:    s.Aliases,
		Banner:     s.Banner,
		FirstAired: s.FirstAired,
		Id:         s.Id,
		Overview:   s.Overview,
		Name:       s.SeriesName,
		Status:     s.Status,
	}
}

func convertEpisode(e TVDBEpisode) provider.Episode {
	return provider.Episode{
		AbsoluteNumber: e.AbsoluteNumber,
		Number:         e.AiredEpisodeNumber,
		Season:         e.AiredSeason,
		Name:           e.EpisodeName,
		Date:           e.FirstAired,
		Id:             e.Id,
		Overview:       e.Overview,
	}
}

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
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", APIToken.Token))

	response, err := httpClient.Do(request)
	if err != nil {
		return http.Response{}, err
	}

	return *response, nil
}

func New(apiKey string, username string, userKey string) TVDBProvider {
	return TVDBProvider{ApiKey: apiKey, Username: username, UserKey: userKey}
}

func (tvdbProvider TVDBProvider) Init() error {
	log.Debug("Beginning authentication process to thetvdb API")

	params := map[string]string{"apikey": tvdbProvider.ApiKey, "username": tvdbProvider.Username, "userkey": tvdbProvider.UserKey}
	response, err := performAPIRequest("POST", "login", params)
	if err != nil {
		return err
	}

	body, readError := ioutil.ReadAll(response.Body)
	if readError != nil {
		return errors.New(fmt.Sprintf("API Response read: %v", readError))
	}

	parseErr := json.Unmarshal(body, &APIToken)
	if parseErr != nil {
		return parseErr
	} else {
		log.Debug("Authentication successful")
		return nil
	}
}

func (tvdbProvider TVDBProvider) FindShow(query string) (provider.Show, error) {
	log.WithFields(log.Fields{
		"name": query,
	}).Debug("Searching show")
	params := map[string]string{"name": query}
	response, err := performAPIRequest("GET", "search/series", params)
	if err != nil {
		return provider.Show{}, err
	}

	body, readError := ioutil.ReadAll(response.Body)
	if readError != nil {
		return provider.Show{}, errors.New(fmt.Sprintf("API Response read: %v", readError))
	}

	var searchResults ShowsSearchResults
	parseErr := json.Unmarshal(body, &searchResults)
	if parseErr != nil {
		return provider.Show{}, errors.New("Cannot find show")
	} else {
		if len(searchResults.Results) == 0 {
			return provider.Show{}, errors.New("No show foun")
		}
		showID := searchResults.Results[0].Id

		log.WithFields(log.Fields{
			"Id":   showID,
			"name": query,
		}).Debug("Found ID for show")

		return tvdbProvider.GetShow(showID)
	}
}

func (tvdbProvider TVDBProvider) GetShow(id int) (provider.Show, error) {
	log.WithFields(log.Fields{
		"Id": id,
	}).Debug("Retrieving info for show")

	response, err := performAPIRequest("GET", fmt.Sprintf("series/%d", id), nil)
	if err != nil {
		return provider.Show{}, err
	}

	body, readError := ioutil.ReadAll(response.Body)
	if readError != nil {
		return provider.Show{}, errors.New(fmt.Sprintf("API Response read: %v", readError))
	}

	var searchResults ShowSearchResult
	parseErr := json.Unmarshal(body, &searchResults)
	if parseErr != nil {
		return provider.Show{}, parseErr
	} else {
		show := searchResults.Result

		return convertShow(show), nil
	}
}

func (tvdbProvider TVDBProvider) getEpisodesByPage(show provider.Show, page int) (EpisodesSearchResults, error) {
	params := map[string]string{"page": strconv.Itoa(page)}
	response, err := performAPIRequest("GET", fmt.Sprintf("series/%d/episodes/query", show.Id), params)
	if err != nil {
		return EpisodesSearchResults{}, err
	}

	body, readError := ioutil.ReadAll(response.Body)
	if readError != nil {
		return EpisodesSearchResults{}, errors.New(fmt.Sprintf("API Response read: %v", readError))
	}

	var searchResults EpisodesSearchResults
	parseErr := json.Unmarshal(body, &searchResults)
	if parseErr != nil {
		return EpisodesSearchResults{}, parseErr
	} else {
		return searchResults, nil
	}
}

func (tvdbProvider TVDBProvider) GetEpisodes(show provider.Show) ([]provider.Episode, error) {
	log.WithFields(log.Fields{
		"Id":   show.Id,
		"name": show.Name,
	}).Debug("Retrieving epîsodes for show")

	episodesSearchResults, err := tvdbProvider.getEpisodesByPage(show, 1)
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Warning("Failed to retrieve episodes")
		return []provider.Episode{}, nil
	}

	var tvdbEpisodesList []TVDBEpisode

	if episodesSearchResults.Pages.Last > 1 {
		tvdbEpisodesList = episodesSearchResults.Results
		for i := 2; i <= episodesSearchResults.Pages.Last; i++ {
			episodesPage, err := tvdbProvider.getEpisodesByPage(show, i)
			if err != nil {
				log.WithFields(log.Fields{
					"Error": err,
				}).Warning("Failed to retrieve episodes")
				return []provider.Episode{}, nil
			}
			tvdbEpisodesList = append(tvdbEpisodesList, episodesPage.Results...)
		}
	} else {
		tvdbEpisodesList = episodesSearchResults.Results
	}

	sort.Slice(tvdbEpisodesList[:], func(i, j int) bool {
		return tvdbEpisodesList[i].AbsoluteNumber > tvdbEpisodesList[j].AbsoluteNumber
	})

	var episodesList []provider.Episode
	for _, e := range tvdbEpisodesList {
		episodesList = append(episodesList, convertEpisode(e))
	}

	return episodesList, nil
}

func (tvdbProvider TVDBProvider) FindNextAiredEpisodes(episodeList []provider.Episode) ([]provider.Episode, error) {
	log.Debug("Looking for next episode in given episodes list")

	sort.Slice(episodeList[:], func(i, j int) bool {
		return episodeList[i].Date < episodeList[j].Date
	})

	var futureEpisodes []provider.Episode
	for _, episode := range episodeList {
		airDate, err := time.Parse("2006-01-02", episode.Date)
		if err != nil {
			continue
		}

		if airDate.After(time.Now()) {
			futureEpisodes = append(futureEpisodes, episode)
		}
	}

	if len(futureEpisodes) == 0 {
		return []provider.Episode{}, errors.New("All episodes in episode list have already aired")
	}

	return futureEpisodes, nil
}

func (tvdbProvider TVDBProvider) FindNextEpisodesForShow(show provider.Show) ([]provider.Episode, error) {
	episodes, err := tvdbProvider.GetEpisodes(show)
	if err != nil {
		return []provider.Episode{}, err
	}

	nextEpisodes, err := tvdbProvider.FindNextAiredEpisodes(episodes)
	if err != nil {
		return []provider.Episode{}, err
	}

	return nextEpisodes, nil
}

func (tvdbProvider TVDBProvider) FindRecentlyAiredEpisodes(episodeList []provider.Episode) ([]provider.Episode, error) {
	log.Debug("Looking for recent episode in given episodes list")

	sort.Slice(episodeList[:], func(i, j int) bool {
		return episodeList[i].Date < episodeList[j].Date
	})

	var recentEpisodes []provider.Episode
	for _, episode := range episodeList {
		airDate, err := time.Parse("2006-01-02", episode.Date)
		if err != nil {
			continue
		}

		if airDate.Before(time.Now()) && airDate.After(time.Now().AddDate(0, 0, -14)) {
			episodeLogString := fmt.Sprintf("S%03dE%03d: %v", episode.Season, episode.Number, episode.Name)
			log.WithFields(log.Fields{
				"airDate": airDate,
				"episode": episodeLogString,
			}).Debug("Found recent episode")

			recentEpisodes = append(recentEpisodes, episode)
		}
	}

	if len(recentEpisodes) == 0 {
		return []provider.Episode{}, errors.New("No recent episodes in episodes list")
	}

	return recentEpisodes, nil
}

func (tvdbProvider TVDBProvider) FindRecentlyAiredEpisodesForShow(show provider.Show) ([]provider.Episode, error) {
	episodes, err := tvdbProvider.GetEpisodes(show)
	if err != nil {
		return []provider.Episode{}, err
	}

	nextEpisodes, err := tvdbProvider.FindRecentlyAiredEpisodes(episodes)
	if err != nil {
		return []provider.Episode{}, err
	}

	return nextEpisodes, nil
}
