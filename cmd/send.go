package cmd

import (
	"github.com/spf13/cobra"

	"github.com/weitbelou/splitter/chunks"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send chunks to remote service",
	Run:   runSend,
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

func runSend(_ *cobra.Command, _ []string) {
	// Open file for reader
	f, err := getInputFile(inputFile)
	if err != nil {
		log.Fatalf("Failed to get input file %s: %v", inputFile, err)
	}
	defer f.Close()
	log.Debugf("Using input file: %s", f.Name())

	// Create chunks reader
	r := chunks.NewLineReader(f, chunkSize)
	log.Debugf("Line reader created: %+v", r)

	// Create chunks processor
	processor := chunks.HTTPSender{
		Log:         log,
		URI:         sendURI,
		Method:      sendMethod,
		ContentType: sendContentType,
	}
	log.Debugf("Chunk processor created: %+v", processor)

	// Process chunks
	err = chunks.Process(r, processor)
	if err != nil {
		log.Errorf("Failed to process chunks: %v", err)
	}
	log.Debugf("Processing completed: %+v", processor)
}
