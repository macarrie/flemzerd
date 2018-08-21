package mock

import (
	"errors"
	"fmt"

	. "github.com/macarrie/flemzerd/objects"
)

type Downloader struct{}
type ErrorDownloader struct{}

var testTorrentsCount int

func (d Downloader) Status() (Module, error) {
	return Module{
		Name: "Downloader",
		Type: "provider",
		Status: ModuleStatus{
			Alive:   true,
			Message: "",
		},
	}, nil
}

func (d ErrorDownloader) Status() (Module, error) {
	var err error = fmt.Errorf("Notifier error")
	return Module{
		Name: "ErrorNotifier",
		Type: "notifier",
		Status: ModuleStatus{
			Alive:   false,
			Message: err.Error(),
		},
	}, err
}

func (w Downloader) GetName() string {
	return "Downloader"
}
func (w ErrorDownloader) GetName() string {
	return "ErrorDownloader"
}

func (d Downloader) AddTorrent(t Torrent) (string, error) {
	testTorrentsCount += 1
	return "id", nil
}

func (d ErrorDownloader) AddTorrent(t Torrent) (string, error) {
	return "id", errors.New("Downloader error")
}

func (d Downloader) AddTorrentMapping(flemzerId string, downloaderId string) {
	return
}

func (d ErrorDownloader) AddTorrentMapping(flemzerId string, downloaderId string) {
	return
}

func (d Downloader) RemoveTorrent(t Torrent) error {
	if testTorrentsCount > 0 {
		testTorrentsCount -= 1
	}
	return nil
}

func (d ErrorDownloader) RemoveTorrent(t Torrent) error {
	return errors.New("Downloader error")
}

func (d Downloader) Init() error {
	return nil
}

func (d ErrorDownloader) Init() error {
	return nil
}

func (d Downloader) GetTorrentStatus(t Torrent) (int, error) {
	return TORRENT_SEEDING, nil
}

func (d ErrorDownloader) GetTorrentStatus(t Torrent) (int, error) {
	return TORRENT_STOPPED, errors.New("Downloader error")
}

func (d Downloader) GetTorrentCount() int {
	return testTorrentsCount
}
func (d ErrorDownloader) GetTorrentCount() int {
	return testTorrentsCount
}
