package units

import (
	"github.com/reeceappling/measurements/units/conv"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBaseUnits(t *testing.T) {
	fakeBaseUnitA, _ := newBaseUnit("fakeBaseUnitA")
	fakeBaseUnitB, _ := newBaseUnit("fakeBaseUnitB")
	allUnits.addNode(fakeBaseUnitA, fakeBaseUnitA.Signature()) // TODO: might mess up some tests?
	allUnits.addNode(fakeBaseUnitB, fakeBaseUnitB.Signature()) // TODO: might mess up some tests?

	t.Run("newBaseUnit()", func(t *testing.T) {
		_, err := newBaseUnit("^badSig^")
		assert.Error(t, err)
		assert.ErrorIs(t, err, errInvalidSymbolChars)
		fakeBaseUnitA, err = newBaseUnit("fakeBaseUnitA")
		assert.NoError(t, err, "creating baseUnit with new string works fine")
	})

	t.Run("NewBaseUnit()", func(t *testing.T) {
		_, err := NewBaseUnit("^badSig^")
		assert.Error(t, err)
		assert.ErrorIs(t, err, errBadNewBaseUnitSymbol)
		nbuA, err := NewBaseUnit("NewBaseUnitTestUnitA")
		assert.NoError(t, err)
		t.Run("option for not erroring on duplicate base unit", func(t *testing.T) {
			oldDuplicateNewBaseOpt := currentServiceOptions.errorOnDuplicateNewBaseUnit
			currentServiceOptions.errorOnDuplicateNewBaseUnit = false
			act, err := NewBaseUnit("NewBaseUnitTestUnitA")
			assert.NoError(t, err)
			assert.Equal(t, nbuA, act)
			currentServiceOptions.errorOnDuplicateNewBaseUnit = oldDuplicateNewBaseOpt
		})
	})

	// Grab necessary base units
	cm, _ := allUnits.Get("cm")
	m, _ := allUnits.Get("m")
	g, _ := allUnits.Get("g")
	s, _ := allUnits.Get("s")
	// Grab necessary non-base units
	kg, _ := allUnits.Get("kg")
	N, _ := allUnits.Get("N")
	doz, _ := allUnits.Get(dozenSym)
	mol, _ := allUnits.Get(molSym)
	t.Run("getConversion", func(t *testing.T) {
		t.Run("nonexistent returns false", func(t *testing.T) {
			_, exists := m.getConversion("a bad Signature")
			assert.False(t, exists)
		})
		t.Run("existent returns true and staticConversion", func(t *testing.T) {
			convers, exists := m.getConversion(cm.Signature())
			assert.True(t, exists)
			assert.NotNil(t, convers)
			ratio, err := convers.rateFrom(cm.Signature()).Resolve(nil)
			assert.NoError(t, err)
			assert.Equal(t, 100.0, ratio)
		})
	})

	t.Run("modifying conversions", func(t *testing.T) {
		testConv := &conversion{
			conversionPair: newConversionPair(fakeBaseUnitA, fakeBaseUnitB),
			rate:           conv.NewStaticRate(37.28),
		}
		sigTo := fakeBaseUnitB.Signature()
		t.Run("setConversion", func(t *testing.T) {
			_, exists := fakeBaseUnitA.getConversion(sigTo)
			assert.False(t, exists, "should not exist")
			assert.NoError(t, fakeBaseUnitA.setConversion(sigTo, testConv))
			_, exists = fakeBaseUnitA.getConversion(sigTo)
			assert.True(t, exists, "setting conversions that did not exist works")
			t.Run("without overwrite active (default)", func(t *testing.T) {
				_, exists = fakeBaseUnitA.getConversion(sigTo)
				assert.True(t, exists, "should still exist")
				assert.Error(t, fakeBaseUnitA.setConversion(sigTo, testConv), "should fail to overwrite")
			})
			t.Run("with overwrite active", func(t *testing.T) {
				_, exists = fakeBaseUnitA.getConversion(sigTo)
				assert.True(t, exists, "should still exist")
				currentServiceOptions.allowOverWritingConversionRate = true // change options to ensure overwrite is allowed
				assert.NoError(t, fakeBaseUnitA.setConversion(sigTo, testConv), "should successFully overwrite")
				resetCurrentOptionsToBuiltIn()
			})
		})
		t.Run("removeConversion", func(t *testing.T) {
			_, exists := fakeBaseUnitA.getConversion(sigTo)
			assert.True(t, exists, "should still exist")
			fakeBaseUnitA.removeConversion(sigTo)
			_, exists = fakeBaseUnitA.getConversion(sigTo)
			assert.False(t, exists, "should no longer exist")
		})
	})

	t.Run("asGroup", func(t *testing.T) {
		ag := fakeBaseUnitA.asGroup()
		assert.True(t, ag.isBaseGroup())
		assert.Equal(t, ag, ag.baseGroup())
		assert.Equal(t, ag, ag.withoutNils())
		assert.Equal(t, UnitSignature("fakeBaseUnitA"), ag.Signature())
		uwes := ag.asUWEs()
		assert.Equal(t, 1, len(uwes))
		uwe := uwes[0]
		assert.Equal(t, 1, uwe.Exp())
		assert.Equal(t, fakeBaseUnitA, uwe.Unit())
		nPow := 3
		powd := ag.pow(nPow).asUWEs()
		assert.Equal(t, len(uwes), len(powd))
		powUwe := powd[0]
		assert.Equal(t, nPow, powUwe.Exp())
		assert.Equal(t, fakeBaseUnitA, powd[0].Unit())
		t.Run("unitless", func(t *testing.T) {
			unitlessGroup := unitless.asGroup()
			assert.True(t, unitlessGroup.isBaseGroup())
			assert.Equal(t, 0, len(unitlessGroup.asUWEs()))
		})
	})

	t.Run("asBaseUnitWithExponent", func(t *testing.T) {
		buwe := fakeBaseUnitA.asBaseUnitWithExponent()
		assert.Equal(t, 1, buwe.Exp())
		assert.True(t, buwe.isBaseGroup())
		assert.Equal(t, fakeBaseUnitA, buwe.Unit())
		assert.Equal(t, fakeBaseUnitA.Signature(), buwe.Signature())
		bg := buwe.baseGroup()
		assert.Equal(t, fakeBaseUnitA.Signature(), bg.Signature())
		assert.Equal(t, 1, len(bg.asUWEs()))
		assert.Equal(t, buwe, bg.asUWEs()[0])
		t.Run("unitless", func(t *testing.T) {
			unitlessUwe := unitless.asBaseUnitWithExponent()
			assert.Equal(t, unitless, unitlessUwe.unit)
			assert.Equal(t, 1, unitlessUwe.exp)
		})
	})

	t.Run("asUnitWithExponent", func(t *testing.T) {
		// Covered by tests for asBaseUnitWithExponent
	})

	t.Run("Per/ConversionRateTo", func(t *testing.T) {
		t.Run("to self", func(t *testing.T) {
			rate, err := fakeBaseUnitA.Per(fakeBaseUnitA).Resolve()
			assert.NoError(t, err)
			assert.Equal(t, 1.0, rate)
			rate, err = unitless.Per(unitless).Resolve()
			assert.NoError(t, err)
			assert.Equal(t, 1.0, rate)
		})
		t.Run("to another base unit (shouldn't work)", func(t *testing.T) {
			_, err := fakeBaseUnitA.Per(fakeBaseUnitB).Resolve()
			assert.Error(t, err)
			_, err = fakeBaseUnitA.ConversionRateTo(fakeBaseUnitB).Resolve()
			assert.Error(t, err)
		})
		t.Run("to a standard unit (works sometimes)", func(t *testing.T) {
			mPerCm, err := m.Per(cm).Resolve()
			assert.NoError(t, err)
			assert.Equal(t, 0.01, mPerCm)
			cmPerM, err := m.ConversionRateTo(cm).Resolve()
			assert.NoError(t, err)
			assert.Equal(t, 100.0, cmPerM)
			_, err = m.Per(kg).Resolve()
			assert.Error(t, err)
			_, err = m.ConversionRateTo(kg).Resolve()
			assert.Error(t, err)
			t.Run("unitless", func(t *testing.T) {
				amt, err := unitless.Per(doz).Resolve()
				assert.NoError(t, err)
				assert.Equal(t, 12.0, amt)
			})
		})
		t.Run("to a compound unit", func(t *testing.T) {
			t.Run("to invalid unit", func(t *testing.T) {
				_, err := m.Per(kg.Mul(fakeBaseUnitB)).Resolve()
				assert.Error(t, err)
				_, err = m.ConversionRateTo(kg.Mul(fakeBaseUnitB)).Resolve()
				assert.Error(t, err)
			})
			t.Run("to valid compound unit", func(t *testing.T) {
				A, _ := allUnits.Get("A")
				V, _ := allUnits.Get("V")
				ohm, _ := allUnits.Get("ohm")
				AEquiv := V.Div(ohm)
				a, err := A.Per(AEquiv).Resolve()
				assert.NoError(t, err)
				assert.Equal(t, 1.0, a)
				a, err = A.ConversionRateTo(AEquiv).Resolve()
				assert.NoError(t, err)
				assert.Equal(t, 1.0, a)
				mV, err := NewUnitWithEquivalence(1000, "mV", 1, V) // TODO: build mV in built_in
				assert.NoError(t, err)
				milliAmp := mV.Div(ohm)
				a, err = A.Per(milliAmp).Resolve()
				assert.NoError(t, err)
				assert.Equal(t, 0.001, a)
				a, err = A.ConversionRateTo(milliAmp).Resolve()
				assert.NoError(t, err)
				assert.Equal(t, 1000.0, a)
			})
		})
	})

	t.Run("Mul/Div", func(t *testing.T) {
		t.Run("unitless", func(t *testing.T) {
			// With itself
			assert.Equal(t, unitless, unitless.Mul(unitless))
			assert.Equal(t, unitless, unitless.Div(unitless))
			// With a base
			assert.Equal(t, fakeBaseUnitA, unitless.Mul(fakeBaseUnitA)) // TODO: compoundUnit result is not simplifying!
			assert.Equal(t, fakeBaseUnitA.Pow(-1), unitless.Div(fakeBaseUnitA))
			// with a std unit
			assert.Equal(t, kg, unitless.Mul(kg))
			assert.Equal(t, kg.Pow(-1), unitless.Div(kg))
			// with a compound unit (incl another unitless)
			kgmPerMol := kg.Mul(m).Div(mol)
			assert.Equal(t, kgmPerMol, unitless.Mul(kgmPerMol))
			a := kgmPerMol.Pow(-1) // TODO: failing here
			b := unitless.Div(kgmPerMol)
			assert.Equal(t, a, b) // TODO: failing here due to nil receiver?
		})
		t.Run("with another base", func(t *testing.T) {
			assert.Equal(t, fakeBaseUnitA.Mul(fakeBaseUnitB), fakeBaseUnitB.Mul(fakeBaseUnitA))
			assert.Equal(t, fakeBaseUnitA.Mul(fakeBaseUnitA), fakeBaseUnitA.Pow(2))
			assert.Equal(t, unitless, fakeBaseUnitA.Div(fakeBaseUnitA), "self-division results in unitless")
			assert.Equal(t, fakeBaseUnitA.Div(fakeBaseUnitB), fakeBaseUnitA.Mul(fakeBaseUnitB.Pow(-1)), "division by another is equal to multiplying by the inverse of that")
		})
		t.Run("with a standard unit", func(t *testing.T) {
			// unrelated units (hz + m)
			hz, _ := allUnits.Get("hz")
			mhz, mphz := m.Mul(hz), m.Div(hz)
			assert.Equal(t, "hz m", string(mhz.Signature()), "that is unrelated")
			assert.Equal(t, "hz^-1 m", string(mphz.Signature()))
			assert.Equal(t, m, mphz.Mul(hz))
			assert.Equal(t, m, mhz.Div(hz))
			// related units (s + hz)
			shz := s.Mul(hz)
			sphz := s.Div(hz)
			assert.Equal(t, UnitSignature("hz s"), shz.Signature())
			assert.Equal(t, UnitSignature("hz^-1 s"), sphz.Signature())
			shzBu, created := shz.baseUnits()
			assert.False(t, created)
			assert.Equal(t, unitless, shzBu)
			sphzBu, created := sphz.baseUnits()
			assert.False(t, created)
			assert.Equal(t, UnitSignature("s^2"), sphzBu.Signature())
			amt, err := shz.Per(unitless).Resolve()
			assert.NoError(t, err)
			assert.Equal(t, 1.0, amt)
			// units that don't cancel out but are related: g -> N
			gN, gpN := g.Mul(N), g.Div(N)
			assert.Equal(t, UnitSignature("N g"), gN.Signature())
			assert.Equal(t, UnitSignature("N^-1 g"), gpN.Signature())

		})
		t.Run("with a compound unit (no overlap)", func(t *testing.T) {
			mps, reciprocal := m.Div(s), s.Div(m)
			gmPerSec := g.Mul(mps)
			assert.Equal(t, UnitSignature("g m s^-1"), gmPerSec.Signature())
			gspm := g.Div(mps)
			assert.Equal(t, UnitSignature("g m^-1 s"), gspm.Signature())
			assert.Equal(t, gspm, g.Mul(reciprocal), "multiplying by the reciprocal is equal to division")
			assert.Equal(t, gmPerSec, g.Div(reciprocal), "dividing by the reciprocal is equal to multiplication")
		})
		t.Run("with a compound unit (with overlapping units)", func(t *testing.T) {
			// with nested same-base units (g + g/m)
			gpm := g.Div(m)
			shouldBeM, recipSame := g.Div(gpm), g.Mul(gpm)
			sbmBu, sbmBuCreated := shouldBeM.baseUnits()
			recipSameBu, recipSameBuCreated := recipSame.baseUnits()
			assert.False(t, sbmBuCreated)       // TODO: ensure correct?
			assert.False(t, recipSameBuCreated) // TODO: ensure correct?
			assert.Equal(t, UnitSignature("g^2 m^-1"), recipSame.Signature())
			assert.Equal(t, m, shouldBeM)
			assert.Equal(t, m.Signature(), shouldBeM.Signature())
			// With same-class units (g + kg/m)
			kgpm, reciprocal := kg.Div(m), m.Div(kg)
			gkgpm, gmpkg := g.Mul(kgpm), g.Div(kgpm)
			gmpkgBu, gmpkgBuCreated := gmpkg.baseUnits()
			gkgpmBu, gkgpmBuCreated := gkgpm.baseUnits()
			assert.False(t, gmpkgBuCreated) // TODO: ensure correct?
			assert.False(t, gkgpmBuCreated) // TODO: ensure correct?
			assert.Equal(t, UnitSignature("g kg m^-1"), gkgpm.Signature())
			assert.Equal(t, UnitSignature("g kg^-1 m"), gmpkg.Signature())
			assert.Equal(t, gmpkg, g.Mul(reciprocal))
			assert.Equal(t, gkgpm, g.Div(reciprocal))
			assert.Equal(t, recipSameBu.Signature(), gkgpmBu.Signature(), "base units should calculate correctly for Mul/Div")
			assert.Equal(t, sbmBu.Signature(), gmpkgBu.Signature(), "base units should calculate correctly for Div/Mul")
		})
		t.Run("with count units", func(t *testing.T) {
			// TODO: this
		})
	})

	t.Run("Pow", func(t *testing.T) {
		posN, negN := 2, -3
		// Positive exponent
		posPow := fakeBaseUnitA.Pow(posN).asGroup().asUWEs()
		assert.Equal(t, 1, len(posPow), "positive result size")
		assert.Equal(t, posN, posPow[0].Exp(), "positive exponent")
		assert.Equal(t, fakeBaseUnitA, posPow[0].Unit(), "positive single unit")
		// Negative exponent
		negPow := fakeBaseUnitA.Pow(negN).asGroup().asUWEs()
		assert.Equal(t, 1, len(negPow), "negative result size")
		assert.Equal(t, negN, negPow[0].Exp(), "negative exponent")
		assert.Equal(t, fakeBaseUnitA, negPow[0].Unit(), "negative single unit")
		// Zero exponent
		assert.Equal(t, unitless, fakeBaseUnitA.Pow(0), "zero power")
		// Unitless
		assert.Equal(t, unitless, unitless.Pow(0))
		assert.Equal(t, unitless, unitless.Pow(posN))
		assert.Equal(t, unitless, unitless.Pow(negN))
	})
}
