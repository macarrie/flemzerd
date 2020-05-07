package mock

import (
	"errors"
	"fmt"
	"time"

	. "github.com/macarrie/flemzerd/objects"
)

type TVProvider struct{}
type MovieProvider struct{}
type DownloadDelayTVProvider struct{}
type DownloadDelayMovieProvider struct{}
type ErrorProvider struct{}

func (p TVProvider) Status() (Module, error) {
	return Module{
		Name: "TVProvider",
		Type: "provider",
		Status: ModuleStatus{
			Alive:   true,
			Message: "",
		},
	}, nil
}
func (p MovieProvider) Status() (Module, error) {
	return Module{
		Name: "MovieProvider",
		Type: "provider",
		Status: ModuleStatus{
			Alive:   true,
			Message: "",
		},
	}, nil
}
func (p DownloadDelayTVProvider) Status() (Module, error) {
	return Module{
		Name: "DownloadDelayTVProvider",
		Type: "provider",
		Status: ModuleStatus{
			Alive:   true,
			Message: "",
		},
	}, nil
}
func (p DownloadDelayMovieProvider) Status() (Module, error) {
	return Module{
		Name: "DownloadDelayMovieProvider",
		Type: "provider",
		Status: ModuleStatus{
			Alive:   true,
			Message: "",
		},
	}, nil
}
func (p ErrorProvider) Status() (Module, error) {
	err := fmt.Errorf("provider error")
	return Module{
		Name: "ErrorProvider",
		Type: "provider",
		Status: ModuleStatus{
			Alive:   false,
			Message: err.Error(),
		},
	}, err
}

func (p TVProvider) GetName() string {
	return "TVProvider"
}
func (p MovieProvider) GetName() string {
	return "MovieProvider"
}
func (p DownloadDelayTVProvider) GetName() string {
	return "DownloadDelayTVProvider"
}
func (p DownloadDelayMovieProvider) GetName() string {
	return "DownloadDelayMovieProvider"
}
func (p ErrorProvider) GetName() string {
	return "ErrorProvider"
}

func (p TVProvider) GetOrder() int {
	return 1
}
func (p MovieProvider) GetOrder() int {
	return 1
}
func (p DownloadDelayTVProvider) GetOrder() int {
	return 1
}
func (p DownloadDelayMovieProvider) GetOrder() int {
	return 1
}
func (p ErrorProvider) GetOrder() int {
	return 1
}

func (p TVProvider) GetShow(tvShow MediaIds) (TvShow, error) {
	return TvShow{
		Title:         "Test show",
		OriginalTitle: "Test show",
	}, nil
}
func (p TVProvider) GetEpisodes(tvShow TvShow) ([]Episode, error) {
	return []Episode{}, nil
}
func (p TVProvider) GetNextEpisodes(tvShow TvShow) ([]Episode, error) {
	return []Episode{}, nil
}
func (p TVProvider) GetRecentlyAiredEpisodes(tvShow TvShow) ([]Episode, error) {
	return []Episode{
		{
			Number: 1,
			Season: 1,
			Title:  "Test episode tvprovider",
		},
	}, nil
}
func (p TVProvider) GetSeasonEpisodeList(show TvShow, seasonNumber int) ([]Episode, error) {
	if seasonNumber == 1000 {
		return []Episode{}, errors.New("unknown season")
	}
	return []Episode{
		{
			Number: 1,
			Season: 1,
			Title:  "Test episode tvprovider",
		},
	}, nil
}
func (p TVProvider) GetEpisode(tvShow MediaIds, seasonNb int, episodeNb int) (Episode, error) {
	return Episode{
		Number: 1,
		Season: 1,
		Title:  "Test episode tvprovider",
	}, nil
}

func (p MovieProvider) GetMovie(movie MediaIds) (Movie, error) {
	return Movie{
		Title:         "Test Movie",
		OriginalTitle: "Test Movie",
	}, nil
}
func (p DownloadDelayMovieProvider) GetMovie(movie MediaIds) (Movie, error) {
	return Movie{
		Title:         "Test Movie",
		OriginalTitle: "Test Movie",
		Date:          time.Now(),
	}, nil
}

func (p DownloadDelayTVProvider) GetShow(tvShow MediaIds) (TvShow, error) {
	return TvShow{
		Title:         "Test show",
		OriginalTitle: "Test show",
	}, nil
}
func (p DownloadDelayTVProvider) GetEpisodes(tvShow TvShow) ([]Episode, error) {
	return []Episode{}, nil
}
func (p DownloadDelayTVProvider) GetNextEpisodes(tvShow TvShow) ([]Episode, error) {
	return []Episode{}, nil
}
func (p DownloadDelayTVProvider) GetRecentlyAiredEpisodes(tvShow TvShow) ([]Episode, error) {
	return []Episode{
		{
			Number: 1,
			Season: 1,
			Title:  "Test episode downloaddelaytvprovider",
			Date:   time.Now(),
		},
	}, nil
}
func (p DownloadDelayTVProvider) GetSeasonEpisodeList(show TvShow, seasonNumber int) ([]Episode, error) {
	if seasonNumber == 1000 {
		return []Episode{}, errors.New("unknown season")
	}
	return []Episode{
		{
			Number: 1,
			Season: 1,
			Title:  "Test episode downloaddelaytvprovider",
			Date:   time.Now(),
		},
	}, nil
}
func (p DownloadDelayTVProvider) GetEpisode(tvShow MediaIds, seasonNb int, episodeNb int) (Episode, error) {
	return Episode{
		Number: 1,
		Season: 1,
		Title:  "Test episode downloaddelaytvprovider",
		Date:   time.Now(),
	}, nil
}

func (p ErrorProvider) GetShow(tvShow MediaIds) (TvShow, error) {
	return TvShow{}, fmt.Errorf("provider error")
}
func (p ErrorProvider) GetEpisodes(tvShow TvShow) ([]Episode, error) {
	return []Episode{}, fmt.Errorf("provider error")
}
func (p ErrorProvider) GetNextEpisodes(tvShow TvShow) ([]Episode, error) {
	return []Episode{}, fmt.Errorf("provider error")
}
func (p ErrorProvider) GetRecentlyAiredEpisodes(tvShow TvShow) ([]Episode, error) {
	return []Episode{}, fmt.Errorf("provider error")
}
func (p ErrorProvider) GetSeasonEpisodeList(show TvShow, seasonNumber int) ([]Episode, error) {
	return []Episode{}, fmt.Errorf("provider error")
}
func (p ErrorProvider) GetMovie(movie MediaIds) (Movie, error) {
	return Movie{}, fmt.Errorf("provider error")
}
func (p ErrorProvider) GetEpisode(tvShow MediaIds, seasonNb int, episodeNb int) (Episode, error) {
	return Episode{}, fmt.Errorf("provider error")
}
