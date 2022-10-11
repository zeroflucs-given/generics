package query

import (
	"context"

	"github.com/zeroflucs-given/generics"
	"github.com/zeroflucs-given/generics/filtering"
)

// Slice creates an instance of a slice-query over a a slice
func Slice[T any](input []T) SliceQuery[T] {
	return sliceQueryInternal[T]{
		input: input,
	}
}

// SliceQuery is our query type for non context-aware slices. You can use WithContext to
// extend with various context aware operations, but the two types are not compatible directly.
type SliceQuery[T any] interface {
	// All returns true if all items in the slice meet our filters
	All(filters ...filtering.Expression[T]) bool

	// Any returns true if any item in the slice matches our filters
	Any(filters ...filtering.Expression[T]) bool

	// Concatenate multiple slices together with this slice
	Concatenate(slices ...[]T) SliceQuery[T]

	// Count the number of items matching the filters
	Count(filters ...filtering.Expression[T]) int

	// Filter the slice
	Filter(filters ...filtering.Expression[T]) SliceQuery[T]

	// First item of the slice matching the filter
	First(filters ...filtering.Expression[T]) T

	// Last item of the slice matching the filter
	Last(filters ...filtering.Expression[T]) T

	// Reverse inverts the order of the slice
	Reverse() SliceQuery[T]

	// Skip n items from the slice
	Skip(n int) SliceQuery[T]

	// Take takes n items from the slice
	Take(n int) SliceQuery[T]

	// TakeUntil takes items from the slice until the first item matching the predicate
	TakeUntil(filters ...filtering.Expression[T]) SliceQuery[T]

	// TakeWhile takes items from the slice until the first item no longer matching
	TakeWhile(filters ...filtering.Expression[T]) SliceQuery[T]

	// ToSlice drops the wrapper type and returns a raw slice
	ToSlice() []T

	// WithContext sets the context used for the slice query
	WithContext(ctx context.Context) SliceQueryContext[T]
}

type sliceQueryInternal[T any] struct {
	input []T
}

// All returns true if all items in the slice meet our filters
func (s sliceQueryInternal[T]) All(filters ...filtering.Expression[T]) bool {
	return generics.All(s.input, filters...)
}

// Any returns true if any item in the slice matches our filters
func (s sliceQueryInternal[T]) Any(filters ...filtering.Expression[T]) bool {
	return generics.Any(s.input, filters...)
}

// Concatenate multiple slices together with this slice
func (s sliceQueryInternal[T]) Concatenate(slices ...[]T) SliceQuery[T] {
	var data [][]T
	data = append(data, s.input)
	data = append(data, slices...)

	return &sliceQueryInternal[T]{
		input: generics.Concatenate(data...),
	}
}

// Count the number of items matching the filters
func (s sliceQueryInternal[T]) Count(filters ...filtering.Expression[T]) int {
	return generics.Count(s.input, filters...)
}

// Filter the slice
func (s sliceQueryInternal[T]) Filter(filters ...filtering.Expression[T]) SliceQuery[T] {
	return sliceQueryInternal[T]{
		input: generics.Filter(s.input, filters...),
	}
}

// First item of the slice matching the filter
func (s sliceQueryInternal[T]) First(filters ...filtering.Expression[T]) T {
	return generics.First(s.input, filters...)
}

// Last item of the slice matching the filter
func (s sliceQueryInternal[T]) Last(filters ...filtering.Expression[T]) T {
	return generics.Last(s.input, filters...)
}

// Reverse inverts the order of the slice
func (s sliceQueryInternal[T]) Reverse() SliceQuery[T] {
	return sliceQueryInternal[T]{
		input: generics.Reverse(s.input),
	}
}

// Skip n items from the slice
func (s sliceQueryInternal[T]) Skip(n int) SliceQuery[T] {
	return sliceQueryInternal[T]{
		input: generics.Skip(s.input, n),
	}
}

/*
// SkipUntil skips items until the first item matching the predicate
func (s Slice[T]) SkipUntil(filters ...filtering.Expression[T]) Slice[T] {
	return SkipUntil(s, filters...)
}

// SkipWhile skips items whilst the predicate is true
func (s Slice[T]) SkipWhile(filters ...filtering.Expression[T]) Slice[T] {
	return SkipWhile(s, filters...)
}
*/

// Take takes n items from the slice
func (s sliceQueryInternal[T]) Take(n int) SliceQuery[T] {
	return sliceQueryInternal[T]{
		input: generics.Take(s.input, n),
	}
}

// TakeUntil takes items from the slice until the first item matching the predicate
func (s sliceQueryInternal[T]) TakeUntil(filters ...filtering.Expression[T]) SliceQuery[T] {
	return sliceQueryInternal[T]{
		input: generics.TakeUntil(s.input, filters...),
	}
}

// TakeWhile takes items from the slice until the first item no longer matching
func (s sliceQueryInternal[T]) TakeWhile(filters ...filtering.Expression[T]) SliceQuery[T] {
	return sliceQueryInternal[T]{
		input: generics.TakeWhile(s.input, filters...),
	}
}

// ToSlice drops the wrapper type and returns a raw slice
func (s sliceQueryInternal[T]) ToSlice() []T {
	return s.input
}

// WithContext sets the context used for the slice
func (s sliceQueryInternal[T]) WithContext(ctx context.Context) SliceQueryContext[T] {
	return sliceQueryContextInternal[T]{
		ctx:   ctx,
		input: s.input,
	}
}
