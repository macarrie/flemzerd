package indexer

import (
	//"errors"
	log "flemzerd/logging"
)

type Indexer interface {
	GetName() string
	GetTorrentForEpisode(show string, season string, episode string) ([]Torrent, error)
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

func GetTorrentForEpisode(show string, season string, episode string) {
	indexers[0].GetTorrentForEpisode(show, season, episode)
}
