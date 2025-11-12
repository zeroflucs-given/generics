package batching_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeroflucs-given/generics"
	"github.com/zeroflucs-given/generics/batching"
)

func TestAbstractBatcher(t *testing.T) {
	// Create a batcher that does a dumb multiply-by-2 with delay at 5msec intervals
	batcher := batching.NewBatcher[int, int, []int, []int](
		func() context.Context {
			return context.TODO()
		},
		func(requests []int) []int {
			return requests
		},
		func(ctx context.Context, req []int) ([]int, error) {
			t.Logf("Executing batch of %d items", len(req))
			return generics.Map(req, func(i int, v int) int {
				return v * 2
			}), nil
		},
		func(response []int) []int {
			return response
		},
		100,
		time.Millisecond*50,
	)
	batcher.Start()
	defer batcher.Stop()

	wg := sync.WaitGroup{}
	t.Log("Queuing items")

	for i := 0; i < 1000; i++ {
		v := i
		wg.Add(1)
		go func() {
			res, err := batcher.Execute(context.TODO(), v)
			assert.NoError(t, err, "Should have no error at index %d", v)
			assert.Equal(t, v*2, res, "Should have got right result for item at index %d", v)
			wg.Done()
		}()

		time.Sleep(time.Millisecond)
	}

	t.Log("Awaiting items")
	wg.Wait()
}

func TestAbstractBatcherErrors(t *testing.T) {
	// Create a batcher that does a dumb multiply-by-2 with delay at 5msec intervals
	errExpect := errors.New("boom")
	batcher := batching.NewBatcher[int, int, []int, []int](
		func() context.Context {
			return context.TODO()
		},
		func(requests []int) []int {
			return requests
		},
		func(ctx context.Context, req []int) ([]int, error) {
			t.Logf("Executing batch of %d items", len(req))
			return nil, errExpect
		},
		func(response []int) []int {
			return response
		},
		100,
		time.Second,
	)
	batcher.Start()
	defer batcher.Stop()

	wg := sync.WaitGroup{}
	t.Log("Queuing items")

	for i := 0; i < 100; i++ {
		v := i
		wg.Add(1)
		go func() {
			res, err := batcher.Execute(context.TODO(), v)
			assert.ErrorIs(t, err, errExpect, "Should have error at index %d", v)
			assert.Zero(t, res, "Should have no data for index %d", v)
			wg.Done()
		}()
	}

	t.Log("Awaiting items")
	wg.Wait()
}
