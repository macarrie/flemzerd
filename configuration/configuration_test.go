package configuration

import (
	"github.com/hashicorp/go-multierror"
	"testing"

	log "github.com/macarrie/flemzerd/logging"
)

func init() {
	ConfigurationFatal = func(code int) {
		log.Error("Configuration fatal error: ", code)
	}
}

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

	testMatrix := []struct {
		ExpectedErrors int
		FilePath       string
	}{
		{
			6,
			"../testdata/test_config.toml",
		},
		{
			18,
			"../testdata/test_config_bad.toml",
		},
		{
			15,
			"../testdata/test_config_mediacenters.toml",
		},
		{
			7,
			"../testdata/test_config_no_notifiers.toml",
		},
	}

	for _, test := range testMatrix {
		t.Run(test.FilePath, func(t *testing.T) {
			UseFile(test.FilePath)
			loadErr := Load()
			if loadErr != nil {
				t.Errorf("Could not load config file: %s", loadErr)
			}
			if err := Check(); err != nil {
				multierr, _ := err.(*multierror.Error)
				if len(multierr.Errors) != test.ExpectedErrors {
					t.Errorf("Expected %d configuration errors, got %d instead", test.ExpectedErrors, len(multierr.Errors))
				}
			}
		})
	}
}

func TestLoad(t *testing.T) {
	UseFile("../testdata/test_config.toml")
	err := Load()
	if err != nil {
		t.Errorf("Error while loading configuration file: %s", err.Error())
	}

	if len(Config.Providers) != 3 {
		t.Error("Providers not correctly loaded")
	}
	if len(Config.Notifiers) != 5 {
		t.Error("Notifiers not correctly loaded")
	}
	if len(Config.Indexers["torznab"]) != 2 {
		t.Error("Indexers not correctly loaded")
	}
	if len(Config.Downloaders) != 2 {
		t.Error("Downloaders not correctly loaded")
	}
	if len(Config.Watchlists) != 3 {
		t.Error("Watchlists not correctly loaded")
	}
}
