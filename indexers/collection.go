package indexer

import (
	"fmt"
	"sort"

	"github.com/macarrie/flemzerd/configuration"
	media_helper "github.com/macarrie/flemzerd/helpers/media"
	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
	"github.com/macarrie/flemzerd/vidocq"

	multierror "github.com/hashicorp/go-multierror"
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
	var errorList *multierror.Error

	for _, indexer := range indexersCollection {
		mod, indexerAliveError := indexer.Status()
		if indexerAliveError != nil {
			log.WithFields(log.Fields{
				"error": indexerAliveError,
			}).Warning("Indexer is not alive")
			errorList = multierror.Append(errorList, indexerAliveError)
		}
		modList = append(modList, mod)
	}

	return modList, errorList.ErrorOrNil()
}

func Reset() {
	indexersCollection = []Indexer{}
}

func GetTorrentForEpisode(show string, season int, episode int) ([]Torrent, error) {
	var torrentList []Torrent
	var errorList *multierror.Error

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
			errorList = multierror.Append(errorList, err)
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

	return torrentList, errorList.ErrorOrNil()
}

func GetTorrentForMovie(movie Movie) ([]Torrent, error) {
	var torrentList []Torrent
	var errorList *multierror.Error

	for _, indexer := range indexersCollection {
		ind, ok := indexer.(MovieIndexer)
		if !ok {
			log.WithFields(log.Fields{
				"indexer": indexer.GetName(),
			}).Debug("Indexer does not support movies, skipping")
			continue
		}

		indexerSearch, err := ind.GetTorrentForMovie(media_helper.GetMovieTitle(movie))
		if err != nil {
			log.WithFields(log.Fields{
				"indexer": indexer.GetName(),
				"movie":   media_helper.GetMovieTitle(movie),
				"error":   err,
			}).Warning("Couldn't get torrents from indexer")
			errorList = multierror.Append(errorList, err)
			continue
		}

		if len(indexerSearch) != 0 {
			torrentList = append(torrentList, indexerSearch...)
			log.WithFields(log.Fields{
				"indexer": ind.GetName(),
				"movie":   media_helper.GetMovieTitle(movie),
			}).Info(len(indexerSearch), " torrents found")
		} else {
			log.WithFields(log.Fields{
				"indexer": ind.GetName(),
				"movie":   media_helper.GetMovieTitle(movie),
			}).Info("No torrents found")
		}
	}

	sort.Slice(torrentList[:], func(i, j int) bool {
		return torrentList[i].Seeders > torrentList[j].Seeders
	})

	torrentList = ApplyUsersPreferencesOnTorrents(torrentList)
	if movie.Date.Year() != 1 {
		torrentList = CheckYearOfTorrents(torrentList, movie.Date.Year())
	}

	return torrentList, errorList.ErrorOrNil()
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
			log.Warning("An error occured during vidocq request", err.Error())
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

// GetIndexer returns the registered indexer with name "name". An non-nil error is returned if no registered indexer are found with the required name
func GetIndexer(name string) (Indexer, error) {
	for _, ind := range indexersCollection {
		mod, _ := ind.Status()
		if mod.Name == name {
			return ind, nil
		}
	}

	return nil, fmt.Errorf("Indexer %s not found in configuration", name)
}
