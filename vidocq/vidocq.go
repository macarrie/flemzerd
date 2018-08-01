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

var localVidocqAvailable bool

func init() {
	_, err := exec.LookPath("vidocq")
	if err != nil {
		log.Debug("Local vidocq executable not found. Quality filters and torrent control will not be done")
		localVidocqAvailable = false
	} else {
		localVidocqAvailable = true
	}
}

func GetInfo(name string) (MediaInfo, error) {
	if !localVidocqAvailable {
		return MediaInfo{}, fmt.Errorf("vidocq not available locally")
	}

	cmd := exec.Command("vidocq", name)
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
