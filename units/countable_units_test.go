package units

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCountableUnits(t *testing.T) {
	dozen, exists := allUnits.Get(dozenSym)
	assert.True(t, exists)
	mol, existsM := allUnits.Get(molSym)
	assert.True(t, existsM)
	t.Run("completely unitless items", func(t *testing.T) {
		assert.False(t, dozen.isBaseUnits())
		dozBase, created := dozen.baseUnits()
		assert.False(t, created)
		dozBaseGroup := dozBase.asGroup()
		assert.True(t, dozBase.isBaseUnits())
		assert.Equal(t, 0, len(dozBaseGroup.asUWEs()), "base group should be empty") // TODO: ensure this is how we want this to be
		nItems, err := dozBase.Per(dozen).Resolve()
		assert.NoError(t, err)
		assert.Equal(t, 12.0, nItems)
		nItems, err = unitless.Per(dozen).Resolve()
		assert.NoError(t, err)
		assert.Equal(t, 12.0, nItems)
		// TODO: base unitless item
		// TODO: hybrid base-standard unitless item
		// TODO: NEED TO DO MORE, MULTI-standard UNITLESS ITEMS, COVER ALL BASES
	})
	grams, existsG := allUnits.Get("g")
	assert.True(t, existsG)
	kg, existskg := allUnits.Get("kg")
	assert.True(t, existskg)
	t.Run("partially unitless", func(t *testing.T) {
		gpm := grams.Div(mol)
		assert.False(t, gpm.isBaseUnits())
		gpmBase, _ := gpm.baseUnits() // TODO: test created or no?
		assert.Equal(t, grams, gpmBase)
		testConversionBetween(t, grams, kg, 1000.0)
		n, err := gpm.Per(grams).Resolve()
		assert.NoError(t, err)
		assert.InDelta(t, AvagadrosNumber, n, AvagadrosNumber/10000000.0)
		n, err = grams.ConversionRateTo(gpm).Resolve()
		assert.NoError(t, err)
		assert.InDelta(t, AvagadrosNumber, n, AvagadrosNumber/10000000.0)
		dozMol := dozen.Mul(mol)
		n, err = mol.Per(dozMol).Resolve()
		assert.NoError(t, err)
		assert.Equal(t, 12.0, n)
		n, err = mol.ConversionRateTo(dozMol).Resolve()
		assert.NoError(t, err)
		assert.Equal(t, 1.0/12.0, n)
		n, err = dozMol.Per(mol).Resolve()
		assert.NoError(t, err)
		assert.Equal(t, 1.0/12.0, n)
		gPerDozMol := grams.Div(dozMol)
		n, err = gPerDozMol.Per(gpm).Resolve()
		assert.NoError(t, err)
		assert.Equal(t, 12.0, n)
		n, err = gpm.ConversionRateTo(gPerDozMol).Resolve()
		assert.NoError(t, err)
		assert.Equal(t, 12.0, n)
	})
}
