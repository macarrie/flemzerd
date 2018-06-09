package indexer

import (
	"bytes"
	"errors"
	"fmt"
	"sort"

	"github.com/macarrie/flemzerd/configuration"
	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
	"github.com/macarrie/flemzerd/vidocq"
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

func Reset() {
	indexersCollection = []Indexer{}
}

func GetTorrentForEpisode(show string, season int, episode int) ([]Torrent, error) {
	var torrentList []Torrent
	var err error

	for _, indexer := range indexersCollection {
		ind, ok := indexer.(TVIndexer)
		if !ok {
			log.WithFields(log.Fields{
				"indexer": indexer.GetName(),
			}).Debug("Indexer does not support movies, skipping")
			continue
		}

		indexerSearch, err := ind.GetTorrentForEpisode(show, season, episode)
		if err != nil {
			log.WithFields(log.Fields{
				"indexer": ind.GetName(),
				"episode": fmt.Sprintf("%v S%03dE%03d", show, season, episode),
				"error":   err,
			}).Warning("Couldn't get torrents from indexer")
			continue
		}

		if len(indexerSearch) != 0 {
			torrentList = append(torrentList, indexerSearch...)
			log.WithFields(log.Fields{
				"indexer": ind.GetName(),
				"episode": fmt.Sprintf("%v S%03dE%03d", show, season, episode),
			}).Info(len(indexerSearch), " torrents found")
		} else {
			log.WithFields(log.Fields{
				"indexer": ind.GetName(),
				"episode": fmt.Sprintf("%v S%03dE%03d", show, season, episode),
			}).Info("No torrents found")
		}
	}

	sort.Slice(torrentList[:], func(i, j int) bool {
		return torrentList[i].Seeders > torrentList[j].Seeders
	})

	torrentList = ApplyUsersPreferencesOnTorrents(torrentList)
	torrentList = FilterBadTorrentsForEpisode(torrentList, season, episode)

	if len(torrentList) == 0 {
		return []Torrent{}, errors.New("No torrents found for episode")
	}

	return torrentList, err
}

func GetTorrentForMovie(movie Movie) ([]Torrent, error) {
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

		indexerSearch, err := ind.GetTorrentForMovie(movie.Title)
		if err != nil {
			log.WithFields(log.Fields{
				"indexer": indexer.GetName(),
				"movie":   movie.Title,
				"error":   err,
			}).Warning("Couldn't get torrents from indexer")
			continue
		}

		if len(indexerSearch) != 0 {
			torrentList = append(torrentList, indexerSearch...)
			log.WithFields(log.Fields{
				"indexer": ind.GetName(),
				"movie":   movie.Title,
			}).Info(len(indexerSearch), " torrents found")
		} else {
			log.WithFields(log.Fields{
				"indexer": ind.GetName(),
				"movie":   movie.Title,
			}).Info("No torrents found")
		}
	}

	if len(torrentList) == 0 {
		return []Torrent{}, errors.New("No torrents found for movie")
	}

	sort.Slice(torrentList[:], func(i, j int) bool {
		return torrentList[i].Seeders > torrentList[j].Seeders
	})

	torrentList = ApplyUsersPreferencesOnTorrents(torrentList)
	log.Warning("After quality filter: ", len(torrentList))
	torrentList = CheckYearOfTorrents(torrentList, movie.Date.Year())

	log.Warning("After year check: ", len(torrentList))

	return torrentList, err
}

func ApplyUsersPreferencesOnTorrents(list []Torrent) []Torrent {
	log.WithFields(log.Fields{
		"quality_filter": configuration.Config.System.PreferredMediaQuality,
	}).Debug("Sorting list according to quality preferences")

	var qualityFilteredList []Torrent
	var otherTorrents []Torrent
	var qualityFilter string

	switch configuration.Config.System.PreferredMediaQuality {
	case "720p", "1080p":
		qualityFilter = configuration.Config.System.PreferredMediaQuality
	default:
		qualityFilter = ""
	}

	if qualityFilter == "" {
		return list
	}

	for _, torrent := range list {
		mediaInfo, err := vidocq.GetInfo(torrent.Name)
		if err != nil {
			log.Warning("An error occured during vidocq request: %s", err.Error())
			otherTorrents = append(otherTorrents, torrent)

			continue
		}

		if mediaInfo.Quality == qualityFilter {
			qualityFilteredList = append(qualityFilteredList, torrent)
		} else {
			otherTorrents = append(otherTorrents, torrent)
		}
	}

	retList := append(qualityFilteredList, otherTorrents...)

	return retList
}

func FilterBadTorrentsForEpisode(list []Torrent, season int, episode int) []Torrent {
	log.Debug("Checking torrent list for bad episodes")
	var returnList []Torrent

	for _, torrent := range list {
		episodeInfo, err := vidocq.GetInfo(torrent.Name)
		if err != nil {
			log.WithFields(log.Fields{
				"torrent": torrent.Name,
			}).Warning("Error while getting media info for torrent: ", err)

			returnList = append(returnList, torrent)
			continue
		}

		if episodeInfo.Season != 0 && episodeInfo.Season == season && episodeInfo.Episode != 0 && episodeInfo.Episode == episode {
			returnList = append(returnList, torrent)
		}
	}

	return returnList
}

func CheckYearOfTorrents(list []Torrent, year int) []Torrent {
	log.Debug("Checking torrent list for bad movie torrents (wrong year)")
	var returnList []Torrent

	for _, torrent := range list {
		movieInfo, err := vidocq.GetInfo(torrent.Name)
		if err != nil {
			log.WithFields(log.Fields{
				"torrent": torrent.Name,
			}).Warning("Error while getting media info for torrent: ", err)

			returnList = append(returnList, torrent)
			continue
		}

		if movieInfo.Year != 0 && movieInfo.Year == year {
			returnList = append(returnList, torrent)
		}
	}

	return returnList
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
