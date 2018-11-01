package provider

import (
	"fmt"
	"sort"

	"github.com/macarrie/flemzerd/db"
	media_helper "github.com/macarrie/flemzerd/helpers/media"
	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"

	watchlist "github.com/macarrie/flemzerd/watchlists"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

var providersCollection []Provider
var TVShows []TvShow
var Movies []Movie

func Status() ([]Module, error) {
	var modList []Module
	var errorList *multierror.Error

	for _, provider := range providersCollection {
		mod, providerAliveError := provider.Status()
		if providerAliveError != nil {
			log.WithFields(log.Fields{
				"error": providerAliveError,
			}).Warning("Provider is not alive")
			errorList = multierror.Append(errorList, providerAliveError)
		}
		modList = append(modList, mod)
	}

	return modList, errorList.ErrorOrNil()
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
		show, err := (*p).GetShow(ids)
		show.MediaIds = mergeMediaIds(ids, show.MediaIds)
		if err != nil {
			return TvShow{}, err
		}
		showReq := TvShow{}
		req := db.Client.Where("media_ids_id = ?", ids.Model.ID).Find(&showReq)
		if req.RecordNotFound() {
			//Look in deleted records too
			unscopedReq := db.Client.Unscoped().Where("media_ids_id = ?", ids.Model.ID).Find(&showReq)
			if unscopedReq.RecordNotFound() {
				db.Client.Create(&show)
				return show, nil
			}
		}

		return showReq, nil
	}

	return TvShow{}, errors.New("Cannot find any TV provider in configuration")
}

func FindMovie(query MediaIds) (Movie, error) {
	p := getMovieProvider()
	if p != nil {
		movie, err := (*p).GetMovie(query)
		movie.MediaIds = query
		if err != nil {
			return Movie{}, err
		}
		movieReq := Movie{}
		req := db.Client.Where("media_ids_id = ?", query.Model.ID).Find(&movieReq)
		if req.RecordNotFound() {
			//Look in deleted records too
			unscopedReq := db.Client.Unscoped().Where("media_ids_id = ?", query.Model.ID).Find(&movieReq)
			if unscopedReq.RecordNotFound() {
				db.Client.Create(&movie)
				return movie, nil
			}
		}

		return movieReq, nil
	}

	return Movie{}, errors.New("Cannot find any movie provider in configuration")
}

func FindRecentlyAiredEpisodesForShow(show TvShow) ([]Episode, error) {
	p := getTVProvider()
	if p != nil {
		recentEpisodes, err := (*p).GetRecentlyAiredEpisodes(show)
		if err != nil {
			return []Episode{}, err
		}

		var retList []Episode
		for _, recentEpisode := range recentEpisodes {
			for _, p := range providers {
				e, err := (*p).GetEpisode(show.MediaIds, recentEpisode.Season, recentEpisode.Number)
				if err != nil {
					continue
				}

				reqEpisode := Episode{}
				req := db.Client.Where(Episode{
					Name:   e.Title,
					Season: e.Season,
					Number: e.Number,
				}).Find(&reqEpisode)
				if req.RecordNotFound() {
					e.TvShow = show
					db.Client.Create(&e)
				} else {
					e = reqEpisode
				}

				if e.TvShow.IsAnime {
					previousSeasonsEpisodeCount := 0
					for _, season := range show.Seasons {
						if season.SeasonNumber < e.Season {
							previousSeasonsEpisodeCount += season.EpisodeCount
						}
					}
					if e.Number > previousSeasonsEpisodeCount {
						e.AbsoluteNumber = e.Number
						e.Number = e.Number - previousSeasonsEpisodeCount
					} else {
						e.AbsoluteNumber = e.Number + previousSeasonsEpisodeCount
					}
				}
				db.Client.Save(&e)
			}

			retList = append(retList, e)
		}
		return retList, nil
	}

	return []Episode{}, errors.New("Cannot find any TV provider in configuration")
}

