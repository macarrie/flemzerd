type Stats = {
	Movies: {
		Tracked     :number,
		Downloading :number,
		Downloaded  :number,
		Removed     :number,
	}
	Shows: {
		Tracked :number,
		Removed :number,
	}
	Episodes: {
		Downloading :number,
		Downloaded  :number,
	}
	Notifications: {
		Read   :number,
		Unread :number,
	}
	Runtime: {
		GoRoutines :number,
		GoMaxProcs :number,
		NumCPU     :number,
	}
};

export default Stats;
