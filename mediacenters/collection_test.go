package mediacenter

import (
	"testing"

	"github.com/macarrie/flemzerd/db"
)

func init() {
	db.DbPath = "/tmp/flemzerd.db"
	db.Load()
}

func TestAddMediaCenter(t *testing.T) {
	mcLength := len(mediaCenterCollection)
	m := MockMediaCenter{}
	AddMediaCenter(m)

	if len(mediaCenterCollection) != mcLength+1 {
		t.Error("Expected ", mcLength+1, " providers, got ", len(mediaCenterCollection))
	}
}

func TestStatus(t *testing.T) {
	m1 := MockMediaCenter{}
	m2 := MockMediaCenter{}

	mediaCenterCollection = []MediaCenter{m1, m2}

	mods, err := Status()
	if err != nil {
		t.Error("Expected not to have error for mediacenter status")
	}
	if len(mods) != 2 {
		t.Errorf("Expected to have 2 mediacenter modules status, got %d instead", len(mods))
	}

	AddMediaCenter(MockErrorMediaCenter{})
	_, err = Status()
	if err == nil {
		t.Error("Expected to have aggregated error for mediacenter status")
	}
}

func TestReset(t *testing.T) {
	m := MockMediaCenter{}
	AddMediaCenter(m)
	Reset()

	if len(mediaCenterCollection) != 0 {
		t.Error("Expected mediacenter collection to be empty after reset")
	}
}

func TestRefreshLibrary(t *testing.T) {
	db.ResetDb()

	refreshCounter = 0
	AddMediaCenter(MockMediaCenter{})
	AddMediaCenter(MockErrorMediaCenter{})
	RefreshLibrary()

	if refreshCounter != 1 {
		t.Errorf("Expected library to have been refresh 1 time, got %d refresh instead", refreshCounter)
	}

}
