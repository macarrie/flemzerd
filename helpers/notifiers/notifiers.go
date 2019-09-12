package notifiers_helper

import (
	"fmt"
	. "github.com/macarrie/flemzerd/objects"
)

func getShowTitle(s TvShow) string {
	if s.ID != 0 {
		return s.GetTitle()
	}

	return "Unknown show"
}

func getMovieTitle(m Movie) string {
	if m.ID != 0 {
		return m.GetTitle()
	}

	return "Unknown movie"
}

func GetNotificationText(notif Notification) (notif_title, notif_content string, err error) {
	title := ""
	content := ""

	switch notif.Type {
	case NOTIFICATION_NEW_EPISODE:
		title = fmt.Sprintf("%v S%03dE%03d: New episode aired", getShowTitle(notif.Episode.TvShow), notif.Episode.Season, notif.Episode.Number)
		content = fmt.Sprintf("New episode aired on %v\n%v Season %03d Episode %03d: %v", notif.Episode.Date, getShowTitle(notif.Episode.TvShow), notif.Episode.Season, notif.Episode.Number, notif.Episode.Title)

	case NOTIFICATION_NEW_MOVIE:
		title = fmt.Sprintf("%s", notif.Movie.GetTitle())
		content = "Movie found in watchlist, adding to tracked movies"

	case NOTIFICATION_DOWNLOAD_START:
		if notif.Episode.ID != 0 {
			title = fmt.Sprintf("%v S%03dE%03d: Download start ", getShowTitle(notif.Episode.TvShow), notif.Episode.Season, notif.Episode.Number)
			content = "Torrents found for episode. Starting download"
		}
		if notif.Movie.ID != 0 {
			title = fmt.Sprintf("%v: Download start", getMovieTitle(notif.Movie))
			content = "Torrents found for movie. Starting download"
		}

	case NOTIFICATION_DOWNLOAD_SUCCESS:
		if notif.Episode.ID != 0 {
			title = fmt.Sprintf("%v S%03dE%03d: Episode downloaded", getShowTitle(notif.Episode.TvShow), notif.Episode.Season, notif.Episode.Number)
			content = fmt.Sprintf("New episode downloaded\n%v Season %03d Episode %03d: %v", getShowTitle(notif.Episode.TvShow), notif.Episode.Season, notif.Episode.Number, notif.Episode.Title)
		}
		if notif.Movie.ID != 0 {
			title = fmt.Sprintf("%v: Movie downloaded", getMovieTitle(notif.Movie))
			content = "New movie downloaded\n"
		}

	case NOTIFICATION_DOWNLOAD_FAILURE:
		if notif.Episode.ID != 0 {
			title = fmt.Sprintf("%v S%03dE%03d: Episode download failed", getShowTitle(notif.Episode.TvShow), notif.Episode.Season, notif.Episode.Number)
			content = fmt.Sprintf("Failed to download episode\n%v Season %03d Episode %03d: %v. Maybe torrents could not be downloaded, or some essential modules for download were not available.", getShowTitle(notif.Episode.TvShow), notif.Episode.Season, notif.Episode.Number, notif.Episode.Title)
		}
		if notif.Movie.ID != 0 {
			title = fmt.Sprintf("%v: Movie download failed", getMovieTitle(notif.Movie))
			content = fmt.Sprintf("Failed to download movie: %v. Maybe torrents could not be downloaded, or some essential modules for download were not available.", getMovieTitle(notif.Movie))
		}

	case NOTIFICATION_TEXT:
		title = notif.Title
		content = notif.Content

	case NOTIFICATION_NO_TORRENTS:
		if notif.Episode.ID != 0 {
			title = fmt.Sprintf("%v S%03dE%03d: No torrents found", getShowTitle(notif.Episode.TvShow), notif.Episode.Season, notif.Episode.Number)
			content = fmt.Sprintf("No torrents found. Try adding indexers or wait for torrents to be available in case of a recent release\n%v Season %03d Episode %03d: %v", getShowTitle(notif.Episode.TvShow), notif.Episode.Season, notif.Episode.Number, notif.Episode.Title)
		}
		if notif.Movie.ID != 0 {
			title = fmt.Sprintf("%v: Movie download failed", getMovieTitle(notif.Movie))
			content = fmt.Sprintf("Failed to download movie: %v", getMovieTitle(notif.Movie))
		}

	default:
		return "", "", fmt.Errorf("Unable to send notification: Unknown notification type (%d)", notif.Type)
	}

	return title, content, nil
}
