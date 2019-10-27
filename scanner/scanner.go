package scanner

import (
	"os"
	"path/filepath"

	"github.com/macarrie/flemzerd/configuration"
	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
	"github.com/macarrie/flemzerd/vidocq"

	"github.com/rs/xid"
)

func scan_dir(directory string, media_type int) ([]MediaInfo, error) {
	fileList := []string{}
	err := filepath.Walk(directory, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			fileList = append(fileList, path)
		}
		return nil
	})

	if err != nil {
		return []MediaInfo{}, err
	}

	media_list := []MediaInfo{}
	for _, file := range fileList {
		info, err := vidocq.GetInfo(file, media_type)
		if err != nil {
			log.WithFields(log.Fields{
				"filepath": file,
			}).Error("Could not get info from vidocq for filename from library")
		} else {
			if info.Container != "" {
				info.Id = xid.New().String()
				media_list = append(media_list, info)
			}
		}
	}

	return media_list, nil
}

func removeDuplicates(array []MediaInfo) []MediaInfo {
	occurences := make(map[string]bool)
	var ret []MediaInfo

	for _, media := range array {
		if !occurences[media.Title] {
			occurences[media.Title] = true
			ret = append(ret, media)
		}
	}

	return ret
}

func groupByTvShow(list []MediaInfo) MediaInfoGroupedByShow {
	var retList = make(MediaInfoGroupedByShow)
	for _, episode := range list {
		if len(retList[episode.Title]) == 0 {
			retList[episode.Title] = make(MediaInfoSeasons)
		}
		if len(retList[episode.Title][episode.Season]) == 0 {
			retList[episode.Title][episode.Season] = make(MediaInfoEpisodes)
		}

		retList[episode.Title][episode.Season][episode.Episode] = episode
	}

	return retList
}

func ScanMovies() ([]MediaInfo, error) {
	movies, err := scan_dir(configuration.Config.Library.MoviePath, MOVIE)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Could not scan movie library")
	}

	return removeDuplicates(movies), err
}

func ScanShows() (MediaInfoGroupedByShow, error) {
	episodes, err := scan_dir(configuration.Config.Library.ShowPath, EPISODE)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Could not scan TV show library")
	}

	return groupByTvShow(episodes), err
}
