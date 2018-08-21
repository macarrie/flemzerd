package mock

import (
	"errors"
	"fmt"

	. "github.com/macarrie/flemzerd/objects"
)

type TVProvider struct{}
type MovieProvider struct{}
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
func (p ErrorProvider) Status() (Module, error) {
	var err error = fmt.Errorf("Provider error")
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
func (p ErrorProvider) GetName() string {
	return "ErrorProvider"
}

func (p TVProvider) GetShow(tvShow MediaIds) (TvShow, error) {
	return TvShow{
		Name:         "Test show",
		OriginalName: "Test show",
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
		Episode{
			Number: 1,
			Season: 1,
			Name:   "Test episode",
		},
	}, nil
}
func (p TVProvider) GetSeasonEpisodeList(show TvShow, seasonNumber int) ([]Episode, error) {
	if seasonNumber == 1000 {
		return []Episode{}, errors.New("Unknown season")
	}
	return []Episode{
		Episode{
			Number: 1,
			Season: 1,
			Name:   "Test episode",
		},
	}, nil
}

func (p MovieProvider) GetMovie(movie MediaIds) (Movie, error) {
	return Movie{
		Title:         "Test Movie",
		OriginalTitle: "Test Movie",
	}, nil
}

func (p ErrorProvider) GetShow(tvShow MediaIds) (TvShow, error) {
	return TvShow{}, fmt.Errorf("Provider error")
}
func (p ErrorProvider) GetEpisodes(tvShow TvShow) ([]Episode, error) {
	return []Episode{}, fmt.Errorf("Provider error")
}
func (p ErrorProvider) GetNextEpisodes(tvShow TvShow) ([]Episode, error) {
	return []Episode{}, fmt.Errorf("Provider error")
}
func (p ErrorProvider) GetRecentlyAiredEpisodes(tvShow TvShow) ([]Episode, error) {
	return []Episode{}, fmt.Errorf("Provider error")
}
func (p ErrorProvider) GetSeasonEpisodeList(show TvShow, seasonNumber int) ([]Episode, error) {
	return []Episode{}, fmt.Errorf("Provider error")
}
func (p ErrorProvider) GetMovie(movie MediaIds) (Movie, error) {
	return Movie{}, fmt.Errorf("Provider error")
}
