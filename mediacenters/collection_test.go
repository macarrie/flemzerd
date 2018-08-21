package mediacenter

import (
	"testing"

	"github.com/macarrie/flemzerd/db"
	mock "github.com/macarrie/flemzerd/mocks"
)

func init() {
	db.DbPath = "/tmp/flemzerd.db"
	db.Load()
}

func TestAddMediaCenter(t *testing.T) {
	mcLength := len(mediaCenterCollection)
	m := mock.MediaCenter{}
	AddMediaCenter(m)

	if len(mediaCenterCollection) != mcLength+1 {
		t.Error("Expected ", mcLength+1, " providers, got ", len(mediaCenterCollection))
	}
}

func TestStatus(t *testing.T) {
	m1 := mock.MediaCenter{}
	m2 := mock.MediaCenter{}

	mediaCenterCollection = []MediaCenter{m1, m2}

	mods, err := Status()
	if err != nil {
		t.Error("Expected not to have error for mediacenter status")
	}
	if len(mods) != 2 {
		t.Errorf("Expected to have 2 mediacenter modules status, got %d instead", len(mods))
	}

	AddMediaCenter(mock.ErrorMediaCenter{})
	_, err = Status()
	if err == nil {
		t.Error("Expected to have aggregated error for mediacenter status")
	}
}

func TestReset(t *testing.T) {
	m := mock.MediaCenter{}
	AddMediaCenter(m)
	Reset()

	if len(mediaCenterCollection) != 0 {
		t.Error("Expected mediacenter collection to be empty after reset")
	}
}

func TestRefreshLibrary(t *testing.T) {
	db.ResetDb()

	mc1 := mock.MediaCenter{}
	count := mc1.GetRefreshCount()
	AddMediaCenter(mc1)
	AddMediaCenter(mock.ErrorMediaCenter{})
	RefreshLibrary()

	if mc1.GetRefreshCount() != count+1 {
		t.Errorf("Expected library to have been refresh 1 time, got %d refresh instead", mc1.GetRefreshCount())
	}

}

func TestGetMediaCenter(t *testing.T) {
	mc1 := mock.MediaCenter{}
	mediaCenterCollection = []MediaCenter{mc1}

	if _, err := GetMediaCenter("Unknown"); err == nil {
		t.Error("Expected to have error when getting unknown mediacenter, got none")
	}

	if _, err := GetMediaCenter("MediaCenter"); err != nil {
		t.Errorf("Got error while retrieving known mediacenter: %s", err.Error())
	}
}
