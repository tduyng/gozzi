package shared

import (
	"context"
	"errors"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
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

func TestTimeoutProcessor(t *testing.T) {
	timeout := 100 * time.Millisecond

	// Test successful processing within timeout
	processor := TimeoutProcessor(timeout, func(ctx context.Context, item any) error {
		// Fast processing
		time.Sleep(10 * time.Millisecond)
		return nil
	})

	err := processor(context.Background(), "test")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestTimeoutProcessor_Timeout(t *testing.T) {
	timeout := 50 * time.Millisecond

	// Test timeout behavior
	processor := TimeoutProcessor(timeout, func(ctx context.Context, item any) error {
		// Slow processing that should timeout
		time.Sleep(200 * time.Millisecond)
		return nil
	})

	err := processor(context.Background(), "test")
	if err == nil {
		t.Error("expected timeout error but got nil")
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected context.DeadlineExceeded in error chain, got: %v", err)
	}
}

func TestBatchFileProcessor(t *testing.T) {
	files := []string{"f1", "f2", "f3", "f4", "f5"}
	var batchesProcessed int64
	var totalFilesProcessed int64

	processor := NewBatchFileProcessor(2, func(ctx context.Context, batch []string) error {
		atomic.AddInt64(&batchesProcessed, 1)
		atomic.AddInt64(&totalFilesProcessed, int64(len(batch)))
		return nil
	})

	err := processor.Process(context.Background(), files)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Should create 3 batches: [f1,f2], [f3,f4], [f5]
	expectedBatches := int64(3)
	if atomic.LoadInt64(&batchesProcessed) != expectedBatches {
		t.Errorf("expected %d batches, got %d", expectedBatches, batchesProcessed)
	}

	if atomic.LoadInt64(&totalFilesProcessed) != int64(len(files)) {
		t.Errorf("expected %d total files processed, got %d", len(files), totalFilesProcessed)
	}
}

func TestBatchFileProcessor_EmptyFiles(t *testing.T) {
	processor := NewBatchFileProcessor(2, func(ctx context.Context, batch []string) error {
		t.Error("processor should not be called for empty files")
		return nil
	})

	err := processor.Process(context.Background(), []string{})
	if err != nil {
		t.Errorf("unexpected error for empty files: %v", err)
	}
}

func TestEnhancedWaitGroup(t *testing.T) {
	ewg := NewEnhancedWaitGroup(context.Background())
	var executed int64

	// Start multiple goroutines
	for i := 0; i < 3; i++ {
		ewg.Go(func(ctx context.Context) {
			atomic.AddInt64(&executed, 1)
		})
	}

	ewg.Wait()

	if atomic.LoadInt64(&executed) != 3 {
		t.Errorf("expected 3 goroutines executed, got %d", executed)
	}
}

func TestEnhancedWaitGroup_WaitWithTimeout(t *testing.T) {
	ewg := NewEnhancedWaitGroup(context.Background())

	// Start a quick goroutine
	ewg.Go(func(ctx context.Context) {
		time.Sleep(10 * time.Millisecond)
	})

	err := ewg.WaitWithTimeout(100 * time.Millisecond)
	if err != nil {
		t.Errorf("unexpected timeout error: %v", err)
	}
}

func TestEnhancedWaitGroup_WaitWithTimeout_Timeout(t *testing.T) {
	ewg := NewEnhancedWaitGroup(context.Background())

	// Start a slow goroutine
	ewg.Go(func(ctx context.Context) {
		time.Sleep(200 * time.Millisecond)
	})

	err := ewg.WaitWithTimeout(50 * time.Millisecond)
	if err == nil {
		t.Error("expected timeout error but got nil")
	}
}

func TestEnhancedWaitGroup_WaitWithTimeout_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ewg := NewEnhancedWaitGroup(ctx)

	// Start a goroutine
	ewg.Go(func(ctx context.Context) {
		time.Sleep(200 * time.Millisecond)
	})

	// Cancel context immediately
	cancel()

	err := ewg.WaitWithTimeout(1 * time.Second)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled error, got: %v", err)
	}
}
