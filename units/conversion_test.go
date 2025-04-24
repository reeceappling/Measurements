package units

import (
	"github.com/reeceappling/measurements/units/conv"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConversions(t *testing.T) {
	var nilConv *conversion = nil
	t.Run("deleteConversionFromNodePairBySigs", func(t *testing.T) {
		// TODO: this
	})
	t.Run("conversion", func(t *testing.T) {
		m := GetUnitBySignature("m")
		g := GetUnitBySignature("g")
		srVal := 37.38
		sr := conv.NewStaticRate(srVal)
		t.Run("rateFrom", func(t *testing.T) {
			var c *conversion = nil
			assert.Panics(t, func() {
				nilConv.rateFrom("")
			}) // TODO: string
			c = &conversion{
				conversionPair: newConversionPair(m, g),
				rate:           sr,
			}
			assert.Equal(t, sr, c.rateFrom(m.Signature()))                   // TODO: string
			assert.Equal(t, conv.RatePow(sr, -1), c.rateFrom(g.Signature())) // TODO: string
			assert.Panics(t, func() {
				c.rateFrom("neither")
			}, "panics if looking for a rate from a signature not on the conversion")
		})
		t.Run("deleteConversionFromNodePair", func(t *testing.T) {
			assert.NotPanics(t, func() {
				nilConv.deleteConversionFromNodePair()
			}) // TODO: string
			// TODO: this
		})

		t.Run("setConversionOnNodePair", func(t *testing.T) {
			assert.Panics(t, func() {
				_ = nilConv.setConversionOnNodePair()
			}) // TODO: string
			// TODO: this
		})
		t.Run("newDynamicConversion", func(t *testing.T) {
			// TODO: only if used
		})
	})

}
