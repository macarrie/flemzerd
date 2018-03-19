package indexer

import (
	"bytes"
	"errors"
	"fmt"
	"sort"

	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
)

var indexersCollection []Indexer

func AddIndexer(indexer Indexer) {
	indexersCollection = append(indexersCollection, indexer)
	log.WithFields(log.Fields{
		"indexer": indexer.GetName(),
	}).Debug("Indexer loaded")
}

func Status() ([]Module, error) {
	var modList []Module
	var aggregatedErrorMessage bytes.Buffer

	for _, indexer := range indexersCollection {
		mod, indexerAliveError := indexer.Status()
		if indexerAliveError != nil {
			log.WithFields(log.Fields{
				"error": indexerAliveError,
			}).Warning("Indexer is not alive")
			aggregatedErrorMessage.WriteString(indexerAliveError.Error())
			aggregatedErrorMessage.WriteString("\n")
		}
		modList = append(modList, mod)
	}

	var retError error
	if aggregatedErrorMessage.Len() == 0 {
		retError = nil
	} else {
		retError = errors.New(aggregatedErrorMessage.String())
	}
	return modList, retError
}

func GetTorrentForEpisode(show string, season int, episode int) ([]Torrent, error) {
	var torrentList []Torrent
	var err error

	for _, indexer := range indexersCollection {
		_, ok := indexer.(TVIndexer)
		if !ok {
			log.WithFields(log.Fields{
				"indexer": indexer.GetName(),
			}).Debug("Indexer does not support movies, skipping")
			continue
		}

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

	if len(torrentList) == 0 {
		return []Torrent{}, errors.New("No torrents found for episode")
	}

	sort.Slice(torrentList[:], func(i, j int) bool {
		return torrentList[i].Seeders > torrentList[j].Seeders
	})

	return torrentList, err
}

func GetTorrentForMovie(movieName string) ([]Torrent, error) {
	var torrentList []Torrent
	var err error

	for _, indexer := range indexersCollection {
		ind, ok := indexer.(MovieIndexer)
		if !ok {
			log.WithFields(log.Fields{
				"indexer": indexer.GetName(),
			}).Debug("Indexer does not support movies, skipping")
			continue
		}

		indexerSearch, err := ind.GetTorrentForMovie(movieName)
		if err != nil {
			log.WithFields(log.Fields{
				"indexer": indexer.GetName(),
				"movie":   movieName,
				"error":   err,
			}).Warning("Couldn't get torrents from indexer")
			continue
		}

		if len(indexerSearch) != 0 {
			torrentList = append(torrentList, indexerSearch...)
			log.WithFields(log.Fields{
				"indexer": ind.GetName(),
				"movie":   movieName,
			}).Info(len(indexerSearch), " torrents found")
		} else {
			log.WithFields(log.Fields{
				"indexer": ind.GetName(),
				"movie":   movieName,
			}).Info("No torrents found")
		}
	}

	if len(torrentList) == 0 {
		return []Torrent{}, errors.New("No torrents found for movie")
	}

	sort.Slice(torrentList[:], func(i, j int) bool {
		return torrentList[i].Seeders > torrentList[j].Seeders
	})

	return torrentList, err
}
