package units

import (
	"fmt"
	"github.com/reeceappling/measurements/units/conv"
)

var (
	_ unitWithExponent = baseUnitWithExponent{}
	_ unitWithExponent = stdUnitWithExponent{}
)

type unitWithExponent interface {
	Signable
	unitsGroup
	Unit() singleUnit
	Exp() int
	rateToBase() conv.Rate
}

func uweSignature(uwe unitWithExponent) UnitSignature {
	u, exp := uwe.Unit(), uwe.Exp()
	if exp == 0 {
		panic("0-exponent unit, should never occur") // TODO: ??? or should we return ""?
	}
	if exp == 1 {
		return u.Signature()
	}
	return UnitSignature(fmt.Sprintf(`%s%s%d`, u.Signature(), exponentSeparator, exp))
}

type baseUnitWithExponent struct {
	unit *baseUnit
	exp  int
}

func (b baseUnitWithExponent) withoutNils() unitsGroup { //TODO TEST ME
	if b.unit == unitless {
		return baseUnitsGroup{}
	}
	return b
}

func (b baseUnitWithExponent) rateToBase() conv.Rate {
	return conv.MultiplicativeRate{}
}

func (b baseUnitWithExponent) pow(i int) unitsGroup {
	return baseUnitWithExponent{
		unit: b.unit,
		exp:  b.exp * i,
	}
}

func (b baseUnitWithExponent) isBaseGroup() bool {
	return true
}

func (b baseUnitWithExponent) baseGroup() baseUnitsGroup {
	return baseUnitsGroup{b}
}

func (b baseUnitWithExponent) asUWEs() []unitWithExponent {
	return b.baseGroup().asUWEs()
}

func (b baseUnitWithExponent) Signature() UnitSignature {
	return uweSignature(b)
}

func (b baseUnitWithExponent) Unit() singleUnit {
	return b.unit
}

func (b baseUnitWithExponent) Exp() int {
	return b.exp
}

type stdUnitWithExponent struct {
	unit singleUnit
	exp  int
}

func (s stdUnitWithExponent) withoutNils() unitsGroup {
	return s
}

func (s stdUnitWithExponent) rateToBase() conv.Rate {
	if s.unit.isBaseUnits() {
		// TODO: remove ignore
		//coverage:ignore
		panic("base units should never be in a stdUnitWithExponent!") // TODO: maybe just get rid of this
	}
	bus, _ := s.unit.baseUnits() // Steamroll bool b/c it will always be false
	//rate, _ := s.unit.getConversion(bus.Signature()) // TODO: reenable when working
	rate, exists := s.unit.getConversion(bus.Signature()) // TODO: for use in debugging only
	if !exists {                                          // TODO: for use in debugging only
		// TODO: remove ignore
		//coverage:ignore
		panic("std unit conversion did not exist")
	}
	return conv.RatePow(rate.rateFrom(s.unit.Signature()), s.exp) // This should be the correct direction
}

func (s stdUnitWithExponent) pow(i int) unitsGroup {
	return stdUnitWithExponent{
		unit: s.unit,
		exp:  s.exp * i,
	}
}

func (s stdUnitWithExponent) isBaseGroup() bool {
	return s.unit.isBaseUnits()
}

func (s stdUnitWithExponent) baseGroup() baseUnitsGroup {
	bus, _ := s.unit.baseUnits() // Steamroll bool b/c std units baseUnits never create new ones
	groupToPow := bus.asGroup().(baseUnitsGroup)
	out := make(baseUnitsGroup, len(groupToPow))
	for i, buwe := range groupToPow {
		out[i] = baseUnitWithExponent{
			unit: buwe.unit,
			exp:  buwe.exp * s.exp,
		}
	}
	return out
}

func (s stdUnitWithExponent) asUWEs() []unitWithExponent {
	return []unitWithExponent{s}
}

func (s stdUnitWithExponent) Signature() UnitSignature {
	return uweSignature(s)
}

func (s stdUnitWithExponent) Unit() singleUnit {
	return s.unit
}

func (s stdUnitWithExponent) Exp() int {
	return s.exp
}
