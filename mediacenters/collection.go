package mediacenter

import (
	"bytes"
	"errors"

	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"
)

var mediaCenterCollection []MediaCenter

func Status() ([]Module, error) {
	var modList []Module
	var aggregatedErrorMessage bytes.Buffer

	for _, mc := range mediaCenterCollection {
		mod, mcAliveError := mc.Status()
		if mcAliveError != nil {
			log.WithFields(log.Fields{
				"error": mcAliveError,
			}).Warning("MediaCenter is not alive")
			aggregatedErrorMessage.WriteString(mcAliveError.Error())
			aggregatedErrorMessage.WriteString("\n")
		}
		modList = append(modList, mod)
	}

	var retError error
	if aggregatedErrorMessage.Len() == 0 {
		retError = nil
	} else {
		retError = errors.New(aggregatedErrorMessage.String())
	}
	return modList, retError
}

func Reset() {
	mediaCenterCollection = []MediaCenter{}
}

func AddMediaCenter(mc MediaCenter) {
	mediaCenterCollection = append(mediaCenterCollection, mc)
	log.WithFields(log.Fields{
		"mediacenter": mc.GetName(),
	}).Debug("Mediacenter loaded")
}

func RefreshLibrary() {
	for _, mc := range mediaCenterCollection {
		err := mc.RefreshLibrary()
		if err != nil {
			log.WithFields(log.Fields{
				"mediacenter": mc.GetName(),
				"error":       err,
			}).Warning("Failed to refresh media center library")
		}
	}
}
