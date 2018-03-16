package provider

import . "github.com/macarrie/flemzerd/objects"

type Provider interface {
	Status() (Module, error)
	GetShow(tvShowName string) (TvShow, error)
	GetRecentlyAiredEpisodes(tvShow TvShow) ([]Episode, error)
}
