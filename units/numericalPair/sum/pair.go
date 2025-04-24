package sum

import (
	"github.com/reeceappling/measurements/units/numericalPair/pairs"
)

// Pairs // TODO: explain
type Pairs[N pairs.Mathable, C comparable, U any] pairs.Pairs[N, C, U]

func CombinePairs[N pairs.Mathable, C comparable, U any](in ...pairs.Pairs[N, C, U]) Pairs[N, C, U] { // TODO: simplify
	return Pairs[N, C, U](pairs.CombinePairs(in...))
}
