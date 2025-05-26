package utils

import (
	"sync"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	t.Run("basic job execution", func(t *testing.T) {
		pool := NewWorkerPool(2, 10)
		pool.Start()
		defer pool.Stop()

		result := 0
		var mu sync.Mutex

		// Create a simple job that increments a counter
		job, err := NewJob(
			func(val int) error {
				mu.Lock()
				result += val
				mu.Unlock()
				return nil
			},
			[]any{1},
			nil,
			nil,
		)

		if err != nil {
			t.Fatalf("Failed to create job: %v", err)
		}

		if err := pool.Enqueue(job); err != nil {
			t.Fatalf("Failed to enqueue job: %v", err)
		}

		// Give some time for the job to execute
		time.Sleep(100 * time.Millisecond)

		if result != 1 {
			t.Errorf("Expected result to be 1, got %d", result)
		}
	})

	t.Run("job with post execution", func(t *testing.T) {
		pool := NewWorkerPool(2, 10)
		pool.Start()
		defer pool.Stop()

		var result int
		var mu sync.Mutex

		// Create a job with a post-execution function
		job, err := NewJob(
			func(x int) (int, error) {
				return x * 2, nil
			},
			[]any{5},
			func(x int, jobReturns ...any) {
				if len(jobReturns) > 0 {
					if val, ok := jobReturns[0].(int); ok {
						mu.Lock()
						result = val
						mu.Unlock()
					}
				}
			},
			[]any{5},
		)

		if err != nil {
			t.Fatalf("Failed to create job: %v", err)
		}

		if err := pool.Enqueue(job); err != nil {
			t.Fatalf("Failed to enqueue job: %v", err)
		}

		// Give some time for the job to execute
		time.Sleep(100 * time.Millisecond)

		if result != 10 {
			t.Errorf("Expected result to be 10, got %d", result)
		}
	})

	t.Run("multiple concurrent jobs", func(t *testing.T) {
		pool := NewWorkerPool(4, 20)
		pool.Start()
		defer pool.Stop()

		var wg sync.WaitGroup
		var counter int
		var mu sync.Mutex

		for i := 0; i < 10; i++ {
			wg.Add(1)
			job, err := NewJob(
				func(val int) error {
					defer wg.Done()
					mu.Lock()
					counter += val
					mu.Unlock()
					return nil
				},
				[]any{1},
				nil,
				nil,
			)

			if err != nil {
				t.Fatalf("Failed to create job: %v", err)
			}

			if err := pool.Enqueue(job); err != nil {
				t.Fatalf("Failed to enqueue job: %v", err)
			}
		}

		wg.Wait()

		if counter != 10 {
			t.Errorf("Expected counter to be 10, got %d", counter)
		}
	})

	t.Run("job validation", func(t *testing.T) {
		// Test with invalid number of arguments
		_, err := NewJob(
			func(x int, y int) error { return nil },
			[]any{1}, // Missing one argument
			nil,
			nil,
		)

		if err == nil {
			t.Error("Expected error for invalid number of arguments, got nil")
		}

		// Test with invalid argument type
		_, err = NewJob(
			func(x int) error { return nil },
			[]any{"string"}, // Wrong type
			nil,
			nil,
		)

		if err == nil {
			t.Error("Expected error for invalid argument type, got nil")
		}

		// Test with nil main job
		_, err = NewJob(
			nil,
			[]any{},
			nil,
			nil,
		)

		if err == nil {
			t.Error("Expected error for nil main job, got nil")
		}
	})
} 