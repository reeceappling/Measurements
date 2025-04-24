package conv

import (
	"errors"
	"fmt"
	"github.com/reeceappling/measurements/units/numericalPair/pairs"
	"github.com/reeceappling/measurements/units/numericalPair/power"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestRates(t *testing.T) {
	aStaticVal, bStaticVal, zStaticVal := 3.0, -4.0, 0.0
	aStaticRate, bStaticRate, zStaticRate := StaticRate(aStaticVal), StaticRate(bStaticVal), StaticRate(zStaticVal)
	aStatic, bStatic, zStatic := &aStaticRate, &bStaticRate, &zStaticRate

	var badDynamicRate *DynamicRate = nil
	okDynamicResolver := RatioFunc(func(inp *Options) (float64, error) {
		return aStaticVal, nil
	})
	okDynamicRate := &DynamicRate{resolver: &okDynamicResolver}

	t.Run("MultiplyRates", func(t *testing.T) {
		t.Run("Works in normal case (non-singular-1-exp)", func(t *testing.T) {
			cancellable := -2 // TODO: rename
			bExpA, bExpB := 7, 3
			pairGroup := []pairs.Pair[int, Rate]{
				{aStatic, -3},
				{bStatic, bExpA},
				{bStatic, bExpB},
			}
			crA := MultiplicativeRate{
				baseInputs:                []baseRate{pairGroup[0].Value.(baseRate), okDynamicRate},
				subRateUncomparableInputs: nil,
				Pairs: power.Pairs[int, Rate, compoundRate]{
					Comparable:   []pairs.Pair[int, Rate]{pairGroup[0], {okDynamicRate, cancellable}},
					Uncomparable: nil,
				},
			}
			crB := MultiplicativeRate{
				baseInputs:                []baseRate{pairGroup[1].Value.(baseRate), okDynamicRate},
				subRateUncomparableInputs: nil,
				Pairs: power.Pairs[int, Rate, compoundRate]{
					Comparable:   []pairs.Pair[int, Rate]{pairGroup[1], {okDynamicRate, 0 - cancellable}},
					Uncomparable: nil,
				},
			}
			crC := MultiplicativeRate{
				baseInputs:                []baseRate{pairGroup[2].Value.(baseRate)},
				subRateUncomparableInputs: nil,
				Pairs: power.Pairs[int, Rate, compoundRate]{
					Comparable:   []pairs.Pair[int, Rate]{pairGroup[2]},
					Uncomparable: nil,
				},
			}
			res := MultiplyRates(crA, crB, crC).AsMultiplicativeRate()
			assert.Equal(t, 3, len(res.Comparable))   // TODO: unsure
			assert.Equal(t, 0, len(res.Uncomparable)) // TODO: unsure
			assert.Contains(t, res.Comparable, pairGroup[0])
			aIndx, bIndx, dynIndx := 0, 0, 0
			for i := 0; i < len(res.Comparable); i++ {
				switch res.Comparable[i].Value {
				case aStatic:
					aIndx = i
				case bStatic:
					bIndx = i
				case okDynamicRate:
					dynIndx = i
				default:
					panic("should never hit this")
				}

			}
			assert.NotNil(t, aIndx, dynIndx) // TODO: more stuff and delete this line
			bResult := res.Comparable[bIndx]
			assert.Equal(t, bExpA+bExpB, bResult.Number)
			assert.Equal(t, bStatic, bResult.Value)
		})
		t.Run("Simplifies for singular results", func(t *testing.T) {
			crA := MultiplicativeRate{
				baseInputs:                rateInputs{aStatic},
				subRateUncomparableInputs: nil,
				Pairs:                     power.Pairs[int, Rate, compoundRate]{Comparable: []pairs.Pair[int, Rate]{{aStatic, 2}}},
			}
			crB := MultiplicativeRate{
				baseInputs:                rateInputs{aStatic},
				subRateUncomparableInputs: nil,
				Pairs:                     power.Pairs[int, Rate, compoundRate]{Comparable: []pairs.Pair[int, Rate]{{aStatic, -1}}},
			}
			assert.Equal(t, aStatic, MultiplyRates(crA, crB))
		})
	})
	// TODO: fix all below this
	t.Run("RatePow", func(t *testing.T) {
		inPowA, inPowB, exp := -3, 72, 8
		ps := power.Pairs[int, Rate, compoundRate]{
			Comparable: []pairs.Pair[int, Rate]{
				{aStatic, inPowA},
				{bStatic, inPowB},
			},
		}
		testRate := MultiplicativeRate{
			baseInputs:                rateInputs{aStatic, bStatic},
			subRateUncomparableInputs: nil,
			Pairs:                     ps,
		}
		t.Run("Zero case returns empty Compound rate", func(t *testing.T) {
			res := RatePow(testRate, 0).AsMultiplicativeRate()
			assert.Equal(t, 0, len(res.Comparable))
			assert.Equal(t, 0, len(res.Uncomparable))
		})
		t.Run("1 case returns input", func(t *testing.T) {
			assert.Equal(t, testRate, RatePow(testRate, 1))
		})
		t.Run("Default case returns pow-d rate as a compound rate", func(t *testing.T) {
			res := RatePow(testRate, exp).AsMultiplicativeRate()
			assert.Equal(t, len(ps.Comparable), len(res.Comparable))
			assert.Equal(t, len(ps.Uncomparable), len(res.Uncomparable))
			for i, inPow := range []int{inPowA, inPowB} {
				assert.Equal(t, ps.Comparable[i].Value, res.Comparable[i].Value)
				assert.Equal(t, exp*inPow, res.Comparable[i].Number)
			}
		})
		t.Run("Default case with an inverse-singular compound rate returns the internal rate of the compound", func(t *testing.T) {
			mr := MultiplicativeRate{
				baseInputs:                rateInputs{aStatic},
				subRateUncomparableInputs: nil,
				Pairs:                     power.Pairs[int, Rate, compoundRate]{Comparable: []pairs.Pair[int, Rate]{{aStatic, -1}}},
			}
			res := RatePow(mr, -1)
			assert.Equal(t, aStatic, res)
		})
	})
	t.Run("Static Rates", func(t *testing.T) {
		t.Run("Resolve", func(t *testing.T) {
			t.Run("Nil rate pointer error", func(t *testing.T) {
				var nilStaticRate *StaticRate = nil
				_, err := nilStaticRate.Resolve(nil)
				assert.Error(t, err)
				assert.ErrorIs(t, err, ErrNilRatePointer)
			})
			t.Run("valid", func(t *testing.T) {
				res, err := aStatic.Resolve(nil)
				assert.NoError(t, err)
				assert.Equal(t, aStaticVal, res)
			})
		})
		t.Run("AsMultiplicativeRate", func(t *testing.T) {
			cr := aStatic.AsMultiplicativeRate()
			assert.Equal(t, 1, len(cr.Comparable))
			assert.Equal(t, 0, len(cr.Uncomparable))
			assert.Equal(t, 1, cr.Comparable[0].Number)
			assert.Equal(t, aStatic, cr.Comparable[0].Value)
		})
	})

	dynamicFieldNameA, dynamicFieldNameB := "A", "B"
	dynamicFloatMultiplierErrTxt := func(fieldName string, dne bool) string {
		if dne {
			return fmt.Sprintf(`%s DID NOT EXIST`, fieldName)
		}
		return fmt.Sprintf(`%s NOT A FLOAT`, fieldName)
	}
	floatMultiplierResolver := RatioFunc(func(inp *Options) (float64, error) {
		fieldIntA, exists := (*inp)[dynamicFieldNameA]
		if !exists {
			return 0.0, errors.New(dynamicFloatMultiplierErrTxt(dynamicFieldNameA, false))
		}
		fieldIntB, exists := (*inp)[dynamicFieldNameB]
		if !exists {
			return 0.0, errors.New(dynamicFloatMultiplierErrTxt(dynamicFieldNameB, false))
		}
		a, ok := fieldIntA.(float64)
		if !ok {
			return 0.0, errors.New(dynamicFloatMultiplierErrTxt(dynamicFieldNameA, false))
		}
		b, ok := fieldIntB.(float64)
		if !ok {
			return 0.0, errors.New(dynamicFloatMultiplierErrTxt(dynamicFieldNameB, false))
		}
		return a * b, nil
	})
	dynamicFloatMultiplierRate := &DynamicRate{resolver: &floatMultiplierResolver}
	t.Run("Dynamic rates", func(t *testing.T) {
		t.Run("Resolve", func(t *testing.T) {
			t.Run("Nil rate pointer error", func(t *testing.T) {
				t.Run("Nil rate pointer error", func(t *testing.T) {
					_, err := badDynamicRate.Resolve(nil)
					assert.Error(t, err)
					assert.ErrorIs(t, err, ErrNilRatePointer)
				})
			})
			t.Run("Nil rate resolver error", func(t *testing.T) {
				nilRate := &DynamicRate{}
				_, err := nilRate.Resolve(nil)
				assert.Error(t, err)
				assert.ErrorIs(t, err, ErrNilRateResolver)
			})
			t.Run("resolver errors propagate up", func(t *testing.T) {
				testResolverErr := errors.New("test resolver error")
				resolver := RatioFunc(func(inp *Options) (float64, error) {
					return 0.0, testResolverErr
				})
				rate := &DynamicRate{resolver: &resolver}
				_, err := rate.Resolve(nil)
				assert.Error(t, err)
				assert.ErrorIs(t, err, testResolverErr)
			})
			t.Run("valid", func(t *testing.T) {
				t.Run("simple rate without options", func(t *testing.T) {
					res, err := okDynamicRate.Resolve(nil)
					assert.NoError(t, err)
					assert.Equal(t, aStaticVal, res)
				})
				t.Run("rate with options", func(t *testing.T) {
					valA, valB := 37.28, 1.123
					exp := valA * valB
					fields := []optionField{
						OptionField(dynamicFieldNameA, valA),
						OptionField(dynamicFieldNameB, valB),
					}

					optsBad := NewOptions().With(fields[0])
					_, err := dynamicFloatMultiplierRate.Resolve(optsBad)
					assert.Error(t, err, "should error without both fields")
					optsGood := optsBad.With(fields[1])
					res, err := dynamicFloatMultiplierRate.Resolve(optsGood)
					assert.NoError(t, err)
					assert.Equal(t, exp, res)

				})
			})
		})
		t.Run("AsMultiplicativeRate", func(t *testing.T) {
			cr := okDynamicRate.AsMultiplicativeRate()
			assert.Equal(t, 1, len(cr.Comparable))
			assert.Equal(t, 0, len(cr.Uncomparable))
			assert.Equal(t, 1, cr.Comparable[0].Number)
			assert.Equal(t, okDynamicRate, cr.Comparable[0].Value)
		})
	})
	t.Run("Compound rates", func(t *testing.T) {
		newCompoundRate := func(baserate Rate, exp int) MultiplicativeRate {
			return MultiplicativeRate{
				baseInputs:                rateInputs{baserate.(baseRate)},
				subRateUncomparableInputs: nil,
				Pairs:                     power.Pairs[int, Rate, compoundRate]{Comparable: []pairs.Pair[int, Rate]{{baserate, exp}}},
			}
		}
		t.Run("Resolve", func(t *testing.T) {
			t.Run("Error case", func(t *testing.T) {
				_, err := newCompoundRate(badDynamicRate, 1).Resolve(nil)
				assert.Error(t, err)
			})
			t.Run("zero returns early", func(t *testing.T) {
				res, err := newCompoundRate(zStatic, 1).Resolve(nil)
				assert.NoError(t, err)
				assert.Equal(t, 0.0, res)
			})
			t.Run("valid", func(t *testing.T) {
				aCRPow, bCRPow := -3, 7
				exp := math.Pow(aStaticVal, float64(aCRPow)) * math.Pow(bStaticVal, float64(bCRPow))
				res, err := MultiplicativeRate{
					baseInputs:                rateInputs{aStatic, bStatic},
					subRateUncomparableInputs: nil,
					Pairs: power.Pairs[int, Rate, compoundRate]{
						Comparable: []pairs.Pair[int, Rate]{
							{aStatic, aCRPow},
							{bStatic, bCRPow},
						},
					},
				}.Resolve(nil)
				assert.NoError(t, err)
				assert.NotContains(t, []float64{math.Inf(1), math.Inf(-1), 0.0, math.NaN()}, exp)
				assert.Equal(t, exp, res)

			})
		})
		t.Run("AsMultiplicativeRate", func(t *testing.T) {
			cr := MultiplicativeRate{
				baseInputs:                rateInputs{aStatic, bStatic},
				subRateUncomparableInputs: nil,
				Pairs: power.Pairs[int, Rate, compoundRate]{
					Comparable: []pairs.Pair[int, Rate]{
						{aStatic, -3},
						{bStatic, 7},
					},
				},
			}
			assert.Equal(t, cr, cr.AsMultiplicativeRate())
		})
	})
}
