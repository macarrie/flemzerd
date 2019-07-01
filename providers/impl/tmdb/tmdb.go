package tmdb

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/macarrie/flemzerd/configuration"
	"github.com/macarrie/flemzerd/db"
	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
	"github.com/ryanbradynd05/go-tmdb"

	"github.com/pkg/errors"
)

type TMDBProvider struct {
	Client *tmdb.TMDb
	Order  int
}

var module Module

// Create new instance of the tmdb info provider
func New(apiKey string, order string) (tmdbProvider *TMDBProvider, err error) {
	client := tmdb.Init(apiKey)
	t := &TMDBProvider{}
	parsedOrder, err := strconv.Atoi(order)
	if err != nil {
		log.WithFields(log.Fields{
			"provider": t.GetName(),
			"order":    order,
		}).Error("Incorrect order defined in configuration.")
		parsedOrder = 1
	}

	module = Module{
		Name: t.GetName(),
		Type: "provider",
		Status: ModuleStatus{
			Alive:   true,
			Message: "",
		},
	}

	return &TMDBProvider{Client: client, Order: parsedOrder}, nil
}

// Check if Provider is alive
func (tmdbProvider *TMDBProvider) Status() (Module, error) {
	log.Debug("Checking TMDB provider status")

	// TODO
	if configuration.TMDB_API_KEY == "" {
		module.Status.Alive = false
		module.Status.Message = "TMDB API key not found"
		return module, errors.New(module.Status.Message)
	}

	return module, nil
}

func (tmdbProvider *TMDBProvider) GetName() string {
	return "TMDB"
}

func (tmdbProvider *TMDBProvider) GetOrder() int {
	return tmdbProvider.Order
}

// Get show from name
func (tmdbProvider *TMDBProvider) GetShow(tvShow MediaIds) (TvShow, error) {
	log.WithFields(log.Fields{
		"title":    tvShow.Title,
		"provider": module.Name,
	}).Debug("Searching show")

	var id int
	if tvShow.Tmdb != 0 {
		id = tvShow.Tmdb
	} else {
		results, err := tmdbProvider.Client.SearchTv(tvShow.Title, nil)
		if err != nil {
			return TvShow{}, errors.Wrap(err, "cannot find show in TMDB")
		}

		if len(results.Results) > 0 {
			id = results.Results[0].ID
		}
	}

	show, err := tmdbProvider.Client.GetTvInfo(id, nil)
	if err != nil {
		return TvShow{}, errors.Wrap(err, "cannot get show from TMDB")
	}

	return convertShow(*show), nil
}

// Get specific episode from tvshow
func (tmdbProvider *TMDBProvider) GetEpisode(tvShowMediaIds MediaIds, seasonNb int, episodeNb int) (Episode, error) {
	log.WithFields(log.Fields{
		"name":     tvShowMediaIds.Title,
		"provider": module.Name,
		"season":   seasonNb,
		"episode":  episodeNb,
	}).Debug("Getting episode")

	var id int
	if tvShowMediaIds.Tmdb != 0 {
		id = tvShowMediaIds.Tmdb
	} else {
		results, err := tmdbProvider.Client.SearchTv(tvShowMediaIds.Title, nil)
		if err != nil {
			return Episode{}, errors.Wrap(err, "cannot find show in TMDB")
		}

		if len(results.Results) > 0 {
			id = results.Results[0].ID
		}
	}

	episode, err := tmdbProvider.Client.GetTvEpisodeInfo(id, seasonNb, episodeNb, nil)
	if err != nil {
		return Episode{}, errors.Wrap(err, "cannot get episode from TMDB")
	}

	showFromDb := TvShow{}
	db.Client.Where("media_ids_id = ?", tvShowMediaIds.Model.ID).Find(&showFromDb)
	convertedEpisode := convertEpisode(*episode)
	getEpisodeAbsoluteNumber(showFromDb, &convertedEpisode)

	return convertedEpisode, nil
}

