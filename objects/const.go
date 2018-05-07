package objects

const RECENTLY_AIRED_EPISODES_INTERVAL = 7
const HTTP_TIMEOUT = 10

const (
	TORRENT_DOWNLOADING = iota
	TORRENT_SEEDING
	TORRENT_STOPPED
	TORRENT_DOWNLOAD_PENDING
	TORRENT_UNKNOWN_STATUS
)

const (
	DB_PATH          = "/var/lib/flemzerd/flemzer.db"
	DOWNLOAD_TMP_DIR = "/var/lib/flemzerd/tmp"
)
