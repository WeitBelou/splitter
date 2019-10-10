package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/weitbelou/splitter/chunks"
)

var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save chunks to local files",
	Run: func(cmd *cobra.Command, args []string) {
		// Open file for reader
		f, err := getInputFile(inputFile)
		if err != nil {
			log.Fatalf("failed to get input file %s: %v", inputFile, err)
		}
		defer f.Close()

		// Create directory for chunks
		err = os.MkdirAll(outputDir, 0700)
		if err != nil {
			log.Fatalf("failed to create dir for result: %v", err)
		}

		// Create chunks reader
		r := chunks.NewLineReader(f, chunkSize)

		// Create chunks processor
		processor := newIndexedChunksWriter(outputDir, inputFile)

		// Process chunks
		err = chunks.Process(r, processor)
		if err != nil {
			log.Fatalf("failed to process chunks: %v", err)
		}
	},
}

var outputDir string

func init() {
	saveCmd.PersistentFlags().StringVarP(&outputDir, "output-dir", "o",
		"out",
		"Path where chunks will be saved (default is 'out'",
	)
}

func getInputFile(name string) (*os.File, error) {
	if name == "-" {
		return os.Stdin, nil
	}
	return os.Open(name)
}

type indexedChunksWriter struct {
	dir      string
	baseName string
	ext      string

	index uint32
}

func newIndexedChunksWriter(outDir string, input string) *indexedChunksWriter {
	if input == "-" {
		input = "chunk.txt"
	}

	ext := filepath.Ext(input)
	baseName := strings.TrimSuffix(filepath.Base(input), ext)

	return &indexedChunksWriter{
		dir:      outDir,
		baseName: baseName,
		ext:      ext,
	}
}

func (w *indexedChunksWriter) ProcessChunk(chunk chunks.Chunk) error {
	// Create file for chunk
	out, err := os.Create(w.getOutputPath())
	if err != nil {
		return fmt.Errorf("failed to create file for result: %w", err)
	}
	defer out.Close()

	// Write chunk
	_, err = out.Write(chunk)
	if err != nil {
		return fmt.Errorf("failed to write chunk to file: %w", err)
	}

	// Increment index
	w.index++

	return nil
}

func (w indexedChunksWriter) getOutputPath() string {
	return filepath.Join(w.dir, fmt.Sprintf("%s-%d.%s", w.baseName, w.index, w.ext))
}
