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
	retentionData.DownloadingEpisodes = make(map[int]*DownloadingEpisode)
	retentionData.FailedEpisodes = make(map[int]Episode)
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
	if IsInDownloadProcess(testEpisode) {
		t.Error("Expected episode not to be in downloading episodes retention")
	}
	if HasDownloadFailed(testEpisode) {
		t.Error("Expected episode not to be in failed torrents retention")
	}

	retentionData = RetentionData{
		NotifiedEpisodes: map[int]Episode{
			testEpisode.Id: testEpisode,
		},
		DownloadedEpisodes: map[int]Episode{
			testEpisode.Id: testEpisode,
		},
		DownloadingEpisodes: map[int]*DownloadingEpisode{
			testEpisode.Id: &DownloadingEpisode{
				Episode: testEpisode,
			},
		},
		FailedEpisodes: map[int]Episode{
			testEpisode.Id: testEpisode,
		},
	}

	if !HasBeenNotified(testEpisode) {
		t.Error("Expected test episode to be present in notified episodes retention")
	}
	if !HasBeenDownloaded(testEpisode) {
		t.Error("Expected test episode to be present in downloaded episodes retention")
	}
	if !IsInDownloadProcess(testEpisode) {
		t.Error("Expected test episode to be present in downloading episodes retention")
	}
	if !HasDownloadFailed(testEpisode) {
		t.Error("Expected test episode to be present in failed torrents retention")
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

func TestRemoveElementFromRetention(t *testing.T) {
	e1 := Episode{
		Id: 1,
	}
	e2 := Episode{
		Id: 2,
	}

	AddNotifiedEpisode(e1)
	AddNotifiedEpisode(e2)
	AddDownloadedEpisode(e1)
	AddDownloadedEpisode(e2)
	AddDownloadingEpisode(e1)
	AddDownloadingEpisode(e2)
	AddFailedEpisode(e1)
	AddFailedEpisode(e2)

	itemToRemove := 2

	RemoveNotifiedEpisode(e2)
	RemoveDownloadingEpisode(e2)
	RemoveDownloadedEpisode(e2)
	RemoveFailedEpisode(e2)

	if HasBeenNotified(e2) {
		t.Error("Expected item \"", itemToRemove, "\" to be removed from notified episodes retention but it is still present")
	}
	if HasBeenDownloaded(e2) {
		t.Error("Expected item \"", itemToRemove, "\" to be removed from downloaded episodes retention but it is still present")
	}
	if IsInDownloadProcess(e2) {
		t.Error("Expected item \"", itemToRemove, "\" to be removed from downloading episodes retention but it is still present")
	}
	if HasDownloadFailed(e2) {
		t.Error("Expected item \"", itemToRemove, "\" to be removed from failed episodes retention but it is still present")
	}
}

func TestTorrentsHandling(t *testing.T) {
	resetRetention()
	testEpisode := Episode{
		Id: 1000,
	}
	testTorrent := Torrent{
		Id: "id",
	}

	retentionData = RetentionData{
		DownloadingEpisodes: map[int]*DownloadingEpisode{
			testEpisode.Id: &DownloadingEpisode{
				Episode:     testEpisode,
				Downloading: true,
			},
		},
	}

	if !IsDownloading(testEpisode) {
		t.Error("Episode is supposed to be downloading")
	}

	ChangeDownloadingState(testEpisode, false)

	if IsDownloading(testEpisode) {
		t.Error("Episode should not be downloading")
	}

	AddFailedTorrent(testEpisode, testTorrent)

	if !IsInFailedTorrents(testEpisode, testTorrent) {
		t.Error("Expected torrent to be in failed torrents")
	}

	failedTorrentsCount := GetFailedTorrentsCount(testEpisode)
	if failedTorrentsCount != 1 {
		t.Error("Expected failed torrents count to be 1, got ", failedTorrentsCount, " instead")
	}
}
