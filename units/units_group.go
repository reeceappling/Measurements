package units

import (
	"github.com/reeceappling/measurements/units/conv"
	"github.com/reeceappling/randomGoStuff/utils/slices"
	"sort"
)

var (
	_ unitsGroup = baseUnitsGroup{}
	_ unitsGroup = stdUnitsGroup{} // TODO: why the hell is this doing this?
)

type unitsGroup interface { // TODO: why the hell is this doing this?
	Signable
	isBaseGroup() bool
	baseGroup() baseUnitsGroup
	asUWEs() []unitWithExponent
	pow(int) unitsGroup // TODO: ensure all pows are tested
	withoutNils() unitsGroup
}

func unitsGroupRateToBase(ug unitsGroup) conv.Rate { // TODO: use this in Mul, all over the place
	allUWEs := ug.asUWEs()
	rateToBaseFromUwe := func(uwe unitWithExponent) conv.Rate {
		return uwe.rateToBase()
	}
	allRates := slices.Map(allUWEs, rateToBaseFromUwe)
	return conv.MultiplyRates(allRates...)
	//return conv.MultiplyRates(
	//	slices.Map(ug.asUWEs(),
	//		func(uwe unitWithExponent) conv.Rate {
	//			return uwe.rateToBase()
	//		},
	//	)...,
	//)

}

func combineGroups(ugSet ...unitsGroup) unitsGroup {
	out := workingUnitsGroup{}
	for _, ug := range ugSet {
		out = append(out, ug.asUWEs()...)
	}
	return out.asProperGroup()
}

// Note: workingUnitsGroup is NOT Signable, as we cannot guarantee simplification or alphabetization
type workingUnitsGroup []unitWithExponent

func (wug workingUnitsGroup) asProperGroup() unitsGroup {
	return wug.
		simplified().
		alphabetized().
		resolveGroupType()
}
func (wug workingUnitsGroup) simplified() workingUnitsGroup {
	amts := map[UnitSymbol]int{}
	ptrs := map[UnitSymbol]singleUnit{}
	for _, uwe := range wug {
		symb := uwe.Unit().symbol()
		if amt, exists := amts[symb]; exists {
			newAmt := amt + uwe.Exp()
			amts[symb] = newAmt
			continue
		}
		amts[symb], ptrs[symb] = uwe.Exp(), uwe.Unit()
	}
	out := workingUnitsGroup{}
	for symb, exp := range amts {
		if exp == 0 {
			continue
		}
		thisUnit := ptrs[symb]
		var newUWE unitWithExponent
		if thisUnit.isBaseUnits() {
			newUWE = baseUnitWithExponent{
				unit: thisUnit.(*baseUnit),
				exp:  exp,
			}
		} else {
			newUWE = stdUnitWithExponent{unit: thisUnit, exp: exp}
		}
		out = append(out, newUWE)
	}
	return out
}
func (wug workingUnitsGroup) alphabetized() workingUnitsGroup {
	sort.Slice(wug, func(i int, j int) bool { // TODO: ensure works
		return wug[i].Unit().symbol() < wug[j].Unit().symbol()
	})
	return wug
}

func (wug workingUnitsGroup) resolveGroupType() unitsGroup {
	bug := make([]baseUnitWithExponent, len(wug))
	for i, uwe := range wug {
		baseUnitI, ok := uwe.Unit().(*baseUnit)
		if !ok {
			return stdUnitsGroup(wug)
		}

		bug[i] = baseUnitWithExponent{
			unit: baseUnitI,
			exp:  uwe.Exp(),
		}
	}
	return baseUnitsGroup(bug)
}

// stdUnitsGroup should be simplified, alphabetized, and contain at least 1 non-base unit
type stdUnitsGroup []unitWithExponent

func (s stdUnitsGroup) withoutNils() unitsGroup {
	out := stdUnitsGroup{} // TODO: FIGURE OUT WHERE ALL WE NEED TO USE withoutNils()
	for _, suwe := range s.asUWEs() {
		if suwe.Unit() != unitless { // TODO: unsure if this will work!
			out = append(out, suwe)
		}
	}
	return out
}

