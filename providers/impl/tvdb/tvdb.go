package tvdb

import (
	"errors"
	"net/url"
	"time"

	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
	"github.com/pioz/tvdb"
)

type TVDBProvider struct {
	Client          tvdb.Client
	LastTokenUpdate time.Time
}

var module Module

// Create new instance of the TVDB info provider
func New(apiKey string) (tvdbProvider *TVDBProvider, err error) {
	module = Module{
		Name: "TVDB",
		Type: "provider",
		Status: ModuleStatus{
			Alive:   true,
			Message: "",
		},
	}

	client := tvdb.Client{Apikey: apiKey}
	err = client.Login()
	log.WithFields(log.Fields{
		"provider": module.Name,
	}).Debug("Checking connection to TheTVDB")

	if err != nil {
		if tvdb.HaveCodeError(401, err) {
			log.Error("Cannot connect to thetvdb (API key not valid). Please check your API key and try again")
		} else {
			log.WithFields(log.Fields{
				"details":  err,
				"provider": module.Name,
				"error":    err,
			}).Error("Cannot connect to thetvdb")
		}
		return &TVDBProvider{}, err
	} else {
		log.Debug("Connection to TheTVDB successful")
		return &TVDBProvider{Client: client, LastTokenUpdate: time.Now()}, nil
	}
}

// Check if Provider is alive
func (tvdbProvider *TVDBProvider) Status() (Module, error) {
	log.Debug("Checking TVDB provider status")

	// Token expires every 24h. Refresh it if needed
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)

	if tvdbProvider.LastTokenUpdate.Before(yesterday) {
		log.Debug("Refreshing TVDB token")
		tvdbProvider.Client.RefreshToken()
		tvdbProvider.LastTokenUpdate = now
	}

	err := tvdbProvider.Client.Login()
	if err != nil {
		module.Status.Alive = false
		module.Status.Message = err.Error()
	} else {
		module.Status.Alive = true
		module.Status.Message = ""
	}

	return module, err
}

// Get show from name
func (tvdbProvider *TVDBProvider) GetShow(tvShowName string) (TvShow, error) {
	log.WithFields(log.Fields{
		"name":     tvShowName,
		"provider": module.Name,
	}).Debug("Searching show")

	tvShow, err := tvdbProvider.Client.BestSearch(tvShowName)
	if err != nil {
		return TvShow{}, handleTvShowNotFoundError(tvShowName, err)
	} else {
		log.WithFields(log.Fields{
			"name":     tvShow.SeriesName,
			"TVDB-ID":  tvShow.ID,
			"provider": module.Name,
		}).Debug("TV show found")
		return convertShow(tvShow), nil

	}
}

// Get list of episodes of a given show
func (tvdbProvider *TVDBProvider) GetEpisodes(tvShow TvShow) ([]Episode, error) {
	log.WithFields(log.Fields{
		"id":       tvShow.Model.ID,
		"name":     tvShow.Name,
		"provider": module.Name,
	}).Debug("Retrieving episodes list for a tv show")

	tvShowSearchResult, err := tvdbProvider.Client.BestSearch(tvShow.Name)
	if err != nil {
		return []Episode{}, handleTvShowNotFoundError(tvShow.Name, err)
	} else {
		err := tvdbProvider.Client.GetSeriesEpisodes(&tvShowSearchResult, url.Values{})
		if err != nil {
			log.WithFields(log.Fields{
				"tvshow_name": tvShow.Name,
				"id":          tvShow.Model.ID,
				"provider":    module.Name,
				"error":       err,
			}).Warn("Cannot retrieve episodes of tv show")
			return []Episode{}, err
		} else {
			log.WithFields(log.Fields{
				"name":      tvShow.Name,
				"tvshow_id": tvShow.Model.ID,
				"provider":  module.Name,
			}).Debug("Episodes found")

			var retVal []Episode
			for _, episode := range tvShowSearchResult.Episodes {
				retVal = append(retVal, convertEpisode(episode))

			}
			return retVal, nil
		}
	}
}

// Get all episodes of a tv show that haven't been aired yet for a given show
func (tvdbProvider *TVDBProvider) GetNextEpisodes(tvShow TvShow) ([]Episode, error) {
	episodes, err := tvdbProvider.GetEpisodes(tvShow)
	if err != nil {
		log.WithFields(log.Fields{
			"tvshow_name": tvShow.Name,
			"id":          tvShow.Model.ID,
			"provider":    module.Name,
			"error":       err,
		}).Warn("Cannot get next aired episodes of the tv show")
		return []Episode{}, err
	} else {
		log.WithFields(log.Fields{
			"tvshow_name": tvShow.Name,
			"id":          tvShow.Model.ID,
			"provider":    module.Name,
		}).Debug("Getting list of episodes that haven't be aired yet")
		now := time.Now()
		futureEpisodes := filterEpisodesAiredBetweenDates(episodes, &now, nil)

		log.WithFields(log.Fields{
			"tvshow_name":    tvShow.Name,
			"id":             tvShow.Model.ID,
			"nb_of_episodes": len(futureEpisodes),
			"provider":       module.Name,
		}).Debug("Successfully filtered list of episodes that haven't be aired yet")
		return futureEpisodes, nil
	}
}

