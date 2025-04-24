package units

import (
	"errors"
	"github.com/reeceappling/measurements/units/conv"
	"golang.org/x/exp/maps"
	"sync"
)

var (
	_ singleUnit = &baseUnit{}
)

type baseUnit struct { // TODO: THIS INCLUDES UNITLESS
	sym         UnitSymbol
	conversions conversionMap
	system      *UnitSystem // TODO: account for this elsewhere
	sync.RWMutex
}

func (b *baseUnit) allConversions() []*conversion {
	return maps.Values(b.conversions)
}

func (b *baseUnit) unitSystem() *UnitSystem {
	// TODO: remove ignore
	//coverage:ignore
	return b.system
}

var (
	errAddingBaseNode = errors.New("failed to add a base node to the global graph")
)

// NewBaseUnit creates a new base unit with a specific symbol. It DOES add it to the graph
func NewBaseUnit(symbol string) (out Units, err error) {
	newSym := UnitSymbol(symbol)
	// TODO: ensure unit does not already exist
	out, err = newBaseUnit(newSym)
	if err != nil {
		return nil, errors.Join(err, errBadNewBaseUnitSymbol)
	}
	if err = allUnits.tryAddNode(out); err != nil {
		out, _ = allUnits.Get(newSym.asSignature())
		if !currentServiceOptions.errorOnDuplicateNewBaseUnit {
			err = nil
		}
		return
	}
	return
}

// newBaseUnit creates a new base unit with a specific symbol. It does NOT add it to the graph
func newBaseUnit(symbol UnitSymbol) (*baseUnit, error) {
	if !symbol.isValid() {
		return nil, errInvalidSymbolChars
	}
	return &baseUnit{
		sym:         symbol,
		conversions: conversionMap{},
		RWMutex:     sync.RWMutex{},
		// TODO: add special processing rules for specific units/classes?
	}, nil
}

func (b *baseUnit) getConversion(otherUnitsSignature UnitSignature) (conv *conversion, exists bool) {
	b.RLock()
	defer b.RUnlock()

	conv, exists = b.conversions[otherUnitsSignature]
	return
}

func (b *baseUnit) removeConversion(otherUnitsSignature UnitSignature) {
	b.Lock()
	defer b.Unlock()

	delete(b.conversions, otherUnitsSignature)
}

func (b *baseUnit) setConversion(otherUnitsSignature UnitSignature, conv *conversion) error {
	// TODO: resolve options, not necessarily always currentServiceOptions
	if !currentServiceOptions.allowOverWritingConversionRate {
		if _, exists := b.getConversion(otherUnitsSignature); exists {
			return errors.New("staticConversion already exists for base unit")
		}
	}

	b.Lock()
	defer b.Unlock()
	b.conversions[otherUnitsSignature] = conv
	return nil
}

func (b *baseUnit) asGroup() unitsGroup {
	if b == unitless {
		return baseUnitsGroup{}
	}
	return baseUnitsGroup{b.asBaseUnitWithExponent()}
}

func (b *baseUnit) asBaseUnitWithExponent() baseUnitWithExponent {
	return baseUnitWithExponent{
		unit: b,
		exp:  1,
	}
}

func (b *baseUnit) asUnitWithExponent() unitWithExponent {
	return b.asBaseUnitWithExponent()
}

var errNoBaseConversion = errors.New("staticConversion rate from base unit does not exist")

func (b *baseUnit) Per(units Units) conv.Rate {
	if b == units {
		return conv.MultiplicativeIdentityRate()
	}
	otherSig := units.Signature() // TODO: ensure this works for countUnits we don't know convs for yet...
	if convers, exists := b.getConversion(otherSig); exists {
		return convers.rateFrom(b.Signature())
	}
	return conv.ErroneousRate{Err: errNoBaseConversion}
}

func (b *baseUnit) ConversionRateTo(units Units) conv.Rate {
	return units.Per(b)
}

func (b *baseUnit) Mul(units ...Units) Units {
	return mulUnits(appendUnitSliceToReciever(b, units)) // TODO: maybe move append into mulUnits
}

func (b *baseUnit) Div(units Units) Units {
	if b == unitless {
		return units.Pow(-1)
	}
	return b.Mul(units.Pow(-1))
}

func (b *baseUnit) Pow(i int) Units {
	if i == 0 || b == unitless { // TODO: unsure if we REALLY want b == unitless
		return unitless
	}
	newBug := baseUnitsGroup{baseUnitWithExponent{unit: b, exp: i}}
	newSig := newBug.Signature()
	if unit, exists := allUnits.Get(newSig); exists {
		return unit
	}
	// if newUWE does not exist, create it and store it
	newUnit := newCompoundUnit(newBug)
	allUnits.addNode(newUnit, newSig)
	return newUnit
}

func (b *baseUnit) Signature() UnitSignature {
	return b.symbol().asSignature()
}

func (b *baseUnit) isBaseUnits() bool {
	return true
}

func (b *baseUnit) baseUnits() (Units, bool) {
	return b, false
}

func (b *baseUnit) symbol() UnitSymbol {
	return b.sym
}
