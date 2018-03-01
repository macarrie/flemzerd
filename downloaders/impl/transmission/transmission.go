package transmission

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
	tr "github.com/odwrtw/transmission"
	"github.com/rs/xid"
)

var transmissionClient tr.Client

// Since torrents in transmission have specific ID and torrent objects in flemzer have their own ID, we need to know which transmission torrent correspond to which flemzer torrent
// This map stores "flemzerd torrent id" -> "transmission_torrent_id" relations
var torrentsMapping map[string]string
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
	return Torrent{
		Id:          id.String(),
		Name:        t.Name,
		Link:        t.TorrentFile,
		DownloadDir: t.DownloadDir,
	}
}

func updateTorrentList() {
	torrentList, _ = transmissionClient.GetTorrents()
	return
}

func getTransmissionTorrent(t Torrent) (tr.Torrent, error) {
	for _, transmissionTorrent := range torrentList {
		if torrentsMapping[t.Id] == transmissionTorrent.HashString {
			return *transmissionTorrent, nil
		}
	}

	return tr.Torrent{}, errors.New("Could not find corresponding transmission torrent")
}

func New(address string, port int, user string, password string) *TransmissionDownloader {
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
		return err
	}

	transmissionClient = *client
	torrentsMapping = make(map[string]string)

	return nil
}

func (d *TransmissionDownloader) Status() (Module, error) {
	log.Debug("Checking transmission downloader status")
	client := &http.Client{
		Timeout: time.Duration(HTTP_TIMEOUT * time.Second),
	}

	msg := ""
	var alive bool
	_, err := client.Get(fmt.Sprintf("http://%s:%d", d.Address, d.Port))
	if err != nil {
		alive = false
		msg = err.Error()
	} else {
		alive = true
	}

	return Module{
		Name: "transmission",
		Type: "downloader",
		Status: ModuleStatus{
			Alive:   alive,
			Message: msg,
		},
	}, err
}

func (d TransmissionDownloader) AddTorrent(t Torrent) (string, error) {
	torrent, err := transmissionClient.AddTorrent(tr.AddTorrentArg{
		Filename:    t.Link,
		DownloadDir: t.DownloadDir,
	})
	if err != nil {
		return "", err
	}

	d.AddTorrentMapping(t.Id, torrent.HashString)
	updateTorrentList()
	return torrent.HashString, nil
}

func (d TransmissionDownloader) AddTorrentMapping(flemzerID string, transmissionID string) {
	torrentsMapping[flemzerID] = transmissionID
}

func (d TransmissionDownloader) RemoveTorrent(t Torrent) error {
	transmissionTorrent, err := getTransmissionTorrent(t)
	if err != nil {
		return err
	}

	removeErr := transmissionClient.RemoveTorrents([]*tr.Torrent{&transmissionTorrent}, false)
	if removeErr != nil {
		return removeErr
	}

	return nil
}

func (d TransmissionDownloader) GetTorrentStatus(t Torrent) (int, error) {
	updateTorrentList()
	for _, torrent := range torrentList {
		if torrentsMapping[t.Id] == torrent.HashString {
			torrent.Update()

			switch (*torrent).Status {
			case tr.StatusDownloading:
				return TORRENT_DOWNLOADING, nil
			case tr.StatusSeeding:
				return TORRENT_SEEDING, nil
			case tr.StatusStopped:
				return TORRENT_STOPPED, nil
			case tr.StatusDownloadPending:
				return TORRENT_DOWNLOAD_PENDING, nil
			default:
				return TORRENT_UNKNOWN_STATUS, nil
			}
		}
	}

	return 0, errors.New("Could not find torrent in transmission")
}
