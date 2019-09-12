package indexer

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/macarrie/flemzerd/downloadable"

	"github.com/macarrie/flemzerd/configuration"
	log "github.com/macarrie/flemzerd/logging"
	"github.com/macarrie/flemzerd/vidocq"

	"github.com/hashicorp/go-multierror"
)

var indexersCollection []Indexer

func AddIndexer(indexer Indexer) {
	indexersCollection = append(indexersCollection, indexer)
	log.WithFields(log.Fields{
		"indexer": indexer.GetName(),
	}).Debug("Indexer loaded")
}

func Status() ([]Module, error) {
	var modChan = make(chan Module, len(indexersCollection))
	var errorList *multierror.Error

	var wg sync.WaitGroup
	var errorListMutex sync.Mutex

	for i, _ := range indexersCollection {
		wg.Add(1)
		go func(indexer Indexer) {
			defer wg.Done()

			mod, indexerAliveError := indexer.Status()
			modChan <- mod
			if indexerAliveError != nil {
				log.WithFields(log.Fields{
					"error":   indexerAliveError,
					"indexer": indexer.GetName(),
				}).Warning("Indexer is not alive")

				errorListMutex.Lock()
				errorList = multierror.Append(errorList, indexerAliveError)
				errorListMutex.Unlock()
			}
		}(indexersCollection[i])
	}

	wg.Wait()
	close(modChan)
	var modList []Module
	for mod := range modChan {
		modList = append(modList, mod)
	}

	return modList, errorList.ErrorOrNil()
}

func Reset() {
	indexersCollection = []Indexer{}
}

func GetTorrents(d downloadable.Downloadable) ([]Torrent, error) {
	var torrentList []Torrent
	var errorList *multierror.Error
	var totalError bool = true

	for _, indexer := range indexersCollection {
		indexerSearch, err := indexer.GetTorrents(d)
		if err != nil {
			d.GetLog().WithFields(log.Fields{
				"indexer": indexer.GetName(),
				"error":   err,
			}).Warning("Couldn't get torrents from indexer")
			errorList = multierror.Append(errorList, err)
			continue
		} else {
			totalError = false
		}

		if len(indexerSearch) != 0 {
			torrentList = append(torrentList, indexerSearch...)
			d.GetLog().WithFields(log.Fields{
				"indexer": indexer.GetName(),
				"nb":      len(indexerSearch),
			}).Info("Torrents found")
		} else {
			d.GetLog().WithFields(log.Fields{
				"indexer": indexer.GetName(),
			}).Info("No torrents found")
		}
	}

	sort.Slice(torrentList[:], func(i, j int) bool {
		return torrentList[i].Seeders > torrentList[j].Seeders
	})

	switch d.(type) {
	case *Movie:
		torrentList = FilterMovieTorrents(*d.(*Movie), torrentList)
	case *Episode:
		torrentList = FilterEpisodeTorrents(*d.(*Episode), torrentList)
	}

	if totalError {
		return torrentList, errorList.ErrorOrNil()
	}

	return torrentList, nil
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
	log.WithFields(log.Fields{
		"strict_check": configuration.Config.System.StrictTorrentCheck,
	}).Debug("Checking torrent list for bad episodes")

	var returnList []Torrent
	var otherTorrents []Torrent

	for _, torrent := range list {
		if episode.TvShow.IsAnime {
			if episode.AbsoluteNumber != 0 {
				match, err := regexp.Match(fmt.Sprintf("%v", episode.AbsoluteNumber), []byte(torrent.Name))
				if match && err == nil {
					returnList = append(returnList, torrent)
				}
			} else {
				otherTorrents = append(otherTorrents, torrent)
			}
		} else {
			episodeInfo, err := vidocq.GetInfo(torrent.Name)
			if err != nil {
				log.WithFields(log.Fields{
					"torrent": torrent.Name,
				}).Warning("Error while getting media info for torrent: ", err)

				otherTorrents = append(otherTorrents, torrent)
				continue
			}

			if episodeInfo.Season != 0 && episodeInfo.Season == episode.Season && episodeInfo.Episode != 0 && episodeInfo.Episode == episode.Number {
				returnList = append(returnList, torrent)
			}
		}
	}

	if configuration.Config.System.StrictTorrentCheck {
		return returnList
	}

	return append(returnList, otherTorrents...)
}

func FilterTorrentYear(list []Torrent, year int) []Torrent {
	log.WithFields(log.Fields{
		"strict_check": configuration.Config.System.StrictTorrentCheck,
	}).Debug("Checking torrent list for bad movie torrents (wrong year)")
	var returnList []Torrent
	var otherTorrents []Torrent

	for _, torrent := range list {
		movieInfo, err := vidocq.GetInfo(torrent.Name)
		if err != nil {
			log.WithFields(log.Fields{
				"torrent": torrent.Name,
			}).Warning("Error while getting media info for torrent: ", err)

			otherTorrents = append(otherTorrents, torrent)
			continue
		}

		if (movieInfo.Year != 0 && movieInfo.Year == year) || movieInfo.Year == 0 {
			returnList = append(returnList, torrent)
		}
	}

	if configuration.Config.System.StrictTorrentCheck {
		return returnList
	}

	return append(returnList, otherTorrents...)
}

func FilterTorrentReleaseType(list []Torrent) []Torrent {
	log.WithFields(log.Fields{
		"excluded_release_types": configuration.Config.System.ExcludedReleaseTypes,
	}).Debug("Excluding release types from torrent list")

	var releaseFilteredList []Torrent
	var otherTorrents []Torrent
	var releaseTypeFilters []string

	releaseTypeFilters = strings.Split(configuration.Config.System.ExcludedReleaseTypes, ",")
	for i := range releaseTypeFilters {
		releaseTypeFilters[i] = strings.TrimSpace(releaseTypeFilters[i])
	}

	for _, torrent := range list {
		mediaInfo, err := vidocq.GetInfo(torrent.Name)
		if err != nil {
			log.Warning("An error occured during vidocq request: ", err.Error())
			otherTorrents = append(otherTorrents, torrent)

			continue
		}

		releaseTypeExcluded := false
		for _, excludedReleaseType := range releaseTypeFilters {
			if excludedReleaseType != "" && excludedReleaseType == mediaInfo.ReleaseType {
				releaseTypeExcluded = true
			}
		}

		if !releaseTypeExcluded {
			releaseFilteredList = append(releaseFilteredList, torrent)
		}
	}

	if configuration.Config.System.StrictTorrentCheck {
		return releaseFilteredList
	}

	return append(releaseFilteredList, otherTorrents...)
}

func FilterTorrentQuality(list []Torrent) []Torrent {
	log.WithFields(log.Fields{
		"quality_filter": configuration.Config.System.PreferredMediaQuality,
		"strict_check":   configuration.Config.System.StrictTorrentCheck,
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

	if configuration.Config.System.StrictTorrentCheck {
		return qualityFilteredList
	}

	return append(qualityFilteredList, otherTorrents...)
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
