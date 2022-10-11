package bplustree

// Package bplustree is a B+ tree implementation in pure Go, using generics.

/**
	On our test with a large data set, 27 order with 2,500,000 inserts we got this time with:

		go test -benchmem -run=^$ -bench ^BenchmarkInsertsSequentialFixedLarge$ -benchtime 10s
		BenchmarkInsertsSequentialFixedLarge-12  	31         350878338 ns/op        166408388 B/op    624992 allocs/op

	We then optimised this to reduce allocations, by using a ring-buffer of pre-allocated blocks of
	records/children. This yielded a 13% improvement:

		BenchmarkInsertsSequentialFixedLarge-12     38         305604530 ns/op        173427506 B/op    212506 allocs/op

	This means that the tree can accomodate approximately ~10 million inserts/second.
 **/
