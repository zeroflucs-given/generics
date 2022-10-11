package weightedrandom

import (
	"fmt"
	"math"
	"math/rand"
	"sync"

	"github.com/zeroflucs-given/generics"
)

// DomainMapper maps the value from the input type to the domain of floats
type DomainMapper[V generics.Numeric] func(v V) float64

// UnsignedFloatMapper clamps all values to the positive float range
var UnsignedFloatMapper = func(v float64) float64 {
	return math.Abs(v)
}

// NewWeightedRandom creates a new weighted random picker
func NewWeightedRandom[T any, W generics.Numeric](initialCapacity int, mapper DomainMapper[W]) *WeightedRandom[T, W] {
	return &WeightedRandom[T, W]{
		mapper:       mapper,
		data:         make([]T, 0, initialCapacity),
		runningTotal: make([]float64, 0, initialCapacity),
	}
}

// WeightedRandom is a structure that performs weighted-random selections of values
// according to their relative frequencies. The maximum value of the weights must
// not exceed te capacity of a Float64
type WeightedRandom[T any, W generics.Numeric] struct {
	mapper       DomainMapper[W]
	data         []T
	runningTotal []float64
	lock         sync.RWMutex
	total        float64
}

// Pick a random value
func (w *WeightedRandom[T, W]) Pick(rnd *rand.Rand) (bool, T) {
	w.lock.RLock()

	if len(w.data) == 0 {
		var blank T
		w.lock.RUnlock()
		return false, blank
	}

	watermark := rnd.Float64() * w.total

	data := w.data
	running := w.runningTotal

	low := 0
	high := len(w.data) - 1

	for low <= high {
		median := (low + high) / 2

		if running[median] < watermark {
			low = median + 1
		} else {
			high = median - 1
		}
	}

	w.lock.RUnlock()
	return true, data[low]
}

// Push a value into the set of weighted random values
func (w *WeightedRandom[T, W]) Push(weight W, value T) {
	w.lock.Lock()

	mappedWeight := w.mapper(weight)

	// Panic not error, due to keeping interface nice and clean (any mapper that does not work
	// as expected is inherently faulty, and this issue is not a runtime recoverable scenario).
	if mappedWeight < 0 {
		panic(fmt.Errorf("the weight %v mappped to %v which is invalid", weight, mappedWeight))
	}

	w.total += mappedWeight
	w.data = append(w.data, value)
	w.runningTotal = append(w.runningTotal, w.total)

	w.lock.Unlock()
}
