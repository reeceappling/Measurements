package conv

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestBaseRates(t *testing.T) {
	t.Run("StaticRate", func(t *testing.T) {
		t.Run("NewStaticRate", func(t *testing.T) {
			val := 37.38
			sr := NewStaticRate(val)
			assert.NotNil(t, sr)
			amt, err := sr.Resolve()
			assert.NoError(t, err)
			assert.Equal(t, val, amt)
		})
		// TODO: this
	})
	t.Run("DynamicRate", func(t *testing.T) {
		t.Run("NewDynamicRate", func(t *testing.T) {
			//val := 37.38 // TODO: delete
			rf := RatioFunc(func(opts *Options) (float64, error) {
				return (*opts)["a"].(float64), nil
			})
			dr := NewDynamicRate(&rf, RatioInputArgument{
				Name:         "a",
				ValueType:    reflect.TypeFor[float64](),
				DefaultValue: nil, // TODO: ok?
			})
			valB := -11.23
			amt, err := dr.Resolve(NewOptions().Add(OptionField("a", valB)))
			assert.NoError(t, err)
			assert.Equal(t, valB, amt)

		})
		// TODO: this
	})
	t.Run("ErroneousRate", func(t *testing.T) {
		// TODO: this
	})
}
