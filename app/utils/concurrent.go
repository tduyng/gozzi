// Package utils provides worker pool and concurrent processing utilities.
package utils

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"time"
)

// WorkerPool represents a pool of workers with Go 1.25.x enhancements.
type WorkerPool struct {
	maxWorkers int
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewWorkerPool creates a new worker pool with container-aware concurrency.
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

// ProcessFiles processes files concurrently using Go 1.25.x patterns.
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

func TimeoutProcessor(
	timeout time.Duration,
	processor func(ctx context.Context, item any) error,
) func(ctx context.Context, item any) error {
	return func(ctx context.Context, item any) error {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		done := make(chan error, 1)
		go func() {
			done <- processor(ctx, item)
		}()

		select {
		case err := <-done:
			return err
		case <-ctx.Done():
			return WrapWithContext(ctx.Err(), ErrServer, ErrorContext{
				Operation: "processor_timeout",
				Component: "concurrent",
				Details:   map[string]any{"timeout": timeout.String()},
			})
		}
	}
}

// BatchFileProcessor processes files in batches concurrently.
type BatchFileProcessor struct {
	batchSize int
	processor func(ctx context.Context, files []string) error
}

// NewBatchFileProcessor creates a new batch file processor.
func NewBatchFileProcessor(
	batchSize int,
	processor func(ctx context.Context, files []string) error,
) *BatchFileProcessor {
	return &BatchFileProcessor{
		batchSize: batchSize,
		processor: processor,
	}
}

// Process processes all files in batches concurrently.
func (bp *BatchFileProcessor) Process(ctx context.Context, files []string) error {
	if len(files) == 0 {
		return nil
	}

	var batches [][]string
	for i := 0; i < len(files); i += bp.batchSize {
		end := i + bp.batchSize
		if end > len(files) {
			end = len(files)
		}
		batches = append(batches, files[i:end])
	}

	wp := NewWorkerPool(ctx)
	defer wp.Close()

	var batchNodes []any
	for _, batch := range batches {
		batchNodes = append(batchNodes, batch)
	}

	return wp.ProcessContentNodes(batchNodes, func(ctx context.Context, node any) error {
		batch, ok := node.([]string)
		if !ok {
			return WrapWithContext(
				errors.New("invalid batch type: expected []string"),
				ErrServer,
				ErrorContext{
					Operation: "validate_batch_type",
					Component: "concurrent",
				})
		}
		return bp.processor(ctx, batch)
	})
}

type EnhancedWaitGroup struct {
	wg  sync.WaitGroup
	ctx context.Context
}

func NewEnhancedWaitGroup(ctx context.Context) *EnhancedWaitGroup {
	return &EnhancedWaitGroup{
		ctx: ctx,
	}
}

func (ewg *EnhancedWaitGroup) Go(f func(ctx context.Context)) {
	ewg.wg.Add(1)
	go func() {
		defer ewg.wg.Done()
		f(ewg.ctx)
	}()
}

func (ewg *EnhancedWaitGroup) Wait() {
	ewg.wg.Wait()
}

func (ewg *EnhancedWaitGroup) WaitWithTimeout(timeout time.Duration) error {
	done := make(chan struct{})
	go func() {
		defer close(done)
		ewg.wg.Wait()
	}()

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return WrapWithContext(
			errors.New("wait group timeout"),
			ErrServer,
			ErrorContext{
				Operation: "wait_group_timeout",
				Component: "concurrent",
				Details:   map[string]any{"timeout": timeout.String()},
			})
	case <-ewg.ctx.Done():
		return ewg.ctx.Err()
	}
}
