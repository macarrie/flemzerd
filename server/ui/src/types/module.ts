type ModuleStatus = {
	Alive   :boolean,
	Message :string,
};

type Module = {
	Name   :string,
	Type   :string,
    Status :ModuleStatus,
};

export default Module;
