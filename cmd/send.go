package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/weitbelou/splitter/chunks"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send chunks to remote service",
	Run: func(cmd *cobra.Command, args []string) {
		// Open file for reader
		f, err := getInputFile(inputFile)
		if err != nil {
			log.Fatalf("failed to get input file %s: %v", inputFile, err)
		}
		defer f.Close()

		// Create chunks reader
		r := chunks.NewLineReader(f, chunkSize)

		// Create chunks processor
		processor := newHTTPPostChunksWriter(sendURI)

		// Process chunks
		err = chunks.Process(r, processor)
		if err != nil {
			log.Fatalf("failed to process chunks: %v", err)
		}
	},
}

var (
	sendURI         string
	sendMethod      string
	sendContentType string
)

func init() {
	sendCmd.PersistentFlags().StringVarP(
		&sendURI,
		"uri", "u",
		"http://devnull-as-a-service.com/dev/null",
		"Where to POST chunks",
	)
	sendCmd.PersistentFlags().StringVarP(
		&sendMethod,
		"method", "m",
		"POST",
		"HTTP method to be used for request",
	)

	sendCmd.PersistentFlags().StringVarP(
		&sendContentType,
		"content-type", "c",
		"application/json",
		"Content type of request body",
	)
}

type httpPostChunksWriter struct {
	uri         string
	method      string
	contentType string
}

func newHTTPPostChunksWriter(uri string) httpPostChunksWriter {
	return httpPostChunksWriter{
		uri: uri,
	}
}

func (w httpPostChunksWriter) ProcessChunk(chunk chunks.Chunk) error {
	req, err := http.NewRequest(w.method, w.uri, bytes.NewReader(chunk))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perfrorm request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		return fmt.Errorf("service responded with status code: %d, body: %s", resp.StatusCode, body)
	}

	return nil
}
