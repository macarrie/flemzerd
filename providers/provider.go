package provider

import . "github.com/macarrie/flemzerd/objects"

type Provider interface {
	Status() (Module, error)
	GetName() string
}

type TVProvider interface {
	Status() (Module, error)
	GetName() string
	GetShow(tvShow MediaIds) (TvShow, error)
	GetRecentlyAiredEpisodes(tvShow TvShow) ([]Episode, error)
	GetSeasonEpisodeList(show TvShow, seasonNumber int) ([]Episode, error)
}

type MovieProvider interface {
	Status() (Module, error)
	GetName() string
	GetMovie(movie MediaIds) (Movie, error)
}
