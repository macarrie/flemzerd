package tvdb

import (
	"net/url"
	"strconv"
	"time"

	"github.com/macarrie/flemzerd/configuration"
	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"

	"github.com/pioz/tvdb"
	"github.com/pkg/errors"
)

type TVDBProvider struct {
	Client          tvdb.Client
	LastTokenUpdate time.Time
	Order           int
}

var module Module

// Create new instance of the TVDB info provider
func New(apiKey string, order string) (tvdbProvider *TVDBProvider, err error) {
	t := &TVDBProvider{}
	parsedOrder, err := strconv.Atoi(order)
	if err != nil {
		log.WithFields(log.Fields{
			"provider": t.GetName(),
			"order":    order,
		}).Error("Incorrect order defined in configuration.")
		parsedOrder = 1
	}

	module = Module{
		Name: t.GetName(),
		Type: "provider",
		Status: ModuleStatus{
			Alive:   true,
			Message: "",
		},
	}

	client := tvdb.Client{Apikey: apiKey}
	log.WithFields(log.Fields{
		"provider": module.Name,
	}).Debug("Checking connection to TheTVDB")
	if err := client.Login(); err != nil {
		//if tvdb.HaveCodeError(401, err) {
		//log.Error("Cannot connect to thetvdb (API key not valid). Please check your API key and try again")
		//} else {
		//log.WithFields(log.Fields{
		//"details":  err,
		//"provider": module.Name,
		//"error":    err,
		//}).Error("Cannot connect to thetvdb")
		//}
		return &TVDBProvider{}, errors.Wrap(err, "cannot connect to TMDB")
	} else {
		log.Debug("Connection to TheTVDB successful")
		return &TVDBProvider{Client: client, LastTokenUpdate: time.Now(), Order: parsedOrder}, nil
	}
}

// Check if Provider is alive
func (tvdbProvider *TVDBProvider) Status() (Module, error) {
	log.Debug("Checking TVDB provider status")

	if configuration.TELEGRAM_BOT_TOKEN == "" {
		module.Status.Alive = false
		module.Status.Message = "TVDB API key not found"
		return module, errors.New(module.Status.Message)
	}

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

	return module, errors.Wrap(err, "cannot login to TVDB")
}

func (tvdbProvider *TVDBProvider) GetName() string {
	return "TVDB"
}

func (tvdbProvider *TVDBProvider) GetOrder() int {
	return tvdbProvider.Order
}

// Get show from name
func (tvdbProvider *TVDBProvider) GetShow(tvShow MediaIds) (TvShow, error) {
	log.WithFields(log.Fields{
		"title":    tvShow.Title,
		"provider": module.Name,
	}).Debug("Searching show")

	show, err := tvdbProvider.Client.BestSearch(tvShow.Title)
	if err != nil {
		return TvShow{}, handleTvShowNotFoundError(tvShow.Title, err)
	} else {
		log.WithFields(log.Fields{
			"title":    tvShow.Title,
			"TVDB-ID":  tvShow.ID,
			"provider": module.Name,
		}).Debug("TV show found")

		return convertShow(show), nil
	}
}

// Get specific episode from tvshow
func (tvdbProvider *TVDBProvider) GetEpisode(tvShow MediaIds, seasonNb int, episodeNb int) (Episode, error) {
	log.WithFields(log.Fields{
		"name":     tvShow.Title,
		"provider": module.Name,
		"season":   seasonNb,
		"episode":  episodeNb,
	}).Debug("Getting episode")

	show, err := tvdbProvider.Client.BestSearch(tvShow.Title)
	if err != nil {
		return Episode{}, errors.Wrap(err, "cannot get episode from TVDB")
	}

	episode := show.GetEpisode(seasonNb, episodeNb)

	if episode == nil {
		return Episode{}, errors.New("Cannot find episode for show from TVDB")
	} else {
		return convertEpisode(*episode), nil
	}
}

