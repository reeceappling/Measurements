package units

import (
	"errors"
	"github.com/reeceappling/measurements/units/conv"
	"golang.org/x/exp/maps"
	"sync"
)

var (
	errStdUnitConvAlreadyExists = errors.New("unit staticConversion already exists for standard unit")
)

var (
	_ singleUnit = &standardUnit{}
)

type standardUnit struct {
	sym         UnitSymbol
	baseGroup   Units
	conversions conversionMap
	system      *UnitSystem
	sync.RWMutex
}

func (s *standardUnit) allConversions() []*conversion {
	return maps.Values(s.conversions)
}

func (s *standardUnit) unitSystem() *UnitSystem {
	// TODO: remove ignore
	//coverage:ignore
	return s.system
}

func (s *standardUnit) getConversion(otherUnitsSignature UnitSignature) (conv *conversion, exists bool) {
	s.RLock()
	defer s.RUnlock()

	conv, exists = s.conversions[otherUnitsSignature]
	return
}

func newStdUnit(symbol UnitSymbol, baseUnits Units) (*standardUnit, error) {
	if !baseUnits.isBaseUnits() {
		return nil, errors.New("cannot create new standard unit with base group that is not all base units")
	}
	if allUnits.contains(symbol.asSignature()) {
		return nil, errors.New("standard unit already exists")
	}
	return &standardUnit{
		sym:         symbol,
		baseGroup:   baseUnits,
		conversions: conversionMap{},
		RWMutex:     sync.RWMutex{},
	}, nil
}

func (s *standardUnit) removeConversion(otherUnitsSignature UnitSignature) {
	s.Lock()
	delete(s.conversions, otherUnitsSignature)
	s.Unlock()
}

func (s *standardUnit) setConversion(otherUnitsSignature UnitSignature, conver *conversion) error {
	// TODO: resolve options, not necessarily always currentServiceOptions
	if !currentServiceOptions.allowOverWritingConversionRate {
		s.RLock()
		_, exists := s.conversions[otherUnitsSignature]
		s.RUnlock()
		if exists {
			return errStdUnitConvAlreadyExists
		}
	}

	s.Lock()
	defer s.Unlock()
	s.conversions[otherUnitsSignature] = conver
	return nil
}

func (s *standardUnit) Signature() UnitSignature {
	return s.symbol().asSignature()
}

func (s *standardUnit) isBaseUnits() bool {
	return false
}

func (s *standardUnit) baseUnits() (Units, bool) {
	return s.baseGroup, false
}

func (s *standardUnit) asGroup() unitsGroup {
	return stdUnitsGroup{s.asUnitWithExponent()}
}

func (s *standardUnit) Per(dst Units) conv.Rate {
	if s == dst {
		return conv.MultiplicativeIdentityRate()
	}
	var err error
	if convers, exists := s.getConversion(dst.Signature()); exists {
		return convers.rateFrom(s.Signature())
	}
	newConv, err := resolveTwoStepConversion(s, dst)
	if err != nil {
		return conv.ErroneousRate{Err: err} // TODO: more err info here or no
	}
	_ = newConv.setConversionOnNodePair()

	return newConv.rateFrom(s.Signature())
}

func (s *standardUnit) ConversionRateTo(units Units) conv.Rate {
	return units.Per(s)
}

func (s *standardUnit) Mul(units ...Units) Units { // TODO: ensure works for unitless base/std
	return mulUnits(appendUnitSliceToReciever(s, units)) // TODO: maybe move append into mulUnits
}

func (s *standardUnit) Div(units Units) Units {
	return s.Mul(units.Pow(-1))
}

func (s *standardUnit) Pow(i int) Units {
	if i == 0 {
		return unitless // TODO: test
	}
	newUWE := stdUnitWithExponent{unit: s, exp: i}
	sig := newUWE.Signature()
	// if unit already exists, return it
	if out, exists := allUnits.Get(sig); exists {
		return out
	}
	// Create new compound unit
	newCu := newCompoundUnit(newUWE)
	// Get base group for new unit. Create if needed
	newBaseGroup := newUWE.baseGroup()
	newBaseSig := newBaseGroup.Signature()
	//createdNewBase := false // TODO: del if unneeded
	newBase, exists := allUnits.Get(newBaseSig)
	if !exists {
		newBase = newCompoundUnit(newBaseGroup)
		allUnits.addNode(newBase, newBaseSig)
		//if err := allUnits.tryAddNode(newBase); err != nil { // TODO: delete if unneeded
		//	// TODO: remove ignore
		//	//coverage:ignore
		//	panic("suPow failed to add new base node")
		//}
		//createdNewBase = true
	}
	// add new compound unit
	allUnits.addNode(newCu, sig)
	//if err := allUnits.tryAddNode(newCu); err != nil { // TODO: delete if unneeded
	//	// TODO: remove ignore
	//	//coverage:ignore
	//	if createdNewBase {
	//		allUnits.removeNode(newBaseSig)
	//	}
	//	panic("failed to add new compound unit in suPow") // TODO: maybe change?
	//}

	// Get new unit's distance from its base, create and set staticConversion
	currentConversionRate, _ := s.getConversion(s.baseGroup.Signature()) // TODO: ok to steamroll exists here?
	newConv := &conversion{
		conversionPair: newConversionPair(newCu, newBase), // TODO: ensure direction ok
		rate:           conv.RatePow(currentConversionRate.rateFrom(s.Signature()), i),
	}
	_ = newConv.setConversionOnNodePair() // Err should be impossible here
	//if err := newConv.setConversionOnNodePair(); err != nil { // TODO: delete if unused
	//	// Rollback all others
	//	if createdNewBase {
	//		allUnits.removeNode(newBaseSig)
	//	}
	//	allUnits.removeNode(sig)
	//	panic("failed to set staticConversion, should never happen") // TODO: something else?
	//}
	return newCu
}

func (s *standardUnit) symbol() UnitSymbol {
	return s.sym
}

func (s *standardUnit) asUnitWithExponent() unitWithExponent {
	return stdUnitWithExponent{unit: s, exp: 1}
}
