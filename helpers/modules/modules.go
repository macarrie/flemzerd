package modules_helper

func AtLeastOneAlive(modList []Module) bool {
	for _, mod := range modList {
		if mod.Status.Alive {
			return true
		}
	}

	return false
}
