package objects

type Module struct {
	Name   string
	Type   string
	Status ModuleStatus
}

type ModuleStatus struct {
	Alive   bool
	Message string
}
