package units

import (
	"errors"
	"github.com/reeceappling/measurements/units/conv"
	"strings"
)

const (
	unitSeparator     = " "
	exponentSeparator = "^"
)

var (
	errInvalidSymbolChars    = errors.New("new unit symbol contained one of the invalid characters: `" + unitSeparator + "` or `" + exponentSeparator + "`")
	errBaseSignatureMismatch = errors.New("staticConversion impossible, units not of same dimensions")
)

var (
	_ Units      = &compoundUnit{}
	_ singleUnit = &baseUnit{}
	_ singleUnit = &standardUnit{}
)

// Units are the interface that is actually exported from this package
// Units are always met by pointers to structs, as they need to modify themselves and be unique in the global graph
//
//go:generate mockery --name Units
type Units interface {
	Signable
	isBaseUnits() bool
	// baseUnits gets the baseUnits for a Units (and creates if nonexistant)returning the resulting base, and whether or not the base was created
	baseUnits() (out Units, createdBase bool)
	asGroup() unitsGroup
	// Per (opposite of ConversionRateTo) Ex: m.Per(km) == 1000 == km.ConversionRateTo(m)
	Per(Units) conv.Rate
	// ConversionRateTo (opposite of Per) Ex: km.ConversionRateTo(m) == 1000 == m.Per(km)
	ConversionRateTo(Units) conv.Rate
	allConversions() []*conversion
	// Mul multiplies two units together. May have weird results like m*cm == m cm  (creates a new unit if it does not exist)
	Mul(...Units) Units
	// Div divides the interface unit by the input unit.  May have weird results like m/cm == m cm^-1 (creates a new unit if it does not exist)
	Div(Units) Units
	// Pow returns the unit that is this unit to the power of the input, creating it if needed
	Pow(int) Units
	// setConversion TODO: this
	setConversion(otherUnitsSignature UnitSignature, conv *conversion) error
	// getConversion TODO: this
	getConversion(otherUnitsSignature UnitSignature) (*conversion, bool)
	// removeConversion TODO: this
	removeConversion(otherUnitsSignature UnitSignature)
	// RLock locks a Units for reading, using its internal sync.RWMutex
	RLock()
	// RUnlock unlocks a Units from a previous RLock, using its internal sync.RWMutex
	RUnlock()
	// Lock locks a Units for writing, using its internal sync.RWMutex. Also blocks reads
	Lock()
	// Unlock unlocks a Units from writing, using its internal sync.RWMutex. Also unblocks reads
	Unlock()
	unitSystem() *UnitSystem
}

//func containsNils(u Units) bool { // TODO: TESTME
//	for _, uwe := range u.asGroup().asUWEs() {
//		if uwe.Unit() == unitless {
//			return true
//		}
//	}
//	return false
//}

type Signable interface {
	Signature() UnitSignature
}

type UnitSymbol UnitSignature

func (sym UnitSymbol) isValid() bool {
	symStr := string(sym)
	containsUnitSep, containsExpSep := strings.Contains(symStr, unitSeparator), strings.Contains(symStr, exponentSeparator)
	// account for currentServiceOptions otherInvalidCharacters
	for _, invalidSym := range currentServiceOptions.extraInvalidUnitSymbols {
		if strings.Contains(symStr, invalidSym) {
			return false
		}
	}
	return !(containsUnitSep || containsExpSep) // TODO: simplify
}

func (sym UnitSymbol) asSignature() UnitSignature { return UnitSignature(sym) }

type UnitSignature string

type singleUnit interface {
	Units
	symbol() UnitSymbol
	asUnitWithExponent() unitWithExponent
}
