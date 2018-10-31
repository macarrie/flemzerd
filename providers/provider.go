package provider

import . "github.com/macarrie/flemzerd/objects"

type Provider interface {
	Status() (Module, error)
	GetName() string
}

type TVProvider interface {
	Status() (Module, error)
	GetName() string
	GetOrder() int
	GetShow(tvShow MediaIds) (TvShow, error)
	GetRecentlyAiredEpisodes(tvShow TvShow) ([]Episode, error)
	GetSeasonEpisodeList(show TvShow, seasonNumber int) ([]Episode, error)
	GetEpisode(tvShow MediaIds, seasonNb int, episodeNb int) (Episode, error)
}

type MovieProvider interface {
	Status() (Module, error)
	GetName() string
	GetOrder() int
	GetMovie(movie MediaIds) (Movie, error)
}
