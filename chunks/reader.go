package chunks

import (
	"bufio"
	"context"
	"fmt"
	"io"
)

// Reader is an interface for reading chunks
type Reader interface {
	Read(ctx context.Context) (Chunk, error)
	ReadToChannel(ctx context.Context) chan ChunkWithErr
}

// ChunkWithErr is a helper struct for sending chunk to channel
type ChunkWithErr struct {
	Chunk Chunk
	Err   error
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
func (r LineReader) Read(ctx context.Context) (Chunk, error) {
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

// ReadToChannel implements ReadToChannel from Reader interface
func (r LineReader) ReadToChannel(ctx context.Context) chan ChunkWithErr {
	ch := make(chan ChunkWithErr, 1)

	go func() {
		chunk, err := r.Read(ctx)
		ch <- ChunkWithErr{
			Chunk: chunk,
			Err:   err,
		}
	}()

	return ch
}
