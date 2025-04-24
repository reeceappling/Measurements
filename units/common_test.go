package units

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

//	func TestMain(m *testing.M) {
//			// TODO: maybe this
//	}
func TestCommon(t *testing.T) {
	m := GetUnitBySignature("m")
	t.Run("NewBaseUnit", func(t *testing.T) {
		t.Run("bad base unit symbol", func(t *testing.T) {
			// TODO: this
		})
		t.Run("tryAddNode errs", func(t *testing.T) {
			t.Run("Does exist", func(t *testing.T) {
				t.Run("errorOnDuplicateNewBaseUnit", func(t *testing.T) {
					// TODO: this
				})
				t.Run("no errorOnDuplicateNewBaseUnit", func(t *testing.T) {
					// TODO: this
				})
			})
			t.Run("Does not exist", func(t *testing.T) {
				// TODO: this
			})
		})
		t.Run("works as intended", func(t *testing.T) {
			// TODO: this
		})
	})
	t.Run("NewUnitWithEquivalence", func(t *testing.T) {
		// TODO: write all other test cases
		t.Run("fails for zeroes", func(t *testing.T) {
			t.Run("first zero", func(t *testing.T) {
				_, err := NewUnitWithEquivalence(0.0, "tempNewUnit", 1.0, m)
				assert.Error(t, err)
				assert.ErrorIs(t, err, errInvalidNewEquivalentUnitZeros)
				// TODO: cleanup
			})
			t.Run("second zero", func(t *testing.T) {
				// TODO: this
			})
		})
		t.Run("fails for invalid symbols", func(t *testing.T) {
			_, err := NewUnitWithEquivalence(1.0, "^bad Symbols^", 1.0, m)
			assert.Error(t, err)
			assert.ErrorIs(t, err, errInvalidSymbolChars)
		})
		t.Run("works as intended", func(t *testing.T) {
			// TODO: this
		})
	})
	t.Run("mulUnits", func(t *testing.T) {
		assert.Equal(t, unitless, mulUnits([]Units{}), "returns unitless for empty input")
		// TODO: other test?
		t.Run("works as intended", func(t *testing.T) {
			// TODO: this
		})
	})
	t.Run("appendUnitSliceToReciever", func(t *testing.T) {
		// TODO: this
	})

	t.Run("unitsForGroupCreateWithBaseIfNonexistent", func(t *testing.T) {
		// TODO: all test cases in here
	})
	t.Run("resolveTwoStepConversion", func(t *testing.T) {
		t.Run("impossible conversion", func(t *testing.T) {
			// TODO: this
		})
		// TODO: all other test cases
	})
	t.Run("GetUnitBySignature", func(t *testing.T) {
		u := GetUnitBySignature("m")
		assert.NotNil(t, u)
		// TODO: invalid signatures?
		assert.Panics(t, func() {
			_ = GetUnitBySignature("aNonExistentUnit")
		}, "panics for nonexistent signatures")
	})
}
