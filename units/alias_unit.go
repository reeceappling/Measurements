package units

import "github.com/reeceappling/measurements/units/conv"

var (
	_ Units = &aliasUnit{}
)

// TODO: everything in here
type aliasUnit struct {
	sym       UnitSymbol
	aliasedTo Units
}

func (a aliasUnit) allConversions() []*conversion { // TODO: do not delete aliased conversions
	return a.aliasedTo.allConversions()
}

func NewAliasUnit(sym UnitSymbol) Units {
	// TODO implement me
	panic("implement me")
}

func (a aliasUnit) symbol() UnitSymbol { // TODO: meet interface???? not singleUnits
	return a.sym
}

func (a aliasUnit) Signature() UnitSignature {
	return UnitSignature(a.sym)
}

func (a aliasUnit) isBaseUnits() bool {
	return a.aliasedTo.isBaseUnits() // TODO: ensure ok
}

func (a aliasUnit) baseUnits() (out Units, createdBase bool) {
	return a.aliasedTo.baseUnits()
}

func (a aliasUnit) asGroup() unitsGroup {
	return a.aliasedTo.asGroup() // TODO: ensure ok
}

func (a aliasUnit) Per(units Units) conv.Rate {
	return a.aliasedTo.Per(units)
}

func (a aliasUnit) ConversionRateTo(units Units) conv.Rate {
	return a.aliasedTo.ConversionRateTo(units)
}

func (a aliasUnit) Mul(units ...Units) Units {
	return a.aliasedTo.Mul(units...)
}

func (a aliasUnit) Div(units Units) Units {
	return a.aliasedTo.Div(units)
}

func (a aliasUnit) Pow(i int) Units {
	return a.aliasedTo.Pow(i)
}

func (a aliasUnit) setConversion(otherUnitsSignature UnitSignature, conv *conversion) error {
	return a.aliasedTo.setConversion(otherUnitsSignature, conv)
}

func (a aliasUnit) getConversion(otherUnitsSignature UnitSignature) (*conversion, bool) {
	return a.aliasedTo.getConversion(otherUnitsSignature)
}

func (a aliasUnit) removeConversion(otherUnitsSignature UnitSignature) {
	a.aliasedTo.removeConversion(otherUnitsSignature)
}

func (a aliasUnit) RLock() {
	a.aliasedTo.RLock() // TODO: is this ok?
}

func (a aliasUnit) RUnlock() {
	a.aliasedTo.RUnlock() // TODO: is this ok?
}

func (a aliasUnit) Lock() {
	a.aliasedTo.Lock() // TODO: is this ok?
}

func (a aliasUnit) Unlock() {
	a.aliasedTo.Unlock() // TODO: is this ok?
}

func (a aliasUnit) unitSystem() *UnitSystem {
	return a.aliasedTo.unitSystem() // TODO: is this ok?
}
