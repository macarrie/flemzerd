package provider

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/macarrie/flemzerd/db"
	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"

	watchlist "github.com/macarrie/flemzerd/watchlists"
)

var providersCollection []Provider
var TVShows []TvShow
var Movies []Movie

func Status() ([]Module, error) {
	var modList []Module
	var aggregatedErrorMessage bytes.Buffer

	for _, provider := range providersCollection {
		mod, providerAliveError := provider.Status()
		if providerAliveError != nil {
			log.WithFields(log.Fields{
				"error": providerAliveError,
			}).Warning("Provider is not alive")
			aggregatedErrorMessage.WriteString(providerAliveError.Error())
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
	providersCollection = []Provider{}
}

func AddProvider(provider Provider) {
	providersCollection = append(providersCollection, provider)
}

func FindShow(ids MediaIds) (TvShow, error) {
	p := getTVProvider()
	if p != nil {
		show, err := p.GetShow(ids)
		show.MediaIds = ids
		if err != nil {
			return TvShow{}, err
		}
		showReq := TvShow{}
		req := db.Client.Where("media_ids_id = ?", ids.Model.ID).Find(&showReq)
		if req.RecordNotFound() {
			fmt.Printf("[DB] Show created\n")
			db.Client.Create(&show)
		}
		return show, nil
	}

	return TvShow{}, errors.New("Cannot find any TV provider in configuration")
}

func FindMovie(query MediaIds) (Movie, error) {
	p := getMovieProvider()
	if p != nil {
		movie, err := p.GetMovie(query)
		if err != nil {
			return Movie{}, err
		}
		movieReq := Movie{}
		req := db.Client.Preload("DownloadingItem").Preload("MediaIds").Where(&movie).First(&movieReq)
		if req.RecordNotFound() {
			fmt.Printf("[DB] Movie created\n")
			db.Client.Create(&movie)
			return movie, nil
		}

		return movieReq, nil
	}

	return Movie{}, errors.New("Cannot find any movie provider in configuration")
}

func FindRecentlyAiredEpisodesForShow(show TvShow) ([]Episode, error) {
	p := getTVProvider()
	if p != nil {
		return p.GetRecentlyAiredEpisodes(show)
	}

	return []Episode{}, errors.New("Cannot find any TV provider in configuration")
}

func GetTVShowsInfoFromConfig() {
	var showObjects []TvShow
	var showList []MediaIds

	showsFromWatchlists, _ := watchlist.GetTvShows()
	showList = append(showList, showsFromWatchlists...)

	for _, showIds := range showList {
		showName := showIds.Name
		show, err := FindShow(showIds)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"show":  showName,
			}).Warning("Unable to get show informations")
		} else {
			showObjects = append(showObjects, show)
		}
	}
	if len(showObjects) == 0 {
		log.Warning("No tvshows found in watchlists")
	}

	TVShows = removeDuplicateShows(showObjects)
}

func GetMoviesInfoFromConfig() {
	var movieObjects []Movie
	var movieList []MediaIds

	moviesFromWatchlists, _ := watchlist.GetMovies()
	movieList = append(movieList, moviesFromWatchlists...)

	for _, movie := range movieList {
		movieName := movie
		movie, err := FindMovie(movie)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"movie": movieName,
			}).Warning("Unable to get movie informations")
		} else {
			fmt.Printf("GET MOVIE INFO FROM PROVIDER: %+v\n", movie)
			movieObjects = append(movieObjects, movie)
		}
	}
	if len(movieObjects) == 0 {
		log.Warning("No movies found in watchlists")
	}

	Movies = removeDuplicateMovies(movieObjects)
}

func removeDuplicateShows(array []TvShow) []TvShow {
	occurences := make(map[string]bool)
	var ret []TvShow

	for _, show := range array {
		if !occurences[show.Name] {
			occurences[show.Name] = true
			ret = append(ret, show)
		}
	}

	return ret
}

func removeDuplicateMovies(array []Movie) []Movie {
	occurences := make(map[string]bool)
	var ret []Movie

	for _, movie := range array {
		if !occurences[movie.Title] {
			occurences[movie.Title] = true
			ret = append(ret, movie)
		}
	}

	return ret
}

func getTVProvider() TVProvider {
	for _, p := range providersCollection {
		tvProvider, ok := p.(TVProvider)
		if ok {
			return tvProvider
		}
	}
	return nil
}

func getMovieProvider() MovieProvider {
	for _, p := range providersCollection {
		movieProvider, ok := p.(MovieProvider)
		if ok {
			return movieProvider
		}
	}
	return nil
}
