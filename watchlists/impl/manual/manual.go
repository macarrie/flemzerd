package manual

import (
	"github.com/macarrie/flemzerd/configuration"
	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
)

type ManualWatchlist struct{}

var module Module

func New() (t *ManualWatchlist, err error) {
	t = &ManualWatchlist{}
	module = Module{
		Name: "Manual",
		Type: "watchlist",
		Status: ModuleStatus{
			Alive:   true,
			Message: "",
		},
	}

	return t, nil
}

func (t *ManualWatchlist) Status() (Module, error) {
	log.Warning("Checking manual watchlist status")

	module.Status.Alive = true
	module.Status.Message = ""

	return module, nil
}

func (t *ManualWatchlist) GetTvShows() ([]string, error) {
	log.WithFields(log.Fields{
		"watchlist": "manual",
	}).Debug("Getting TV shows from watchlist")

	var shows []string
	for _, show := range configuration.Config.Watchlists["manual"].([]interface{}) {
		shows = append(shows, show.(string))
	}

	return shows, nil
}
func (t *ManualWatchlist) GetMovies() ([]string, error) {
	log.WithFields(log.Fields{
		"watchlist": "manual",
	}).Debug("Getting movies from watchlist")

	return []string{}, nil
}
