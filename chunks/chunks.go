package chunks

import (
	"context"
	"io"
	"time"

	"github.com/sirupsen/logrus"
)

// Chunk is a slice of bytes for processing
type Chunk []byte

// Dispatcher is an interface for reading chunks from a reader and passing them to processor
type Dispatcher interface {
	Dispatch(r Reader, processor Processor) error
}

// ConcurrentDispatcher is a Dispatcher that performs reading and processing in concurrent mode
type ConcurrentDispatcher struct {
	// Logger instance
	Log *logrus.Logger

	// ReadTimeout is a timeout for reading of single chunk
	ReadTimeout time.Duration

	// ProcessTimeout is a timeout for processing of single chunk
	ProcessTimeout time.Duration

	// Concurrency is a maximum count of simultaneously running processors
	Concurrency int
}

// Dispatch implements Dispatch method of Dispatcher interface
func (d ConcurrentDispatcher) Dispatch(r Reader, processor Processor) error {
	// Create a channel with buffer size equal to concurrency
	chunks := make(chan Chunk, d.Concurrency)

	// Root context
	ctx := context.Background()

	// Read chunks to channel
	for {
		chunk, err := readWithTimeout(ctx, r, d.ReadTimeout)
		if err == io.EOF {
			chunks <- chunk
			close(chunks)
			break
		}
		if err != nil {
			d.Log.Errorf("Failed to read chunk: %v", err)
			continue
		}
		chunks <- chunk
	}

	// Process chunks from channel
	for chunk := range chunks {
		err := processWithTimeout(ctx, processor, chunk, d.ProcessTimeout)
		if err != nil {
			d.Log.Errorf("Failed to process chunk: %v", err)
			continue
		}
	}

	return nil
}

func readWithTimeout(ctx context.Context, r Reader, timeout time.Duration) (Chunk, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-r.ReadToChannel(ctx):
		return res.Chunk, res.Err
	}
}

func processWithTimeout(ctx context.Context, processor Processor, chunk Chunk, timeout time.Duration) error {
	panic("not implemented")
}
