package retention

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
)

type DownloadingEpisode struct {
	Episode         Episode
	Downloading     bool
	FailedTorrents  map[string]*Torrent
	CurrentDownload struct {
		Torrent      Torrent
		DownloaderId string
	}
}

type DownloadingMovie struct {
	Movie           Movie
	Downloading     bool
	FailedTorrents  map[string]*Torrent
	CurrentDownload struct {
		Torrent      Torrent
		DownloaderId string
	}
}

type RetentionData struct {
	NotifiedEpisodes    map[int]Episode
	DownloadedEpisodes  map[int]Episode
	DownloadingEpisodes map[int]*DownloadingEpisode
	FailedEpisodes      map[int]Episode

	NotifiedMovies    map[int]Movie
	DownloadedMovies  map[int]Movie
	DownloadingMovies map[int]*DownloadingMovie
	FailedMovies      map[int]Movie

	TraktToken string
	Watchlists map[string][]string
}

var retentionData RetentionData

func InitStruct() {
	retentionData.NotifiedEpisodes = make(map[int]Episode)
	retentionData.DownloadedEpisodes = make(map[int]Episode)
	retentionData.DownloadingEpisodes = make(map[int]*DownloadingEpisode)
	retentionData.FailedEpisodes = make(map[int]Episode)

	retentionData.NotifiedMovies = make(map[int]Movie)
	retentionData.DownloadedMovies = make(map[int]Movie)
	retentionData.DownloadingMovies = make(map[int]*DownloadingMovie)
	retentionData.FailedMovies = make(map[int]Movie)
}

func Load() error {
	InitStruct()

	file, err := os.Open(RETENTION_FILE_PATH)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	if !json.Valid(data) {
		return errors.New("Retention file contains invalid data, aborting retention load")
	}

	unmarshalError := json.Unmarshal(data, &retentionData)
	if unmarshalError != nil {
		return unmarshalError
	}

	log.Info("Retention data successfully loaded")
	return nil
}

