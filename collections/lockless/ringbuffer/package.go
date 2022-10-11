package ringbuffer

// Package ringbuffer contains a non-thread safe implementation of the
// buffer from the ringbuffer package, optimised for scenarios that do not require
// the weight associated with sync.Mutex.
//
// Uses Go generics. It allows for storage of N items of a type T with FIFO
// semantics. When pushing more than the buffer can hold, an error will be
// generated.
//
// On a MacBook Pro i9-8950HK the benchmarks included in this repository can
// push/pop cycle the buffer 379 million items/second. If you want a
// thread-safe version, use the ringbuffer package, however the mutexes
// account for ~95% of time.