// Get list of episodes of a show aired less than RECENTLY_AIRED_EPISODES_INTERVAL days ago
func (tmdbProvider *TMDBProvider) GetRecentlyAiredEpisodes(tvShow TvShow) ([]Episode, error) {
	if tvShow.MediaIds.Tmdb == 0 {
		results, err := tmdbProvider.Client.SearchTv(tvShow.GetTitle(), nil)
		if err != nil {
			return []Episode{}, errors.Wrap(err, "cannot find show in TMDB")
		}

		if results.TotalResults != 0 {
			tvShow.MediaIds.Tmdb = results.Results[0].ID
		} else {
			return []Episode{}, errors.New("Cannot find show in TMDB")
		}
	}

	show, err := tmdbProvider.Client.GetTvInfo(int(tvShow.MediaIds.Tmdb), nil)
	if err != nil {
		return []Episode{}, errors.Wrap(err, "cannot get show info from TMDB")
	}

	season, err := tmdbProvider.Client.GetTvSeasonInfo(show.ID, show.NumberOfSeasons, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"tvshow":   tvShow.GetTitle(),
			"id":       tvShow.Model.ID,
			"provider": module.Name,
			"season":   show.NumberOfSeasons,
			"error":    err,
		}).Warning("Cannot get recently aired episodes of the tv show")
		return []Episode{}, errors.Wrap(err, "cannot get season info from TMDB")
	}

	var episodeList []Episode
	for _, episode := range season.Episodes {
		episodeList = append(episodeList, convertEpisode(episode))
	}

	oldestDate := time.Now().AddDate(0, 0, -RECENTLY_AIRED_EPISODES_INTERVAL)
	now := time.Now()
	filteredEpisodes := filterEpisodesAiredBetweenDates(episodeList, &oldestDate, &now)

	for i, ep := range filteredEpisodes {
		externalids, err := tmdbProvider.Client.GetTvEpisodeExternalIds(show.ID, ep.Season, ep.Number, nil)
		if err != nil {
			log.WithFields(log.Fields{
				"tvshow":   tvShow.GetTitle(),
				"id":       tvShow.Model.ID,
				"provider": module.Name,
				"season":   ep.Season,
				"number":   ep.Number,
				"error":    err,
			}).Warning("Cannot get external ids from Tmdb for episode")
		}
		filteredEpisodes[i].MediaIds = MediaIds{
			Tmdb: externalids.ID,
			Imdb: externalids.ImdbID,
			Tvdb: externalids.TvdbID,
		}
	}

	return filteredEpisodes, nil
}

func (tmdbProvider *TMDBProvider) GetMovie(m MediaIds) (Movie, error) {
	log.WithFields(log.Fields{
		"title":    m.Title,
		"provider": module.Name,
	}).Debug("Searching movie")

	var id int
	if m.Tmdb != 0 {
		id = m.Tmdb
	} else {
		results, err := tmdbProvider.Client.SearchMovie(m.Title, nil)
		if err != nil {
			return Movie{}, errors.Wrap(err, "cannot find movie in TMDB")
		}
		if len(results.Results) > 0 {
			id = results.Results[0].ID
		}

	}

	movie, err := tmdbProvider.Client.GetMovieInfo(id, nil)
	if err != nil {
		return Movie{}, errors.Wrap(err, "cannot get movie info from TMDB")
	}

	return convertMovie(*movie), nil
}

func (tmdbProvider *TMDBProvider) GetSeasonEpisodeList(show TvShow, seasonNumber int) ([]Episode, error) {
	if show.MediaIds.Tmdb == 0 {
		return []Episode{}, errors.New("No tmdb id found on show")
	}

	results, err := tmdbProvider.Client.GetTvSeasonInfo(show.MediaIds.Tmdb, seasonNumber, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"show":   show.GetTitle(),
			"season": seasonNumber,
		}).Warning("Encountered error when querying season details from TMDB")
		return []Episode{}, errors.Wrap(err, "cannot get show season info from TMDB")
	}

	var retList []Episode
	for _, ep := range results.Episodes {
		convertedEpisode := convertEpisode(ep)
		getEpisodeAbsoluteNumber(show, &convertedEpisode)
		retList = append(retList, convertedEpisode)
	}

	return retList, nil
}

