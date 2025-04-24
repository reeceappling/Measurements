package units

import (
	"fmt"
	conv2 "github.com/reeceappling/measurements/units/conv"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestUnitsGroups(t *testing.T) {
	// Grab necessary base units
	g, _ := allUnits.Get("g")
	// Grab necessary non-base units
	km, _ := allUnits.Get("km")

	baseUnitA, _ := newBaseUnit("bA")
	baseUnitB, _ := newBaseUnit("bB")
	baseUnitC, _ := newBaseUnit("bC")
	testStdUnit := func(symb string, baseU *baseUnit) *standardUnit {
		stu := &standardUnit{
			sym:         UnitSymbol(symb),
			baseGroup:   baseU,
			conversions: map[UnitSignature]*conversion{},
			RWMutex:     sync.RWMutex{},
		}
		conv := &conversion{
			conversionPair: newConversionPair(stu, baseU),
			rate:           conv2.NewStaticRate(1),
		}
		if err := conv.setConversionOnNodePair(); err != nil {
			panic(err.Error())
		}

		return stu
	}
	stdUnitA := testStdUnit("sA", baseUnitA)
	stdUnitB := testStdUnit("sB", baseUnitB)

	t.Run("unitsGroupRateToBase", func(t *testing.T) {
		// TODO: this
	})
	t.Run("workingUnitsGroup", func(t *testing.T) {
		t.Run("asProperGroup", func(t *testing.T) {
			// TODO: thiss
		})
		t.Run("simplified", func(t *testing.T) {
			// TODO: thiss
		})
		t.Run("alphabetized", func(t *testing.T) {
			wugInit := workingUnitsGroup([]unitWithExponent{
				stdUnitWithExponent{
					unit: km.(singleUnit),
					exp:  -1,
				},
				baseUnitWithExponent{
					unit: g.(*baseUnit),
					exp:  1,
				},
			})
			alph := workingUnitsGroup([]unitWithExponent{
				stdUnitWithExponent{
					unit: km.(singleUnit),
					exp:  -1,
				},
				baseUnitWithExponent{
					unit: g.(*baseUnit),
					exp:  1,
				},
			}).alphabetized()
			assert.Equal(t, wugInit[0], alph[1])
			assert.Equal(t, wugInit[1], alph[0])
		})
		t.Run("resolveGroupType", func(t *testing.T) {
			// TODO: this
		})
	})

	spow := 2
	sug := stdUnitsGroup{
		baseUnitA.asUnitWithExponent(),
		stdUnitA.asUnitWithExponent(),
		stdUnitB.Pow(spow).asGroup().asUWEs()[0],
	}
	t.Run("stdUnitsGroup", func(t *testing.T) {
		assert.False(t, sug.isBaseGroup())
		t.Run("pow", func(t *testing.T) {
			uwePow := -1
			sp3 := sug.pow(uwePow)
			assert.False(t, sp3.isBaseGroup())
			sugUwes, sp3Uwes := sug.asUWEs(), sp3.asUWEs()
			assert.Equal(t, len(sug), len(sugUwes), "asUWEs should result in correct number of UWEs")
			assert.Equal(t, len(sugUwes), len(sp3Uwes), "asUWEs")
			for i, buwe := range sugUwes {
				assert.Equal(t, buwe.Unit().Signature(), sp3Uwes[i].Unit().Signature(), "sub-unit signatures should match")
				assert.Equal(t, buwe.Exp()*uwePow, sp3Uwes[i].Exp())
			}
			// TODO: this
		})
		t.Run("isBaseGroup", func(t *testing.T) {
			// TODO: this
		})
		t.Run("asUWEs", func(t *testing.T) {
			// TODO: this
		})
		t.Run("baseGroup", func(t *testing.T) {
			// TODO: this
		})
		t.Run("Signature", func(t *testing.T) {
			// TODO: this
			assert.Equal(t, UnitSignature(""), stdUnitsGroup{unitless.asUnitWithExponent()}.Signature(), "only unitless, or empty, should return empty string")
			assert.Equal(t, UnitSignature(""), stdUnitsGroup{}.Signature(), "only unitless, or empty, should return empty string")
		})
		t.Run("withoutNils", func(t *testing.T) {
			uweA := stdUnitWithExponent{unit: stdUnitA, exp: -3}
			uweB := baseUnitWithExponent{unit: baseUnitB, exp: 7}
			uweC := stdUnitWithExponent{unit: Unitless()}
			exp := stdUnitsGroup{uweA, uweB}
			for testName, inp := range map[string]unitsGroup{
				"properly removes from beginning, middle, and end": stdUnitsGroup{uweC, uweA, uweC, uweB, uweC},
				"does nothing when no unitless entries exist":      exp,
			} {
				assert.Equal(t, exp, inp.withoutNils(), testName)
			}
		})
	})

	sbuPwr := 3
	bug, power := baseUnitsGroup{
		baseUnitA.asBaseUnitWithExponent(),
		baseUnitWithExponent{unit: baseUnitB, exp: sbuPwr},
	}, -2
	t.Run("baseUnitsGroup", func(t *testing.T) {

		assert.True(t, bug.isBaseGroup(), "should be base group")
		bp3 := bug.pow(power)
		assert.True(t, bp3.isBaseGroup(), "pow'd base group should still be a base group")

		t.Run("pow", func(t *testing.T) {
			assert.True(t, bp3.isBaseGroup())
			bugUwes, bp3Uwes := bug.asUWEs(), bp3.asUWEs()
			assert.Equal(t, len(bug), len(bugUwes), "asUWEs should result in correct number of UWEs")
			assert.Equal(t, len(bugUwes), len(bp3Uwes), "asUWEs")
			for i, buwe := range bugUwes {
				assert.Equal(t, buwe.Unit().Signature(), bp3Uwes[i].Unit().Signature(), "sub-unit signatures should match")
				assert.Equal(t, buwe.Exp()*power, bp3Uwes[i].Exp())
			}
		})
		t.Run("baseGroup", func(t *testing.T) {
			compareGroups(t, bug, bug.baseGroup())
		})
		t.Run("Signature", func(t *testing.T) {
			assert.Equal(t, bug.Signature(), UnitSignature(fmt.Sprintf(`%s %s`, bug[0].Signature(), bug[1].Signature())))
		})
		t.Run("withoutNils", func(t *testing.T) {
			uweA := baseUnitWithExponent{unit: baseUnitA, exp: 2}
			uweB := baseUnitWithExponent{unit: baseUnitB, exp: -3}
			uweC := baseUnitWithExponent{unit: Unitless()}
			exp := baseUnitsGroup{uweA, uweB}
			for testName, inp := range map[string]unitsGroup{
				"properly removes from beginning, middle, and end": baseUnitsGroup{uweC, uweA, uweC, uweB, uweC},
				"does nothing when no unitless entries exist":      exp,
			} {
				assert.Equal(t, exp, inp.withoutNils(), testName)
			}
		})
	})

	bugB := baseUnitsGroup{baseUnitB.asBaseUnitWithExponent(), baseUnitC.asBaseUnitWithExponent()}
	t.Run("combineGroups", func(t *testing.T) {
		t.Run("combining all base groups", func(t *testing.T) {
			// bA, bB, bC
			combined := combineGroups(bug, bugB)
			assert.True(t, combined.isBaseGroup())
			assert.Equal(t, 3, len(combined.asUWEs()))
		})
		t.Run("combining mixed isBase groups", func(t *testing.T) {
			// bA, bB, bC, sA, sB == 5
			combined := combineGroups(sug, bug, bugB)
			assert.False(t, combined.isBaseGroup())
			assert.Equal(t, 5, len(combined.asUWEs()))
		})
	})
	t.Run("differentialUnitsGroup", func(t *testing.T) {
		cancellable := 2
		var us = []*baseUnit{
			GetUnitBySignature("m").(*baseUnit),
			GetUnitBySignature("g").(*baseUnit),
			GetUnitBySignature("s").(*baseUnit),
		}
		uga := baseUnitsGroup{
			{unit: us[0], exp: 3},
			{unit: us[1], exp: cancellable},
		}
		ugb := baseUnitsGroup{
			{unit: us[1], exp: cancellable},
			{unit: us[2], exp: 5},
		}
		res := differentialUnitsGroup(uga, ugb)
		for _, uwe := range res.asUWEs() {
			switch uwe.Unit() {
			case us[0]:
				assert.Equal(t, -uga[0].exp, uwe.Exp())
			case us[1]:
				assert.Fail(t, "us[1] should not exist in the result")
			case us[2]:
				assert.Equal(t, ugb[1].exp, uwe.Exp())
			default:
				t.Fail()
			}
		}
	})
}

func compareGroups(t *testing.T, a, b unitsGroup) {
	aUwes, bUwes := a.asUWEs(), b.asUWEs()
	assert.Equal(t, len(aUwes), len(bUwes), "asUWEs")
	for i, buwe := range aUwes {
		assert.Equal(t, buwe.Unit().Signature(), bUwes[i].Unit().Signature(), "sub-unit signatures should match")
		assert.Equal(t, buwe.Exp(), bUwes[i].Exp())
	}
}
