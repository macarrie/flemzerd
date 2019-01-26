package notifiers_helper

import (
	"fmt"

	media_helper "github.com/macarrie/flemzerd/helpers/media"
	. "github.com/macarrie/flemzerd/objects"
)

func GetNotificationText(notif Notification) (notif_title, notif_content string, err error) {
	title := ""
	content := ""

	switch notif.Type {
	case NOTIFICATION_NEW_EPISODE:
		title = fmt.Sprintf("%v S%03dE%03d: New episode aired", media_helper.GetShowTitle(notif.Episode.TvShow), notif.Episode.Season, notif.Episode.Number)
		content = fmt.Sprintf("New episode aired on %v\n%v Season %03d Episode %03d: %v", notif.Episode.Date, media_helper.GetShowTitle(notif.Episode.TvShow), notif.Episode.Season, notif.Episode.Number, notif.Episode.Title)

	case NOTIFICATION_NEW_MOVIE:
		title = fmt.Sprintf("%s", media_helper.GetMovieTitle(notif.Movie))
		content = "Movie found in watchlist, adding to tracked movies"

	case NOTIFICATION_DOWNLOAD_START:
		if notif.Episode.ID != 0 {
			title = fmt.Sprintf("%v S%03dE%03d: Download start ", media_helper.GetShowTitle(notif.Episode.TvShow), notif.Episode.Season, notif.Episode.Number)
			content = "Torrents found for episode. Starting download"
		}
		if notif.Movie.ID != 0 {
			title = fmt.Sprintf("%v: Download start", media_helper.GetMovieTitle(notif.Movie))
			content = "Torrents found for movie. Starting download"
		}

	case NOTIFICATION_DOWNLOAD_SUCCESS:
		if notif.Episode.ID != 0 {
			title = fmt.Sprintf("%v S%03dE%03d: Episode downloaded", media_helper.GetShowTitle(notif.Episode.TvShow), notif.Episode.Season, notif.Episode.Number)
			content = fmt.Sprintf("New episode downloaded\n%v Season %03d Episode %03d: %v", media_helper.GetShowTitle(notif.Episode.TvShow), notif.Episode.Season, notif.Episode.Number, notif.Episode.Title)
		}
		if notif.Movie.ID != 0 {
			title = fmt.Sprintf("%v: Movie downloaded", media_helper.GetMovieTitle(notif.Movie))
			content = "New movie downloaded\n"
		}

	case NOTIFICATION_DOWNLOAD_FAILURE:
		if notif.Episode.ID != 0 {
			title = fmt.Sprintf("%v S%03dE%03d: Episode download failed", media_helper.GetShowTitle(notif.Episode.TvShow), notif.Episode.Season, notif.Episode.Number)
			content = fmt.Sprintf("Failed to download episode\n%v Season %03d Episode %03d: %v", media_helper.GetShowTitle(notif.Episode.TvShow), notif.Episode.Season, notif.Episode.Number, notif.Episode.Title)
		}
		if notif.Movie.ID != 0 {
			title = fmt.Sprintf("%v: Movie download failed", media_helper.GetMovieTitle(notif.Movie))
			content = fmt.Sprintf("Failed to download movie: %v", media_helper.GetMovieTitle(notif.Movie))
		}

	case NOTIFICATION_TEXT:
		title = notif.Title
		content = notif.Content

	case NOTIFICATION_NO_TORRENTS:
		if notif.Episode.ID != 0 {
			title = fmt.Sprintf("%v S%03dE%03d: No torrents found", notif.Episode.TvShow.OriginalTitle, notif.Episode.Season, notif.Episode.Number)
			content = fmt.Sprintf("No torrents found. Try adding indexers or wait for torrents to be available in case of a recent release\n%v Season %03d Episode %03d: %v", notif.Episode.TvShow.OriginalTitle, notif.Episode.Season, notif.Episode.Number, notif.Episode.Title)
		}
		if notif.Movie.ID != 0 {
			title = fmt.Sprintf("%v: Movie download failed", notif.Movie.OriginalTitle)
			content = fmt.Sprintf("Failed to download movie: %v", notif.Movie.OriginalTitle)
		}

	default:
		return "", "", fmt.Errorf("Unable to send notification: Unknown notification type (%d)", notif.Type)
	}

	return title, content, nil
}

func GetNotificationType(notif Notification) string {
	switch notif.Type {
	case NOTIFICATION_NEW_EPISODE:
		return "new_episode"
	case NOTIFICATION_NEW_MOVIE:
		return "new_movie"
	case NOTIFICATION_DOWNLOAD_START:
		return "download_start"
	case NOTIFICATION_DOWNLOAD_SUCCESS:
		return "download_success"
	case NOTIFICATION_DOWNLOAD_FAILURE:
		return "download_failure"
	case NOTIFICATION_TEXT:
		return "text"
	case NOTIFICATION_NO_TORRENTS:
		return "no_torrents"
	default:
		return "unknown"
	}
}
