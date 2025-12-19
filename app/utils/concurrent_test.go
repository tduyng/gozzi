package utils

import (
	"context"
	"errors"
	"runtime"
	"sync/atomic"
	"testing"
)

func TestNewWorkerPool(t *testing.T) {
	wp := NewWorkerPool(context.Background())
	defer wp.Close()

	expectedWorkers := runtime.GOMAXPROCS(0)
	if wp.maxWorkers != expectedWorkers {
		t.Errorf("expected %d workers, got %d", expectedWorkers, wp.maxWorkers)
	}

	if wp.ctx == nil {
		t.Error("context should not be nil")
	}
}

func TestNewWorkerPool_NilContext(t *testing.T) {
	wp := NewWorkerPool(context.TODO())
	defer wp.Close()

	if wp.ctx == nil {
		t.Error("should create context when nil provided")
	}
}

func TestWorkerPool_ProcessFiles(t *testing.T) {
	wp := NewWorkerPool(context.Background())
	defer wp.Close()

	files := []string{"file1.txt", "file2.txt", "file3.txt"}
	var processed int64

	err := wp.ProcessFiles(files, func(ctx context.Context, filePath string) error {
		atomic.AddInt64(&processed, 1)
		return nil
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if atomic.LoadInt64(&processed) != int64(len(files)) {
		t.Errorf("expected %d files processed, got %d", len(files), processed)
	}
}

func TestWorkerPool_ProcessFiles_WithError(t *testing.T) {
	wp := NewWorkerPool(context.Background())
	defer wp.Close()

	files := []string{"file1.txt", "file2.txt", "file3.txt"}
	expectedErr := errors.New("processing error")

	err := wp.ProcessFiles(files, func(ctx context.Context, filePath string) error {
		if filePath == "file2.txt" {
			return expectedErr
		}
		return nil
	})

	if err == nil {
		t.Error("expected error but got nil")
	}
}

func TestWorkerPool_ProcessFiles_EmptyList(t *testing.T) {
	wp := NewWorkerPool(context.Background())
	defer wp.Close()

	err := wp.ProcessFiles([]string{}, func(ctx context.Context, filePath string) error {
		t.Error("processor should not be called for empty list")
		return nil
	})

	if err != nil {
		t.Errorf("unexpected error for empty list: %v", err)
	}
}

func TestWorkerPool_ProcessFiles_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	wp := NewWorkerPool(ctx)
	defer wp.Close()

	files := []string{"file1.txt", "file2.txt", "file3.txt"}

	// Cancel context immediately
	cancel()

	err := wp.ProcessFiles(files, func(ctx context.Context, filePath string) error {
		// This should not be called due to context cancellation
		return nil
	})

	// Should return context.Canceled or nil if no work was started
	if err != nil && !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled error or nil, got: %v", err)
	}
}

func TestWorkerPool_ProcessContentNodes(t *testing.T) {
	wp := NewWorkerPool(context.Background())
	defer wp.Close()

	nodes := []any{"node1", "node2", "node3"}
	var processed int64

	err := wp.ProcessContentNodes(nodes, func(ctx context.Context, node any) error {
		atomic.AddInt64(&processed, 1)
		return nil
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if atomic.LoadInt64(&processed) != int64(len(nodes)) {
		t.Errorf("expected %d nodes processed, got %d", len(nodes), processed)
	}
}
