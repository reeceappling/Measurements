package units

import (
	"errors"
	"github.com/reeceappling/measurements/units/conv"
	sliceutils "github.com/reeceappling/randomGoStuff/utils/slices"
)

var (
	errBadNewBaseUnitSymbol          = errors.New("invalid symbol provided for new base unit")
	errUnitDoesNotExist              = errors.New("unit does not exist")
	errInvalidNewEquivalentUnitZeros = errors.New("cannot create unit with 0.0 as either equivalence")
	errFailedToAddNewNonBaseNode     = errors.New("failed to add new non-base node")
	errEquivUnitToBaseConv           = errors.New("failed to find known staticConversion to base units for new equivalent unit")
	errSetConvToBaseOfKnown          = errors.New("failed to set new staticConversion to base units of a known equivalent")
	errUnitAlreadyExists             = errors.New("unit already exists")
	errNoBaseMatch                   = errors.New("staticConversion impossible, units not of same dimensions") // TODO: make different?
)

// TODO: allow this to do dynamic rates as well?
func NewUnitWithEquivalence(amtNew float64, unitSymb string, amtKnown float64, knownUnits Units) (newUnit Units, err error) { // TODO: ISSUE HERE
	if amtNew*amtKnown == 0.0 {
		return nil, errInvalidNewEquivalentUnitZeros
	}
	// ensure newSymb has no invalid chars
	newSymb := UnitSymbol(unitSymb)
	if !newSymb.isValid() {
		return nil, errInvalidSymbolChars
	}
	if u, exists := allUnits.Get(newSymb.asSignature()); exists { // TODO: do not allow new units that already exist (unless unit replacement is active?)
		return u, errUnitAlreadyExists
	}
	// Try to create new unit
	knownUnitsBaseUnits, _ := knownUnits.baseUnits()      // Should create/add base units if nonexistent
	newUnit, _ = newStdUnit(newSymb, knownUnitsBaseUnits) // steamroll error because we know it is base units

	// create staticConversion Factor (from * knownPerNew -> to, to/knownPerNew -> from)
	newStaticRatio := conv.StaticRate(amtNew / amtKnown)
	newToKnownConv := &conversion{
		conversionPair: newConversionPair(newUnit, knownUnits),
		rate:           &newStaticRatio,
	}
	if err = newToKnownConv.setConversionOnNodePair(); err != nil {
		// TODO: remove ignore, remove entire block if no err is possible
		//coverage:ignore
		return nil, err
	}

	if knownUnits.isBaseUnits() {
		allUnits.addNode(newUnit, UnitSignature(newSymb))
		return newUnit, nil
	}

	// If knownUnits are not baseUnits, also create staticConversion from new to base
	// newConversions := []*conversion{newToKnownConv} // TODO: remove if unneeded
	knownToBase := unitsGroupRateToBase(knownUnits.asGroup())

	newToBase := conv.MultiplyRates(newToKnownConv.rate, knownToBase) // TODO: ensure correct direction
	//newToBase := conv.MultiplyRates(newToKnownConv.rate, conv.RatePow(knownToBase, -1)) // TODO: ensure correct direction
	newConversionToBase := &conversion{
		conversionPair: newConversionPair(newUnit, knownUnitsBaseUnits),
		rate:           newToBase,
	}
	_ = newConversionToBase.setConversionOnNodePair() // Steamroll error because error cannot exist if conversion is not nil and the pair has not existed before
	//if err != nil { // TODO: delete if unused
	//	// revert already-done tasks
	//	newToKnownConv.deleteConversionFromNodePair()
	//	newConversionToBase.deleteConversionFromNodePair()
	//	// TODO: delete new nodes if needed
	//	err = errors.Join(err, errSetConvToBaseOfKnown)
	//	return nil, err
	//}
	// newConversions = append(newConversions, newConversionToBase) // TODO: remove if unneeded

	allUnits.addNode(newUnit, UnitSignature(newSymb))
	//err = allUnits.tryAddNode(newUnit) // TODO: repeated above, write once? is "Try" needed? // TODO: remove if unneeded
	//if err != nil {
	//	for convers := range slices.Values(newConversions) {
	//		convers.deleteConversionFromNodePair()
	//	}
	//	// TODO: delete new nodes if needed
	//	err = errors.Join(err, errFailedToAddNewNonBaseNode)
	//}
	return newUnit, err
}

