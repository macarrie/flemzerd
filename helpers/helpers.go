package helpers

import (
	"crypto/tls"
	"fmt"
	log "github.com/macarrie/flemzerd/logging"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"os"
	"time"
)

func DownloadImage(url string, destinationPath string, timeout int) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(timeout) * time.Second,
	}

	var request *http.Request

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return errors.Wrap(err, "Could not build HTTP request object")
	}
	request.Close = true

	response, err := httpClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "Could not perform HTTP request")
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("Could not download image file (status %d)", response.StatusCode)
	}

	file, err := os.Create(destinationPath)
	if err != nil {
		return errors.Wrap(err, "Could not create destination file")
	}

	size, err := io.Copy(file, response.Body)
	if err != nil {
		return errors.Wrap(err, "Could not save downloaded image into destination path")
	}

	log.WithFields(log.Fields{
		"size": size,
		"url":  url,
		"path": destinationPath,
	}).Debug("Image download successful")

	return nil
}
