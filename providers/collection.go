package provider

import (
	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
)

var providersCollection []Provider

func Status() ([]Module, error) {
	mod, err := providersCollection[0].Status()
	return []Module{mod}, err
}

func AddProvider(provider Provider) {
	providersCollection = append(providersCollection, provider)
	log.Debug("The TVDB provider loaded")
}

func FindShow(query string) (TvShow, error) {
	return providersCollection[0].GetShow(query)
}

func FindRecentlyAiredEpisodesForShow(show TvShow) ([]Episode, error) {
	return providersCollection[0].GetRecentlyAiredEpisodes(show)
}
