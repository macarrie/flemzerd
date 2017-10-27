package provider

import (
	log "github.com/macarrie/flemzerd/logging"
)

const RECENTLY_AIRED_EPISODES_INTERVAL = 14

type Show struct {
	Aliases    []string
	Banner     string
	FirstAired string
	Id         int
	Network    string
	Overview   string
	Name       string
	Status     string
}

type Episode struct {
	AbsoluteNumber int
	Number         int
	Season         int
	Name           string
	Date           string
	Id             int
	Overview       string
}

type Provider interface {
	GetShow(tvShowName string) (Show, error)
	GetEpisodes(tvShow Show) ([]Episode, error)
	GetNextEpisodes(tvShow Show) ([]Episode, error)
	GetRecentlyAiredEpisodes(tvShow Show) ([]Episode, error)
}

var providers []Provider

func AddProvider(provider Provider) {
	providers = append(providers, provider)
	log.Debug("The TVDB provider loaded")
}

func FindShow(query string) (Show, error) {
	return providers[0].GetShow(query)
}

func FindRecentlyAiredEpisodesForShow(show Show) ([]Episode, error) {
	return providers[0].GetRecentlyAiredEpisodes(show)
}
