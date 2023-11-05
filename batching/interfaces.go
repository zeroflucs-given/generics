package batching

import "context"

// Batcher is a contract that describes a batcher of operations
type Batcher[TRequest any, TResponse any] interface {
	// Start the batcher
	Start()

	// Stop the batcher
	Stop()

	// Execute a request, batching where possible
	Execute(ctx context.Context, req TRequest) (TResponse, error)
}

// BatchContextProvider is a function that provides the context to execute batches with
type BatchContextProvider func() context.Context

// BatchRequestBuilder is a function that builds a batch
type BatchRequestBuilder[TRequest any, TBatchRequest any] func(requests []TRequest) TBatchRequest

// BatchExecutor is a function that executes a batched request
type BatchExecutor[TBatchRequest any, TBatchResponse any] func(ctx context.Context, req TBatchRequest) (TBatchResponse, error)

// BatchResponseUnpacker unpacks a batched response
type BatchResponseUnpacker[TBatchResponse any, TResponse any] func(response TBatchResponse) []TResponse
