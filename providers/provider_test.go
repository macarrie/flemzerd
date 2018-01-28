package provider

import . "github.com/macarrie/flemzerd/objects"

type MockProvider struct{}

func (p MockProvider) IsAlive() error {
	return nil
}

func (p MockProvider) GetShow(tvShowName string) (TvShow, error) {
	return TvShow{
		Id:   1000,
		Name: "Test show",
	}, nil
}
func (p MockProvider) GetEpisodes(tvShow TvShow) ([]Episode, error) {
	return []Episode{}, nil
}
func (p MockProvider) GetNextEpisodes(tvShow TvShow) ([]Episode, error) {
	return []Episode{}, nil
}
func (p MockProvider) GetRecentlyAiredEpisodes(tvShow TvShow) ([]Episode, error) {
	return []Episode{
		Episode{
			Number: 1,
			Season: 1,
			Name:   "Test episode",
			Id:     1000,
		},
	}, nil
}