func Save() error {
	file, err := os.OpenFile(RETENTION_FILE_PATH, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(retentionData)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	log.Debug("Retention data saved")

	return nil
}

func EpisodeHasBeenNotified(e Episode) bool {
	_, ok := retentionData.NotifiedEpisodes[e.Id]
	return ok
}

func MovieHasBeenNotified(m Movie) bool {
	_, ok := retentionData.NotifiedMovies[m.Id]
	return ok
}

func EpisodeHasBeenDownloaded(e Episode) bool {
	_, ok := retentionData.DownloadedEpisodes[e.Id]
	return ok
}

func MovieHasBeenDownloaded(m Movie) bool {
	_, ok := retentionData.DownloadedMovies[m.Id]
	return ok
}

func EpisodeIsInDownloadProcess(e Episode) bool {
	_, ok := retentionData.DownloadingEpisodes[e.Id]
	return ok
}

func MovieIsInDownloadProcess(m Movie) bool {
	_, ok := retentionData.DownloadingMovies[m.Id]
	return ok
}

func EpisodeIsDownloading(e Episode) bool {
	if !EpisodeIsInDownloadProcess(e) {
		return false
	}

	ep, _ := retentionData.DownloadingEpisodes[e.Id]
	return ep.Downloading
}

func MovieIsDownloading(m Movie) bool {
	if !MovieIsInDownloadProcess(m) {
		return false
	}

	movie, _ := retentionData.DownloadingMovies[m.Id]
	return movie.Downloading
}

func EpisodeIsInFailedTorrents(e Episode, t Torrent) bool {
	if !EpisodeIsDownloading(e) {
		return false
	}

	_, ok := retentionData.DownloadingEpisodes[e.Id].FailedTorrents[t.Id]
	return ok
}

func MovieIsInFailedTorrents(m Movie, t Torrent) bool {
	if !MovieIsDownloading(m) {
		return false
	}

	_, ok := retentionData.DownloadingMovies[m.Id].FailedTorrents[t.Id]
	return ok
}

func EpisodeGetFailedTorrentsCount(e Episode) int {
	if !EpisodeIsDownloading(e) {
		return 0
	}

	return len(retentionData.DownloadingEpisodes[e.Id].FailedTorrents)
}

func MovieGetFailedTorrentsCount(m Movie) int {
	if !MovieIsDownloading(m) {
		return 0
	}

	return len(retentionData.DownloadingMovies[m.Id].FailedTorrents)
}

func EpisodeHasDownloadFailed(e Episode) bool {
	_, present := retentionData.FailedEpisodes[e.Id]

	return present
}

func MovieHasDownloadFailed(m Movie) bool {
	_, present := retentionData.FailedEpisodes[m.Id]

	return present
}

func AddNotifiedEpisode(e Episode) {
	if !EpisodeHasBeenNotified(e) {
		retentionData.NotifiedEpisodes[e.Id] = e
	}
}

func AddNotifiedMovie(m Movie) {
	if !MovieHasBeenNotified(m) {
		fmt.Printf("RetentionData before: %+v\n", retentionData.NotifiedMovies)
		retentionData.NotifiedMovies[m.Id] = m
		fmt.Printf("RetentionData after: %+v\n", retentionData.NotifiedMovies)
		return
	}
}

func AddDownloadedEpisode(e Episode) {
	if !EpisodeHasBeenDownloaded(e) {
		retentionData.DownloadedEpisodes[e.Id] = e
	}

	delete(retentionData.DownloadingEpisodes, e.Id)
}

func AddDownloadedMovie(m Movie) {
	if !MovieHasBeenDownloaded(m) {
		retentionData.DownloadedMovies[m.Id] = m
	}

	delete(retentionData.DownloadingMovies, m.Id)
}

func AddDownloadingEpisode(e Episode) {
	if !EpisodeIsDownloading(e) {
		retentionData.DownloadingEpisodes[e.Id] = &DownloadingEpisode{
			Episode:        e,
			Downloading:    true,
			FailedTorrents: make(map[string]*Torrent),
		}
	}
}

func AddDownloadingMovie(m Movie) {
	if !MovieIsDownloading(m) {
		retentionData.DownloadingMovies[m.Id] = &DownloadingMovie{
			Movie:          m,
			Downloading:    true,
			FailedTorrents: make(map[string]*Torrent),
		}
	}
}

func AddFailedEpisode(e Episode) {
	if !EpisodeHasDownloadFailed(e) {
		retentionData.FailedEpisodes[e.Id] = e
	}
}

func AddFailedMovie(m Movie) {
	if !MovieHasDownloadFailed(m) {
		retentionData.FailedMovies[m.Id] = m
	}
}

func AddFailedEpisodeTorrent(e Episode, t Torrent) {
	if !EpisodeIsDownloading(e) {
		AddDownloadingEpisode(e)
	}
	retentionData.DownloadingEpisodes[e.Id].FailedTorrents[t.Id] = &t
}

func AddFailedMovieTorrent(m Movie, t Torrent) {
	if !MovieIsDownloading(m) {
		AddDownloadingMovie(m)
	}
	retentionData.DownloadingMovies[m.Id].FailedTorrents[t.Id] = &t
}

func SetCurrentEpisodeDownload(e Episode, t Torrent, downloaderId string) {
	if !EpisodeIsDownloading(e) {
		return
	}
	retentionData.DownloadingEpisodes[e.Id].CurrentDownload.Torrent = t
	retentionData.DownloadingEpisodes[e.Id].CurrentDownload.DownloaderId = downloaderId
}

func SetCurrentMovieDownload(m Movie, t Torrent, downloaderId string) {
	if !MovieIsDownloading(m) {
		return
	}
	retentionData.DownloadingMovies[m.Id].CurrentDownload.Torrent = t
	retentionData.DownloadingMovies[m.Id].CurrentDownload.DownloaderId = downloaderId
}

func GetDownloadingEpisodes() ([]Episode, error) {
	var episodesList []Episode
	for _, val := range retentionData.DownloadingEpisodes {
		episodesList = append(episodesList, val.Episode)
	}
	return episodesList, nil
}

func GetDownloadingMovies() ([]Movie, error) {
	var moviesList []Movie
	for _, val := range retentionData.DownloadingMovies {
		moviesList = append(moviesList, val.Movie)
	}
	return moviesList, nil
}

func GetCurrentDownloadEpisodeFromRetention(e Episode) (Torrent, string, error) {
	if !EpisodeIsDownloading(e) {
		return Torrent{}, "", errors.New("Episode is not getting downloaded, no torrents to fetch")
	}
	currentDownload := retentionData.DownloadingEpisodes[e.Id].CurrentDownload
	return currentDownload.Torrent, currentDownload.DownloaderId, nil
}

func GetCurrentDownloadMovieFromRetention(m Movie) (Torrent, string, error) {
	if !MovieIsDownloading(m) {
		return Torrent{}, "", errors.New("Movie is not getting downloaded, no torrents to fetch")
	}
	currentDownload := retentionData.DownloadingMovies[m.Id].CurrentDownload
	return currentDownload.Torrent, currentDownload.DownloaderId, nil
}

func ChangeEpisodeDownloadingState(e Episode, state bool) {
	if !EpisodeIsInDownloadProcess(e) {
		return
	}

	downloadingEpisode := retentionData.DownloadingEpisodes[e.Id]
	downloadingEpisode.Downloading = state
}

func ChangeMovieDownloadingState(m Movie, state bool) {
	if !MovieIsInDownloadProcess(m) {
		return
	}

	downloadingMovie := retentionData.DownloadingMovies[m.Id]
	downloadingMovie.Downloading = state
}

func CleanOldNotifiedEpisodes() {
	for _, ep := range retentionData.NotifiedEpisodes {
		// Remove notified episodes older than 14 days from retention to avoid
		if ep.Date.Before(time.Now().AddDate(0, 0, -14)) {
			RemoveNotifiedEpisode(ep)
		}
	}
}

func CleanOldNotifiedMovies() {
	for _, m := range retentionData.NotifiedMovies {
		// Remove notified movies older than 14 days from retention to avoid
		if m.Date.Before(time.Now().AddDate(0, 0, -14)) {
			RemoveNotifiedMovie(m)
		}
	}
}

func RemoveNotifiedEpisode(e Episode) {
	delete(retentionData.NotifiedEpisodes, e.Id)
}

func RemoveNotifiedMovie(m Movie) {
	delete(retentionData.NotifiedMovies, m.Id)
}

func RemoveDownloadedEpisode(e Episode) {
	delete(retentionData.DownloadedEpisodes, e.Id)
}

func RemoveDownloadedMovie(m Movie) {
	delete(retentionData.DownloadedMovies, m.Id)
}

func RemoveDownloadingEpisode(e Episode) {
	delete(retentionData.DownloadingEpisodes, e.Id)
}

func RemoveDownloadingMovie(m Movie) {
	delete(retentionData.DownloadingMovies, m.Id)
}

func RemoveFailedEpisode(e Episode) {
	delete(retentionData.FailedEpisodes, e.Id)
}

func RemoveFailedMovie(m Movie) {
	delete(retentionData.FailedMovies, m.Id)
}

func SaveTraktToken(token string) {
	retentionData.TraktToken = token
}

func LoadTraktToken() string {
	return retentionData.TraktToken
}

func AddShowToWatchlistRetention(watchlistName string, show string) {
	retentionData.Watchlists[watchlistName] = append(retentionData.Watchlists[watchlistName], show)
}

func GetShowsFromWatchlist(watchlistName string) []string {
	if _, ok := retentionData.Watchlists[watchlistName]; ok {
		return retentionData.Watchlists[watchlistName]
	} else {
		return []string{}
	}
}
