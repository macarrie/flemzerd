package mediacenter

import (
	"errors"
	"fmt"

	. "github.com/macarrie/flemzerd/objects"
)

type MockMediaCenter struct{}
type MockErrorMediaCenter struct{}

func (m MockMediaCenter) Status() (Module, error) {
	return Module{
		Name: "MockMediaCenter",
		Type: "mediacenter",
		Status: ModuleStatus{
			Alive:   true,
			Message: "",
		},
	}, nil
}

func (m MockErrorMediaCenter) Status() (Module, error) {
	var err error = fmt.Errorf("MediaCenter error")
	return Module{
		Name: "MockErrorMediaCenter",
		Type: "mediacenter",
		Status: ModuleStatus{
			Alive:   false,
			Message: err.Error(),
		},
	}, err
}

func (m MockMediaCenter) GetName() string {
	return "MockMediaCenter"
}
func (m MockErrorMediaCenter) GetName() string {
	return "MockErrorMediaCenter"
}

func (m MockMediaCenter) RefreshLibrary() error {
	return nil
}
func (m MockErrorMediaCenter) RefreshLibrary() error {
	return errors.New("Error when refreshing library")
}
