package provider

import (
	"testing"
)

type MockProvider struct{}

func (p MockProvider) Init() error {
	return nil
}

func (p MockProvider) FindShow(query string) (Show, error) {
	return Show{}, nil
}

func (p MockProvider) GetShow(id int) (Show, error) {
	return Show{}, nil
}

func (p MockProvider) GetEpisodes(show Show) ([]Episode, error) {
	return []Episode{}, nil
}

func (p MockProvider) FindNextAiredEpisodes(episodeList []Episode) ([]Episode, error) {
	return []Episode{}, nil
}

func (p MockProvider) FindNextEpisodesForShow(show Show) ([]Episode, error) {
	return []Episode{}, nil
}

func (p MockProvider) FindRecentlyAiredEpisodes(episodeList []Episode) ([]Episode, error) {
	return []Episode{}, nil
}

func (p MockProvider) FindRecentlyAiredEpisodesForShow(show Show) ([]Episode, error) {
	return []Episode{}, nil
}

func TestAddProvider(t *testing.T)                      {}
func TestFindShow(t *testing.T)                         {}
func TestFindRecentlyAiredEpisodesForShow(t *testing.T) {}
