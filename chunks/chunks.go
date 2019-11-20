package chunks

import (
	"bufio"
	"fmt"
	"io"
	"log"
)

// Chunk is a slice of bytes for processing
type Chunk []byte

// Reader is an interface for reading chunks
type Reader interface {
	Read() (Chunk, error)
}

// LineReader is a Reader that splits reader into chunks of no more than N lines
type LineReader struct {
	n   uint32
	buf *bufio.Reader
}

// NewLineReader creates new LineReader from Reader and chunk size
func NewLineReader(r io.Reader, linesInChunk uint32) LineReader {
	buf := bufio.NewReader(r)

	return LineReader{
		n:   linesInChunk,
		buf: buf,
	}
}

// Read reads chunks from underlining reader
func (r LineReader) Read() (Chunk, error) {
	chunk := make(Chunk, 0)

	for i := uint32(0); i < r.n; i++ {
		line, isPrefix, err := r.buf.ReadLine()
		// Last line
		if err == io.EOF {
			chunk = appendLine(chunk, line)

			return chunk, io.EOF
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read line: %w", err)
		}
		if isPrefix {
			return nil, fmt.Errorf("failed to read full line: line too long")
		}

		chunk = appendLine(chunk, line)
	}

	return chunk, nil
}

func appendLine(chunk Chunk, line []byte) Chunk {
	chunk = append(chunk, line...)
	if len(line) != 0 {
		chunk = append(chunk, '\n')
	}

	return chunk
}

// Processor processes chunks
type Processor interface {
	ProcessChunk(chunk Chunk) error
}

// ProcessorFunc is an adapter that helps with using functions as chunk processors
type ProcessorFunc func(chunk Chunk) error

// ProcessChunk calls underlining function with chunk
func (f ProcessorFunc) ProcessChunk(chunk Chunk) error {
	return f(chunk)
}

// Process reads chunks from chunk reader and process them with given processor
func Process(r Reader, processor Processor) {
	chunks := make(chan Chunk, 8)
	// Send chunks to channel
	go func() {
		for {
			chunk, err := r.Read()
			// Last chunk
			if err == io.EOF {
				chunks <- chunk
				close(chunks)
			}
			if err != nil {
				log.Printf("[ERROR] Failed to read chunk: %v", err)
			}
		}
	}()

	for chunk := range chunks {
		err := processor.ProcessChunk(chunk)
		if err != nil {
			log.Printf("[ERROR] Failed to process chunk: %v", err)
		}
	}
}
