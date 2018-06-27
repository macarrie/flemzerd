package kodi_helper

import (
	"fmt"
	"time"

	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"

	"github.com/pdf/kodirpc"
)

func CreateKodiClient(address string, port string) (*kodirpc.Client, error) {
	log.WithFields(log.Fields{
		"address": address,
		"port":    port,
	}).Debug("Trying to create kodi client")

	fullAddress := fmt.Sprintf("%s:%s", address, port)
	return kodirpc.NewClient(fullAddress, &kodirpc.Config{
		ReadTimeout:         time.Duration(HTTP_TIMEOUT * time.Second),
		ConnectTimeout:      time.Duration(HTTP_TIMEOUT * time.Second),
		Reconnect:           false,
		ConnectBackoffScale: 1,
	})
}
