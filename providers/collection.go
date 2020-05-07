package provider

import (
	"fmt"
	"github.com/macarrie/flemzerd/cache"
	"github.com/macarrie/flemzerd/db"
	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
	"github.com/odwrtw/fanarttv"
	"sort"
	"sync"

	watchlist "github.com/macarrie/flemzerd/watchlists"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

var providersCollection []Provider
var TVShows []TvShow
var Movies []Movie
var watchlistRefreshMutex sync.Mutex

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
		req := db.Client.Unscoped().Where("media_ids_id = ?", ids.Model.ID).Find(&showReq)
		if req.RecordNotFound() {
			db.Client.Create(&show)

			posterCachePath, err := cache.CacheShowPoster(&show)
			if err != nil {
				show.GetLog().WithFields(log.Fields{
					"err": err,
				}).Error("Could not cache poster")
			}

			if show.Background == "" {
				url, err := GetShowFanart(&show)
				if err != nil {
					show.GetLog().WithFields(log.Fields{
						"err": err,
					}).Error("Could not get fanart for movie")
				}
				show.Background = url
			}
			fanartCachePath, err := cache.CacheShowFanart(&show)
			if err != nil {
				show.GetLog().WithFields(log.Fields{
					"err": err,
				}).Error("Could not cache fanart")
			}

			show.Poster = posterCachePath
			show.Background = fanartCachePath
			db.Client.Save(&show)

			return show, nil
		}

		if showReq.DeletedAt != nil {
			return showReq, nil
		}

		// Update show saved in db with info retrieved from provider
		//mergeRecentShowProperties(&showReq, &show)
		posterCachePath, err := cache.CacheShowPoster(&showReq)
		if err != nil {
			showReq.GetLog().WithFields(log.Fields{
				"err": err,
			}).Error("Could not cache poster")
		}

		if showReq.Background == "" {
			url, err := GetShowFanart(&showReq)
			if err != nil {
				showReq.GetLog().WithFields(log.Fields{
					"err": err,
				}).Error("Could not get fanart")
			}
			showReq.Background = url
		}
		fanartCachePath, err := cache.CacheShowFanart(&showReq)
		if err != nil {
			showReq.GetLog().WithFields(log.Fields{
				"err": err,
			}).Error("Could not cache fanart")
		}

		showReq.Poster = posterCachePath
		showReq.Background = fanartCachePath
		db.Client.Save(&showReq)

		return showReq, nil
	}

	return TvShow{}, errors.New("Cannot find any TV provider in configuration")
}

