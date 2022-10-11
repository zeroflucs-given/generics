package query

import (
	"context"

	"github.com/zeroflucs-given/generics"
	"github.com/zeroflucs-given/generics/filtering"
)

// SliceQueryContext is context-aware operations. Any operation that materializes back to a slice
// or scalar will trigger any failures accumulated to be propegated to the caller. You can use
// filtering.Wrap to convert any filters to be context aware.
type SliceQueryContext[T any] interface {
	// All returns true if all items in the slice meet our filters
	All(filters ...filtering.ExpressionWithContext[T]) (bool, error)

	// Any returns true if any item in the slice matches our filters
	Any(filters ...filtering.ExpressionWithContext[T]) (bool, error)

	// Concatenate multiple slices together with this slice
	Concatenate(slices ...[]T) SliceQueryContext[T]

	// Count the number of items matching the filters
	Count(filters ...filtering.ExpressionWithContext[T]) (int, error)

	// Filter the slice
	Filter(filters ...filtering.ExpressionWithContext[T]) SliceQueryContext[T]

	// First item of the slice matching the filter
	First(filters ...filtering.ExpressionWithContext[T]) (T, error)

	// Last item of the slice matching the filter
	Last(filters ...filtering.ExpressionWithContext[T]) (T, error)

	// Reverse inverts the order of the slice
	Reverse() SliceQueryContext[T]

	// Skip n items from the slice
	Skip(n int) SliceQueryContext[T]

	// Take takes n items from the slice
	Take(n int) SliceQueryContext[T]

	// TakeUntil takes items from the slice until the first item matching the predicate
	TakeUntil(filters ...filtering.ExpressionWithContext[T]) SliceQueryContext[T]

	// TakeWhile takes items from the slice until the first item no longer matching
	TakeWhile(filters ...filtering.ExpressionWithContext[T]) SliceQueryContext[T]

	// ToSlice drops the wrapper type and returns a raw slice
	ToSlice() ([]T, error)
}

type sliceQueryContextInternal[T any] struct {
	ctx   context.Context
	input []T
	err   error
}

// All returns true if all items in the slice match our filters
func (s sliceQueryContextInternal[T]) All(filters ...filtering.ExpressionWithContext[T]) (bool, error) {
	if s.err != nil {
		return false, s.err
	}

	return generics.AllWithContext(s.ctx, s.input, filters...)
}

// Any returns true if any item in the slice matches our filters
func (s sliceQueryContextInternal[T]) Any(filters ...filtering.ExpressionWithContext[T]) (bool, error) {
	if s.err != nil {
		return false, s.err
	}

	return generics.AnyWithContext(s.ctx, s.input, filters...)
}

// Concatenate multiple slices together with this slice
func (s sliceQueryContextInternal[T]) Concatenate(slices ...[]T) SliceQueryContext[T] {
	if s.err != nil {
		return s
	}

	var data [][]T
	data = append(data, s.input)
	data = append(data, slices...)

	return sliceQueryContextInternal[T]{
		ctx:   s.ctx,
		input: generics.Concatenate(data...),
	}
}

// Count the number of items matching the filters
func (s sliceQueryContextInternal[T]) Count(filters ...filtering.ExpressionWithContext[T]) (int, error) {
	if s.err != nil {
		return 0, s.err
	}

	return generics.CountWithContext(s.ctx, s.input, filters...)
}

// Filter the slice
func (s sliceQueryContextInternal[T]) Filter(filters ...filtering.ExpressionWithContext[T]) SliceQueryContext[T] {
	if s.err != nil {
		return s
	}

	filtered, err := generics.FilterWithContext(s.ctx, s.input, filters...)
	return &sliceQueryContextInternal[T]{
		ctx:   s.ctx,
		input: filtered,
		err:   err,
	}
}

// First item of the slice matching the filter
func (s sliceQueryContextInternal[T]) First(filters ...filtering.ExpressionWithContext[T]) (T, error) {
	var def T
	if s.err != nil {
		return def, s.err
	}

	return generics.FirstWithContext(s.ctx, s.input, filters...)
}

// Last item of the slice matching the filter
func (s sliceQueryContextInternal[T]) Last(filters ...filtering.ExpressionWithContext[T]) (T, error) {
	var def T
	if s.err != nil {
		return def, s.err
	}

	return generics.LastWithContext(s.ctx, s.input, filters...)
}

// Reverse inverts the order of the slice
func (s sliceQueryContextInternal[T]) Reverse() SliceQueryContext[T] {
	if s.err != nil {
		return s
	}

	return sliceQueryContextInternal[T]{
		ctx:   s.ctx,
		input: generics.Reverse(s.input),
	}
}

// Skip n items from the slice
func (s sliceQueryContextInternal[T]) Skip(n int) SliceQueryContext[T] {
	if s.err != nil {
		return s
	}

	return sliceQueryContextInternal[T]{
		ctx:   s.ctx,
		input: generics.Skip(s.input, n),
	}
}

// Take n items from the slice
func (s sliceQueryContextInternal[T]) Take(n int) SliceQueryContext[T] {
	if s.err != nil {
		return s
	}

	return sliceQueryContextInternal[T]{
		ctx:   s.ctx,
		input: generics.Take(s.input, n),
	}
}

// TakeUntil takes items from the slice until the first item matching the predicate
func (s sliceQueryContextInternal[T]) TakeUntil(filters ...filtering.ExpressionWithContext[T]) SliceQueryContext[T] {
	if s.err != nil {
		return s
	}

	taken, err := generics.TakeUntilWithContext(s.ctx, s.input, filters...)
	return sliceQueryContextInternal[T]{
		ctx:   s.ctx,
		input: taken,
		err:   err,
	}
}

// TakeWhile takes items from the slice until the first item no longer matching
func (s sliceQueryContextInternal[T]) TakeWhile(filters ...filtering.ExpressionWithContext[T]) SliceQueryContext[T] {
	if s.err != nil {
		return s
	}

	taken, err := generics.TakeWhileWithContext(s.ctx, s.input, filters...)
	return sliceQueryContextInternal[T]{
		ctx:   s.ctx,
		input: taken,
		err:   err,
	}
}

// ToSlice drops the wrapper type and returns a raw slice
func (s sliceQueryContextInternal[T]) ToSlice() ([]T, error) {
	return s.input, s.err
}
