package retention

import (
	"testing"
	"time"

	. "github.com/macarrie/flemzerd/objects"
)

func TestLoad(t *testing.T) {
	// TODO
}

func TestSave(t *testing.T) {
	// TODO
}

func TestHasBeenNotified(t *testing.T) {
	retentionData = RetentionData{}
	testEpisode := Episode{
		Id: 1000,
	}

	if HasBeenNotified(testEpisode) {
		t.Error("Expected episode not to be in notified episodes retention")
	}

	retentionData = RetentionData{
		NotifiedEpisodes: []Episode{
			testEpisode,
		},
	}

	if !HasBeenNotified(testEpisode) {
		t.Error("Expected test episode to be present in notified episodes retention")
	}
}

func TestAddNotifiedEpisode(t *testing.T) {
	retentionData = RetentionData{}
	testEpisode := Episode{
		Id: 1000,
	}

	AddNotifiedEpisode(testEpisode)

	if len(retentionData.NotifiedEpisodes) != 1 {
		t.Error("Expected to have 1 item in notified episodes retention, got ", len(retentionData.NotifiedEpisodes), " items instead")
	}
}

func TestCleanOldNotifiedEpisodes(t *testing.T) {
	retentionData = RetentionData{
		NotifiedEpisodes: []Episode{
			Episode{
				Id:   1000,
				Date: time.Time{},
			},
		},
	}

	CleanOldNotifiedEpisodes()

	if len(retentionData.NotifiedEpisodes) > 0 {
		t.Error("Expected old episode to have been removed from notified episodes retention but episode is still present in retention")
	}
}

func TestRemoveNotifiedEpisode(t *testing.T) {
	retentionData = RetentionData{
		NotifiedEpisodes: []Episode{
			Episode{
				Id: 1,
			},
			Episode{
				Id: 2,
			},
		},
	}
	itemToRemove := 2

	RemoveNotifiedEpisode(Episode{Id: itemToRemove})

	if HasBeenNotified(Episode{Id: itemToRemove}) {
		t.Error("Expected item \"", itemToRemove, "\" to be removed from retention but it is still present")
	}
}
