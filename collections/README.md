# collections
Generic structures for Go 1.18+ developers

## Interfaces
The following core interfaces are defined in `collections` package:

| Interface | Role | Notes |
|-----------|------|-------|
| Queue[T]  | Queue, Order Invariant | This interface defines a any queue where you can _Push_ a value, _Pop_ a value and _Count_ the contents of the object. |
| 

## Included Packages
These packages are thread-safe, using mutexes to maintain consistency.

| Package | Thread Safety | Interfaces | Notes |
|---------|-------------|------------|-------|
| `collections/bplustree` | Concurrent Reads & Single Writer | TreeMap[K, V] | A B+ tree implementation that implements a seekable list of key-values. |
| `collections/linkedlist` | Concurrent Reads & Single Writer | Queue[T] | A linked list that implements Queue[T] with FIFO semantics. Capacity limited by system resources. |
| `collections/ringbuffer` | Concurrent Reads & Single Writer | Queue[T] | A linked list with a fixed upper size that implements Queue[T] with FIFO semantics, optimised for fixed sets of data. Attempts tow write data when full will return errors. |
| `collections/stack` | Concurrent Reads & Single Writer | Queue[T] | A fixed size stack that implements Queue[T] with LIFO semantics. Attempts to exceed stack capacity will return errors. |
| `collections/weightedrandom` | Concurrent Reads & Single Writer | N/A | Allows selection of a value from a set of values in accordance with their relative weights/frequencies. Weights can be any `Comparable` type, but you must supply a mapper function that reduces these values to the space of float64(0>maxFloat64)

## Non-Thread Safe
The `lockless` sub-package contains variants of the existing packages. These 
packages are not thread safe, and concurrent operations must be controlled by 
consuming code.

| Package | Notes |
|---------|-------|
| `collections/lockless/ringbuffer` | A non-locking version of the circular ring buffer. Assumes that it is used only in contexts that prevent concurrent operations. |