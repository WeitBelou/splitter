package chunks

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

// Processor processes chunks
type Processor interface {
	ProcessChunk(ctx context.Context, chunk Chunk) error
}

// ProcessorFunc is an adapter that helps with using functions as chunk processors
type ProcessorFunc func(ctx context.Context, chunk Chunk) error

// ProcessChunk calls underlining function with chunk
func (f ProcessorFunc) ProcessChunk(ctx context.Context, chunk Chunk) error {
	return f(ctx, chunk)
}

// HTTPSender is a processor that sends chunks via HTTP
type HTTPSender struct {
	// Logger instance
	Log *logrus.Logger

	// Request params
	URI         string
	Method      string
	ContentType string
}

// ProcessChunk implements ProcessChunk method of Processor interface
func (s HTTPSender) ProcessChunk(ctx context.Context, chunk Chunk) error {
	s.Log.Debugf("Chunk: %q", chunk)
	req, err := s.newRequest(ctx, chunk)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	s.Log.Debugf("Request object created: %+v", req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perfrorm request: %w", err)
	}
	defer resp.Body.Close()
	s.Log.Debugf("Response from the server: %+v", resp)

	if resp.StatusCode >= http.StatusBadRequest {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		return fmt.Errorf("service responded with status code: %d, body: %s", resp.StatusCode, body)
	}
	s.Log.Debug("Chunk processing completed")

	return nil
}

func (s HTTPSender) newRequest(ctx context.Context, chunk Chunk) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, s.Method, s.URI, bytes.NewBuffer(chunk))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", s.ContentType)

	return req, nil
}
