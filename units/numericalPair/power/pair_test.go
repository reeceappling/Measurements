package power

import (
	"github.com/reeceappling/measurements/units/numericalPair/pairs"
	"github.com/reeceappling/measurements/units/numericalPair/pairs/inefficientComparable"
	"github.com/reeceappling/measurements/units/numericalPair/pairs/inefficientComparable/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestPowerPairs(t *testing.T) {
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
	t.Run("Pow", func(t *testing.T) {
		t.Run("0 exp returns empty pairs", func(t *testing.T) {
			act := testPairs.Pow(0)
			assert.Equal(t, 0, len(act.Comparable))
			assert.Equal(t, 0, len(act.Uncomparable))
		})
		t.Run("non-0 exp works as intended", func(t *testing.T) {
			tPow := 5
			act := testPairs.Pow(tPow)
			assert.Equal(t, len(testPairs.Comparable), len(act.Comparable))
			for i := range testPairs.Comparable {
				assert.Equal(t, testPairs.Comparable[i].Value, act.Comparable[i].Value)
				assert.Equal(t, testPairs.Comparable[i].Number*tPow, act.Comparable[i].Number)
			}
			assert.Equal(t, len(testPairs.Uncomparable), len(act.Uncomparable))
			for i := range testPairs.Uncomparable {
				assert.Equal(t, testPairs.Uncomparable[i].Value, act.Uncomparable[i].Value)
				assert.Equal(t, testPairs.Uncomparable[i].Number*tPow, act.Uncomparable[i].Number)
			}
		})
	})
	t.Run("CombinePairs", func(t *testing.T) {
		pm1.On("EqualTo", mock.Anything).Return(false)
		pm2.On("EqualTo", mock.Anything).Return(false)
		tp := pairs.Pairs[int, string, string](testPairs)
		exp := Pairs[int, string, string](pairs.CombinePairs[int, string, string](tp, tp))
		assert.Equal(t, exp, CombinePairs[int, string, string](tp, tp))
	})
}
