package torznab

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
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
}

type TorrentSearchResults struct {
	Torrents []TorznabTorrent `xml:"channel>item"`
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
	return TorznabIndexer{Name: name, Url: url, ApiKey: apikey}
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
	urlObject.RawQuery = params.Encode()

	request, err := http.NewRequest("GET", urlObject.String(), nil)
	if err != nil {
		returnStruct.Status.Message = err.Error()
		return returnStruct, errors.Wrap(err, "error while constructing HTTP request to torznab indexer")
	}
	request.Close = true

	response, err := httpClient.Do(request)
	if err != nil {
		returnStruct.Status.Message = err.Error()
		return returnStruct, errors.Wrap(err, "error while performing HTTP request to torznab indexer")
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		err := errors.New(fmt.Sprintf("Torznab indexer request returned %d status code", response.StatusCode))
		returnStruct.Status.Message = err.Error()
		return returnStruct, err
	}

	returnStruct.Status.Alive = true
	return returnStruct, nil
}

func (torznabIndexer TorznabIndexer) GetTorrents(d downloadable.Downloadable) ([]Torrent, error) {
	if !torznabIndexer.CheckCapabilities(d) {
		d.GetLog().Info("Torznab indexer does not support torrent search for this item. Skipping")
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
	case *Episode:
		episode := *(d.(*Episode))
		params.Add("t", "tvsearch")
		if episode.TvShow.IsAnime && episode.AbsoluteNumber != 0 {
			params.Add("q", fmt.Sprintf("%v %v", episode.TvShow.GetTitle(), episode.AbsoluteNumber))
		} else {
			params.Add("q", episode.TvShow.GetTitle())
			params.Add("season", strconv.Itoa(episode.Season))
			params.Add("episode", strconv.Itoa(episode.Number))
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

func (torznabIndexer TorznabIndexer) CheckCapabilities(d downloadable.Downloadable) bool {
	switch d.(type) {
	case *Movie:
		return true
	case *Episode:
		return true
	default:
		return false
	}
}
