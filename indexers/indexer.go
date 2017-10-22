package indexer

import (
	"fmt"
	log "github.com/macarrie/flemzerd/logging"
	"sort"
	"strconv"
)

type Indexer interface {
	GetName() string
	GetTorrentForEpisode(show string, season int, episode int) ([]Torrent, error)
}

type Torrent struct {
	Title       string
	Description string
	Link        string
	Attributes  map[string]string
}

var indexers []Indexer

//func Init() {
//for _, indexer := range indexers {
//indexer.Init()
//}
//}

func AddIndexer(indexer Indexer) {
	indexers = append(indexers, indexer)
	log.WithFields(log.Fields{
		"indexer": indexer.GetName(),
	}).Debug("Indexer loaded")
}

func GetTorrentForEpisode(show string, season int, episode int) ([]Torrent, error) {
	var torrentList []Torrent
	var err error

	for _, indexer := range indexers {
		indexerSearch, err := indexer.GetTorrentForEpisode(show, season, episode)
		if err != nil {
			log.WithFields(log.Fields{
				"indexer": indexer.GetName(),
				"episode": fmt.Sprintf("%v S%03dE%03d", show, season, episode),
				"error":   err,
			}).Warning("Couldn't get torrents from indexer")
			continue
		}

		if len(indexerSearch) != 0 {
			torrentList = append(torrentList, indexerSearch...)
			log.WithFields(log.Fields{
				"indexer": indexer.GetName(),
				"episode": fmt.Sprintf("%v S%03dE%03d", show, season, episode),
			}).Info(len(indexerSearch), " torrents found")
		} else {
			log.WithFields(log.Fields{
				"indexer": indexer.GetName(),
				"episode": fmt.Sprintf("%v S%03dE%03d", show, season, episode),
			}).Info("No torrents found")
		}
	}

	sort.Slice(torrentList[:], func(i, j int) bool {
		iValue, _ := strconv.Atoi(torrentList[i].Attributes["seeders"])
		jValue, _ := strconv.Atoi(torrentList[j].Attributes["seeders"])

		return iValue > jValue
	})

	return torrentList, err
}
