package objects

const CACHE_BASE_PATH = "/var/lib/flemzerd/cache"

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
	NOTIFICATION_NEW_EPISODE = iota
	NOTIFICATION_NEW_MOVIE
	NOTIFICATION_DOWNLOAD_START
	NOTIFICATION_DOWNLOAD_SUCCESS
	NOTIFICATION_DOWNLOAD_FAILURE
	NOTIFICATION_TEXT
	NOTIFICATION_NO_TORRENTS
)

const (
	MOVIE = iota
	EPISODE
)

const (
	DB_PATH = "/var/lib/flemzerd/db/flemzer.db"
)

const (
	TRAKT_CLIENT_ID = "4d5d84f3a95adc528d08b1e2f6df6c8197548d128e72e68b4c33fc0dfa634926"
	FANART_TV_KEY   = "648ff4214a1eea4416ad51417fc8a4e4"
)