func (s stdUnitsGroup) pow(pwr int) unitsGroup {
	out := make(stdUnitsGroup, len(s))
	for i, uwe := range s {
		uweUnit := uwe.Unit()
		switch uwe.Unit().(type) { // TODO: move this out elsewhere
		case *baseUnit:
			out[i] = baseUnitWithExponent{ // TODO: should not always be stdUnitsWithExponent
				unit: uweUnit.(*baseUnit),
				exp:  uwe.Exp() * pwr,
			}
		default:
			out[i] = stdUnitWithExponent{ // TODO: should not always be stdUnitsWithExponent
				unit: uweUnit,
				exp:  uwe.Exp() * pwr,
			}
		}
	}
	return out
}

func (s stdUnitsGroup) isBaseGroup() bool {
	return false
}

func (s stdUnitsGroup) asUWEs() []unitWithExponent {
	return s
}

func (s stdUnitsGroup) baseGroup() baseUnitsGroup {
	out := workingUnitsGroup{}
	for _, uwe := range s {
		baseUwesFor := func(nonBaseUwe unitWithExponent) []unitWithExponent {
			bug, _ := nonBaseUwe.Unit().baseUnits()
			return bug.asGroup().pow(nonBaseUwe.Exp()).asUWEs()
		}
		out = append(out, baseUwesFor(uwe)...)
	}

	resultGenericGroup := out.asProperGroup()
	result := resultGenericGroup.(baseUnitsGroup) // TODO: collapse in return if kept
	//result, ok := resultGenericGroup.(baseUnitsGroup) // TODO: get rid of if unused
	//if !ok {
	//	panic("failed to turn new proper group from standard to base, should never happen") // TODO: do something here
	//}
	return result
}

func (s stdUnitsGroup) Signature() UnitSignature {
	uweToSig := func(uwe unitWithExponent) UnitSignature {
		return uwe.Signature()
	}
	sigIsNonNil := func(sig UnitSignature) bool {
		return sig != ""
	}
	nonNilSigs := slices.Filter(slices.Map(s, uweToSig), sigIsNonNil)
	if len(nonNilSigs) == 0 {
		return "" // TODO: maybe panic here instead?
	}
	out := nonNilSigs[0]
	for i := 1; i < len(s); i++ {
		out = out + unitSeparator + nonNilSigs[i]
	}
	return out
}

// baseUnitsGroup should be simplified, alphabetized, and contain only base units
type baseUnitsGroup []baseUnitWithExponent

func (b baseUnitsGroup) withoutNils() unitsGroup { // TODO: TESTME
	out := baseUnitsGroup{} // TODO: FIGURE OUT WHERE ALL WE NEED TO USE withoutNils()
	for _, buwe := range b {
		if buwe.unit.sym != "" {
			out = append(out, buwe)
		}
	}
	return out
}

func (b baseUnitsGroup) pow(pwr int) unitsGroup {
	out := make(baseUnitsGroup, len(b))
	for i, baseUwe := range b {
		out[i] = baseUnitWithExponent{
			unit: baseUwe.unit,
			exp:  baseUwe.Exp() * pwr,
		}
	}
	return out
}

func (b baseUnitsGroup) isBaseGroup() bool {
	return true
}

func (b baseUnitsGroup) asUWEs() []unitWithExponent {
	return slices.Map(b, func(in baseUnitWithExponent) unitWithExponent {
		return in
	})
}

func (b baseUnitsGroup) baseGroup() baseUnitsGroup {
	return b
}

func (b baseUnitsGroup) Signature() UnitSignature { // TODO: test
	sigs := slices.Map(b, func(buwe baseUnitWithExponent) UnitSignature {
		return buwe.Signature()
	})
	sigsWithoutNils := slices.Filter(sigs, func(sig UnitSignature) bool {
		return sig != ""
	})
	if len(sigsWithoutNils) == 0 {
		return ""
	}
	out := sigsWithoutNils[0]
	for i := 1; i < len(sigsWithoutNils); i++ {
		out = out + unitSeparator + sigsWithoutNils[i]
	}
	return out
}

func differentialUnitsGroup(from, to unitsGroup) unitsGroup { // TODO: TESTME
	differentialWug := workingUnitsGroup{}
	for _, uwe := range from.pow(-1).asUWEs() {
		differentialWug = append(differentialWug, uwe)
	}
	for _, uwe := range to.asUWEs() {
		differentialWug = append(differentialWug, uwe)
	}
	return differentialWug.asProperGroup()
}
