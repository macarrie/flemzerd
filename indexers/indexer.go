package indexer

import (
	log "flemzerd/logging"
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
	}).Info("Indexer loaded")
}

func GetTorrentForEpisode(show string, season int, episode int) ([]Torrent, error) {
	torrentList, err := indexers[0].GetTorrentForEpisode(show, season, episode)
	if err != nil {
		return []Torrent{}, err
	}

	sort.Slice(torrentList[:], func(i, j int) bool {
		iValue, _ := strconv.Atoi(torrentList[i].Attributes["seeders"])
		jValue, _ := strconv.Atoi(torrentList[j].Attributes["seeders"])

		return iValue > jValue
	})

	return torrentList, err
}
