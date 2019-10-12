package vidocq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"

	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"

	"github.com/pkg/errors"
)

var LocalVidocqAvailable bool

func init() {
	CheckVidocq()
}

func CheckVidocq() {
	_, err := exec.LookPath("vidocq")
	if err != nil {
		log.Error("Local vidocq executable not found. Quality filters and torrent control will not be done")
		LocalVidocqAvailable = false
		return
	}

	cmd := exec.Command("vidocq", "--help")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Vidocq found but could not be executed correctly")
		LocalVidocqAvailable = false
		return
	}

	LocalVidocqAvailable = true
}

func GetInfo(name string, media_type ...int) (MediaInfo, error) {
	if !LocalVidocqAvailable {
		return MediaInfo{}, fmt.Errorf("vidocq not available locally")
	}

	var cmd *exec.Cmd
	if len(media_type) > 0 && media_type[0] == MOVIE {
		cmd = exec.Command("vidocq", "--type", "movie", name)
	} else if len(media_type) > 0 && media_type[0] == EPISODE {
		cmd = exec.Command("vidocq", "--type", "episode", name)
	} else {
		cmd = exec.Command("vidocq", name)
	}

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		log.Debug("[Vidocq] Could not get execute vidocq: ", err)
		return MediaInfo{}, errors.Wrap(err, "cannot execute vidocq")
	}

	var mediaInfo MediaInfo
	parseErr := json.Unmarshal(stdout.Bytes(), &mediaInfo)
	if parseErr != nil {
		log.Error("[Vidocq] Could not parse request result:", parseErr)
		return MediaInfo{}, errors.Wrap(err, "cannot parse vidocq request result")
	}

	return mediaInfo, nil
}
