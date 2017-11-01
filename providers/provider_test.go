package provider

import . "github.com/macarrie/flemzerd/objects"

type MockProvider struct{}

func (p MockProvider) GetShow(tvShowName string) (TvShow, error) {
	return TvShow{}, nil
}
func (p MockProvider) GetEpisodes(tvShow TvShow) ([]Episode, error) {
	return []Episode{}, nil
}
func (p MockProvider) GetNextEpisodes(tvShow TvShow) ([]Episode, error) {
	return []Episode{}, nil
}
func (p MockProvider) GetRecentlyAiredEpisodes(tvShow TvShow) ([]Episode, error) {
	return []Episode{}, nil
}
