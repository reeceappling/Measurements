package units

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnitWithExponent(t *testing.T) {
	t.Run("uweSignature", func(t *testing.T) {
		assert.Panics(t, func() {
			uweSignature(baseUnitWithExponent{
				unit: unitless,
				exp:  0,
			})
		}, "0-exponent should panic")
	})
	bua, _ := newBaseUnit("bua")
	t.Run("withoutNils", func(t *testing.T) {
		t.Run("baseUnitWithExponent", func(t *testing.T) {
			buweReal := baseUnitWithExponent{unit: bua, exp: 2}
			buweNil := baseUnitWithExponent{unit: Unitless(), exp: -2}
			assert.Equal(t, buweReal, buweReal.withoutNils(), "non-nil should be no-op")
			assert.Equal(t, baseUnitsGroup{}, buweNil.withoutNils(), "nil should return empty base group")
		})
		t.Run("stdUnitWithExponent", func(t *testing.T) {
			suweReal := baseUnitWithExponent{unit: bua, exp: 2}
			//suweNil := stdUnitWithExponent{unit: Unitless(), exp: -2} // TODO: see below
			assert.Equal(t, suweReal, suweReal.withoutNils(), "non-nil should be no-op")
			//assert.Equal(t, stdUnitsGroup{}, suweNil.withoutNils(), "nil should return empty base group") // TODO: ensure we will never have a case of a nil sUWE
		})
	})

	t.Run("baseUnitWithExponent", func(t *testing.T) {
		t.Run("pow", func(t *testing.T) {
			exp, pow := 2, -3
			powd := baseUnitWithExponent{unit: bua, exp: exp}.pow(pow)
			uwes := powd.asUWEs()
			assert.Equal(t, 1, len(uwes))
			assert.Equal(t, bua, uwes[0].Unit())
			assert.Equal(t, exp*pow, uwes[0].Exp())
		})
		// TODO: this
	})
	t.Run("", func(t *testing.T) {
		// TODO: this
	})

	// TODO: ADD MORE TESTS FOR UNITS WITH EXPONENTS
}
