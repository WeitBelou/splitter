package chunks

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

type HTTPWriter struct {
	log *logrus.Logger

	uri         string
	method      string
	contentType string
}

func NewHTTPWriter(log *logrus.Logger, uri string, method string, contentType string) HTTPWriter {
	return HTTPWriter{
		log:         log,
		uri:         uri,
		method:      method,
		contentType: contentType,
	}
}

func (w HTTPWriter) ProcessChunk(chunk Chunk) error {
	w.log.Debugf("Chunk: %q", chunk)
	req, err := http.NewRequest(w.method, w.uri, bytes.NewReader(chunk))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	w.log.Debugf("Request object created: %+v", req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perfrorm request: %w", err)
	}
	defer resp.Body.Close()
	w.log.Debugf("Response from server: %+v", resp)

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		return fmt.Errorf("service responded with status code: %d, body: %s", resp.StatusCode, body)
	}
	w.log.Debug("Chunk processing completed")

	return nil
}
