package provider

import . "github.com/macarrie/flemzerd/objects"

type Provider interface {
	IsAlive() error
	GetShow(tvShowName string) (TvShow, error)
	GetEpisodes(tvShow TvShow) ([]Episode, error)
	GetNextEpisodes(tvShow TvShow) ([]Episode, error)
	GetRecentlyAiredEpisodes(tvShow TvShow) ([]Episode, error)
}
