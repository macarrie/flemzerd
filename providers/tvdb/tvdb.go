package tvdb

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	log "flemzerd/logging"
	provider "flemzerd/providers"
	"fmt"
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
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", APIToken.Token))

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

func New(apiKey string, username string, userKey string) TVDBProvider {
	return TVDBProvider{ApiKey: apiKey, Username: username, UserKey: userKey}
}

func (tvdbProvider TVDBProvider) Init() bool {
	log.Debug("Beginning authentication process to thetvdb API")

	params := map[string]string{"apikey": tvdbProvider.ApiKey, "username": tvdbProvider.Username, "userkey": tvdbProvider.UserKey}
	_, body := performAPIRequest("POST", "login", params)

	err := json.Unmarshal(body, &APIToken)
	if err != nil {
		log.Fatal(err)
		return false
	} else {
		log.Debug("Authentication successful")
		return true
	}
}

func (tvdbProvider TVDBProvider) FindShow(query string) (provider.Show, error) {
	log.WithFields(log.Fields{
		"name": query,
	}).Info("Searching show")
	params := map[string]string{"name": query}
	_, body := performAPIRequest("GET", "search/series", params)
	//log.Info("Status: ", status)
	//fmt.Println("Content: ", string(body[:]))

	var searchResults ShowsSearchResults
	err := json.Unmarshal(body, &searchResults)
	if err != nil {
		log.Fatal(err)
		return provider.Show{}, errors.New("Cannot find show")
	} else {
		if len(searchResults.Results) == 0 {
			log.WithFields(log.Fields{
				"search": query,
			}).Warning("No show found")
			return provider.Show{}, errors.New(fmt.Sprintf("No show found for query \"%s\"", query))
		}
		showID := searchResults.Results[0].Id

		log.WithFields(log.Fields{
			"Id":   showID,
			"name": query,
		}).Info("Found ID for show")

		return tvdbProvider.GetShow(showID)
	}
}

func (tvdbProvider TVDBProvider) GetShow(id int) (provider.Show, error) {
	log.WithFields(log.Fields{
		"Id": id,
	}).Info("Retrieving info for show")

	_, body := performAPIRequest("GET", fmt.Sprintf("series/%d", id), nil)

	var searchResults ShowSearchResult
	err := json.Unmarshal(body, &searchResults)
	if err != nil {
		log.Fatal(err)
		return provider.Show{}, errors.New(err.Error())
	} else {
		show := searchResults.Result

		return convertShow(show), nil
	}
}

func (tvdbProvider TVDBProvider) getEpisodesByPage(show provider.Show, page int) (EpisodesSearchResults, error) {
	params := map[string]string{"page": strconv.Itoa(page)}
	_, body := performAPIRequest("GET", fmt.Sprintf("series/%d/episodes/query", show.Id), params)

	var searchResults EpisodesSearchResults
	err := json.Unmarshal(body, &searchResults)
	if err != nil {
		log.Fatal(err)
		return EpisodesSearchResults{}, errors.New(err.Error())
	} else {
		return searchResults, nil
	}
}

func (tvdbProvider TVDBProvider) GetEpisodes(show provider.Show) ([]provider.Episode, error) {
	log.WithFields(log.Fields{
		"Id":   show.Id,
		"name": show.Name,
	}).Info("Retrieving epîsodes for show")

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
			}).Info("Found recent episode")

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
