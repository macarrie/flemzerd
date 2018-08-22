package mock

import (
	"errors"
	"fmt"

	. "github.com/macarrie/flemzerd/objects"
)

type MediaCenter struct{}
type ErrorMediaCenter struct{}

var refreshCounter int

func (m MediaCenter) Status() (Module, error) {
	return Module{
		Name: "MediaCenter",
		Type: "mediacenter",
		Status: ModuleStatus{
			Alive:   true,
			Message: "",
		},
	}, nil
}

func (m ErrorMediaCenter) Status() (Module, error) {
	var err error = fmt.Errorf("MediaCenter error")
	return Module{
		Name: "ErrorMediaCenter",
		Type: "mediacenter",
		Status: ModuleStatus{
			Alive:   false,
			Message: err.Error(),
		},
	}, err
}

func (m MediaCenter) GetName() string {
	return "MediaCenter"
}
func (m ErrorMediaCenter) GetName() string {
	return "ErrorMediaCenter"
}

func (m MediaCenter) RefreshLibrary() error {
	refreshCounter += 1
	return nil
}
func (m ErrorMediaCenter) RefreshLibrary() error {
	return errors.New("Error when refreshing library")
}

func (m MediaCenter) GetRefreshCount() int {
	return refreshCounter
}
func (m ErrorMediaCenter) GetRefreshCount() int {
	return refreshCounter
}
