package cache

import (
	"fmt"
	"github.com/macarrie/flemzerd/cachable"
	"github.com/pkg/errors"
	"path/filepath"
	"time"

	"github.com/macarrie/flemzerd/helpers"
	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
)

func Clear(c cachable.Cachable) {
	c.ClearCache()
}

func CacheMoviePoster(m *Movie) (string, error) {
	return cacheItem(m, "poster", m.Poster, m.DeletedAt)
}

func CacheMovieFanart(m *Movie) (string, error) {
	return cacheItem(m, "fanart", m.Background, m.DeletedAt)
}

func CacheShowPoster(s *TvShow) (string, error) {
	return cacheItem(s, "poster", s.Poster, s.DeletedAt)
}

func CacheShowFanart(s *TvShow) (string, error) {
	return cacheItem(s, "fanart", s.Background, s.DeletedAt)
}

func cacheItem(item cachable.Cachable, cacheType string, source string, deleted *time.Time) (cachePath string, err error) {
	if item.IsCached(source) || deleted != nil {
		return source, nil
	}

	var itemType string
	switch item.(type) {
	case *Movie:
		itemType = "movie"
	case *TvShow:
		itemType = "show"
	}

	item.GetLog().WithFields(log.Fields{
		"source": source,
	}).Info(fmt.Sprintf("Saving %s to cache", cacheType))

	if source == "" {
		item.GetLog().Debug(fmt.Sprintf("No %s to download because source path is empty", cacheType))

		return "", nil
	}

	fileExtension := filepath.Ext(source)
	cacheFilePath := fmt.Sprintf("%s/%ss/%d_%s%s", CACHE_BASE_PATH, itemType, item.GetId(), cacheType, fileExtension)

	if err := helpers.DownloadImage(source, cacheFilePath, HTTP_TIMEOUT); err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Could not download %s %s to cache", itemType, cacheType))
	}

	item.GetLog().WithFields(log.Fields{
		"source": source,
	}).Info("Item saved to cache")

	return fmt.Sprintf("/cache/%ss/%d_%s%s", itemType, item.GetId(), cacheType, fileExtension), nil
}
