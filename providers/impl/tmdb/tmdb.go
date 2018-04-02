package tmdb

import (
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

	module = Module{
		Name: "TMDB",
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
	show, _ := tmdbProvider.Client.GetTvInfo(tvShow.Id, nil)

	season, err := tmdbProvider.Client.GetTvSeasonInfo(show.ID, show.NumberOfSeasons, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"tvshow_name": tvShow.Name,
			"id":          tvShow.Id,
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

	// Set season number on episodes
	for i := range filteredEpisodes {
		filteredEpisodes[i].Season = season.SeasonNumber
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
		Id:         tvShow.ID,
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
		Name:     episode.Name,
		Date:     firstAired,
		Id:       episode.ID,
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
		Id:            movie.ID,
	}
}
