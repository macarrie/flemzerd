package torznab

import (
	"crypto/tls"
	"encoding/xml"
	//"errors"
	"flemzerd/indexers"
	log "flemzerd/logging"
	//"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
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

func New(name string, url string, apikey string) TorznabIndexer {
	return TorznabIndexer{Name: name, Url: url, ApiKey: apikey}
}

func (torznabIndexer TorznabIndexer) GetTorrentForEpisode(show string, season int, episode int) ([]indexer.Torrent, error) {
	baseURL := torznabIndexer.Url

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}

	urlObject, _ := url.ParseRequestURI(baseURL)

	var request *http.Request

	params := url.Values{}
	params.Add("apikey", torznabIndexer.ApiKey)
	params.Add("t", "tvsearch")
	params.Add("q", show)
	params.Add("season", strconv.Itoa(season))
	params.Add("episode", strconv.Itoa(episode))
	urlObject.RawQuery = params.Encode()

	request, err := http.NewRequest("GET", urlObject.String(), nil)
	if err != nil {
		return []indexer.Torrent{}, err
	}

	response, err := httpClient.Do(request)
	if err != nil {
		return []indexer.Torrent{}, err
	}

	body, readError := ioutil.ReadAll(response.Body)
	if readError != nil {
		return []indexer.Torrent{}, err
	}

	var searchResults TorrentSearchResults
	parseErr := xml.Unmarshal(body, &searchResults)
	if parseErr != nil {
		log.Debug(parseErr)
		return []indexer.Torrent{}, parseErr
	}

	// Construct Attributes map
	var results []indexer.Torrent
	for _, torrent := range searchResults.Torrents {
		resultTorrent := &indexer.Torrent{
			Title:       torrent.Title,
			Description: torrent.Description,
			Link:        torrent.Link,
		}

		resultTorrent.Attributes = make(map[string]string, len(torrent.Attr))
		for _, attr := range torrent.Attr {
			resultTorrent.Attributes[attr.Name] = attr.Value
		}
		results = append(results, *resultTorrent)
	}

	return results, nil
}

func (torznabIndexer TorznabIndexer) GetName() string {
	return torznabIndexer.Name
}
