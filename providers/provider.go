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
	Init() bool
	FindShow(query string) (Show, error)
	GetShow(id int) (Show, error)
	GetEpisodes(show Show) ([]Episode, error)
	FindNextAiredEpisode(episodeList []Episode) (Episode, error)
	FindNextEpisodeForShow(show Show) (Episode, error)
	FindRecentlyAiredEpisode(episodeList []Episode) (Episode, error)
	FindRecentlyAiredEpisodeForShow(show Show) (Episode, error)
}

var providers []Provider

func AddProvider(provider Provider) {
	providers = append(providers, provider)
	log.WithFields(log.Fields{
		"provider": provider,
	}).Info("Provider loaded")
}

func FindShow(query string) (Show, error) {
	return providers[0].FindShow(query)
}

func FindRecentlyAiredEpisodeForShow(show Show) (Episode, error) {
	return providers[0].FindRecentlyAiredEpisodeForShow(show)
}
