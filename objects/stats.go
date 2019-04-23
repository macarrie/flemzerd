package objects

type StatsFields struct {
	Movies struct {
		Tracked     int
		Downloading int
		Downloaded  int
		Removed     int
	}
	Shows struct {
		Tracked int
		Removed int
	}
	Episodes struct {
		Downloading int
		Downloaded  int
	}
	Notifications struct {
		Read   int
		Unread int
	}
	Runtime struct {
		GoRoutines int
		GoMaxProcs int
		NumCPU     int
	}
}

