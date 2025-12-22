// Package utils provides worker pool and concurrent processing utilities.
package utils

import (
	"context"
	"runtime"
	"sync"
)

// WorkerPool represents a pool of workers with Go 1.25.x enhancements.
type WorkerPool struct {
	maxWorkers int
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewWorkerPool creates a new worker pool.
func NewWorkerPool(ctx context.Context) *WorkerPool {
	if ctx == nil {
		ctx = context.Background()
	}

	ctxWithCancel, cancel := context.WithCancel(ctx)

	maxWorkers := runtime.GOMAXPROCS(0)

	return &WorkerPool{
		maxWorkers: maxWorkers,
		ctx:        ctxWithCancel,
		cancel:     cancel,
	}
}

// ProcessFiles processes files concurrently.
func (wp *WorkerPool) ProcessFiles(files []string, processor func(ctx context.Context, filePath string) error) error {
	if len(files) == 0 {
		return nil
	}

	// Create a buffered channel to limit concurrent workers
	semaphore := make(chan struct{}, wp.maxWorkers)
	errChan := make(chan error, len(files))

	for _, file := range files {
		select {
		case semaphore <- struct{}{}:
		case <-wp.ctx.Done():
			return wp.ctx.Err()
		}

		// Start worker goroutine
		wp.wg.Add(1)
		go func(filePath string) {
			defer wp.wg.Done()
			defer func() { <-semaphore }()

			if err := processor(wp.ctx, filePath); err != nil {
				select {
				case errChan <- err:
				case <-wp.ctx.Done():
				}
			}
		}(file)
	}

	go func() {
		wp.wg.Wait()
		close(errChan)
	}()

	// Collect errors
	var firstErr error
	for err := range errChan {
		if firstErr == nil {
			firstErr = err
			wp.cancel() // Cancel remaining work on first error
		}
	}

	return firstErr
}

// ProcessContentNodes processes content nodes concurrently.
func (wp *WorkerPool) ProcessContentNodes(nodes []any, processor func(ctx context.Context, node any) error) error {
	if len(nodes) == 0 {
		return nil
	}

	// Create a buffered channel to limit concurrent workers
	semaphore := make(chan struct{}, wp.maxWorkers)
	errChan := make(chan error, len(nodes))

	for _, node := range nodes {
		select {
		case semaphore <- struct{}{}:
		case <-wp.ctx.Done():
			return wp.ctx.Err()
		}

		// Start worker goroutine
		wp.wg.Add(1)
		go func(item any) {
			defer wp.wg.Done()
			defer func() { <-semaphore }()

			if err := processor(wp.ctx, item); err != nil {
				select {
				case errChan <- err:
				case <-wp.ctx.Done():
				}
			}
		}(node)
	}

	go func() {
		wp.wg.Wait()
		close(errChan)
	}()

	// Collect errors
	var firstErr error
	for err := range errChan {
		if firstErr == nil {
			firstErr = err
			wp.cancel() // Cancel remaining work on first error
		}
	}

	return firstErr
}

func (wp *WorkerPool) Close() {
	wp.cancel()
	wp.wg.Wait()
}
