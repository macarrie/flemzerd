package provider

import (
	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"

	watchlist "github.com/macarrie/flemzerd/watchlists"
)

var providersCollection []Provider
var TVShows []TvShow

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

func removeDuplicates(array []TvShow) []TvShow {
	occurences := make(map[int]bool)
	var ret []TvShow

	for _, show := range array {
		if !occurences[show.Id] {
			occurences[show.Id] = true
			ret = append(ret, show)
		}
	}

	return ret
}

func GetTVShowsInfoFromConfig() {
	var showObjects []TvShow
	var showList []string

	showsFromWatchlists, _ := watchlist.GetTvShows()
	showList = append(showList, showsFromWatchlists...)

	for _, show := range showList {
		showName := show
		show, err := FindShow(show)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"show":  showName,
			}).Warning("Unable to get show informations")
		} else {
			showObjects = append(showObjects, show)
		}
	}
	if len(showObjects) == 0 {
		log.Error("Impossible to get show informations for shows defined in configuration. Shutting down")
	}

	TVShows = removeDuplicates(showObjects)
}
