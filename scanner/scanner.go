package scanner

import (
	"fmt"
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
				fmt.Println(file)
				fmt.Printf("File info: %+v\n", info)
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

func ScanMovies() ([]MediaInfo, error) {
	movies, err := scan_dir(configuration.Config.Library.MoviePath, MOVIE)
	fmt.Printf("Movies: %+v\n", movies)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Could not scan movie library")
	}

	movies = removeDuplicates(movies)
	for _, movie := range movies {
		fmt.Printf("Movie: %v\n", movie.Title)
	}

	return movies, err
}

func ScanShows() ([]MediaInfo, error) {
	episodes, err := scan_dir(configuration.Config.Library.ShowPath, EPISODE)
	fmt.Printf("Episodes: %+v\n", episodes)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Could not scan TV show library")
	}

	episodes = removeDuplicates(episodes)
	for _, ep := range episodes {
		fmt.Printf("Episode: %v\n", ep.Title)
	}

	return episodes, err
}
