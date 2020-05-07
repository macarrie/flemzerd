package cachable

import log "github.com/macarrie/flemzerd/logging"

type Cachable interface {
	IsCached(value string) bool
	ClearCache()
	GetId() uint
	GetTitle() string
	GetLog() *log.Entry
}