// Get list of episodes of a given show
func (tvdbProvider *TVDBProvider) GetEpisodes(tvShow TvShow) ([]Episode, error) {
	log.WithFields(log.Fields{
		"id":       tvShow.Model.ID,
		"title":    tvShow.GetTitle(),
		"provider": module.Name,
	}).Debug("Retrieving episodes list for a tv show")

	tvShowSearchResult, err := tvdbProvider.Client.BestSearch(tvShow.GetTitle())
	if err != nil {
		return []Episode{}, handleTvShowNotFoundError(tvShow.GetTitle(), err)
	} else {
		if err := tvdbProvider.Client.GetSeriesEpisodes(&tvShowSearchResult, url.Values{}); err != nil {
			log.WithFields(log.Fields{
				"tvshow":   tvShow.GetTitle(),
				"id":       tvShow.Model.ID,
				"provider": module.Name,
				"error":    err,
			}).Warn("Cannot retrieve episodes of tv show")
			return []Episode{}, errors.Wrap(err, "cannot get show episodes from TVDB")
		} else {
			log.WithFields(log.Fields{
				"title":     tvShow.GetTitle(),
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
			"tvshow":   tvShow.GetTitle(),
			"id":       tvShow.Model.ID,
			"provider": module.Name,
			"error":    err,
		}).Warn("Cannot get next aired episodes of the tv show")
		return []Episode{}, errors.Wrap(err, "cannot get episodes from TVDB")
	} else {
		log.WithFields(log.Fields{
			"tvshow":   tvShow.GetTitle(),
			"id":       tvShow.Model.ID,
			"provider": module.Name,
		}).Debug("Getting list of episodes that haven't be aired yet")
		now := time.Now()
		futureEpisodes := filterEpisodesAiredBetweenDates(episodes, &now, nil)

		log.WithFields(log.Fields{
			"tvshow":         tvShow.GetTitle(),
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
			"tvshow":   tvShow.GetTitle(),
			"id":       tvShow.Model.ID,
			"provider": module.Name,
			"error":    err,
		}).Warn("Cannot get recently aired episodes of the tv show")
		return []Episode{}, errors.Wrap(err, "cannot get recent episodes from TVDB")
	}
	log.WithFields(log.Fields{
		"tvshow":   tvShow.GetTitle(),
		"id":       tvShow.Model.ID,
		"provider": module.Name,
	}).Debug("Getting list of episodes that haven recently been aired")

	// Current date minus RECENTLY_AIRED_EPISODES_INTERVAL days
	oldestDate := time.Now().AddDate(0, 0, -RECENTLY_AIRED_EPISODES_INTERVAL)
	now := time.Now()
	recentlyAiredEpisodes := filterEpisodesAiredBetweenDates(episodes, &oldestDate, &now)

	log.WithFields(log.Fields{
		"tvshow":         tvShow.GetTitle(),
		"id":             tvShow.Model.ID,
		"nb_of_episodes": len(recentlyAiredEpisodes),
		"provider":       module.Name,
	}).Debug("Successfully filtered list of episodes that have been recently aired")

	return recentlyAiredEpisodes, nil
}

func (tvdbProvider *TVDBProvider) GetSeasonEpisodeList(show TvShow, seasonNumber int) ([]Episode, error) {
	// Method used for UI
	// Not implemented for TVDB provider
	return []Episode{}, nil
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

func handleTvShowNotFoundError(tvShowTitle string, err error) error {
	if tvdb.HaveCodeError(404, err) {
		// The request response is a 404: this means no results have been found
		log.WithFields(log.Fields{
			"tvshow":   tvShowTitle,
			"provider": module.Name,
		}).Warning("TV show not found")
		return errors.New("TV show '" + tvShowTitle + "' not found")
	} else if tvdb.HaveCodeError(401, err) {
		// The request response is a 401: Authentication failure
		log.WithFields(log.Fields{
			"provider": module.Name,
		}).Warning("Could not connect to provider, authentication failed for unkonwn reasons")
		return errors.Wrap(err, "cannot connect to TVDB (auth error)")
	} else {
		log.WithFields(log.Fields{
			"provider": module.Name,
		}).Error("Unknown error encountered while getting show information from The TVDB")
		return errors.Wrap(err, "cannot get show from TVDB")
	}
}

// Convert a github.com/pioz/tvdb series object to flemzerd tvShow object
func convertShow(tvShow tvdb.Series) TvShow {
	firstAired, err := time.Parse("2006-01-02", tvShow.FirstAired)
	if err != nil {
		firstAired = time.Time{}
	}

	return TvShow{
		Banner:          tvShow.Banner,
		FirstAired:      firstAired,
		Overview:        tvShow.Overview,
		Title:           tvShow.SeriesName,
		UseDefaultTitle: true,
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
		AbsoluteNumber: episode.AbsoluteNumber,
		Number:         episode.AiredEpisodeNumber,
		Season:         episode.AiredSeason,
		Title:          episode.EpisodeName,
		Date:           firstAired,
		Overview:       episode.Overview,
	}
}
