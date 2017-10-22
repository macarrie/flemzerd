package provider

import (
	log "flemzerd/logging"
)

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
	Init() error
	FindShow(query string) (Show, error)
	GetShow(id int) (Show, error)
	GetEpisodes(show Show) ([]Episode, error)
	FindNextAiredEpisodes(episodeList []Episode) ([]Episode, error)
	FindNextEpisodesForShow(show Show) ([]Episode, error)
	FindRecentlyAiredEpisodes(episodeList []Episode) ([]Episode, error)
	FindRecentlyAiredEpisodesForShow(show Show) ([]Episode, error)
}

var providers []Provider

func AddProvider(provider Provider) {
	providers = append(providers, provider)
	log.WithFields(log.Fields{
		"provider": provider,
	}).Debug("Provider loaded")
}

func FindShow(query string) (Show, error) {
	return providers[0].FindShow(query)
}

func FindRecentlyAiredEpisodesForShow(show Show) ([]Episode, error) {
	return providers[0].FindRecentlyAiredEpisodesForShow(show)
}