func FindMovie(query MediaIds) (Movie, error) {
	p := getMovieProvider()
	if p != nil {
		movie, err := (*p).GetMovie(query)
		if err != nil {
			return Movie{}, err
		}
		movie.MediaIds = query
		movieReq := Movie{}
		req := db.Client.Unscoped().Where("media_ids_id = ?", query.Model.ID).Find(&movieReq)
		if req.RecordNotFound() {
			db.Client.Create(&movie)
			posterCachePath, err := cache.CacheMoviePoster(&movie)
			if err != nil {
				movie.GetLog().WithFields(log.Fields{
					"err": err,
				}).Error("Could not cache poster")
			}

			if movie.Background == "" {
				url, err := GetMovieFanart(&movie)
				if err != nil {
					movie.GetLog().WithFields(log.Fields{
						"err": err,
					}).Error("Could not get fanart")
				}
				movie.Background = url
			}
			fanartCachePath, err := cache.CacheMovieFanart(&movie)
			if err != nil {
				movie.GetLog().WithFields(log.Fields{
					"err": err,
				}).Error("Could not cache fanart")
			}

			movie.Background = fanartCachePath
			movie.Poster = posterCachePath
			db.Client.Save(&movie)
			return movie, nil
		}

		if movieReq.DeletedAt != nil {
			return movieReq, nil
		}

		// Update movie saved in db with info retrieved from provider
		//mergeRecentMovieProperties(&movieReq, &movie)
		if movieReq.Background == "" {
			url, err := GetMovieFanart(&movieReq)
			if err != nil {
				movieReq.GetLog().WithFields(log.Fields{
					"err": err,
				}).Error("Could not get fanart for movie")
			}
			movieReq.Background = url
		}
		fanartCachePath, err := cache.CacheMovieFanart(&movieReq)
		if err != nil {
			movieReq.GetLog().WithFields(log.Fields{
				"err": err,
			}).Error("Could not cache fanart")
		}
		posterCachePath, err := cache.CacheMoviePoster(&movieReq)
		if err != nil {
			movieReq.GetLog().WithFields(log.Fields{
				"err": err,
			}).Error("Could not cache poster")
		}

		movieReq.Poster = posterCachePath
		movieReq.Background = fanartCachePath
		db.Client.Save(&movieReq)

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
			e, err := (*p).GetEpisode(show.MediaIds, recentEpisode.Season, recentEpisode.Number)
			if err != nil {
				continue
			}

			reqEpisode := Episode{}
			req := db.Client.Where(Episode{
				TvShowID: show.ID,
				Season:   e.Season,
				Number:   e.Number,
			}).Find(&reqEpisode)
			if req.RecordNotFound() {
				log.WithFields(log.Fields{
					"name":    show.GetTitle(),
					"season":  recentEpisode.Season,
					"episode": recentEpisode.Number,
				}).Debug("Storing new episode in DB")

				e.TvShow = show
				db.Client.Create(&e)
			} else {
				log.WithFields(log.Fields{
					"name":    show.GetTitle(),
					"season":  recentEpisode.Season,
					"episode": recentEpisode.Number,
				}).Debug("Episode already stored in DB. Updating informations")

				// Episodes did not have absolute number computed. Absolute Number is updated when querying episode
				reqEpisode.AbsoluteNumber = e.AbsoluteNumber
				e = reqEpisode
			}
			db.Client.Save(&e)

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
	watchlistRefreshMutex.Lock()
	defer watchlistRefreshMutex.Unlock()

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
	watchlistRefreshMutex.Lock()
	defer watchlistRefreshMutex.Unlock()

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
		if !occurences[movie.GetTitle()] {
			occurences[movie.GetTitle()] = true
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

func GetShowFanart(s *TvShow) (url string, err error) {
	s.GetLog().Debug("Looking for fanart")

	client := fanarttv.New(FANART_TV_KEY)
	if client == nil {
		return "", errors.New("could not create fanart API client")
	}

	res, err := client.GetShowImages(fmt.Sprintf("%d", s.MediaIds.Tvdb))
	if err != nil {
		return "", errors.Wrap(err, "could not find fanart for show")
	}

	if len(res.Backgrounds) == 0 {
		return "", errors.New("no background fanart available for show")
	}

	return fanarttv.Best(res.Backgrounds).URL, nil
}

func GetMovieFanart(m *Movie) (url string, err error) {
	m.GetLog().Debug("Looking for fanart")

	client := fanarttv.New(FANART_TV_KEY)
	if client == nil {
		return "", errors.New("could not create fanart API client")
	}

	res, err := client.GetMovieImages(m.MediaIds.Imdb)
	if err != nil {
		return "", errors.Wrap(err, "could not find fanart for movie")
	}

	if len(res.Backgrounds) == 0 {
		return "", errors.New("no background fanart available for movie")
	}

	return fanarttv.Best(res.Backgrounds).URL, nil
}

// Updates currentMovie (saved in DB) with updates retrieved from Provider
func mergeRecentShowProperties(currentShow *TvShow, updatedShow *TvShow) {
	currentShow.Title = updatedShow.Title
	currentShow.OriginalTitle = updatedShow.OriginalTitle
	currentShow.Overview = updatedShow.Overview
	currentShow.Banner = updatedShow.Banner
	currentShow.NumberOfSeasons = updatedShow.NumberOfSeasons
	currentShow.NumberOfEpisodes = updatedShow.NumberOfEpisodes
	currentShow.FirstAired = updatedShow.FirstAired
	currentShow.Status = updatedShow.Status
	if !currentShow.IsCached(currentShow.Poster) {
		currentShow.Poster = updatedShow.Poster
	}
}

// Updates currentMovie (saved in DB) with updates retrieved from Provider
func mergeRecentMovieProperties(currentMovie *Movie, updatedMovie *Movie) {
	currentMovie.Title = updatedMovie.Title
	currentMovie.OriginalTitle = updatedMovie.OriginalTitle
	currentMovie.Overview = updatedMovie.Overview
	if !currentMovie.IsCached(currentMovie.Poster) {
		currentMovie.Poster = updatedMovie.Poster
	}
}
