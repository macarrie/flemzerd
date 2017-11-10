package transmission

import (
	"fmt"

	. "github.com/macarrie/flemzerd/objects"
	tr "github.com/odwrtw/transmission"
)

var transmissionClient tr.Client
var torrentList []tr.Torrent

type TransmissionDownloader struct {
	Address   string
	Port      int
	User      string
	Password  string
	SessionId string
}

func convertTorrent(t tr.Torrent) Torrent {
	return Torrent{
		Name:        t.Name,
		Link:        t.TorrentFile,
		DownloadDir: t.DownloadDir,
	}
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

	return nil
}

func (d TransmissionDownloader) AddTorrent(t Torrent) error {
	torrent, err := transmissionClient.Add(t.Link)
	if err != nil {
		return err
	}

	torrentList = append(torrentList, *torrent)
	return nil
}