func filterEpisodesAiredBetweenDates(episodes []Episode, beginning *time.Time, end *time.Time) []Episode {
	// Set beginning date to zero time if beginning date is nil
	if beginning == nil {
		*beginning = time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC)
	}
	// Set end date to year 100,000 if end date is nil
	if end == nil {
		*end = time.Date(100000, time.January, 1, 0, 0, 0, 0, time.UTC)
	}

	var retVal []Episode
	for _, episode := range episodes {
		if episode.Date.After(*beginning) && episode.Date.Before(*end) {
			retVal = append(retVal, episode)
		}
	}

	return retVal
}

func convertShow(tvShow tmdb.TV) TvShow {
	firstAired, err := time.Parse("2006-01-02", tvShow.FirstAirDate)
	if err != nil {
		firstAired = time.Time{}
	}

	var status int

	switch tvShow.Status {
	case "Returning Series":
		status = TVSHOW_RETURNING
	case "Planned", "In Production":
		status = TVSHOW_PLANNED
	case "Ended", "Canceled":
		status = TVSHOW_ENDED
	default:
		status = TVSHOW_UNKNOWN
	}

	seasons := []TvSeason{}
	for _, season := range tvShow.Seasons {
		airDate, err := time.Parse("2006-01-02", season.AirDate)
		if err != nil {
			airDate = time.Time{}
		}
		seasons = append(seasons, TvSeason{
			AirDate:      &airDate,
			EpisodeCount: season.EpisodeCount,
			SeasonNumber: season.SeasonNumber,
			PosterPath:   season.PosterPath,
		})
	}
	sort.Slice(seasons, func(i, j int) bool {
		return seasons[i].SeasonNumber < seasons[j].SeasonNumber
	})

	isAnime := false
	for _, genre := range tvShow.Genres {
		if genre.ID == 16 {
			isAnime = true
		}
	}

	return TvShow{
		Poster:           fmt.Sprintf("https://image.tmdb.org/t/p/w500%s", tvShow.PosterPath),
		FirstAired:       firstAired,
		Overview:         tvShow.Overview,
		Title:            tvShow.Name,
		OriginalTitle:    tvShow.OriginalName,
		Status:           status,
		Seasons:          seasons,
		NumberOfSeasons:  tvShow.NumberOfSeasons,
		NumberOfEpisodes: tvShow.NumberOfEpisodes,
		MediaIds: MediaIds{
			Tmdb: tvShow.ID,
		},
		UseDefaultTitle: true,
		IsAnime:         isAnime,
	}
}

func convertEpisode(episode tmdb.TvEpisode) Episode {
	firstAired, err := time.Parse("2006-01-02", episode.AirDate)
	if err != nil {
		firstAired = time.Time{}
	}

	return Episode{
		Number:   episode.EpisodeNumber,
		Season:   episode.SeasonNumber,
		Title:    episode.Name,
		Date:     firstAired,
		Overview: episode.Overview,
	}
}

func convertMovie(movie tmdb.Movie) Movie {
	releaseDate, err := time.Parse("2006-01-02", movie.ReleaseDate)
	if err != nil {
		releaseDate = time.Time{}
	}

	return Movie{
		Poster:          fmt.Sprintf("https://image.tmdb.org/t/p/w500%s", movie.PosterPath),
		Title:           movie.Title,
		OriginalTitle:   movie.OriginalTitle,
		Overview:        movie.Overview,
		Date:            releaseDate,
		UseDefaultTitle: true,
	}
}

func getEpisodeAbsoluteNumber(tvshow TvShow, e *Episode) {
	previousSeasonsEpisodeCount := 0
	for _, season := range tvshow.Seasons {
		if season.SeasonNumber < e.Season && season.SeasonNumber != 0 {
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
