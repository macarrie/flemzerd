package media_helper

import (
	. "github.com/macarrie/flemzerd/objects"
)

func GetMovieTitle(m Movie) string {
	if m.CustomTitle != "" {
		return m.CustomTitle
	}

	if m.UseDefaultTitle {
		return m.Title
	}

	return m.OriginalTitle
}

func GetShowTitle(s TvShow) string {
	if s.CustomTitle != "" {
		return s.CustomTitle
	}

	if s.UseDefaultTitle {
		return s.Title
	}

	return s.OriginalTitle
}
