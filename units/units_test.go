package units

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnits(t *testing.T) {
	t.Run("unitSymbol", func(t *testing.T) {
		t.Run("isValid", func(t *testing.T) {
			t.Run("withDefaultOptions", func(t *testing.T) {
				// TODO: this
			})
			t.Run("with Extra Invalidating strings", func(t *testing.T) {
				oldInvalidSymbols := currentServiceOptions.extraInvalidUnitSymbols
				currentServiceOptions.extraInvalidUnitSymbols = []string{"q"}
				assert.False(t, UnitSymbol("q").isValid())
				currentServiceOptions.extraInvalidUnitSymbols = oldInvalidSymbols
				// TODO: this
			})
		})
		t.Run("asSignature", func(t *testing.T) {
			// TODO: this
		})
	})

	testUnitStr := "Test!"
	testMName := testUnitStr + "m"
	testM, errGlobal := NewBaseUnit(testMName)
	assert.NoError(t, errGlobal, "should create nonexistent base units")
	assert.True(t, testM.isBaseUnits())
	t.Run("Creating base units that already exist", func(t *testing.T) {
		sameTestUnit, err := NewBaseUnit(testMName)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errGraphAlreadyHas)
		assert.Equal(t, testM, sameTestUnit)
	})

	testCmName := testUnitStr + "cm"
	cmPerM := 100.0
	testCm, errGlobal := NewUnitWithEquivalence(cmPerM, testCmName, 1, testM)
	t.Run("Creating non-base units", func(t *testing.T) {
		assert.NoError(t, errGlobal)
		assert.Equal(t, UnitSignature(testCmName), testCm.Signature())
		assert.False(t, testCm.isBaseUnits())
		t.Run("Creating a non-base unit again", func(t *testing.T) {
			_, err := NewUnitWithEquivalence(cmPerM, testCmName, 1, testM)
			assert.Error(t, err)
			assert.ErrorIs(t, err, errUnitAlreadyExists, "error for existing units should be error already exists")
		})
		// Base Units are correct
		baseUs, created := testCm.baseUnits()
		assert.False(t, created)
		assert.True(t, baseUs.isBaseUnits())
		assert.Equal(t, testM, baseUs)
		// Conversions to base units are correct
		distOneWay, err := testCm.Per(testM).Resolve()
		assert.NoError(t, err)
		assert.Equal(t, cmPerM, distOneWay)
		distOtherWay, err := testM.Per(testCm).Resolve()
		assert.NoError(t, err)
		assert.Equal(t, 1/cmPerM, distOtherWay)
	})

	testPerMStr := testUnitStr + "perM"
	expPerM := testM.Pow(-1)
	assert.True(t, expPerM.isBaseUnits())
	testPerM, errGlobal := NewUnitWithEquivalence(1, testPerMStr, 1, expPerM)
	t.Run("Creating non-base units that equate to compound base units", func(t *testing.T) {
		assert.NoError(t, errGlobal)
		assert.False(t, testPerM.isBaseUnits())
		_, err := testPerM.Per(testM).Resolve()
		assert.Error(t, err)
		// TODO: test type of error
		rate, err := testPerM.Per(expPerM).Resolve()
		assert.NoError(t, err)
		assert.Equal(t, 1.0, rate)
	})

	testBarnName := testUnitStr + "barn"
	testCm2 := testCm.Pow(2)
	testActM2, created := testCm2.baseUnits()
	assert.False(t, created)
	cm2PerBarn := 1e-24
	testBarn, errGlobal := NewUnitWithEquivalence(1, testBarnName, cm2PerBarn, testCm2)
	t.Run("Creating non-base units that equate to compound non-base units", func(t *testing.T) {
		assert.NoError(t, errGlobal)
		testM2 := testM.Pow(2)
		assert.Equal(t, testM2, testActM2, "move this test, base units of squared units should match")
		_, err := testM.Per(testBarn).Resolve()
		assert.Error(t, err)
		assert.ErrorIs(t, err, errNoBaseConversion)
		actCm2PerM2, err := testCm2.Per(testM2).Resolve()
		assert.NoError(t, err)
		expCm2PerM2, err := testCm2.Per(testActM2).Resolve()
		assert.NoError(t, err)
		assert.Equal(t, expCm2PerM2, actCm2PerM2)
		assert.Equal(t, cmPerM*cmPerM, expCm2PerM2)
		expBarnPerCm2, err := testBarn.Per(testCm2).Resolve()
		assert.NoError(t, err)
		assert.Equal(t, cm2PerBarn, 1/expBarnPerCm2)
		expBarnPerM2 := expBarnPerCm2 * expCm2PerM2
		testConversionBetween(t, testBarn, testM2, expBarnPerM2)
	})
}

func testConversionBetween(t *testing.T, a, b Units, expectedAPerB float64) {
	aPerB, err := a.Per(b).Resolve() // a per b
	assert.NoError(t, err, "staticConversion rate should exist from base")
	bPerA, err := b.Per(a).Resolve() // b per a
	assert.NoError(t, err, "staticConversion rate should exist to base")
	assert.InDelta(t, expectedAPerB, aPerB, min(aPerB, expectedAPerB)*.000000001, "a per b")
	assert.InDelta(t, 1/aPerB, bPerA, min(aPerB, bPerA)*.000000001, "b per a")
	aPerB, err = b.ConversionRateTo(a).Resolve() // a per b
	assert.NoError(t, err, "staticConversion rate should exist from base")
	bPerA, err = a.ConversionRateTo(b).Resolve() // b per a
	assert.NoError(t, err, "staticConversion rate should exist from base")
	assert.InDelta(t, expectedAPerB, aPerB, min(aPerB, expectedAPerB)*.000000001, "Conversion rate from b to a")
	assert.InDelta(t, 1/aPerB, bPerA, min(aPerB, bPerA)*.000000001, "Conversion rate from a to b")

}
