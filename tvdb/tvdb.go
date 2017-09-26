package tvdb

import (
    "fmt"
    "errors"
    "net/http"
    "net/url"
    "crypto/tls"
    "encoding/json"
    "io/ioutil"
    "bytes"
    log "flemzerd/logging"
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

type Token struct {
    Token string `json:"token"`
}

type Show struct {
    Aliases       []string  `json:"aliases"`
    Banner        string    `json:"banner"`
    FirstAired    string    `json:"firstAired"`
    Id            int       `json:"id"`
    Network       string    `json:"network"`
    Overview      string    `json:"overview"`
    SeriesName    string    `json:"seriesName"`
    Status        string    `json:"status"`
}

type Episode struct {
    AbsoluteNumber      int         `json:"absoluteNumber"`
    AiredEpisodeNumber  int         `json:"airedEpisodeNumber"`
    AiredSeason         int         `json:"airedSeason"`
    //DvdEpisodeNumber    int         `json:"dvdEpisodeNumber"`
    //DvdSeason           int         `json:"dvdSeason"`
    EpisodeName         string      `json:"episodeName"`
    FirstAired          string      `json:"firstAired"`
    Id                  int         `json:"id"`
    LastUpdated         int         `json:"lastUpdated"`
    Overview            string      `json:"overview"`
}

type ShowsSearchResults struct {
    //data interface{} `json:"data"`
    Results []Show `json:"data"`
}

type ShowSearchResult struct {
    //data interface{} `json:"data"`
    Result Show `json:"data"`
}

type EpisodesSearchResults struct {
    Results []Episode `json:"data"`
    Pages   struct {
        First   int     `json:"first"`
        Last    int     `json:"last"`
    } `json:"links"`
}

//////////////////
// Package vars //
//////////////////

var APIToken Token

func performAPIRequest(method string, path string, paramsMap map[string]string) (string,[]byte) {
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify : true},
    }
    httpClient := &http.Client{ Transport: tr}

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

        log.Debug("GET ", urlObject.String())

        request, err = http.NewRequest(method, urlObject.String(), nil)
        if err != nil {
            log.Fatal("NewRequest: ", err)
        }
    } else if method == "POST" {
        jsonParams, _ := json.Marshal(paramsMap)

        log.Debug("POST ", urlObject.String(), paramsMap)
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

    return response.Status, body
}

func Authenticate(apiKey string, username string, userKey string) bool {
    log.Info("Beginning authentication process to thetvdb API")

    params := map[string]string{"apikey": apiKey, "username": username, "userkey": userKey}
    _, body := performAPIRequest("POST", "login", params)
    //fmt.Println("Status: ", status)
    //fmt.Println("Content: ", string(body[:]))

    err := json.Unmarshal(body, &APIToken)
    log.Debug("API Token: ", APIToken.Token)
    if err != nil {
        log.Fatal(err)
        return false
    } else {
        log.Info("Authentication successful")
        return true
    }
}

func FindShow(query string) (Show, error) {
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
        return Show{}, errors.New("Cannot find show")
    } else {
        if len(searchResults.Results) == 0 {
            log.WithFields(log.Fields{
                "search": query,
            }).Warning("No show found")
            return Show{}, errors.New(fmt.Sprintf("No show found for query \"%s\"", query))
        }
        showID := searchResults.Results[0].Id

        log.WithFields(log.Fields{
            "Id": showID,
            "name": query,
        }).Info("Found ID for show")

        return GetShow(showID)
    }
}

func GetShow(id int) (Show, error) {
    log.WithFields(log.Fields{
        "Id": id,
    }).Info("Retrieving info for show")

    _, body := performAPIRequest("GET", fmt.Sprintf("series/%d", id), nil)
    //log.Info("Status: ", status)
    //fmt.Println("Content: ", string(body[:]))

    var searchResults ShowSearchResult
    err := json.Unmarshal(body, &searchResults)
    if err != nil {
        log.Fatal(err)
        return Show{}, errors.New(err.Error())
    } else {
        show := searchResults.Result

        return show, nil
    }
}

func getEpisodesByPage(showID int, page int) (EpisodesSearchResults, error) {
    params := map[string]string{"page": strconv.Itoa(page)}
    _, body := performAPIRequest("GET", fmt.Sprintf("series/%d/episodes/query", showID), params)

    var searchResults EpisodesSearchResults
    err := json.Unmarshal(body, &searchResults)
    if err != nil {
        log.Fatal(err)
        return EpisodesSearchResults{}, errors.New(err.Error())
    } else {
        return searchResults, nil
    }
}

func GetEpisodes(showID int) ([]Episode, error) {
    log.WithFields(log.Fields{
        "Id": showID,
    }).Info("Retrieving epîsodes for show")

    episodesSearchResults, err := getEpisodesByPage(showID, 1)
    if err != nil {
        log.WithFields(log.Fields{
            "Error": err,
        }).Warning("Failed to retrieve episodes")
        return []Episode{}, nil
    }

    var episodesList []Episode

    if episodesSearchResults.Pages.Last > 1 {
        episodesList = episodesSearchResults.Results
        for i := 2; i <= episodesSearchResults.Pages.Last; i++ {
            episodesPage, err := getEpisodesByPage(showID, i)
            if err != nil {
                log.WithFields(log.Fields{
                    "Error": err,
                }).Warning("Failed to retrieve episodes")
                return []Episode{}, nil
            }
            episodesList = append(episodesList, episodesPage.Results...)
        }
    } else {
        episodesList = episodesSearchResults.Results
    }

    sort.Slice(episodesList[:], func (i, j int) bool {
         return episodesList[i].AbsoluteNumber > episodesList[j].AbsoluteNumber
    })

    return episodesList, nil
}

func FindNextAiredEpisode(episodeList []Episode) (Episode, error) {
    log.Debug("Looking for next episode in given episodes list")

    sort.Slice(episodeList[:], func (i, j int) bool {
        return episodeList[i].FirstAired < episodeList[j].FirstAired
    })

    var futureEpisodes []Episode
    for _, episode := range episodeList {
        airDate, err := time.Parse("2006-01-02", episode.FirstAired)
        if err != nil {
            continue
        }


        if airDate.After(time.Now()) {
            //log.WithFields(log.Fields{
                //"airDate": airDate,
                //"episode": episode,
            //}).Info("Episode not aired yet")

            futureEpisodes = append(futureEpisodes, episode)
        }
    }

    if len(futureEpisodes) == 0 {
         return Episode{}, errors.New("All episodes in episode list have already aired")
    }

    return futureEpisodes[0], nil
}

func FindNextEpisodeForShow(name string) (Episode, error) {
    show, err := FindShow(name)
    if err != nil {
        return Episode{}, err
    }

    episodes, err := GetEpisodes(show.Id)
    if err != nil {
        return Episode{}, err
    }

    nextEpisode, err := FindNextAiredEpisode(episodes)
    if err != nil {
        return Episode{}, err
    }

    return nextEpisode, nil
}
