package chunks

import (
	"context"
	"time"
)

// chunkWithErr is a helper struct for sending chunk to channel
type chunkWithErr struct {
	Chunk Chunk
	Err   error
}

// readAsync reads chunk from reader async and returns channel
func readAsync(ctx context.Context, r Reader) chan chunkWithErr {
	ch := make(chan chunkWithErr)

	go func() {
		chunk, err := r.Read(ctx)
		ch <- chunkWithErr{
			Chunk: chunk,
			Err:   err,
		}
	}()

	return ch
}

// processAsync process chunk in background and returns error to channel
func processAsync(ctx context.Context, processor Processor, chunk Chunk) chan error {
	ch := make(chan error)

	go func() {
		ch <- processor.ProcessChunk(ctx, chunk)
	}()

	return ch
}

// readWithTimeout reads from reader with timeout
func readWithTimeout(ctx context.Context, r Reader, timeout time.Duration) (Chunk, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-readAsync(ctx, r):
		return res.Chunk, res.Err
	}
}

// readWithTimeout process chunk with timeout
func processWithTimeout(ctx context.Context, processor Processor, chunk Chunk, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-processAsync(ctx, processor, chunk):
		return err
	}
}