// Get list of episodes of a show aired less than RECENTLY_AIRED_EPISODES_INTERVAL days ago
func (tvdbProvider *TVDBProvider) GetRecentlyAiredEpisodes(tvShow TvShow) ([]Episode, error) {
	episodes, err := tvdbProvider.GetEpisodes(tvShow)
	if err != nil {
		log.WithFields(log.Fields{
			"tvshow_name": tvShow.Name,
			"id":          tvShow.Model.ID,
			"provider":    module.Name,
			"error":       err,
		}).Warn("Cannot get recently aired episodes of the tv show")
		return []Episode{}, err
	}
	log.WithFields(log.Fields{
		"tvshow_name": tvShow.Name,
		"id":          tvShow.Model.ID,
		"provider":    module.Name,
	}).Debug("Getting list of episodes that haven recently been aired")

	// Current date minus RECENTLY_AIRED_EPISODES_INTERVAL days
	oldestDate := time.Now().AddDate(0, 0, -RECENTLY_AIRED_EPISODES_INTERVAL)
	now := time.Now()
	recentlyAiredEpisodes := filterEpisodesAiredBetweenDates(episodes, &oldestDate, &now)

	log.WithFields(log.Fields{
		"tvshow_name":    tvShow.Name,
		"id":             tvShow.Model.ID,
		"nb_of_episodes": len(recentlyAiredEpisodes),
		"provider":       module.Name,
	}).Debug("Successfully filtered list of episodes that have been recently aired")
	return recentlyAiredEpisodes, nil
}

func filterEpisodesAiredBetweenDates(episodes []Episode, beginning *time.Time, end *time.Time) []Episode {
	// Set beginning date to zero time if beginning date is nil
	if beginning == nil {
		*beginning = time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC)
	}
	// Set end date to year 100,000 if end date is nil
	if end == nil {
		*end = time.Date(100000, time.January, 1, 0, 0, 0, 0, time.UTC)
	}

	log.WithFields(log.Fields{
		"start_date": beginning,
		"end_date":   end,
	}).Debug("Filtering episodes list by airing date")
	var retVal []Episode
	for _, episode := range episodes {
		if episode.Date.After(*beginning) && episode.Date.Before(*end) {
			retVal = append(retVal, episode)
		}
	}

	log.WithFields(log.Fields{
		"start_date":     beginning,
		"end_date":       end,
		"nb_of_episodes": len(retVal),
	}).Debug("Successfully filtered episodes list by airing date")

	return retVal
}

func handleTvShowNotFoundError(tvShowName string, err error) error {
	if tvdb.HaveCodeError(404, err) {
		// The request response is a 404: this means no results have been found
		log.WithFields(log.Fields{
			"tvshow_name": tvShowName,
			"provider":    module.Name,
		}).Warning("TV show not found")
		return errors.New("TV show '" + tvShowName + "' not found")
	} else if tvdb.HaveCodeError(401, err) {
		// The request response is a 401: Authentication failure
		log.WithFields(log.Fields{
			"provider": module.Name,
		}).Warning("Could not connect to provider, authentication failed for unkonwn reasons")
		return err
	} else {
		log.WithFields(log.Fields{
			"provider": module.Name,
		}).Error("Unknown error encountered while getting show information from The TVDB")
		return err
	}
}

// Convert a github.com/pioz/tvdb series object to flemzerd tvShow object
func convertShow(tvShow tvdb.Series) TvShow {
	firstAired, err := time.Parse("2006-01-02", tvShow.FirstAired)
	if err != nil {
		firstAired = time.Time{}
	}

	return TvShow{
		Banner:     tvShow.Banner,
		FirstAired: firstAired,
		Overview:   tvShow.Overview,
		Name:       tvShow.SeriesName,
		// TODO: Update with correct status
		Status: TVSHOW_UNKNOWN,
	}
}

// Convert a github.com/pioz/tvdb episode object to flemzerd episode object
func convertEpisode(episode tvdb.Episode) Episode {
	firstAired, err := time.Parse("2006-01-02", episode.FirstAired)
	if err != nil {
		firstAired = time.Time{}
	}

	return Episode{
		//AbsoluteNumber: episode.AbsoluteNumber,
		Number:   episode.AiredEpisodeNumber,
		Season:   episode.AiredSeason,
		Name:     episode.EpisodeName,
		Date:     firstAired,
		Overview: episode.Overview,
	}
}
