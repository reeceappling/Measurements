package sum

import (
	"github.com/reeceappling/measurements/units/numericalPair/pairs"
	"github.com/reeceappling/measurements/units/numericalPair/pairs/inefficientComparable"
	"github.com/reeceappling/measurements/units/numericalPair/pairs/inefficientComparable/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestSumPair(t *testing.T) {
	pm1 := mocks.NewIC[string](t)
	pm2 := mocks.NewIC[string](t)
	testPairs := Pairs[int, string, string]{
		Comparable: []pairs.Pair[int, string]{
			{"a", 1},
			{"b", 2},
		},
		Uncomparable: []pairs.Pair[int, inefficientComparable.IC[string]]{
			{pm1, -3},
			{pm2, -4},
		},
	}
	t.Run("CombinePairs", func(t *testing.T) {
		pm1.On("EqualTo", mock.Anything).Return(false)
		pm2.On("EqualTo", mock.Anything).Return(false)
		tp := pairs.Pairs[int, string, string](testPairs)
		exp := Pairs[int, string, string](pairs.CombinePairs[int, string, string](tp, tp))
		assert.Equal(t, exp, CombinePairs[int, string, string](tp, tp))
	})
}
