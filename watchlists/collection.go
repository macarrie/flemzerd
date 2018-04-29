package watchlist

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/macarrie/flemzerd/db"
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

func Reset() {
	watchlistsCollection = []Watchlist{}
}

func AddWatchlist(watchlist Watchlist) {
	watchlistsCollection = append(watchlistsCollection, watchlist)
	log.WithFields(log.Fields{
		"watchlist": watchlist,
	}).Debug("Watchlist loaded")
}

func GetTvShows() ([]MediaIds, error) {
	tvshows := []MediaIds{}
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

	idsFromDb := []MediaIds{}
	showsFromDb := []TvShow{}
	db.Client.Find(&showsFromDb)
	for _, show := range showsFromDb {
		idsFromDb = append(idsFromDb, show.MediaIds)
	}

	tvshows = append(tvshows, idsFromDb...)

	tvshows = removeDuplicates(tvshows)
	// TODO: Sort mediaids
	//sort.Strings(tvshows)

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

func GetMovies() ([]MediaIds, error) {
	movieWatchlist := []MediaIds{}
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
	// TODO: Sort movieIds
	//sort.Strings(movieWatchlist)

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

func GetWatchlist(name string) (Watchlist, error) {
	for _, w := range watchlistsCollection {
		mod, _ := w.Status()
		if mod.Name == name {
			return w, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("Watchlist %s not found in configuration", name))
}
