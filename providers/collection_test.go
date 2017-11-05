package provider

import (
	"testing"

	. "github.com/macarrie/flemzerd/objects"
)

func TestAddProvider(t *testing.T) {
	providersLength := len(providersCollection)
	p := MockProvider{}
	AddProvider(p)

	if len(providersCollection) != providersLength+1 {
		t.Error("Expected ", providersLength+1, " providers, got ", len(providersCollection))
	}
}

func TestFindShow(t *testing.T) {
	p := MockProvider{}
	providersCollection = []Provider{p}

	show, err := FindShow("Test show")
	if err != nil {
		t.Error("Got error during FindShow: ", err)
	}
	if show.Id != 1000 {
		t.Errorf("Expected show with id 1000, got id %v instead\n", show.Id)
	}
}

func TestFindRecentlyAiredEpisodesForShow(t *testing.T) {
	p := MockProvider{}
	providersCollection = []Provider{p}

	episodeList, err := FindRecentlyAiredEpisodesForShow(TvShow{
		Name: "Test show",
	})
	if err != nil {
		t.Error("Got error during FindRecentlyAiredEpisodesForShow: ", err)
	}
	if episodeList[0].Id != 1000 {
		t.Errorf("Expected episode with id 1000, got id %v instead\n", episodeList[0].Id)
	}
}
