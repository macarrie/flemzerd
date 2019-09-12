package modules_helper

import (
	. "github.com/macarrie/flemzerd/objects"
)

func AtLeastOneAlive(modList []Module) bool {
	for _, mod := range modList {
		if mod.Status.Alive {
			return true
		}
	}

	return false
}
