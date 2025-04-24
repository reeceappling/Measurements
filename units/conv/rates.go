package conv

import (
	"errors"
	"github.com/reeceappling/measurements/units/numericalPair/pairs"
	"github.com/reeceappling/measurements/units/numericalPair/pairs/inefficientComparable"
	"github.com/reeceappling/measurements/units/numericalPair/power"
	"github.com/reeceappling/randomGoStuff/utils"
	"reflect"
)

var (
	ErrNilRatePointer  = errors.New("rate pointer was nil")
	ErrNilRateResolver = errors.New("rate resolver was nil")
)

var (
	NegativeRate = utils.Pointer(StaticRate(-1))
)

var (
	_ Rate                                   = MultiplicativeRate{} // TODO: likely unnecessary
	_ inefficientComparable.IC[compoundRate] = MultiplicativeRate{}
	_ Rate                                   = AdditiveRate{}
	_ inefficientComparable.IC[compoundRate] = AdditiveRate{}
	_ compoundRate                           = MultiplicativeRate{}
	_ compoundRate                           = AdditiveRate{}
	_ baseRate                               = utils.Pointer(StaticRate(0))
	_ baseRate                               = &DynamicRate{}
	_ baseRate                               = ErroneousRate{}
)

type Rate interface {
	// Resolve // TODO: more info
	Resolve(...*Options) (float64, error) // TODO: change to Resolve(Options)
	// asCompoundRate // TODO: more info
	AsMultiplicativeRate() MultiplicativeRate
}

// RatePow // TODO: more info
func RatePow(rate Rate, exp int) Rate {
	if err, isErr := rate.(ErroneousRate); isErr { // TODO: testMe
		return err
	}
	switch exp {
	case 0:
		return MultiplicativeRate{}
	case 1:
		return rate
	default:
		return rate.AsMultiplicativeRate().pow(exp).simplestForm()
	}
}

// MultiplyRates // TODO: more info
func MultiplyRates(conversionRates ...Rate) Rate {
	multiplicativeRates := make([]MultiplicativeRate, len(conversionRates))
	for i := 0; i < len(conversionRates); i++ {
		currentRate := conversionRates[i]
		if _, isErr := currentRate.(ErroneousRate); isErr {
			return currentRate
		}
		multiplicativeRates[i] = currentRate.AsMultiplicativeRate()
	}
	// create final base inputs
	finalBaseInputs := []baseRate{}
	// Go through all comparables
	spotsUsed := 0
	baseInputSpots := map[baseRate]int{}
	finalComparableNumbers := []int{}
	for _, mr := range multiplicativeRates {
		for _, bi := range mr.Comparable {
			br, num := bi.Value.(baseRate), bi.Number
			if i, exists := baseInputSpots[br]; exists { // TODO: this should be exists, right?
				currentNum := finalComparableNumbers[i]
				finalComparableNumbers[i] = currentNum + num
				continue
			}
			finalBaseInputs = append(finalBaseInputs, br)
			finalComparableNumbers = append(finalComparableNumbers, num)
			baseInputSpots[br] = spotsUsed
			spotsUsed++
		}
	}
	finalComparablePairs := make([]pairs.Pair[int, Rate], spotsUsed)
	for base, i := range baseInputSpots {
		finalComparablePairs[i] = pairs.Pair[int, Rate]{
			Value:  base,
			Number: finalComparableNumbers[i],
		}
	}
	newUncRateNumbers := []int{}
	newUncomparableRates := []compoundRate{}
	indexForRate := func(known []compoundRate, toCheck compoundRate) int {
		for i, cr := range known {
			if toCheck.EqualTo(cr) {
				return i
			}
		}
		return -1
	}
	// combine like uncomparables
	for _, mr := range multiplicativeRates {
		for _, uncP := range mr.Uncomparable { // TODO: deal with negative additionRates (-1*-1 == 1) as well as varying ratios
			uncVal := uncP.Value.(compoundRate)
			// get list of all required base rates, then see if they need adding to the system
			indexInKnownRates := indexForRate(newUncomparableRates, uncVal)
			if indexInKnownRates != -1 {
				current := newUncRateNumbers[indexInKnownRates]
				newUncRateNumbers[indexInKnownRates] = current + uncP.Number
			} else {
				newUncomparableRates = append(newUncomparableRates, uncVal)
				newUncRateNumbers = append(newUncRateNumbers, uncP.Number)
			}
		}
	}
	// go through simplified uncomparables and map out all baseRates
	finalUncompPairs := []pairs.Pair[int, inefficientComparable.IC[compoundRate]]{}
	nUncRates := len(newUncRateNumbers)
	finalSubRateUncomparableInputs := make([][]int, nUncRates)
	for i := 0; i < nUncRates; i++ {
		r := newUncomparableRates[i]
		newPair := pairs.Pair[int, inefficientComparable.IC[compoundRate]]{
			Value:  r,
			Number: newUncRateNumbers[i],
		}
		finalUncompPairs = append(finalUncompPairs, newPair)
		uncBaseInps := r.getBaseInputs()
		thisRateInputs := make([]int, len(uncBaseInps))
		for j, bi := range uncBaseInps {
			// go through all base inputs, add any new ones
			if index, exists := baseInputSpots[bi]; exists {
				thisRateInputs[j] = index
			} else {
				finalBaseInputs = append(finalBaseInputs, bi)
				thisRateInputs[j] = len(finalBaseInputs) - 1 // TODO: ensure correct: len(finalBaseInputs) or len(finalBaseInputs) - 1?
			}
		}
		finalSubRateUncomparableInputs[i] = thisRateInputs
	}
	// combine into output
	return MultiplicativeRate{
		baseInputs:                finalBaseInputs,
		subRateUncomparableInputs: finalSubRateUncomparableInputs,
		Pairs: power.Pairs[int, Rate, compoundRate]{
			Comparable:   finalComparablePairs,
			Uncomparable: finalUncompPairs,
		},
	}.simplestForm()
}

// RatioInputArgument is only exported for internal use. // TODO: add more info
type RatioInputArgument struct {
	Name         string
	ValueType    reflect.Type
	DefaultValue *reflect.Value
} // TODO: add ratioInputArgument creator

type rateInputs []baseRate

func (ris rateInputs) resolve(inp ...*Options) ([]float64, error) {
	out := make([]float64, len(ris))
	for i, baseInp := range ris {
		baseResult, err := baseInp.Resolve(inp...)
		if err != nil {
			return out, err // TODO: maybe more stuff here?
		}
		out[i] = baseResult
	}
	return out, nil
}

type subRateUncomparableInputs [][]int // TODO: rename
