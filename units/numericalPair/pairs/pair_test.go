package pairs

import (
	"github.com/reeceappling/measurements/units/numericalPair/pairs/inefficientComparable"
	"github.com/reeceappling/measurements/units/numericalPair/pairs/inefficientComparable/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestPairs(t *testing.T) {
	pm1 := mocks.NewIC[string](t)
	pm2 := mocks.NewIC[string](t)
	pm3 := mocks.NewIC[string](t)
	comparables := []Pair[int, string]{
		{"a", 1},
		{"b", -2},
		{"c", 3},
	}
	uncomparables := []Pair[int, inefficientComparable.IC[string]]{
		{pm1, -4},
		{pm2, 5},
		{pm3, -6},
	}
	testPairs := Pairs[int, string, string]{
		Comparable: []Pair[int, string]{
			comparables[0],
			comparables[1],
		},
		Uncomparable: []Pair[int, inefficientComparable.IC[string]]{
			uncomparables[0],
			uncomparables[1],
		},
	}

	// TODO: EVERYTHING IN THIS FILE
	t.Run("IndexForInefficientlyComparable", func(t *testing.T) {
		t.Run("Returns index of found item", func(t *testing.T) {
			pm := mocks.NewIC[string](t)
			pm.On("EqualTo", pm1).Return(false).Once()
			pm.On("EqualTo", pm2).Return(true).Once()
			assert.Equal(t, 1, IndexForInefficientlyComparable[int, string](testPairs.Uncomparable, pm))
		})
		t.Run("Returns -1 for unfound index", func(t *testing.T) {
			pm := mocks.NewIC[string](t)
			pm.On("EqualTo", mock.Anything).Return(false)
			assert.Equal(t, -1, IndexForInefficientlyComparable[int, string](testPairs.Uncomparable, pm))
		})
	})
	t.Run("CombinePairs", func(t *testing.T) {
		t.Run("Empty input -> empty output", func(t *testing.T) {
			act := CombinePairs[int, string, string]()
			assert.Empty(t, act.Comparable)
			assert.Empty(t, act.Uncomparable)
		})
		t.Run("Works as intended", func(t *testing.T) {
			tp1 := Pairs[int, string, string]{
				Comparable: []Pair[int, string]{
					comparables[0],
					comparables[1],
				},
				Uncomparable: []Pair[int, inefficientComparable.IC[string]]{
					{pm1, -3},
					{pm2, -4},
				},
			}
			tp2 := Pairs[int, string, string]{
				Comparable: []Pair[int, string]{
					comparables[1],
					comparables[2],
				},
				Uncomparable: []Pair[int, inefficientComparable.IC[string]]{
					uncomparables[2],
					uncomparables[1],
				},
			}
			argMatcher := func(pma *mocks.IC[string]) func(interface{}) bool {
				return func(pmb interface{}) bool {
					return pma == pmb
				}
			}
			argMatcherOtherwise := func(pma *mocks.IC[string]) func(interface{}) bool {
				return func(pmb interface{}) bool {
					return pma != pmb
				}
			}
			pm1.On("EqualTo", mock.MatchedBy(argMatcher(pm1))).Return(true).Maybe()
			pm2.On("EqualTo", mock.MatchedBy(argMatcher(pm2))).Return(true).Maybe()
			pm3.On("EqualTo", mock.MatchedBy(argMatcher(pm3))).Return(true).Maybe()
			pm1.On("EqualTo", mock.MatchedBy(argMatcherOtherwise(pm1))).Return(false).Maybe()
			pm2.On("EqualTo", mock.MatchedBy(argMatcherOtherwise(pm2))).Return(false).Maybe()
			pm3.On("EqualTo", mock.MatchedBy(argMatcherOtherwise(pm3))).Return(false).Maybe()
			act := CombinePairs[int, string, string](tp1, tp2)
			assert.NotEmpty(t, act.Comparable)
			assert.NotEmpty(t, act.Uncomparable)
			// TODO: this (one case should cover all of the following)
			// TODO: same comparables in different pairs,
			// TODO: different comparables in different pairs
			// TODO: equivalent inComps in different pairs,
			// TODO: different inComps in different pairs
		})
	})
}
