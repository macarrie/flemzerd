package transmission

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	tr "github.com/macarrie/transmission"

	"github.com/macarrie/flemzerd/configuration"
	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"

	"github.com/pkg/errors"
	"github.com/rs/xid"
	"github.com/upgear/go-kit/retry"
)

var module Module
var transmissionClient tr.Client

// Since torrents in transmission have specific ID and torrent objects in flemzer have their own ID, we need to know which transmission torrent correspond to which flemzer torrent
// This map stores "flemzerd torrent id" -> "transmission_torrent_id" relations
var torrentsMapping map[string]string
var torrentsMappingMutex sync.Mutex

var torrentList []*tr.Torrent

type TransmissionDownloader struct {
	Address   string
	Port      int
	User      string
	Password  string
	SessionId string
}

func convertTorrent(t tr.Torrent) Torrent {
	id := xid.New()
	// TODO: Compute ETA into Time format
	return Torrent{
		TorrentId:    id.String(),
		Name:         t.Name,
		Link:         t.TorrentFile,
		DownloadDir:  t.DownloadDir,
		PercentDone:  t.PercentDone,
		TotalSize:    t.TotalSize,
		RateDownload: t.RateDownload,
		RateUpload:   t.RateUpload,
	}
}

func updateTorrentList() error {
	torrents, err := transmissionClient.GetTorrents()
	if err != nil {
		return errors.Wrap(err, "cannot get torrent list from transmission")
	}

	torrentList = torrents

	return nil
}

func getTransmissionTorrent(t Torrent) (tr.Torrent, error) {
	for _, transmissionTorrent := range torrentList {
		torrentsMappingMutex.Lock()
		id := torrentsMapping[t.TorrentId]
		torrentsMappingMutex.Unlock()

		if id == transmissionTorrent.HashString {
			return *transmissionTorrent, nil
		}
	}

	return tr.Torrent{}, errors.New("Could not find corresponding transmission torrent")
}

func New(address string, port int, user string, password string) *TransmissionDownloader {
	module = Module{
		Name: "transmission",
		Type: "downloader",
		Status: ModuleStatus{
			Alive:   true,
			Message: "",
		},
	}

	return &TransmissionDownloader{
		Address:  address,
		Port:     port,
		User:     user,
		Password: password,
	}
}

func (d *TransmissionDownloader) Init() error {
	client, err := tr.New(tr.Config{
		Address:  fmt.Sprintf("http://%s:%d/transmission/rpc", d.Address, d.Port),
		User:     d.User,
		Password: d.Password,
	})
	if err != nil {
		return errors.Wrap(err, "cannot create transmission client")
	}

	transmissionClient = *client

	torrentsMappingMutex.Lock()
	torrentsMapping = make(map[string]string)
	torrentsMappingMutex.Unlock()

	return nil
}

func (d *TransmissionDownloader) GetName() string {
	return "transmission"
}

func (d *TransmissionDownloader) Status() (Module, error) {
	log.Debug("Checking transmission downloader status")
	client := &http.Client{
		Timeout: time.Duration(HTTP_TIMEOUT * time.Second),
	}
	url := fmt.Sprintf("http://%s:%d/transmission/rpc", d.Address, d.Port)
	params := map[string]string{
		"method":    "port-test",
		"arguments": "",
	}

	jsonParams, _ := json.Marshal(params)

	request, err := http.NewRequest("POST", url, bytes.NewReader(jsonParams))
	if err != nil {
		return module, errors.Wrap(err, "cannot not create HTTP request for transmission")
	}
	request.SetBasicAuth(d.User, d.Password)

	res, err := client.Do(request)
	if err != nil {
		module.Status.Alive = false
		module.Status.Message = err.Error()

		return module, errors.Wrap(err, "cannot perform HTTP request to transmission")
	} else if res.StatusCode == http.StatusUnauthorized {
		module.Status.Alive = false
		module.Status.Message = "Credentials refused when attempting to connect to transmission daemon"

		return module, fmt.Errorf(module.Status.Message)
	} else if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusConflict {
		module.Status.Alive = true
		module.Status.Message = ""

		var retError error = nil
		if _, err := transmissionClient.FreeSpace(configuration.Config.Library.CustomTmpPath); err != nil {
			retError = errors.Wrapf(err, "Error while checking access to show library path %s", configuration.Config.Library.CustomTmpPath)
			module.Status.Alive = false
			module.Status.Message = retError.Error()
		}

		return module, retError
	} else {
		module.Status.Alive = false
		module.Status.Message = fmt.Sprintf("Unknow return status code: %d", res.StatusCode)

		return module, fmt.Errorf(module.Status.Message)
	}
}

func (d TransmissionDownloader) AddTorrent(t Torrent) (string, error) {
	torrent, err := transmissionClient.AddTorrent(tr.AddTorrentArg{
		Filename:    t.Link,
		DownloadDir: t.DownloadDir,
	})
	if err != nil {
		return "", errors.Wrap(err, "cannot add torrent to transmission")
	}

	d.AddTorrentMapping(t.TorrentId, torrent.HashString)
	if err = retry.Double(3).Run(func() error {
		return updateTorrentList()
	}); err != nil {
		return "", errors.Wrap(err, "cannot update transmission torrent list")
	}

	return torrent.HashString, nil
}

func (d TransmissionDownloader) AddTorrentMapping(flemzerID string, transmissionID string) {
	torrentsMappingMutex.Lock()
	torrentsMapping[flemzerID] = transmissionID
	torrentsMappingMutex.Unlock()
}

func (d TransmissionDownloader) RemoveTorrent(t Torrent) error {
	transmissionTorrent, err := getTransmissionTorrent(t)
	if err != nil {
		return errors.Wrap(err, "cannot get torrent list from transmission")
	}

	removeErr := transmissionClient.RemoveTorrents([]*tr.Torrent{&transmissionTorrent}, false)
	if removeErr != nil {
		return errors.Wrap(removeErr, "cannot remove torrent from transmission")
	}

	return nil
}

func (d TransmissionDownloader) GetTorrentStatus(t *Torrent) error {
	if err := retry.Double(3).Run(func() error {
		return updateTorrentList()
	}); err != nil {
		t.Status = TORRENT_UNKNOWN_STATUS
		return errors.Wrap(err, "cannot update transmission torrent list")
	}

	for _, torrent := range torrentList {
		torrentsMappingMutex.Lock()
		id := torrentsMapping[t.TorrentId]
		torrentsMappingMutex.Unlock()

		if id == torrent.HashString {
			err := retry.Double(3).Run(func() error {
				return torrent.Update()
			})
			if err != nil {
				t.Status = TORRENT_UNKNOWN_STATUS
				return nil
			}

			t.ETA = time.Unix(torrent.Eta, 0)
			t.PercentDone = torrent.PercentDone
			t.TotalSize = torrent.TotalSize
			t.RateDownload = torrent.RateDownload
			t.RateUpload = torrent.RateUpload

			switch (*torrent).Status {
			case tr.StatusDownloading:
				t.Status = TORRENT_DOWNLOADING
				return nil
			case tr.StatusSeeding:
				t.Status = TORRENT_SEEDING
				return nil
			case tr.StatusStopped:
				t.Status = TORRENT_STOPPED
				return nil
			case tr.StatusDownloadPending:
				t.Status = TORRENT_DOWNLOAD_PENDING
				return nil
			default:
				t.Status = TORRENT_UNKNOWN_STATUS
				return nil
			}
		}
	}

	t.Status = TORRENT_UNKNOWN_STATUS
	return errors.New("Could not find torrent in transmission")
}
