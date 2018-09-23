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
	if s.CustomName != "" {
		return s.CustomName
	}

	if s.UseDefaultTitle {
		return s.Name
	}

	return s.OriginalName
}
