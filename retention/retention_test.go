package retention

import (
	"testing"
	"time"

	. "github.com/macarrie/flemzerd/objects"
)

func init() {
	resetRetention()
}

func resetRetention() {
	retentionData.NotifiedEpisodes = make(map[int]Episode)
	retentionData.DownloadedEpisodes = make(map[int]Episode)
	retentionData.DownloadingEpisodes = make(map[int]DownloadingEpisode)
}

func TestLoad(t *testing.T) {
	// TODO
}

func TestSave(t *testing.T) {
	// TODO
}

func TestElementInRetention(t *testing.T) {
	resetRetention()
	testEpisode := Episode{
		Id: 1000,
	}

	if HasBeenNotified(testEpisode) {
		t.Error("Expected episode not to be in notified episodes retention")
	}
	if HasBeenDownloaded(testEpisode) {
		t.Error("Expected episode not to be in downloaded episodes retention")
	}
	if IsDownloading(testEpisode) {
		t.Error("Expected episode not to be in downloading episodes retention")
	}

	retentionData = RetentionData{
		NotifiedEpisodes: map[int]Episode{
			testEpisode.Id: testEpisode,
		},
		DownloadedEpisodes: map[int]Episode{
			testEpisode.Id: testEpisode,
		},
		DownloadingEpisodes: map[int]DownloadingEpisode{
			testEpisode.Id: DownloadingEpisode{
				Episode: testEpisode,
			},
		},
	}

	if !HasBeenNotified(testEpisode) {
		t.Error("Expected test episode to be present in notified episodes retention")
	}
	if !HasBeenDownloaded(testEpisode) {
		t.Error("Expected test episode to be present in downloaded episodes retention")
	}
	if !IsDownloading(testEpisode) {
		t.Error("Expected test episode to be present in downloading episodes retention")
	}
}

func TestAddElementInRetention(t *testing.T) {
	resetRetention()
	testEpisode := Episode{
		Id: 1000,
	}

	AddNotifiedEpisode(testEpisode)
	AddDownloadedEpisode(testEpisode)
	AddDownloadingEpisode(testEpisode)

	if len(retentionData.NotifiedEpisodes) != 1 {
		t.Error("Expected to have 1 item in notified episodes retention, got ", len(retentionData.NotifiedEpisodes), " items instead")
	}
	if len(retentionData.DownloadedEpisodes) != 1 {
		t.Error("Expected to have 1 item in downloaded episodes retention, got ", len(retentionData.DownloadedEpisodes), " items instead")
	}
	if len(retentionData.DownloadingEpisodes) != 1 {
		t.Error("Expected to have 1 item in downloading episodes retention, got ", len(retentionData.DownloadingEpisodes), " items instead")
	}
}

func TestCleanOldNotifiedEpisodes(t *testing.T) {
	resetRetention()
	testEpisode := Episode{
		Id:   1000,
		Date: time.Time{},
	}

	AddNotifiedEpisode(testEpisode)

	CleanOldNotifiedEpisodes()

	if len(retentionData.NotifiedEpisodes) > 0 {
		t.Error("Expected old episode to have been removed from notified episodes retention but episode is still present in retention")
	}
}

func TestRemoveNotifiedEpisode(t *testing.T) {
	e1 := Episode{
		Id: 1,
	}
	e2 := Episode{
		Id: 2,
	}

	AddNotifiedEpisode(e1)
	AddNotifiedEpisode(e2)

	itemToRemove := 2

	RemoveNotifiedEpisode(e2)

	if HasBeenNotified(e2) {
		t.Error("Expected item \"", itemToRemove, "\" to be removed from retention but it is still present")
	}
}
