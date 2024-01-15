package batching

import (
	"context"
	"sync"
	"time"
)

// NewBatcher creates a new batcher instance
func NewBatcher[TRequest any, TResponse any, TBatchRequest any, TBatchResponse any](
	batchContextProvider BatchContextProvider,
	batchBuilder BatchRequestBuilder[TRequest, TBatchRequest],
	batchExecutor BatchExecutor[TBatchRequest, TBatchResponse],
	batchUnpacker BatchResponseUnpacker[TBatchResponse, TResponse],
	maxItemsPerBatch int,
	maxWait time.Duration,
) Batcher[TRequest, TResponse] {
	return &abstractBatcher[TRequest, TResponse, TBatchRequest, TBatchResponse]{
		maxItemsPerBatch: maxItemsPerBatch,
		maxWait:          maxWait,
		batchContext:     batchContextProvider,
		batchBuilder:     batchBuilder,
		batchExecutor:    batchExecutor,
		batchUnpacker:    batchUnpacker,
		done:             make(chan struct{}),
	}
}

// AbstractBatcher is a batcher of arbitrary things
type abstractBatcher[TRequest any, TResponse any, TBatchRequest any, TBatchResponse any] struct {
	maxItemsPerBatch int           // Maximum items per batch
	maxWait          time.Duration // Maximum wait Msec
	batchContext     BatchContextProvider
	batchBuilder     BatchRequestBuilder[TRequest, TBatchRequest]
	batchExecutor    BatchExecutor[TBatchRequest, TBatchResponse]
	batchUnpacker    BatchResponseUnpacker[TBatchResponse, TResponse]

	mtx          sync.Mutex
	currentBatch *batch[TRequest, TResponse, TBatchRequest, TBatchResponse]
	done         chan struct{}
}

// Start commences a goroutine that dispatces the items
func (a *abstractBatcher[TRequest, TResponse, TBatchRequest, TBatchResponse]) Start() {
	go func() {
		ticker := time.NewTicker(time.Duration(a.maxWait))
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				a.mtx.Lock()
				a.dispatch()
				a.mtx.Unlock()
			case <-a.done:
				return
			}
		}
	}()
}

// Stop the dispatcher loop
func (a *abstractBatcher[TRequest, TResponse, TBatchRequest, TBatchResponse]) Stop() {
	close(a.done)
}

// Execute a batched request
func (a *abstractBatcher[TRequest, TResponse, TBatchRequest, TBatchResponse]) Execute(ctx context.Context, req TRequest) (TResponse, error) {
	a.mtx.Lock()

	// Start a new batch if needed
	if a.currentBatch == nil {
		a.currentBatch = &batch[TRequest, TResponse, TBatchRequest, TBatchResponse]{
			batchContext:  a.batchContext,
			batchBuilder:  a.batchBuilder,
			batchExecutor: a.batchExecutor,
			batchUnpacker: a.batchUnpacker,
		}
		a.currentBatch.wg.Add(1)
	}

	batchRef := a.currentBatch
	itemIndex := len(batchRef.requests)
	batchRef.requests = append(batchRef.requests, req)

	// Eager dispatch
	if len(batchRef.requests) >= a.maxItemsPerBatch {
		a.dispatch()
	}

	a.mtx.Unlock()

	return batchRef.AwaitAndReturn(ctx, itemIndex)
}

// Dispatch a batch of items and reset
func (a *abstractBatcher[TRequest, TResponse, TBatchRequest, TBatchResponse]) dispatch() {
	if a.currentBatch != nil {
		ref := a.currentBatch
		go func() {
			ref.dispatch()
		}()
	}

	a.currentBatch = nil
}

type batch[TRequest any, TResponse any, TBatchRequest any, TBatchResponse any] struct {
	batchContext  BatchContextProvider
	batchBuilder  BatchRequestBuilder[TRequest, TBatchRequest]
	batchExecutor BatchExecutor[TBatchRequest, TBatchResponse]
	batchUnpacker BatchResponseUnpacker[TBatchResponse, TResponse]
	requests      []TRequest
	responses     []TResponse
	wg            sync.WaitGroup
	err           error
}

// dispatch a single batced request
func (b *batch[TRequest, TResponse, TBatchRequest, TBatchResponse]) dispatch() {
	defer b.wg.Done()

	// Build our object
	ctx := b.batchContext()
	req := b.batchBuilder(b.requests)

	// Execute the work
	res, err := b.batchExecutor(ctx, req)

	// Unpack
	if err == nil {
		b.responses = b.batchUnpacker(res)
	}
	b.err = err

}

// Await the batch to complete, and return our request by ordinal
func (b *batch[TRequest, TResponse, TBatchReq, TBatchRes]) AwaitAndReturn(ctx context.Context, index int) (TResponse, error) {
	var result TResponse
	b.wg.Wait()

	if len(b.responses) > index {
		result = b.responses[index]
	}

	return result, b.err
}
