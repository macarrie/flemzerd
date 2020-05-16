package torznab

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"github.com/macarrie/flemzerd/downloadable"

	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
	"github.com/rs/xid"

	"github.com/pkg/errors"
)

type TorznabIndexer struct {
	Name   string
	Url    string
	ApiKey string
	Caps   TorznabCaps
}

type TorznabError struct {
	Code        int
	Description string
}

func (e TorznabError) Error() string {
	return fmt.Sprintf("torznab error (code %d): %s", e.Code, e.Description)
}

type CapBool bool

func (b *CapBool) UnmarshalText(text []byte) error {
	switch strings.ToLower(string(text)) {
	case "yes":
		*b = true
	default:
		*b = false
	}
	return nil
}

type TorznabSupportedParams struct {
	SeasonParam  string
	EpisodeParam string
	Imdb         bool
	Tvdb         bool
}

func (p *TorznabSupportedParams) UnmarshalText(text []byte) error {
	lower := strings.ToLower(string(text))
	for _, param := range strings.Split(lower, ",") {
		switch param {
		case "season":
			(*p).SeasonParam = "season"
		case "imdbid":
			(*p).Imdb = true
		case "tvdbid":
			(*p).Tvdb = true
		}

		if strings.HasPrefix(param, "ep") {
			(*p).EpisodeParam = param
		}
	}

	return nil
}

type TorznabSearchCaps struct {
	Available       CapBool                `xml:"available,attr"`
	SupportedParams TorznabSupportedParams `xml:"supportedParams,attr"`
}

type TorznabCaps struct {
	Server struct {
		Title string `xml:"title,attr"`
	} `xml:"server"`
	Searching struct {
		Search      TorznabSearchCaps `xml:"search"`
		TVSearch    TorznabSearchCaps `xml:"tv-search"`
		MovieSearch TorznabSearchCaps `xml:"movie-search"`
	} `xml:"searching"`
	Categories []struct {
		ID   int    `xml:"id,attr"`
		Name string `xml:"name,attr"`
	} `xml:"categories>category"`
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
}

func NewTorznabCaps() TorznabCaps {
	t := TorznabCaps{}
	t.Searching.Search = TorznabSearchCaps{
		Available: true,
	}
	t.Searching.TVSearch = TorznabSearchCaps{
		Available: true,
		SupportedParams: TorznabSupportedParams{
			SeasonParam:  "season",
			EpisodeParam: "ep",
		},
	}
	t.Searching.MovieSearch = TorznabSearchCaps{
		Available: true,
	}

	return t
}

type TorrentSearchResults struct {
	Torrents []TorznabTorrent `xml:"channel>item"`
	XMLName  xml.Name
	Attrs    []xml.Attr `xml:",any,attr"`
}

type TorznabTorrent struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Guid        string `xml:"guid"`
	Comments    string `xml:"comments"`
	Link        string `xml:"link"`
	Category    string `xml:"category"`
	PubDate     string `xml:"pubDate"`
	Attr        []struct {
		Name  string `xml:"name,attr"`
		Value string `xml:"value,attr"`
	} `xml:"attr"`
}

func convertTorrent(t TorznabTorrent) Torrent {
	id := xid.New()
	return Torrent{
		TorrentId: id.String(),
		Name:      t.Title,
		Link:      t.Link,
	}
}

func New(name string, url string, apikey string) TorznabIndexer {
	t := TorznabIndexer{Name: name, Url: url, ApiKey: apikey}

	caps, err := t.GetCapabilities()
	if err != nil {
		t.Caps = NewTorznabCaps()
	}
	t.Caps = caps

	return t
}

func (torznabIndexer TorznabIndexer) Status() (Module, error) {
	returnStruct := Module{
		Name: torznabIndexer.GetName(),
		Type: "indexer",
		Status: ModuleStatus{
			Alive: false,
		},
	}

	log.WithFields(log.Fields{
		"name": torznabIndexer.GetName(),
	}).Debug("Checking torznab indexer status")

	testMovie := Movie{
		CustomTitle: "Big Buck Bunny",
	}
	_, err := torznabIndexer.GetTorrents(&testMovie)
	if err != nil {
		// Perhaps the indexer does not support movie search, try again with an episode
		tzErr, ok := err.(TorznabError)
		if ok && tzErr.Code == 201 {
			fmt.Println("TZ INDEXER STATUS 201")
			// Capabilities error returned by indexer, try again with an episode instead of a movie
			testEpisode := Episode{
				TvShow: TvShow{
					CustomTitle: "test",
				},
				Title:  "Big Buck Bunny",
				Season: 1,
				Number: 1,
			}
			_, episodeErr := torznabIndexer.GetTorrents(&testEpisode)
			if episodeErr != nil {
				returnStruct.Status.Message = episodeErr.Error()
				return returnStruct, episodeErr
			}
		}
		returnStruct.Status.Message = err.Error()
		return returnStruct, err
	}

	returnStruct.Status.Alive = true
	return returnStruct, nil
}

