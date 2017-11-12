package transmission

import (
	"errors"
	"fmt"

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

func addTorrentMapping(flemzerID string, transmissionID string) {
	torrentsMapping[flemzerID] = transmissionID
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

func (d TransmissionDownloader) AddTorrent(t Torrent) error {
	torrent, err := transmissionClient.Add(t.Link)
	if err != nil {
		return err
	}

	addTorrentMapping(t.Id, torrent.HashString)
	updateTorrentList()
	return nil
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
