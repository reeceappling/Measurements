package units

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuiltIn(t *testing.T) {
	t.Run("newMetricBaseUnits", func(t *testing.T) {
		// TODO: this
		assert.Panics(t, func() {
			_ = newMetricBaseUnits("m", "meter", "length")
		}, "fails when NewUnitWithEquivalence fails")
	})
	// TODO: this
}
