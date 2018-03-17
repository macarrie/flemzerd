package provider

import . "github.com/macarrie/flemzerd/objects"

type Provider interface {
	Status() (Module, error)
}

type TVProvider interface {
	Status() (Module, error)
	GetShow(tvShowName string) (TvShow, error)
	GetRecentlyAiredEpisodes(tvShow TvShow) ([]Episode, error)
}

type MovieProvider interface {
	Status() (Module, error)
	GetMovie(movieName string) (Movie, error)
}
