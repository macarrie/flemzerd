// Package watchlist groups methods for retrieving items (shows and movies) form multiple watchlists
// Multiple watchlist types can be registered. Retrieving items from watchlists will aggregate results from all registered watchlists.
package watchlist

import (
	"fmt"

	"github.com/macarrie/flemzerd/db"
	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"

	multierror "github.com/hashicorp/go-multierror"
)

var watchlistsCollection []Watchlist

// Status checks registered watchlists status. A module list is returned, each module corresponds to a registered watchlist. A non nil error is returned if at least one registered watchlist is in error
func Status() ([]Module, error) {
	var modList []Module
	var errorList *multierror.Error

	for _, watchlist := range watchlistsCollection {
		mod, watchlistAliveError := watchlist.Status()
		if watchlistAliveError != nil {
			log.WithFields(log.Fields{
				"error": watchlistAliveError,
			}).Warning("Watchlist is not alive")
			errorList = multierror.Append(errorList, watchlistAliveError)
		}
		modList = append(modList, mod)
	}

	return modList, errorList.ErrorOrNil()
}

// Reset empties registered watchlists list
func Reset() {
	watchlistsCollection = []Watchlist{}
}

// AddWatchlist registers a new watchlist
func AddWatchlist(watchlist Watchlist) {
	watchlistsCollection = append(watchlistsCollection, watchlist)
	log.WithFields(log.Fields{
		"watchlist": watchlist.GetName(),
	}).Debug("Watchlist loaded")
}

// GetTvShows retrieves TV Shows from all watchlists then aggregates the result. Retrieved shows are added to the database to ease future requests.
// Duplicates are removed.
// Results are returned as an array of MediaIds structs
func GetTvShows() ([]MediaIds, error) {
	tvshows := []MediaIds{}
	for _, watchlist := range watchlistsCollection {
		shows, err := watchlist.GetTvShows()
		if err != nil {
			log.WithFields(log.Fields{
				"watchlist": watchlist.GetName(),
				"error":     err,
			}).Warning("Couldn't get TV shows from watchlist")
			continue
		}
		tvshows = append(tvshows, shows...)
	}

	idsFromDb := []MediaIds{}
	showsFromDb := []TvShow{}
	db.Client.Find(&showsFromDb)
	for _, show := range showsFromDb {
		idsFromDb = append(idsFromDb, show.MediaIds)
	}

	tvshows = append(tvshows, idsFromDb...)

	tvshows = removeDuplicates(tvshows)

	//Return elements saved into Db
	retList := []MediaIds{}
	for _, showIds := range tvshows {
		idsFromDb := MediaIds{}
		req := db.Client.Where("name = ?", showIds.Name).Find(&idsFromDb)
		if req.RecordNotFound() {
			db.Client.Create(&showIds)
			retList = append(retList, showIds)
		} else {
			retList = append(retList, idsFromDb)
		}
	}

	return retList, nil
}

// GetMovies retrieves movies from all watchlists then aggregates the result. Retrieved movies are added to the database to ease future requests.
// Duplicates are removed.
// Results are returned as an array of MediaIds structs
func GetMovies() ([]MediaIds, error) {
	movieWatchlist := []MediaIds{}
	for _, watchlist := range watchlistsCollection {
		movies, err := watchlist.GetMovies()
		if err != nil {
			log.WithFields(log.Fields{
				"watchlist": watchlist.GetName(),
				"error":     err,
			}).Warning("Couldn't get movies from watchlist")
			continue
		}
		movieWatchlist = append(movieWatchlist, movies...)
	}

	movieWatchlist = removeDuplicates(movieWatchlist)

	//Return elements saved into Db
	retList := []MediaIds{}
	for _, movieIds := range movieWatchlist {
		idsFromDb := MediaIds{}
		req := db.Client.Where("name = ?", movieIds.Name).Find(&idsFromDb)
		if req.RecordNotFound() {
			db.Client.Create(&movieIds)
			retList = append(retList, movieIds)
		} else {
			retList = append(retList, idsFromDb)
		}
	}

	return retList, nil
}

func removeDuplicates(array []MediaIds) []MediaIds {
	occurences := make(map[string]bool)
	var ret []MediaIds

	for _, media := range array {
		if !occurences[media.Name] {
			occurences[media.Name] = true
			ret = append(ret, media)
		}
	}

	return ret
}

// GetWatchlist returns the registered watchlist with name "name". An non-nil error is returned if no registered watchlists are found with the required name
func GetWatchlist(name string) (Watchlist, error) {
	for _, w := range watchlistsCollection {
		mod, _ := w.Status()
		if mod.Name == name {
			return w, nil
		}
	}

	return nil, fmt.Errorf("Watchlist %s not found in configuration", name)
}