func GetSeasonEpisodeList(show TvShow, seasonNumber int) ([]Episode, error) {
	p := getTVProvider()
	if p != nil {
		episodes, err := (*p).GetSeasonEpisodeList(show, seasonNumber)
		if err != nil {
			return []Episode{}, err
		}
		var retList []Episode
		for _, ep := range episodes {
			var epFromDb Episode
			req := db.Client.Where(ep).Find(&epFromDb)
			if req.RecordNotFound() {
				ep.TvShow = show
				db.Client.Create(&ep)
				retList = append(retList, ep)
			} else {
				retList = append(retList, epFromDb)
			}
		}

		return retList, nil
	}

	return []Episode{}, errors.New("Cannot find any TV provider in configuration")
}

func GetTVShowsInfoFromConfig() {
	var showObjects []TvShow
	var showList []MediaIds

	showsFromWatchlists, _ := watchlist.GetTvShows()
	showList = append(showList, showsFromWatchlists...)

	for _, showIds := range showList {
		showTitle := showIds.Title
		show, err := FindShow(showIds)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"show":  showTitle,
			}).Warning("Unable to get show informations")
		} else {
			if show.DeletedAt == nil {
				showObjects = append(showObjects, show)
			}
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
		movieTitle := movie
		movie, err := FindMovie(movie)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"movie": movieTitle,
			}).Warning("Unable to get movie informations")
		} else {
			if movie.DeletedAt == nil {
				movieObjects = append(movieObjects, movie)
			}
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
		if !occurences[show.Title] {
			occurences[show.Title] = true
			ret = append(ret, show)
		}
	}

	return ret
}

func removeDuplicateMovies(array []Movie) []Movie {
	occurences := make(map[string]bool)
	var ret []Movie

	for _, movie := range array {
		if !occurences[media_helper.GetMovieTitle(movie)] {
			occurences[media_helper.GetMovieTitle(movie)] = true
			ret = append(ret, movie)
		}
	}

	return ret
}

func getTVProvider() *TVProvider {
	var list []*TVProvider
	for _, p := range providersCollection {
		tvProvider, ok := p.(TVProvider)
		if ok {
			list = append(list, &tvProvider)
		}
	}

	sort.Slice(list, func(i, j int) bool {
		return (*list[i]).GetOrder() < (*list[j]).GetOrder()
	})
	if len(list) > 0 {
		return list[0]
	}

	return nil
}

func getMovieProvider() *MovieProvider {
	for _, p := range providersCollection {
		movieProvider, ok := p.(MovieProvider)
		if ok {
			return &movieProvider
		}
	}
	return nil
}

func mergeMediaIds(m1, m2 MediaIds) MediaIds {
	merged := MediaIds{}

	if m1.ID != 0 {
		merged.ID = m1.ID
	}
	if m1.Title != "" {
		merged.Title = m1.Title
	}
	if m1.Trakt != 0 {
		merged.Trakt = m1.Trakt
	}
	if m1.Tmdb != 0 {
		merged.Tmdb = m1.Tmdb
	}
	if m1.Imdb != "" {
		merged.Imdb = m1.Imdb
	}
	if m1.Tvdb != 0 {
		merged.Tvdb = m1.Tvdb
	}

	if merged.ID != 0 && m2.ID != 0 {
		merged.ID = m2.ID
	}
	if merged.Title != "" && m2.Title != "" {
		merged.Title = m2.Title
	}
	if merged.Trakt == 0 && m2.Trakt != 0 {
		merged.Trakt = m2.Trakt
	}
	if merged.Tmdb == 0 && m2.Tmdb != 0 {
		merged.Tmdb = m2.Tmdb
	}
	if merged.Imdb == "" && m2.Imdb != "" {
		merged.Imdb = m2.Imdb
	}
	if merged.Tvdb == 0 && m2.Tvdb != 0 {
		merged.Tvdb = m2.Tvdb
	}

	return merged
}

// GetProvider returns the registered provider with name "name". An non-nil error is returned if no registered provider are found with the required name
func GetProvider(name string) (Provider, error) {
	for _, p := range providersCollection {
		mod, _ := p.Status()
		if mod.Name == name {
			return p, nil
		}
	}

	return nil, fmt.Errorf("Provider %s not found in configuration", name)
}