func mulUnits(units []Units) Units { // TODO: testMe. Also use me in base, std, compound
	if len(units) == 0 {
		return unitless // TODO: testMe, may be unnecessary
	}
	//TODO: from std/compound
	resultGroup := combineGroups(sliceutils.Map(units, func(u Units) unitsGroup { // TODO: I don't think this gets rid of unitless?
		return u.asGroup().withoutNils() // TODO: avoid using withoutNils here if possible
	})...)
	out := unitsForGroupCreateWithBaseIfNonexistent(resultGroup)
	return out
}

func appendUnitSliceToReciever(rec Units, uSlice []Units) []Units {
	out := make([]Units, len(uSlice)+1) // TODO: do we want to do anything with unitless in here?
	out[0] = rec
	for i, unit := range uSlice {
		out[i+1] = unit
	}
	return out
}

func unitsForGroupCreateWithBaseIfNonexistent(finalUnitsGroup unitsGroup) Units { // TODO: rename, not yet tested
	newSig := finalUnitsGroup.Signature()
	// if Signature exists, return units for it
	u, exists := allUnits.Get(newSig) // Covers all base units (incl. unitless), std units, and known compound units
	if exists {
		return u
	}
	var unitsOut Units = newCompoundUnit(finalUnitsGroup)
	allUnits.addNode(unitsOut, newSig)
	if finalUnitsGroup.isBaseGroup() {
		return unitsOut
	}
	var newBaseUnits Units

	newBase := finalUnitsGroup.baseGroup().withoutNils() // TODO: unsure if withoutNils is a no-op here or not
	newBaseSig := newBase.Signature()
	// if new base Signature does not exist, create it
	newBaseUnits, exists = allUnits.Get(newBaseSig)
	if !exists {
		newBaseUnits = newCompoundUnit(newBase)
		allUnits.addNode(newBaseUnits, newBaseSig)
	}
	// create conversion from new unit to base
	rateToBase := unitsGroupRateToBase(finalUnitsGroup)
	newConv := &conversion{
		conversionPair: newConversionPair(unitsOut, newBaseUnits),
		rate:           rateToBase, // TODO:  conv.RatePow(rateToBase, -1)
	}
	_ = newConv.setConversionOnNodePair() // Error should be impossible because conversion exists and conversion between pair should not already exist
	return unitsOut
}

func resolveTwoStepConversion(src, dst Units) (*conversion, error) { // TODO: test me
	// if base sigs not equal, cannot convert
	baseU, baseCreated := src.baseUnits()      // TODO: ensure this works for countUnits we don't know convs for uet...
	dstBase, dstBaseCreated := dst.baseUnits() // TODO: ensure this works for countUnits we don't know convs for uet...
	baseSig, dstBaseSig := baseU.Signature(), dstBase.Signature()
	if baseSig != dstBaseSig {
		// TODO: remove ignore later. Remove block if impossible to trigger
		//coverage:ignore
		if baseCreated { // TODO: unsure if necessary
			// TODO: remove ignore later. Remove block if impossible to trigger
			//coverage:ignore
			panic("resolveTwoStepConversion baseCreated was triggered")
		}
		if dstBaseCreated { // TODO: unsure if necessary
			// TODO: remove ignore later. Remove block if impossible to trigger
			//coverage:ignore
			panic("resolveTwoStepConversion dstBaseCreated was triggered")
		}
		return nil, errNoBaseMatch // TODO: move and change
	}

	srcToBase, existsA := src.getConversion(baseSig)
	baseToDst, existsB := baseU.getConversion(dst.Signature())
	var rate conv.Rate
	if !existsA || !existsB {
		// find both's distance to base, combine for real rate in new staticConversion
		differentialGroup := differentialUnitsGroup(src.asGroup(), dst.asGroup()) // TODO: remove or switch back
		rate = unitsGroupRateToBase(differentialGroup)
	} else {
		rate = conv.MultiplyRates(srcToBase.rateFrom(src.Signature()), baseToDst.rateFrom(baseSig))
	}

	return &conversion{
		conversionPair: newConversionPair(src, dst),
		rate:           rate,
	}, nil
}

func GetUnitBySignature(sig string) Units { // TODO: testMe
	result, exists := allUnits.Get(UnitSignature(sig))
	if !exists {
		panic("unit not found: " + sig) // TODO: is it ok to panic here?
	}
	return result
}