func (torznabIndexer TorznabIndexer) GetTorrents(d downloadable.Downloadable) ([]Torrent, error) {
	if !torznabIndexer.CheckCapabilities(d) {
		d.GetLog().WithFields(log.Fields{
			"indexer": torznabIndexer.Name,
		}).Info("Torznab indexer does not support torrent search for this item. Skipping")
		return []Torrent{}, nil
	}

	baseURL := torznabIndexer.Url

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(HTTP_TIMEOUT * time.Second),
	}

	urlObject, _ := url.ParseRequestURI(baseURL)

	var request *http.Request

	params := url.Values{}
	params.Add("apikey", torznabIndexer.ApiKey)
	switch d.(type) {
	case *Movie:
		params.Add("t", "movie")
		params.Add("q", d.GetTitle())
		if torznabIndexer.Caps.Searching.MovieSearch.SupportedParams.Imdb && d.GetMediaIds().Imdb != "" {
			params.Add("imdbid", d.GetMediaIds().Imdb)
		}
	case *Episode:
		episode := *(d.(*Episode))
		params.Add("t", "tvsearch")
		if torznabIndexer.Caps.Searching.MovieSearch.SupportedParams.Tvdb && d.GetMediaIds().Tvdb != 0 {
			params.Add("imdbid", strconv.Itoa(d.GetMediaIds().Tvdb))
		}
		if episode.TvShow.IsAnime && episode.AbsoluteNumber != 0 {
			params.Add("q", fmt.Sprintf("%v %v", episode.TvShow.GetTitle(), episode.AbsoluteNumber))
		} else {
			params.Add("q", episode.TvShow.GetTitle())
			params.Add(torznabIndexer.Caps.Searching.TVSearch.SupportedParams.SeasonParam, strconv.Itoa(episode.Season))
			params.Add(torznabIndexer.Caps.Searching.TVSearch.SupportedParams.EpisodeParam, strconv.Itoa(episode.Number))
		}
	default:
		return []Torrent{}, errors.New("Unknown downloadable type")
	}

	urlObject.RawQuery = params.Encode()

	request, err := http.NewRequest("GET", urlObject.String(), nil)
	if err != nil {
		return []Torrent{}, errors.Wrap(err, "error while constructing HTTP request to torznab indexer")
	}
	request.Close = true

	response, err := httpClient.Do(request)
	if err != nil {
		return []Torrent{}, errors.Wrap(err, "error while performing HTTP request to torznab indexer")
	}
	defer response.Body.Close()

	body, readError := ioutil.ReadAll(response.Body)
	if readError != nil {
		return []Torrent{}, errors.Wrap(readError, "error while reading HTTP result from torznab indexer request")
	}

	if len(body) == 0 {
		return []Torrent{}, errors.New("Empty result")
	}

	var searchResults TorrentSearchResults
	parseErr := xml.Unmarshal(body, &searchResults)
	if parseErr != nil {
		return []Torrent{}, errors.Wrap(parseErr, "cannot parse search results xml")
	}
	if searchResults.XMLName.Local == "error" {
		var code int
		var desc string
		for _, attr := range searchResults.Attrs {
			if attr.Name.Local == "description" {
				desc = attr.Value
			}
			if attr.Name.Local == "code" {
				code, _ = strconv.Atoi(attr.Value)
			}
		}

		return []Torrent{}, TorznabError{Code: code, Description: desc}
	}

	// Get seeders count for each torrent
	var results []Torrent
	for _, torrent := range searchResults.Torrents {
		resultTorrent := convertTorrent(torrent)

		for _, attr := range torrent.Attr {
			if attr.Name == "seeders" {
				seedersNb, _ := strconv.Atoi(attr.Value)
				resultTorrent.Seeders = seedersNb
			}
		}

		results = append(results, resultTorrent)
	}

	return results, nil
}

func (torznabIndexer TorznabIndexer) GetName() string {
	return torznabIndexer.Name
}

func (torznabIndexer TorznabIndexer) GetCapabilities() (TorznabCaps, error) {
	baseURL := torznabIndexer.Url

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(HTTP_TIMEOUT * time.Second),
	}

	urlObject, _ := url.ParseRequestURI(baseURL)

	var request *http.Request

	params := url.Values{}
	params.Add("apikey", torznabIndexer.ApiKey)
	params.Add("t", "caps")

	urlObject.RawQuery = params.Encode()

	request, err := http.NewRequest("GET", urlObject.String(), nil)
	if err != nil {
		return TorznabCaps{}, errors.Wrap(err, "error while constructing HTTP request to torznab indexer")
	}
	request.Close = true

	response, err := httpClient.Do(request)
	if err != nil {
		return TorznabCaps{}, errors.Wrap(err, "error while performing HTTP request to torznab indexer")
	}
	defer response.Body.Close()

	body, readError := ioutil.ReadAll(response.Body)
	if readError != nil {
		return TorznabCaps{}, errors.Wrap(readError, "error while reading HTTP result from torznab indexer request")
	}

	if len(body) == 0 {
		return TorznabCaps{}, errors.New("Empty result")
	}

	var capsResults TorznabCaps
	parseErr := xml.Unmarshal(body, &capsResults)
	if parseErr != nil {
		return TorznabCaps{}, errors.Wrap(parseErr, "cannot parse caps results xml")
	}
	if capsResults.XMLName.Local == "error" {
		var code int
		var desc string
		for _, attr := range capsResults.Attrs {
			if attr.Name.Local == "description" {
				desc = attr.Value
			}
			if attr.Name.Local == "code" {
				code, _ = strconv.Atoi(attr.Value)
			}
		}

		return TorznabCaps{}, TorznabError{Code: code, Description: desc}
	}

	return capsResults, nil
}

func (torznabIndexer TorznabIndexer) CheckCapabilities(d downloadable.Downloadable) bool {
	switch d.(type) {
	case *Movie:
		return bool(torznabIndexer.Caps.Searching.MovieSearch.Available)
	case *Episode:
		return bool(torznabIndexer.Caps.Searching.TVSearch.Available)
	default:
		return false
	}
}
