package downloader

import (
	"testing"
)

type MockDownloader struct{}

func (d MockDownloader) AddTorrent(url string) error {
	return nil
}

func (d MockDownloader) Init() error {
	return nil
}

func TestAddDownloader(t *testing.T) {}
func TestAddTorrent(t *testing.T)    {}
