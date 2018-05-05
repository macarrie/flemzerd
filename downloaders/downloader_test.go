package downloader

import (
	"errors"
	"fmt"
	"strconv"

	. "github.com/macarrie/flemzerd/objects"
)

type MockDownloader struct{}
type MockErrorDownloader struct{}

var testTorrentsCount int

func (d MockDownloader) Status() (Module, error) {
	return Module{
		Name: "MockDownloader",
		Type: "provider",
		Status: ModuleStatus{
			Alive:   true,
			Message: "",
		},
	}, nil
}

func (d MockErrorDownloader) Status() (Module, error) {
	var err error = fmt.Errorf("Notifier error")
	return Module{
		Name: "MockErrorNotifier",
		Type: "notifier",
		Status: ModuleStatus{
			Alive:   false,
			Message: err.Error(),
		},
	}, err
}

func (w MockDownloader) GetName() string {
	return "MockDownloader"
}
func (w MockErrorDownloader) GetName() string {
	return "MockErrorDownloader"
}

func (d MockDownloader) AddTorrent(t Torrent) (string, error) {
	testTorrentsCount += 1
	return "id", nil
}

func (d MockErrorDownloader) AddTorrent(t Torrent) (string, error) {
	return "id", errors.New("Downloader error")
}

func (d MockDownloader) AddTorrentMapping(flemzerId string, downloaderId string) {
	return
}

func (d MockErrorDownloader) AddTorrentMapping(flemzerId string, downloaderId string) {
	return
}

func (d MockDownloader) RemoveTorrent(t Torrent) error {
	if testTorrentsCount > 0 {
		testTorrentsCount -= 1
	}
	return nil
}

func (d MockErrorDownloader) RemoveTorrent(t Torrent) error {
	return errors.New("Downloader error")
}

func (d MockDownloader) Init() error {
	return nil
}

func (d MockErrorDownloader) Init() error {
	return nil
}

func (d MockDownloader) GetTorrentStatus(t Torrent) (int, error) {
	status, _ := strconv.Atoi(t.TorrentId)
	return status, nil
}

func (d MockErrorDownloader) GetTorrentStatus(t Torrent) (int, error) {
	return 0, errors.New("Downloader error")
}
