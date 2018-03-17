package watchlist

import (
	"bytes"
	"errors"
	"fmt"
	"sort"

	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
)

var watchlistsCollection []Watchlist

func Status() ([]Module, error) {
	var modList []Module
	var aggregatedErrorMessage bytes.Buffer

	for _, watchlist := range watchlistsCollection {
		mod, watchlistAliveError := watchlist.Status()
		if watchlistAliveError != nil {
			log.WithFields(log.Fields{
				"error": watchlistAliveError,
			}).Warning("Watchlist is not alive")
			aggregatedErrorMessage.WriteString(watchlistAliveError.Error())
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

func AddWatchlist(watchlist Watchlist) {
	watchlistsCollection = append(watchlistsCollection, watchlist)
	log.WithFields(log.Fields{
		"watchlist": watchlist,
	}).Debug("Watchlist loaded")
}

func GetTvShows() ([]string, error) {
	tvshows := []string{}
	for _, watchlist := range watchlistsCollection {
		shows, err := watchlist.GetTvShows()
		if err != nil {
			log.WithFields(log.Fields{
				"watchlist": watchlist,
				"error":     err,
			}).Warning("Couldn't get TV shows from watchlist")
			continue
		}
		tvshows = append(tvshows, shows...)
	}

	tvshows = removeDuplicates(tvshows)
	sort.Strings(tvshows)

	return tvshows, nil
}

func GetMovies() ([]string, error) {
	movieWatchlist := []string{}
	for _, watchlist := range watchlistsCollection {
		movies, err := watchlist.GetMovies()
		if err != nil {
			log.WithFields(log.Fields{
				"watchlist": watchlist,
				"error":     err,
			}).Warning("Couldn't get movies from watchlist")
			continue
		}
		movieWatchlist = append(movieWatchlist, movies...)
	}

	movieWatchlist = removeDuplicates(movieWatchlist)
	sort.Strings(movieWatchlist)

	return movieWatchlist, nil
}

func removeDuplicates(array []string) []string {
	occurences := make(map[string]bool)
	var ret []string

	for _, str := range array {
		if !occurences[str] {
			occurences[str] = true
			ret = append(ret, str)
		}
	}

	return ret
}

func GetWatchlist(name string) (Watchlist, error) {
	for _, w := range watchlistsCollection {
		mod, _ := w.Status()
		if mod.Name == name {
			return w, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("Watchlist %s not found in configuration", name))
}
