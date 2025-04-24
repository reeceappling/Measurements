package units

import (
	"github.com/reeceappling/measurements/units/conv"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompoundUnits(t *testing.T) {
	m, _ := allUnits.Get("m")
	s, _ := allUnits.Get("s")
	g, _ := allUnits.Get("g")
	kg, _ := allUnits.Get("kg")
	mps := m.Div(s)
	mps2 := m.Div(s.Pow(2))
	doz, _ := allUnits.Get(dozenSym)
	mol, _ := allUnits.Get(molSym)

	accGroup := mps2.asGroup()
	newCU := newCompoundUnit(accGroup)
	t.Run("newCompoundUnit", func(t *testing.T) {
		assert.Equal(t, 0, len(newCU.conversions))
		assert.Equal(t, accGroup, newCU.group)
	})

	t.Run("conversionMap fxns", func(t *testing.T) {
		// Covers get/set/has/remove
		nonexistentSig, newTestSig := UnitSignature("nonexistent"), UnitSignature("aNewConversion")
		assert.False(t, newCU.hasConversion(nonexistentSig), "nonexistent staticConversion should not have staticConversion")
		_, exists := newCU.getConversion(nonexistentSig)
		assert.False(t, exists, "nonexistent staticConversion should return false existence")
		newConv := &conversion{
			conversionPair: newConversionPair(newCU, kg),
			rate:           conv.NewStaticRate(2),
		}
		assert.NoError(t, newCU.setConversion(newTestSig, newConv))
		assert.True(t, newCU.hasConversion(newTestSig), "new staticConversion should exist")
		accConv, newExists := newCU.getConversion(newTestSig)
		assert.True(t, newExists, "newly set staticConversion should exist")
		assert.Equal(t, newConv, accConv, "newly set staticConversion should match")
		// TODO: test staticConversion overwriting (pass and fail)
		t.Run("setConversion when one already exists", func(t *testing.T) {
			err := newCU.setConversion(newTestSig, newConv)
			assert.Error(t, err)
			assert.ErrorIs(t, err, errCompoundConvExists)
			// TODO: when overwrite is available
		})
		newCU.removeConversion(newTestSig)
		assert.False(t, newCU.hasConversion(newTestSig), "newly removed staticConversion should no longer exist")

	})

	t.Run("Signature", func(t *testing.T) {
		assert.Equal(t, UnitSignature("m s^-2"), mps2.Signature())
	})

	t.Run("isBaseUnits", func(t *testing.T) {
		assert.True(t, mps2.isBaseUnits(), "with only base UWEs, should be considered base units")
		assert.False(t, kg.Div(m.Pow(2)).isBaseUnits(), "compound units containing non-base units should not be base units")
		assert.False(t, m.Mul(mol).isBaseUnits(), "compound units with count units are not considered base units")
	})

	t.Run("baseUnits", func(t *testing.T) {
		bu, created := mps2.baseUnits()
		assert.False(t, created) // TODO: ensure correct
		assert.Equal(t, mps2, bu, "all-base compound units should return themselves")
		bu, created = kg.Mul(m).Div(s.Pow(2)).baseUnits()
		assert.False(t, created) // TODO: ensure correct
		assert.Equal(t, g.Mul(m).Div(s.Pow(2)), bu, "nonbase compound units should properly resolve their base")
		bu, created = g.Div(mol).baseUnits()
		assert.False(t, created) // TODO: ensure correct
		assert.Equal(t, g, bu, "compound units with count units")
	})

	t.Run("asGroup", func(t *testing.T) {
		assert.Equal(t, mps2.(*compoundUnit).group, mps2.asGroup(), "should return internal group")
		assert.Equal(t, kg.Mul(m).Div(s.Pow(2)).asGroup(), s.Pow(-2).Mul(m.Mul(kg)).asGroup(), "keeps order consistent")
	})

	t.Run("Per/ConversionRateTo", func(t *testing.T) {
		t.Run("same rate returns multiplicative identity rate", func(t *testing.T) {
			cu := GetUnitBySignature("m s^-2")
			assert.Equal(t, conv.MultiplicativeIdentityRate(), cu.Per(cu))
		})
		t.Run("count units with base units and standard units", func(t *testing.T) {
			// TODO: count units to base units!
			dozgkg := doz.Mul(g).Mul(kg).Mul(m)
			gPerKg, err := dozgkg.Per(doz.Mul(kg).Mul(kg).Mul(m)).Resolve()
			assert.NoError(t, err)
			assert.Equal(t, 1000.0, gPerKg)
		})
		// TODO: this
	})

	t.Run("Mul", func(t *testing.T) {
		assert.Equal(t, mps2, mps2.Mul(unitless), "identity multiplication with unitless")
		assert.Equal(t, "g m s^-1", string(mps.Mul(g).Signature()), "with unrelated base")
		assert.Equal(t, mps2, mps.Mul(s.Pow(-1)), "with related base")
		kgmps := mps.Mul(kg)
		assert.Equal(t, "kg m s^-1", string(kgmps.Signature()), "with unrelated std")
		assert.Equal(t, "kg^2 m s^-1", string(kgmps.Mul(kg).Signature()), "with related std")
		assert.Equal(t, "kg m mol s^-1", string(kgmps.Mul(mol).Signature()), "with count")
	})

	t.Run("Div", func(t *testing.T) {
		assert.Equal(t, mps2, mps2.Div(unitless), "identity division with unitless")
		assert.Equal(t, unitless, mps2.Div(mps2), "self division results in unitless")
		assert.Equal(t, "g^-1 m s^-2", string(mps2.Div(g).Signature()), "with unrelated base")
		assert.Equal(t, mps2, mps.Div(s), "with related base")
		mpskg, mpskg2 := mps.Div(kg), mps.Div(kg.Pow(2))
		assert.Equal(t, "kg^-1 m s^-1", string(mpskg.Signature()), "with unrelated std")
		assert.Equal(t, "kg^-2 m s^-1", string(mpskg2.Signature()), "with related std")
		assert.Equal(t, "kg^-2 m mol^-1 s^-1", string(mpskg2.Div(mol).Signature()), "with related count")
	})

	t.Run("Pow", func(t *testing.T) {
		assert.Equal(t, unitless, mps2.Pow(0))
		assert.Equal(t, mps2.Mul(mps2).Mul(mps2), mps2.Pow(3))
		assert.Equal(t, unitless.Div(mps2).Div(mps2), mps2.Pow(-2))
	})
}
