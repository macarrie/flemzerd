package indexer

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

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

func GetTorrentForEpisode(episode Episode) ([]Torrent, error) {
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

		indexerSearch, err := ind.GetTorrentForEpisode(episode)
		if err != nil {
			log.WithFields(log.Fields{
				"indexer": ind.GetName(),
				"show":    media_helper.GetShowTitle(episode.TvShow),
				"season":  episode.Season,
				"number":  episode.Number,
				"episode": episode.Title,
				"error":   err,
			}).Warning("Couldn't get torrents from indexer")
			errorList = multierror.Append(errorList, err)
			continue
		}

		if len(indexerSearch) != 0 {
			torrentList = append(torrentList, indexerSearch...)
			log.WithFields(log.Fields{
				"indexer": ind.GetName(),
				"show":    media_helper.GetShowTitle(episode.TvShow),
				"season":  episode.Season,
				"number":  episode.Number,
				"episode": episode.Title,
			}).Info(len(indexerSearch), " torrents found")
		} else {
			log.WithFields(log.Fields{
				"indexer": ind.GetName(),
				"show":    media_helper.GetShowTitle(episode.TvShow),
				"season":  episode.Season,
				"number":  episode.Number,
				"episode": episode.Title,
			}).Info("No torrents found")
		}
	}

	sort.Slice(torrentList[:], func(i, j int) bool {
		return torrentList[i].Seeders > torrentList[j].Seeders
	})

	torrentList = FilterEpisodeTorrents(episode, torrentList)

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

	torrentList = FilterMovieTorrents(movie, torrentList)

	return torrentList, errorList.ErrorOrNil()
}

func FilterEpisodeTorrents(episode Episode, torrentList []Torrent) []Torrent {
	torrentList = FilterTorrentEpisodeNumber(torrentList, episode)
	torrentList = FilterTorrentQuality(torrentList)
	torrentList = FilterTorrentReleaseType(torrentList)

	return torrentList
}

func FilterMovieTorrents(movie Movie, torrentList []Torrent) []Torrent {
	torrentList = FilterTorrentQuality(torrentList)
	if movie.Date.Year() != 1 {
		torrentList = FilterTorrentYear(torrentList, movie.Date.Year())
	}
	torrentList = FilterTorrentReleaseType(torrentList)

	return torrentList
}

func FilterTorrentEpisodeNumber(list []Torrent, episode Episode) []Torrent {
	log.Debug("Checking torrent list for bad episodes")
	var returnList []Torrent

	for _, torrent := range list {
		if episode.TvShow.IsAnime {
			if episode.AbsoluteNumber != 0 {
				match, err := regexp.Match(fmt.Sprintf("%v", episode.AbsoluteNumber), []byte(torrent.Name))
				if match && err == nil {
					returnList = append(returnList, torrent)
				}
			} else {
				returnList = append(returnList, torrent)
			}
		} else {
			episodeInfo, err := vidocq.GetInfo(torrent.Name)
			if err != nil {
				log.WithFields(log.Fields{
					"torrent": torrent.Name,
				}).Warning("Error while getting media info for torrent: ", err)

				returnList = append(returnList, torrent)
				continue
			}

			if episodeInfo.Season != 0 && episodeInfo.Season == episode.Season && episodeInfo.Episode != 0 && episodeInfo.Episode == episode.Number {
				returnList = append(returnList, torrent)
			}
		}
	}

	return returnList
}

func FilterTorrentYear(list []Torrent, year int) []Torrent {
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

		if (movieInfo.Year != 0 && movieInfo.Year == year) || movieInfo.Year == 0 {
			returnList = append(returnList, torrent)
		}
	}

	return returnList
}

func FilterTorrentReleaseType(list []Torrent) []Torrent {
	log.WithFields(log.Fields{
		"excluded_release_types": configuration.Config.System.ExcludedReleaseTypes,
	}).Debug("Excluding release types from torrent list")

	var releaseFilteredList []Torrent
	var releaseTypeFilters []string

	releaseTypeFilters = strings.Split(configuration.Config.System.ExcludedReleaseTypes, ",")
	for i := range releaseTypeFilters {
		releaseTypeFilters[i] = strings.TrimSpace(releaseTypeFilters[i])
	}

	for _, torrent := range list {
		mediaInfo, err := vidocq.GetInfo(torrent.Name)
		if err != nil {
			log.Warning("An error occured during vidocq request: ", err.Error())
			releaseFilteredList = append(releaseFilteredList, torrent)

			continue
		}

		releaseTypeExcluded := false
		for _, releaseType := range releaseTypeFilters {
			if releaseType != "" && releaseType == mediaInfo.ReleaseType {
				releaseTypeExcluded = true
			}
		}

		if !releaseTypeExcluded {
			releaseFilteredList = append(releaseFilteredList, torrent)
		}
	}

	return releaseFilteredList
}

func FilterTorrentQuality(list []Torrent) []Torrent {
	log.WithFields(log.Fields{
		"quality_filter": configuration.Config.System.PreferredMediaQuality,
	}).Debug("Sorting list according to quality preferences")

	var qualityFilteredList []Torrent
	var otherTorrents []Torrent
	var qualityFilters []string

	qualityFilters = strings.Split(configuration.Config.System.PreferredMediaQuality, ",")
	for i := range qualityFilters {
		qualityFilters[i] = strings.TrimSpace(qualityFilters[i])
	}

	for _, torrent := range list {
		mediaInfo, err := vidocq.GetInfo(torrent.Name)
		if err != nil {
			log.Warning("An error occured during vidocq request: ", err.Error())
			otherTorrents = append(otherTorrents, torrent)

			continue
		}

		qualityMatches := false
		for _, quality := range qualityFilters {
			if quality == mediaInfo.Quality {
				qualityMatches = true
			}
		}

		if qualityMatches {
			qualityFilteredList = append(qualityFilteredList, torrent)
		} else {
			otherTorrents = append(otherTorrents, torrent)
		}
	}

	retList := append(qualityFilteredList, otherTorrents...)

	return retList
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
