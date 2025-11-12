package generics_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zeroflucs-given/generics"
)

func TestExecuteOnce(t *testing.T) {
	t.Parallel()

	t.Run("simple", func(t *testing.T) {
		type XX struct {
			X int64
		}

		counter := atomic.Int32{}
		var task = generics.ExecuteOnce[XX](func(ctx context.Context) (*XX, error) {
			// _ = <-time.After(5 * time.Second)
			counter.Add(1)
			return &XX{X: 42}, nil
		})

		ctx := context.Background()

		expected := &XX{X: 42}

		ch := make(chan struct{})
		go func() {
			x1, err := task.Get(ctx)
			assert.NoError(t, err)
			assert.Equal(t, expected, x1)
			ch <- struct{}{}
		}()

		x2, err := task.Get(ctx)
		require.NoError(t, err)
		require.Equal(t, expected, x2)

		// Check our counter (i.e. the job) has only been executed once.
		require.Equal(t, int32(1), counter.Load())
		<-ch
	})

	t.Run("panic_recovery", func(t *testing.T) {
		task := generics.ExecuteOnce[struct{}](func(ctx context.Context) (*struct{}, error) {
			panic("something went wrong")
		})

		ctx := context.Background()

		go func() {
			_, err := task.Get(ctx)
			assert.Error(t, err)
		}()

		_, err := task.Get(ctx)
		assert.Error(t, err)
	})

	t.Run("context_cancellation", func(t *testing.T) {
		taskWait := make(chan struct{})
		taskStarted := make(chan struct{})

		task := generics.ExecuteOnce[struct{}](func(ctx context.Context) (*struct{}, error) {
			taskStarted <- struct{}{}
			<-taskWait
			return &struct{}{}, nil
		})

		// Start a blocked task and wait for it to start.
		go func() {
			_, err := task.Get(context.Background())
			assert.NoError(t, err)
		}()

		<-taskStarted

		// Start another task waiting for the first to finish, and cancel it.
		ctx1, cancel1 := context.WithCancel(context.Background())
		cancel1()

		_, err := task.Get(ctx1)
		assert.ErrorIs(t, err, context.Canceled)

		taskWait <- struct{}{}
	})

	t.Run("big", func(t *testing.T) {
		type XX struct {
			X int64
		}

		task := generics.ExecuteOnce[XX](func(ctx context.Context) (*XX, error) {
			<-time.After(5 * time.Second)
			return &XX{X: 42}, nil
		})

		ctx := context.Background()

		expected := &XX{X: 42}
		//     100,000 =        5s 140ms
		//   1,000,000 =        6s  40ms
		//  10,000,000 =       25s 100ms
		// 100,000,000 = 5 min 24s
		wg := sync.WaitGroup{}
		wg.Add(1_000_000)
		for i := 0; i < 1_000_000; i += 1 {
			go func() {
				val, err := task.Get(ctx)
				assert.NoError(t, err)
				assert.Equal(t, expected, val)
				wg.Done()

			}()
		}

		wg.Wait()
	})
}
