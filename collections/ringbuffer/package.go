package ringbuffer

// Package ringbuffer contains a thread-safe ring-buffer implementation that
// uses Go generics. It allows for storage of N items of a type T with FIFO
// semantics. When pushing more than the buffer can hold, an error will be
// generated.
//
// On a MacBook Pro i9-8950HK the benchmarks included in this repository can
// push/pop cycle the buffer 25.1 million items/second. If you want a non
// thread-safe version, use the ringbuffer package.
