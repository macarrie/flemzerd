package tmdb

import (
	"errors"
	"fmt"
	"time"

	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
	tmdb "github.com/ryanbradynd05/go-tmdb"
)

type TMDBProvider struct {
	Client *tmdb.TMDb
}

var module Module

// Create new instance of the tmdb info provider
func New(apiKey string) (tmdbProvider *TMDBProvider, err error) {
	client := tmdb.Init(apiKey)
	t := &TMDBProvider{}

	module = Module{
		Name: t.GetName(),
		Type: "provider",
		Status: ModuleStatus{
			Alive:   true,
			Message: "",
		},
	}

	return &TMDBProvider{Client: client}, nil
}

// Check if Provider is alive
func (tmdbProvider *TMDBProvider) Status() (Module, error) {
	log.Debug("Checking TMDB provider status")
	// TODO

	return module, nil
}

func (tmdbProvider *TMDBProvider) GetName() string {
	return "TMDB"
}

// Get show from name
func (tmdbProvider *TMDBProvider) GetShow(tvShow MediaIds) (TvShow, error) {
	log.WithFields(log.Fields{
		"name":     tvShow.Name,
		"provider": module.Name,
	}).Debug("Searching show")

	var id int
	if tvShow.Tmdb != 0 {
		id = tvShow.Tmdb
	} else {
		results, err := tmdbProvider.Client.SearchTv(tvShow.Name, nil)
		if err != nil {
			return TvShow{}, err
		}

		if len(results.Results) > 0 {
			id = results.Results[0].ID
		}
	}

	show, err := tmdbProvider.Client.GetTvInfo(id, nil)
	if err != nil {
		return TvShow{}, err
	}

	return convertShow(*show), nil
}

// Get list of episodes of a show aired less than RECENTLY_AIRED_EPISODES_INTERVAL days ago
func (tmdbProvider *TMDBProvider) GetRecentlyAiredEpisodes(tvShow TvShow) ([]Episode, error) {
	if tvShow.MediaIds.Tmdb == 0 {
		results, err := tmdbProvider.Client.SearchTv(tvShow.Name, nil)
		if err != nil {
			return []Episode{}, err
		}

		if results.TotalResults != 0 {
			tvShow.MediaIds.Tmdb = results.Results[0].ID
		} else {
			return []Episode{}, errors.New("Cannot find show in TMDB")
		}
	}

	show, _ := tmdbProvider.Client.GetTvInfo(int(tvShow.MediaIds.Tmdb), nil)

	season, err := tmdbProvider.Client.GetTvSeasonInfo(show.ID, show.NumberOfSeasons, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"tvshow_name": tvShow.Name,
			"id":          tvShow.Model.ID,
			"provider":    module.Name,
			"season":      show.NumberOfSeasons,
			"error":       err,
		}).Warning("Cannot get recently aired episodes of the tv show")
		return []Episode{}, err
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
				"tvshow_name": tvShow.Name,
				"id":          tvShow.Model.ID,
				"provider":    module.Name,
				"season":      ep.Season,
				"number":      ep.Number,
				"error":       err,
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
		"name":     m.Name,
		"provider": module.Name,
	}).Debug("Searching movie")

	var id int
	if m.Tmdb != 0 {
		id = m.Tmdb
	} else {
		results, err := tmdbProvider.Client.SearchMovie(m.Name, nil)
		if err != nil {
			return Movie{}, err
		}
		if len(results.Results) > 0 {
			id = results.Results[0].ID
		}

	}

	movie, _ := tmdbProvider.Client.GetMovieInfo(id, nil)
	return convertMovie(*movie), nil
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

	return TvShow{
		Poster:     fmt.Sprintf("https://image.tmdb.org/t/p/w500%s", tvShow.PosterPath),
		FirstAired: firstAired,
		Overview:   tvShow.Overview,
		Name:       tvShow.Name,
		Status:     status,
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
		Name:     episode.Name,
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
		Poster:        fmt.Sprintf("https://image.tmdb.org/t/p/w500%s", movie.PosterPath),
		Title:         movie.Title,
		OriginalTitle: movie.OriginalTitle,
		Overview:      movie.Overview,
		Date:          releaseDate,
	}
}
