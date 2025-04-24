package measurement

import (
	"errors"
	"github.com/reeceappling/measurements/units"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMeasurements(t *testing.T) {
	m := units.GetUnitBySignature("m")
	cm := units.GetUnitBySignature("cm")
	expSmVal := 37.28
	sm := NewStaticMeasurement(expSmVal, m)
	t.Run("NewStaticMeasurement()", func(t *testing.T) {
		assert.Equal(t, m, sm.units)
		actVal, err := sm.val.Resolve()
		assert.NoError(t, err)
		assert.Equal(t, expSmVal, actVal)
	})

	dynRateFunc := conv.RatioFunc(func(inp *conv.Options) (float64, error) {
		return 372.8, nil
	})
	dynRate := conv.NewDynamicRate(&dynRateFunc)
	basicDynamicMeasurement := NewDynamicMeasurement(dynRate, m)
	t.Run("NewDynamicMeasurement()", func(t *testing.T) {
		assert.Equal(t, m, basicDynamicMeasurement.units)
		expVal, err := dynRateFunc(nil)
		assert.NoError(t, err)
		actVal, err := basicDynamicMeasurement.val.Resolve()
		assert.NoError(t, err)
		assert.Equal(t, expVal, actVal)
	})

	t.Run("Measurement", func(t *testing.T) {
		t.Run("rate", func(t *testing.T) {
			assert.Equal(t, dynRate, basicDynamicMeasurement.rate())
		})
		t.Run("Value", func(t *testing.T) {
			act, err := basicDynamicMeasurement.Value()
			assert.NoError(t, err)
			exp, err := basicDynamicMeasurement.val.Resolve()
			assert.NoError(t, err)
			assert.Equal(t, exp, act)
		})
		errRate := conv.ErroneousRate{Err: errors.New("A TEST ERROR!")}
		errMmt := Measurement{
			val:   errRate,
			units: m,
		}
		t.Run("ConvertTo", func(t *testing.T) {
			t.Run("with existing error", func(t *testing.T) {
				cmErr := errMmt.ConvertTo(cm)
				_, isErrRate := cmErr.rate().(conv.ErroneousRate)
				assert.True(t, isErrRate)
			})
			t.Run("new error", func(t *testing.T) {
				// TODO: baseUnits don't match
			})
			t.Run("without error", func(t *testing.T) {
				smAsCm := sm.ConvertTo(cm)
				assert.Equal(t, cm, smAsCm.units)
				resultValue, err := smAsCm.Value(conv.NewOptions()) // Ensures options get passed through // TODO: kinda hacky for coverage
				assert.NoError(t, err)
				assert.Equal(t, expSmVal*100, resultValue)
			})
		})
		t.Run("Mul/Div", func(t *testing.T) {
			t.Run("with existing error", func(t *testing.T) {
				smcmVal := 11.23
				smcm := NewStaticMeasurement(smcmVal, cm)
				mErr := errMmt.Mul(smcm, sm)
				expUnits := errMmt.units.Mul(smcm.units, sm.units)
				assert.Equal(t, expUnits, mErr.units)
				_, isErrRate := mErr.rate().(conv.ErroneousRate)
				assert.True(t, isErrRate)
				mErr = errMmt.Div(sm)
				expUnits = errMmt.units.Div(sm.units)
				assert.Equal(t, expUnits, mErr.units)
				_, isErrRate = mErr.rate().(conv.ErroneousRate)
				assert.True(t, isErrRate)
			})
			t.Run("new error", func(t *testing.T) {
				// TODO: something that causes a new error
			})
			t.Run("without error", func(t *testing.T) {
				// TODO: this
			})
		})
		t.Run("Plus/Minus", func(t *testing.T) {
			smcmVal := 11.23
			smcm := NewStaticMeasurement(smcmVal, cm)
			t.Run("with existing error", func(t *testing.T) { // TODO: fixMe
				mErr := errMmt.Minus(smcm, sm, Measurement{
					val:   errRate,
					units: units.GetUnitBySignature("g"),
				})
				expUnits := errMmt.units
				assert.Equal(t, expUnits, mErr.units)
				_, isErrRate := mErr.rate().(conv.ErroneousRate)
				assert.True(t, isErrRate)
			})
			t.Run("new error", func(t *testing.T) {
				// TODO: something that causes a new error
			})
			t.Run("without error", func(t *testing.T) {
				sm2 := sm.Plus(sm)
				assert.Equal(t, errMmt.units, sm2.units)
				actVal, err := sm2.rate().Resolve()
				assert.NoError(t, err)
				assert.Equal(t, 2*expSmVal, actVal)
			})
		})
		t.Run("Pow", func(t *testing.T) {
			t.Run("with existing error", func(t *testing.T) {
				powBy := -3
				mErr := errMmt.Pow(powBy)
				assert.Equal(t, errMmt.units.Pow(powBy), mErr.units.Pow(powBy))
				_, isErrRate := mErr.rate().(conv.ErroneousRate)
				assert.True(t, isErrRate)
			})
			t.Run("new error", func(t *testing.T) {
				// TODO: something that causes a new error
			})
			t.Run("without error", func(t *testing.T) {
				// TODO: this
			})
		})
		// TODO: this
	})
}
