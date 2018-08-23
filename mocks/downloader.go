package mock

import (
	"errors"
	"fmt"

	. "github.com/macarrie/flemzerd/objects"
)

type Downloader struct{}
type ErrorDownloader struct{}
type StalledDownloader struct{}
type DLErrorDownloader struct{}

var testTorrentsCount int

// Status
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
func (d StalledDownloader) Status() (Module, error) {
	var err error = fmt.Errorf("Notifier error")
	return Module{
		Name: "StalledDownloader",
		Type: "notifier",
		Status: ModuleStatus{
			Alive:   false,
			Message: err.Error(),
		},
	}, err
}
func (d DLErrorDownloader) Status() (Module, error) {
	var err error = fmt.Errorf("Notifier error")
	return Module{
		Name: "DLErrorDownloader",
		Type: "notifier",
		Status: ModuleStatus{
			Alive:   false,
			Message: err.Error(),
		},
	}, err
}

// GetName
func (d Downloader) GetName() string {
	return "Downloader"
}
func (d ErrorDownloader) GetName() string {
	return "ErrorDownloader"
}
func (d StalledDownloader) GetName() string {
	return "StalledDownloader"
}
func (d DLErrorDownloader) GetName() string {
	return "DLErrorDownloader"
}

// AddTorrent
func (d Downloader) AddTorrent(t Torrent) (string, error) {
	testTorrentsCount += 1
	return "id", nil
}
func (d ErrorDownloader) AddTorrent(t Torrent) (string, error) {
	return "id", errors.New("Downloader error")
}
func (d StalledDownloader) AddTorrent(t Torrent) (string, error) {
	testTorrentsCount += 1
	return "id", nil
}
func (d DLErrorDownloader) AddTorrent(t Torrent) (string, error) {
	testTorrentsCount += 1
	return "id", nil
}

// AddTorrentMapping
func (d Downloader) AddTorrentMapping(flemzerId string, downloaderId string) {
	return
}
func (d ErrorDownloader) AddTorrentMapping(flemzerId string, downloaderId string) {
	return
}
func (d StalledDownloader) AddTorrentMapping(flemzerId string, downloaderId string) {
	return
}
func (d DLErrorDownloader) AddTorrentMapping(flemzerId string, downloaderId string) {
	return
}

// RemoveTorrent
func (d Downloader) RemoveTorrent(t Torrent) error {
	if testTorrentsCount > 0 {
		testTorrentsCount -= 1
	}
	return nil
}
func (d ErrorDownloader) RemoveTorrent(t Torrent) error {
	return errors.New("Downloader error")
}
func (d StalledDownloader) RemoveTorrent(t Torrent) error {
	if testTorrentsCount > 0 {
		testTorrentsCount -= 1
	}
	return nil
}
func (d DLErrorDownloader) RemoveTorrent(t Torrent) error {
	if testTorrentsCount > 0 {
		testTorrentsCount -= 1
	}
	return nil
}

// Init
func (d Downloader) Init() error {
	return nil
}
func (d ErrorDownloader) Init() error {
	return nil
}
func (d StalledDownloader) Init() error {
	return nil
}
func (d DLErrorDownloader) Init() error {
	return nil
}

// GetTorrentStatus
func (d Downloader) GetTorrentStatus(t Torrent) (int, error) {
	return TORRENT_SEEDING, nil
}
func (d ErrorDownloader) GetTorrentStatus(t Torrent) (int, error) {
	return TORRENT_STOPPED, errors.New("Downloader error")
}
func (d StalledDownloader) GetTorrentStatus(t Torrent) (int, error) {
	return TORRENT_DOWNLOADING, nil
}
func (d DLErrorDownloader) GetTorrentStatus(t Torrent) (int, error) {
	return TORRENT_STOPPED, nil
}

// GetTorrentCount
func (d Downloader) GetTorrentCount() int {
	return testTorrentsCount
}
func (d ErrorDownloader) GetTorrentCount() int {
	return testTorrentsCount
}
func (d StalledDownloader) GetTorrentCount() int {
	return testTorrentsCount
}
func (d DLErrorDownloader) GetTorrentCount() int {
	return testTorrentsCount
}
