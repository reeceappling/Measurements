package units

import (
	"errors"
	"github.com/reeceappling/measurements/units/conv"
	"golang.org/x/exp/maps"
	"slices"
	"sync"
)

type compoundUnit struct {
	group       unitsGroup
	conversions conversionMap
	system      *UnitSystem
	sync.RWMutex
}

func (c *compoundUnit) allConversions() []*conversion {
	return maps.Values(c.conversions)
}

func (c *compoundUnit) unitSystem() *UnitSystem {
	// TODO: remove ignore later when we actually test this
	//coverage:ignore
	return c.system
}

func (c *compoundUnit) getConversion(toSig UnitSignature) (conv *conversion, exists bool) {
	c.RLock()
	defer c.RUnlock()

	conv, exists = c.conversions[toSig]
	return
}

func newCompoundUnit(group unitsGroup) *compoundUnit {
	return &compoundUnit{
		group:       group.withoutNils(), // TODO: unsure if withoutNils is the best call here
		conversions: conversionMap{},
		RWMutex:     sync.RWMutex{},
	}
}

func (c *compoundUnit) hasConversion(toSig UnitSignature) bool {
	_, exists := c.getConversion(toSig)
	return exists
}

func (c *compoundUnit) removeConversion(otherUnitsSignature UnitSignature) {
	c.Lock()
	delete(c.conversions, otherUnitsSignature)
	c.Unlock()
}

var errCompoundConvExists = errors.New("staticConversion already exists for compound unit")

func (c *compoundUnit) setConversion(otherUnitsSignature UnitSignature, conv *conversion) error {
	// TODO: resolve options, not necessarily always currentServiceOptions
	if !currentServiceOptions.allowOverWritingConversionRate {
		if c.hasConversion(otherUnitsSignature) {
			return errCompoundConvExists
		}
	}

	c.Lock()
	defer c.Unlock()
	c.conversions[otherUnitsSignature] = conv
	return nil
}

func (c *compoundUnit) Signature() UnitSignature {
	return c.group.Signature()
}

func (c *compoundUnit) isBaseUnits() bool {
	return c.group.isBaseGroup()
}

func (c *compoundUnit) baseUnits() (out Units, createdBase bool) {
	if c.isBaseUnits() /*&& !containsNils(c)*/ { // TODO: likely no longer need to check that this contains nils if newCompoundUnit returns withoutNils()
		return c, false
	}
	// if base group exists, return it
	baseGroup := c.group.baseGroup().withoutNils() // TODO: this is failing on compoundUnit
	if len(baseGroup.asUWEs()) == 0 {              // TODO: this thing here may be necessary
		return unitless, false
	}
	baseSig := baseGroup.Signature()
	if unit, exists := allUnits.Get(baseSig); exists {
		return unit, false
	}
	// If base droup does not exist, create it and add it to the graph // TODO: this may be rough for count units
	newBaseUnits := newCompoundUnit(baseGroup)
	allUnits.addNode(newBaseUnits, baseSig)
	// Resolve group distance from theoretical base

	var preMergeGroups []conv.Rate
	for uwe := range slices.Values(c.group.asUWEs()) {
		subUnit := uwe.Unit()
		if subUnit.isBaseUnits() {
			continue
		}
		subUnitBaseUnits, _ := subUnit.baseUnits() // Steamroll the bool because base/std units will always return false
		// If not base units, find distance to base
		preMergeGroups = append(preMergeGroups, unitsGroupRateToBase(subUnitBaseUnits.asGroup()))
	}
	// Create and set new staticConversion
	convToBase := &conversion{
		conversionPair: newConversionPair(c, newBaseUnits),
		rate:           conv.MultiplyRates(preMergeGroups...),
	}
	_ = convToBase.setConversionOnNodePair() // Should be impossible to error here
	//if err := convToBase.setConversionOnNodePair(); err != nil { // TODO: delete if unneeded
	//	convToBase.deleteConversionFromNodePair()
	//	allUnits.removeNode(baseSig)
	//	panic("failed to set new staticConversion to base units, new base node has been removed")
	//}
	return newBaseUnits, true
}

func (c *compoundUnit) asGroup() unitsGroup {
	return c.group
}

func (c *compoundUnit) Per(dst Units) conv.Rate {
	if c == dst {
		return conv.MultiplicativeIdentityRate()
	}
	var err error
	convers, exists := c.getConversion(dst.Signature())
	if exists {
		return convers.rateFrom(c.Signature())
	}

	newConv, err := resolveTwoStepConversion(c, dst)
	if err != nil {
		return conv.ErroneousRate{Err: err} // TODO: more err info here or no?
	}

	//err = newConv.setConversionOnNodePair() // TODO: may only be able to error if node pair is nil
	//if err != nil {
	//	// TODO: delete something????
	//	err = errors.Join(err, errors.New("failed to set staticConversion in compoundunit Per")) // TODO: move
	//	return conv.ErroneousRate{Err: err}
	//}
	_ = newConv.setConversionOnNodePair()

	return newConv.rateFrom(c.Signature())
}

func (c *compoundUnit) ConversionRateTo(units Units) conv.Rate {
	return units.Per(c)
}

func (c *compoundUnit) Mul(units ...Units) Units { // TODO: not tested yet, need to confirm works properly
	return mulUnits(appendUnitSliceToReciever(c, units)) // TODO: maybe move append into mulUnits
}

func (c *compoundUnit) Div(units Units) Units {
	return c.Mul(units.Pow(-1))
}

func (c *compoundUnit) Pow(p int) Units {
	if p == 0 {
		return unitless
	}
	newGroup := c.group.pow(p)                                // TODO: simplify by un-variable-ing if this works properly
	out := unitsForGroupCreateWithBaseIfNonexistent(newGroup) // TODO: should we use this or can we shortcut parts of it?
	return out
}
