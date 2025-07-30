package generics

import (
	"context"
	"fmt"
	"iter"
	"slices"

	"github.com/zeroflucs-given/generics/filtering"
)

// First item in a slice that passes the filters. If multiple filters are set, they are treated
// as a logical AND. If no filters are set, will return the first item in the slice. If no items
// match, the type default is returned.
func First[T any](items []T, filters ...filtering.Expression[T]) T {
	return Must(FirstWithContext(context.Background(), items, filtering.AndWrapWithContext(filters...)))
}

// FirstWithContext gets the first item in a slice that passes the filters. If multiple filters are set, they are treated
// as a logical AND. If no filters are set, will return the first item in the slice. If no items
// match, the type default is returned.
func FirstWithContext[T any](ctx context.Context, items []T, filters ...filtering.ExpressionWithContext[T]) (T, error) {
	return FirstSeqWithContext(ctx, slices.Values(items), filters...)
}

func FirstSeq[T any](seq iter.Seq[T], filters ...filtering.Expression[T]) T {
	return Must(FirstSeqWithContext(context.Background(), seq, filtering.AndWrapWithContext(filters...)))
}

func FirstSeqWithContext[T any](ctx context.Context, items iter.Seq[T], filters ...filtering.ExpressionWithContext[T]) (T, error) {
	filter := filtering.AndWithContext(filters...)

	var def T
	i := 0
	for v := range items {
		match, err := filter(ctx, i, v)
		if err != nil {
			return def, fmt.Errorf("error applying filter to item %d: %w", i, err)
		}

		if match {
			return v, nil
		}
		i += 1
	}

	return def, nil
}

// Except
func Except[T comparable](items []T, exclusions ...T) []T {
	result := make([]T, 0, len(items))
	for _, item := range items {
		if Contains(exclusions, item) {
			continue
		}
		result = append(result, item)
	}

	return result
}

// Intersect
func Intersect[T comparable](items []T, others ...T) []T {
	result := make([]T, 0, len(items))

	for _, item := range items {
		if Contains(others, item) {
			result = append(result, item)
		}
	}

	return result
}

// Filter filters item in a list
func Filter[T any](items []T, filters ...filtering.Expression[T]) []T {
	return Must(FilterWithContext(context.Background(), items, filtering.AndWrapWithContext(filters...)))
}

// FilterWithContext filters item in a list
func FilterWithContext[T any](ctx context.Context, items []T, filters ...filtering.ExpressionWithContext[T]) ([]T, error) {
	filter := filtering.AndWithContext(filters...)

	output := make([]T, 0, len(items))
	for i, v := range items {
		ok, err := filter(ctx, i, v)
		if err != nil {
			return nil, fmt.Errorf("error applying filter to item %d: %w", i, err)
		}
		if ok {
			output = append(output, v)
		}
	}

	return output, nil
}

func FilterSeq[T any](items iter.Seq[T], filters ...filtering.Expression[T]) []T {
	return Must(FilterSeqWithContext(context.Background(), items, filtering.AndWrapWithContext(filters...)))
}

func FilterSeqWithContext[T any](ctx context.Context, items iter.Seq[T], filters ...filtering.ExpressionWithContext[T]) ([]T, error) {
	filter := filtering.AndWithContext(filters...)

	var output []T
	i := 0
	for v := range items {
		ok, err := filter(ctx, i, v)
		if err != nil {
			return nil, fmt.Errorf("error applying filter to item %d: %w", i, err)
		}
		if ok {
			output = append(output, v)
		}

		i += 1
	}

	return output, nil
}

// Last item in a slice that matches the specified filters. Returns the type
// default if none found.
func Last[T any](items []T, filters ...filtering.Expression[T]) T {
	return Must(LastWithContext(context.Background(), items, filtering.AndWrapWithContext(filters...)))
}

func LastSeq[T any](seq iter.Seq[T], filters ...filtering.Expression[T]) T {
	return Must(LastSeqWithContext(context.Background(), seq, filtering.AndWrapWithContext(filters...)))
}

func LastSeqWithContext[T any](ctx context.Context, seq iter.Seq[T], filters ...filtering.ExpressionWithContext[T]) (T, error) {
	// We have to collect it here, unfortunately.
	return LastWithContext(ctx, slices.Collect(seq), filters...)
}

// LastWithContext item in a slice that matches the specified filters. Returns the type
// default if none found.
func LastWithContext[T any](ctx context.Context, items []T, filters ...filtering.ExpressionWithContext[T]) (T, error) {
	filter := filtering.AndWithContext(filters...)

	var def T
	for reverseIndex := len(items) - 1; reverseIndex >= 0; reverseIndex-- {
		match, err := filter(ctx, reverseIndex, items[reverseIndex])
		if err != nil {
			return def, fmt.Errorf("error applying filter to index %d: %w", reverseIndex, err)
		}
		if match {
			return items[reverseIndex], nil
		}
	}

	return def, nil
}
