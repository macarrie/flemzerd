package configuration

import (
	"testing"
)

func TestUseFile(t *testing.T) {
	UseFile("test")

	if !customConfigFile {
		if customConfigFilePath != "test" {
			t.Errorf("Expected custom config file path to be 'test', got '%s' instead", customConfigFilePath)
		} else {
			t.Error("Not using custom file")
		}
	}
}

func TestCheck(t *testing.T) {
	customConfigFile = false
	customConfigFilePath = ""
	Load()
	Check()

	UseFile("../testdata/test_config.yaml")
	Load()
	Check()

	UseFile("../testdata/test_config_bad.yaml")
	Load()
	Check()
}

func TestLoad(t *testing.T) {
	UseFile("../testdata/test_config.yaml")
	err := Load()
	if err != nil {
		t.Errorf("Error while loading configuration file: %s", err.Error())
	}

	if len(Config.Providers) != 2 {
		t.Error("Providers not correctly loaded")
	}
	if len(Config.Notifiers) != 2 {
		t.Error("Notifiers not correctly loaded")
	}
	if len(Config.Indexers["torznab"]) != 2 {
		t.Error("Providers not correctly loaded")
	}
	if len(Config.Downloaders) != 1 {
		t.Error("Downloaders not correctly loaded")
	}
	if len(Config.Watchlists) != 2 {
		t.Error("Watchlists not correctly loaded")
	}
}
