package weightedrand

import (
	"errors"
	"math/rand"
	"sort"
)

// Choice is a generic wrapper that can be used to add weights for any item.
type Choice[T any, W integer] struct {
	Item   T
	Weight W
}

type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// NewChoice creates a new Choice with specified item and weight.
func NewChoice[T any, W integer](item T, weight W) Choice[T, W] {
	return Choice[T, W]{Item: item, Weight: weight}
}

// A Chooser caches many possible Choices in a structure designed to improve
// performance on repeated calls for weighted random selection.
type Chooser[T any, W integer] struct {
	data   []Choice[T, W]
	totals []int
	max    int
}

// NewChooser initializes a new Chooser for picking from the provided choices.
func NewChooser[T any, W integer](choices ...Choice[T, W]) (*Chooser[T, W], error) {
	sort.Slice(choices, func(i, j int) bool {
		return choices[i].Weight < choices[j].Weight
	})

	totals := make([]int, len(choices))
	runningTotal := 0
	for i, c := range choices {
		if c.Weight < 0 {
			continue // ignore negative weights, can never be picked
		}

		// case of single ~uint64 or similar value that exceeds maxInt on its own
		if uint64(c.Weight) >= maxInt {
			return nil, errWeightOverflow
		}

		weight := int(c.Weight) // convert weight to int for internal counter usage
		if (maxInt - runningTotal) <= weight {
			return nil, errWeightOverflow
		}
		runningTotal += weight
		totals[i] = runningTotal
	}

	if runningTotal < 1 {
		return nil, errNoValidChoices
	}

	return &Chooser[T, W]{data: choices, totals: totals, max: runningTotal}, nil
}

const (
	intSize   = 32 << (^uint(0) >> 63) // cf. strconv.IntSize
	maxInt    = 1<<(intSize-1) - 1
	maxUint64 = 1<<64 - 1
)

// Possible errors returned by NewChooser, preventing the creation of a Chooser
// with unsafe runtime states.
var (
	// If the sum of provided Choice weights exceed the maximum integer value
	// for the current platform (e.g. math.MaxInt32 or math.MaxInt64), then
	// the internal running total will overflow, resulting in an imbalanced
	// distribution generating improper results.
	errWeightOverflow = errors.New("sum of Choice Weights exceeds max int")
	// If there are no Choices available to the Chooser with a weight >= 1,
	// there are no valid choices and Pick would produce a runtime panic.
	errNoValidChoices = errors.New("zero Choices with Weight >= 1")
)

// Pick returns a single weighted random Choice.Item from the Chooser.
//
// Utilizes global rand as the source of randomness. Safe for concurrent usage.
func (c Chooser[T, W]) Pick() T {
	r := rand.Intn(c.max) + 1
	i := searchInts(c.totals, r)
	return c.data[i].Item
}

// PickSource returns a single weighted random Choice.Item from the Chooser,
// utilizing the provided *rand.Rand source rs for randomness.
//
// The primary use-case for this is avoid lock contention from the global random
// source if utilizing Chooser(s) from multiple goroutines in extremely
// high-throughput situations.
//
// It is the responsibility of the caller to ensure the provided rand.Source is
// free from thread safety issues.
//
// Deprecated: Since go1.21 global rand no longer suffers from lock contention
// when used in multiple high throughput goroutines, as long as you don't
// manually seed it. Use [Chooser.Pick] instead.
func (c Chooser[T, W]) PickSource(rs *rand.Rand) T {
	r := rs.Intn(c.max) + 1
	i := searchInts(c.totals, r)
	return c.data[i].Item
}

// The standard library sort.SearchInts() just wraps the generic sort.Search()
// function, which takes a function closure to determine truthfulness. However,
// since this function is utilized within a for loop, it cannot currently be
// properly inlined by the compiler, resulting in non-trivial performance
// overhead.
//
// Thus, this is essentially manually inlined version.  In our use case here, it
// results in a significant throughput increase for Pick.
//
// See also github.com/mroth/xsort.
func searchInts(a []int, x int) int {
	// Possible further future optimization for searchInts via SIMD if we want
	// to write some Go assembly code: http://0x80.pl/articles/simd-search.html
	i, j := 0, len(a)
	for i < j {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		if a[h] < x {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}
