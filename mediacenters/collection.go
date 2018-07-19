package mediacenter

import (
	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"

	multierror "github.com/hashicorp/go-multierror"
)

var mediaCenterCollection []MediaCenter

func Status() ([]Module, error) {
	var modList []Module
	var errorList *multierror.Error

	for _, mc := range mediaCenterCollection {
		mod, mcAliveError := mc.Status()
		if mcAliveError != nil {
			log.WithFields(log.Fields{
				"error": mcAliveError,
			}).Warning("MediaCenter is not alive")
			errorList = multierror.Append(errorList, mcAliveError)
		}
		modList = append(modList, mod)
	}

	return modList, errorList.ErrorOrNil()
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
