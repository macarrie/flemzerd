package provider

import (
	"testing"
)

type MockProvider struct{}

func (p MockProvider) GetShow(tvShowName string) (Show, error) {
	return Show{}, nil
}
func (p MockProvider) GetEpisodes(tvShow Show) ([]Episode, error) {
	return []Episode{}, nil
}
func (p MockProvider) GetNextEpisodes(tvShow Show) ([]Episode, error) {
	return []Episode{}, nil
}
func (p MockProvider) GetRecentlyAiredEpisodes(tvShow Show) ([]Episode, error) {
	return []Episode{}, nil
}

func TestAddProvider(t *testing.T)                      {}
func TestFindShow(t *testing.T)                         {}
func TestFindRecentlyAiredEpisodesForShow(t *testing.T) {}
