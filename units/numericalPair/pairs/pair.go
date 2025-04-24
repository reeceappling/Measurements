package pairs

import (
	"github.com/reeceappling/measurements/units/numericalPair/pairs/inefficientComparable"
	"golang.org/x/exp/constraints"
	"slices"
)

// Pairs // TODO: explain
type Pairs[N Mathable, C comparable, T any] struct {
	Comparable   []Pair[N, C]
	Uncomparable []Pair[N, inefficientComparable.IC[T]] // TODO: fixMe
}

type UncomparablePair[N Mathable, T any] Pair[N, inefficientComparable.IC[T]]

type ComparablePair[N Mathable, T comparable] Pair[N, T]

type Pair[N Mathable, T any] struct {
	Value  T
	Number N
}

type Mathable interface { // TODO: rename
	constraints.Float | constraints.Integer // TODO: consider adding constraints.Signed/Complex/Ordered/Unsigned
}

func IndexForInefficientlyComparable[N Mathable, T any](s []Pair[N, inefficientComparable.IC[T]], v inefficientComparable.IC[T]) int { // TODO: stuff
	for i, pair := range s {
		if v.EqualTo(pair.Value) {
			return i
		}
	}
	return -1
}

// CombinePairs is the adds the numbers of like-terms
func CombinePairs[N Mathable, C comparable, U any](in ...Pairs[N, C, U]) Pairs[N, C, U] { // TODO: simplify (less efficient than if we assume each pairSet has no duplicates), test
	out := Pairs[N, C, U]{}
	if len(in) == 0 {
		return out
	}

	// combine comparable
	knownValues := []C{}
	finalComparable := []Pair[N, C]{}
	for _, group := range in { // TODO: ensure this is most efficient way to do this
		for _, comp := range group.Comparable {
			num := comp.Number
			existingIndex := slices.Index(knownValues, comp.Value)
			if existingIndex != -1 {
				finalComparable[existingIndex].Number += num
				continue
			}
			knownValues = append(knownValues, comp.Value) // TODO: is this most efficient?
			finalComparable = append(finalComparable, comp)
		}
	}
	// combine uncomparable
	finalUncomparable := []Pair[N, inefficientComparable.IC[U]]{}
	for _, group := range in {
		for _, uncomp := range group.Uncomparable {
			val, num := uncomp.Value, uncomp.Number
			existingIndex := IndexForInefficientlyComparable(finalUncomparable, val)
			if existingIndex != -1 {
				finalUncomparable[existingIndex].Number += num
				continue
			}
			finalUncomparable = append(finalUncomparable, uncomp)
		}
	}
	return Pairs[N, C, U]{
		Comparable:   finalComparable,
		Uncomparable: finalUncomparable,
	}
}
