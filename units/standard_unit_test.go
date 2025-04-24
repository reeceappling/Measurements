package units

import (
	conv2 "github.com/reeceappling/measurements/units/conv"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStandardUnits(t *testing.T) {
	// standard units
	kg, _ := allUnits.Get("kg")
	km, _ := allUnits.Get("km")
	N, _ := allUnits.Get("N")
	// base units
	g, _ := allUnits.Get("g")
	m, _ := allUnits.Get("m")
	s, _ := allUnits.Get("s")

	t.Run("modifying conversions", func(t *testing.T) {
		conv, exists := kg.getConversion("g")
		assert.True(t, exists, "staticConversion from kg to g should exist")
		_, exists = kg.getConversion("km")
		assert.False(t, exists, "staticConversion from kg to km should not exist")
		assert.Error(t, kg.setConversion(g.Signature(), &conversion{
			conversionPair: newConversionPair(kg, g),
			rate:           conv2.NewStaticRate(-1),
		}), "without overwriting enabled, should error for existing conversions")
		newConvRatio := 3728.0
		newConversion := &conversion{
			conversionPair: newConversionPair(kg, km),
			rate:           conv2.NewStaticRate(newConvRatio),
		}
		assert.NoError(t, kg.setConversion(km.Signature(), newConversion), "nonexistent conversions should be settable")
		conv, exists = kg.getConversion("km")
		assert.True(t, exists, "recently set staticConversion should now exist")
		assert.Equal(t, newConversion, conv)
		kg.removeConversion("km")
		_, exists = kg.getConversion("km")
		assert.False(t, exists, "staticConversion should not exist post-deletion")
		t.Run("setConversion when allowing staticConversion rate overwriting", func(t *testing.T) {
			// TODO: this
		})
	})

	t.Run("newStdUnit", func(t *testing.T) {
		_, err := newStdUnit("bad", kg)
		assert.Error(t, err, "should error when provided with non-base unit(s)")
		_, err = newStdUnit("kg", kg)
		assert.Error(t, err, "should error when provided with existing unit symbol")
		newOkSym := UnitSymbol("newStdUnitTest")
		newUnit, err := newStdUnit(newOkSym, g)
		assert.NoError(t, err)
		assert.Equal(t, newOkSym, newUnit.sym, "symbol should match input")
		assert.Equal(t, g, newUnit.baseGroup, "base group should match input")
		assert.Equal(t, 0, len(newUnit.conversions), "conversions should initially have 0 entries")
	})

	t.Run("symbol/Signature", func(t *testing.T) {
		assert.Equal(t, kg.(*standardUnit).sym, kg.(*standardUnit).symbol(), "symbol should match")
		assert.Equal(t, UnitSignature("kg"), kg.Signature(), "Signature should be kg")
		assert.Equal(t, UnitSymbol(kg.Signature()), kg.(*standardUnit).sym, "unitSymbol should match Signature")
	})

	assert.False(t, kg.isBaseUnits(), "isBaseUnits should always be false")

	kgBu, created := kg.baseUnits()
	assert.False(t, created)
	NBu, _ := N.baseUnits()
	t.Run("baseUnits", func(t *testing.T) {
		assert.Equal(t, g, kgBu, "base units should always be the internal base units group")
		assert.Equal(t, kg.(*standardUnit).baseGroup, kgBu, "base units should always be the internal base units group")
		assert.Equal(t, N.(*standardUnit).baseGroup, NBu, "base units should always be the internal base units group, even for more complex standard units")
	})

	t.Run("asGroup", func(t *testing.T) {
		stdGroup := kg.asGroup()
		groupUwes := stdGroup.asUWEs()
		assert.Equal(t, 1, len(groupUwes), "should always have only 1 UWE")
		assert.Equal(t, 1, groupUwes[0].Exp(), "should always have exp of 1")
		assert.Equal(t, kg, groupUwes[0].Unit(), "UWE unit should match stdUnit")
		assert.Equal(t, kg.Signature(), stdGroup.Signature(), "group Signature should match stdUnit Signature")
		assert.Equal(t, kgBu.asGroup(), stdGroup.baseGroup(), "base groups should match")
	})

	t.Run("Per/ConversionRateTo", func(t *testing.T) {
		t.Run("unconvertable unit", func(t *testing.T) {
			_, err := kg.Per(km).Resolve()
			assert.Error(t, err)
			_, err = kg.ConversionRateTo(km).Resolve()
			assert.Error(t, err)
		})

		t.Run("convertable unit", func(t *testing.T) {
			t.Run("simple", func(t *testing.T) {
				convTo, err := kg.ConversionRateTo(g).Resolve()
				assert.NoError(t, err)
				assert.Equal(t, 1000.0, convTo)
				convPer, err := kg.Per(g).Resolve()
				assert.NoError(t, err)
				assert.Equal(t, 1.0/convTo, convPer)
			})

			t.Run("complex", func(t *testing.T) {
				convPer, err := N.Per(NBu).Resolve()
				assert.NoError(t, err)
				assert.Equal(t, 0.001, convPer)
				convTo, err := N.ConversionRateTo(NBu).Resolve()
				assert.NoError(t, err)
				assert.Equal(t, 1.0/convPer, convTo)
			})

		})

		t.Run("convertable count unit", func(t *testing.T) {
			// TODO: this
		})
	})

	t.Run("Mul/Div", func(t *testing.T) {
		t.Run("with baseUnit", func(t *testing.T) {
			t.Run("related", func(t *testing.T) {
				kgg, kgpg := kg.Mul(g), kg.Div(g)
				kggBu, _ := kgg.baseUnits()
				kgpgBu, _ := kgpg.baseUnits()
				assert.Equal(t, UnitSignature("g kg"), kgg.Signature())
				assert.Equal(t, UnitSignature("g^-1 kg"), kgpg.Signature())
				assert.Equal(t, UnitSignature("g^2"), kggBu.Signature())
				assert.Equal(t, unitless, kgpgBu)
			})

			t.Run("unrelated", func(t *testing.T) {
				kgm, kgpm := kg.Mul(m), kg.Div(m)
				assert.Equal(t, UnitSignature("kg m"), kgm.Signature())
				assert.Equal(t, UnitSignature("kg m^-1"), kgpm.Signature())
			})
		})

		t.Run("with compound unit", func(t *testing.T) {
			t.Run("related", func(t *testing.T) {
				kgpm := kg.Div(m)
				kg2pm, shouldBem := kg.Mul(kgpm), kg.Div(kgpm)
				kg2pmBu, _ := kg2pm.baseUnits()
				assert.Equal(t, UnitSignature("kg^2 m^-1"), kg2pm.Signature())
				assert.Equal(t, m, shouldBem)
				gpm := g.Div(m)
				kggpm, kgmpg := kg.Mul(gpm), kg.Div(gpm)
				kggpmBu, _ := kggpm.baseUnits()
				kgmpgBu, _ := kgmpg.baseUnits()
				assert.Equal(t, UnitSignature("g kg m^-1"), kggpm.Signature())
				assert.Equal(t, UnitSignature("g^-1 kg m"), kgmpg.Signature())
				assert.Equal(t, kg2pmBu, kggpmBu)
				assert.Equal(t, m, kgmpgBu)
			})

			t.Run("unrelated", func(t *testing.T) {
				mps := m.Div(s)
				kgmps, kgspm := kg.Mul(mps), kg.Div(mps)
				assert.Equal(t, UnitSignature("kg m s^-1"), kgmps.Signature())
				assert.Equal(t, UnitSignature("kg m^-1 s"), kgspm.Signature())
			})
		})

		t.Run("with stdUnit", func(t *testing.T) {
			t.Run("self", func(t *testing.T) {
				kg2, shouldBeUnitless := kg.Mul(kg), kg.Div(kg)
				kg2Bu, _ := kg2.baseUnits()
				assert.Equal(t, UnitSignature("kg^2"), kg2.Signature())
				assert.Equal(t, UnitSignature("g^2"), kg2Bu.Signature())
				assert.Equal(t, unitless, shouldBeUnitless)
			})

			t.Run("related", func(t *testing.T) {
				kgg, kgpg := kg.Mul(g), kg.Div(g)
				kgpgBu, _ := kgpg.baseUnits()
				assert.Equal(t, UnitSignature("g kg"), kgg.Signature())
				assert.Equal(t, UnitSignature("g^-1 kg"), kgpg.Signature())
				assert.Equal(t, unitless, kgpgBu)
			})

			t.Run("unrelated", func(t *testing.T) {
				kgkm, kgpkm := kg.Mul(km), kg.Div(km)
				kgkmBu, _ := kgkm.baseUnits()
				kgpkmBu, _ := kgpkm.baseUnits()
				assert.Equal(t, UnitSignature("kg km"), kgkm.Signature())
				assert.Equal(t, UnitSignature("kg km^-1"), kgpkm.Signature())
				assert.Equal(t, UnitSignature("g m"), kgkmBu.Signature())
				assert.Equal(t, UnitSignature("g m^-1"), kgpkmBu.Signature())
			})
		})

		t.Run("with count unit", func(t *testing.T) {
			// TODO: this
		})
	})

	t.Run("Pow", func(t *testing.T) {
		t.Run("positive", func(t *testing.T) {
			posExp := 3
			kgPos3 := kg.Pow(posExp)
			assert.Equal(t, 1, len(kgPos3.asGroup().asUWEs()))
			assert.Equal(t, posExp, kgPos3.asGroup().asUWEs()[0].Exp())
			assert.Equal(t, kg, kgPos3.asGroup().asUWEs()[0].Unit())
		})

		t.Run("negative", func(t *testing.T) {
			negExp := -2
			kgNeg2 := kg.Pow(negExp)
			assert.Equal(t, 1, len(kgNeg2.asGroup().asUWEs()))
			assert.Equal(t, negExp, kgNeg2.asGroup().asUWEs()[0].Exp())
			assert.Equal(t, kg, kgNeg2.asGroup().asUWEs()[0].Unit())
		})

		t.Run("zero", func(t *testing.T) {
			assert.Equal(t, unitless, kg.Pow(0))
		})
	})

	t.Run("asUnitWithExponent", func(t *testing.T) {
		kgAsUwe := kg.(*standardUnit).asUnitWithExponent()
		assert.Equal(t, 1, kgAsUwe.Exp())
		assert.Equal(t, kg, kgAsUwe.Unit())
	})
}
