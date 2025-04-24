package power

import (
	"github.com/reeceappling/measurements/units/numericalPair/pairs"
	"github.com/reeceappling/measurements/units/numericalPair/pairs/inefficientComparable"
	"github.com/reeceappling/randomGoStuff/utils/slices"
)

// Pairs // TODO: explain
type Pairs[N pairs.Mathable, C comparable, U any] pairs.Pairs[N, C, U] // TODO: do power generic.Pairs things

// Pow // TODO: explain
func (thesePairs Pairs[N, C, U]) Pow(exp N) Pairs[N, C, U] { // TODO: figure out where used. Figure out how to not just take int
	if exp == 0 { // TODO: ensure ok with floats
		return Pairs[N, C, U]{}
	}
	pairToExp := func(pair pairs.Pair[N, C]) pairs.Pair[N, C] {
		return pairs.Pair[N, C]{pair.Value, pair.Number * exp}
	}
	uncompPairToExp := func(pair pairs.Pair[N, inefficientComparable.IC[U]]) pairs.Pair[N, inefficientComparable.IC[U]] {
		return pairs.Pair[N, inefficientComparable.IC[U]]{pair.Value, pair.Number * exp}
	}
	return Pairs[N, C, U]{
		Comparable:   slices.Map(thesePairs.Comparable, pairToExp),
		Uncomparable: slices.Map(thesePairs.Uncomparable, uncompPairToExp),
	}
}

// CombinePairs is the equivalent of multiplying(?) like-terms // TODO: explain
func CombinePairs[N pairs.Mathable, C comparable, U any](in ...pairs.Pairs[N, C, U]) Pairs[N, C, U] { // TODO: simplify
	return Pairs[N, C, U](pairs.CombinePairs(in...))
}
