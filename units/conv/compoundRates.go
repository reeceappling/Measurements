package conv

import (
	"github.com/reeceappling/measurements/units/numericalPair/pairs"
	"github.com/reeceappling/measurements/units/numericalPair/pairs/inefficientComparable"
	"github.com/reeceappling/measurements/units/numericalPair/power"
	"github.com/reeceappling/measurements/units/numericalPair/sum"
	"github.com/reeceappling/randomGoStuff/utils/slices"
	"math"
	sliceutils "slices"
)

type compoundRate interface {
	Rate
	resolveWithKnownValues(...float64) (float64, error)
	inefficientComparable.IC[compoundRate]
	getBaseInputs() rateInputs
}

func MultiplicativeIdentityRate() MultiplicativeRate {
	// TODO: ensure this is ok in conjunction with additiveRate
	return MultiplicativeRate{} // This is the multiplicative identity rate (always resolves to 1.0)
}

// TODO: would numerator/denominator work better?
// MultiplicativeRate // TODO: more info.
// can contain base rates, as well as addition rates
type MultiplicativeRate struct {
	baseInputs                rateInputs // TODO: maybe remove comparables from this and just reference them in pairs
	subRateUncomparableInputs subRateUncomparableInputs
	power.Pairs[int, Rate, compoundRate]
} // TODO: is type aliasing ok here or do we need a struct?

func (mr MultiplicativeRate) simplestForm() Rate { // TODO: test
	if len(mr.Pairs.Comparable) == 1 && mr.Pairs.Comparable[0].Number == 1 {
		return mr.Pairs.Comparable[0].Value
	}
	if len(mr.Pairs.Uncomparable) == 1 && mr.Pairs.Uncomparable[0].Number == 1 {
		return mr.Pairs.Uncomparable[0].Value.InefficientlyComparableValue()
	}
	return mr
}

func (mr MultiplicativeRate) getBaseInputs() rateInputs {
	return mr.baseInputs
}

func (mr MultiplicativeRate) resolveWithKnownValues(f ...float64) (float64, error) {
	if len(f) != len(mr.baseInputs) {
		// TODO: err?
	}
	out := 1.0
	for i, comp := range mr.Comparable {

		out *= math.Pow(f[i], float64(comp.Number))
	}
	for i, uncomp := range mr.Uncomparable {
		inputs := slices.Map(mr.subRateUncomparableInputs[i], func(indx int) float64 {
			return f[indx]
		})
		value, err := uncomp.Value.InefficientlyComparableValue().resolveWithKnownValues(inputs...)
		if err != nil {
			return 0, err // TODO: this
		}
		out *= math.Pow(value, float64(uncomp.Number))
	}
	return out, nil
}

func (rates MultiplicativeRate) InefficientlyComparableValue() compoundRate {
	return rates
}

func (rates MultiplicativeRate) EqualTo(ucv inefficientComparable.IC[compoundRate]) bool { // TODO: reimagine this if necesssary
	// Weed out the easy ones
	rate, isCompound := ucv.InefficientlyComparableValue().(MultiplicativeRate)
	if !isCompound {
		return false
	}
	if len(rate.baseInputs) != len(rates.baseInputs) { // TODO: This second or third?
		return false
	}
	if len(rate.Comparable) != len(rates.Comparable) || len(rate.Uncomparable) != len(rates.Uncomparable) { // TODO: This second or third?
		return false
	}
	// Work on the comparables first (less compute-intensive)
	for _, val := range rates.Comparable {
		if !sliceutils.Contains(rate.Comparable, val) { // TODO: more efficient way to do this? likely to small to warrant a map-based solution
			return false
		}
	}
	// then do the uncomparables
	left := make([]pairs.Pair[int, inefficientComparable.IC[compoundRate]], len(rate.Uncomparable))
	copy(left, rate.Uncomparable)
	for i := 0; i < len(rates.Uncomparable); i++ {
		// TODO: ensure checking numbers too
		pairIndex := pairs.IndexForInefficientlyComparable(left, rates.Uncomparable[i].Value)
		if pairIndex == -1 {
			return false
		}
		// ensure numbers are also the same
		if left[pairIndex].Number != rates.Uncomparable[i].Number {
			return false
		}
		// TODO: consider tracking found with an array of bools
		left = append(left[:pairIndex], left[pairIndex+1:]...) // Remove this index from the slice of those left
	}
	return true
}

// asCompoundRate returns the MultiplicativeRate it is a method on
func (rates MultiplicativeRate) AsMultiplicativeRate() MultiplicativeRate {
	return rates
}

func (mr MultiplicativeRate) pow(exp int) MultiplicativeRate {
	return MultiplicativeRate{
		baseInputs:                mr.baseInputs,
		subRateUncomparableInputs: mr.subRateUncomparableInputs,
		Pairs:                     mr.Pairs.Pow(exp),
	}
}

// Resolve // TODO: more info
func (mr MultiplicativeRate) Resolve(opts ...*Options) (float64, error) {
	if len(mr.Pairs.Comparable)+len(mr.Pairs.Uncomparable) == 0 { // TODO: kinda weird way to do it, but works for now
		return 1.0, nil
	}
	baseValues, err := mr.baseInputs.resolve(opts...)
	if err != nil {
		return 0, err // TODO: maybe more stuff here?
	}
	return mr.resolveWithKnownValues(baseValues...)

}

// AdditiveRate // TODO: needs thorough explanation
// AdditiveRate can contain base rates, as well as other AdditiveRate(s), and MultiplicativeRate(s)
type AdditiveRate struct {
	baseInputs                rateInputs
	subRateUncomparableInputs subRateUncomparableInputs // TODO: comparable should probably just be an int
	sum.Pairs[int, Rate, compoundRate]
}

func (ar AdditiveRate) getBaseInputs() rateInputs {
	return ar.baseInputs
}

func (ar AdditiveRate) resolveWithKnownValues(f ...float64) (float64, error) {
	if len(f) != len(ar.baseInputs) {
		// TODO: err?
	}
	out := 0.0
	for i, comp := range ar.Comparable {
		out += f[i] * float64(comp.Number)
	}
	for i, uncomp := range ar.Uncomparable {
		inputs := slices.Map(ar.subRateUncomparableInputs[i], func(indx int) float64 {
			return f[indx]
		})
		value, err := uncomp.Value.InefficientlyComparableValue().resolveWithKnownValues(inputs...)
		if err != nil {
			// TODO: this
		}
		out += value * float64(uncomp.Number)
	}
	return out, nil
}

func (ar AdditiveRate) InefficientlyComparableValue() compoundRate {
	return ar
}

func NewAdditionRate(conversionRates ...Rate) Rate { // TODO: test the hell out of this
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
	return AdditiveRate{
		baseInputs:                finalBaseInputs,
		subRateUncomparableInputs: finalSubRateUncomparableInputs,
		Pairs: sum.Pairs[int, Rate, compoundRate]{
			Comparable:   finalComparablePairs,
			Uncomparable: finalUncompPairs,
		},
	}.simplestForm()
}

func (ar AdditiveRate) EqualTo(ucv inefficientComparable.IC[compoundRate]) bool { // TODO: reimagine this if necesssary
	// Weed out the easy ones first
	rate, isAr := ucv.InefficientlyComparableValue().(AdditiveRate)
	if !isAr {
		return false
	}
	if len(ar.baseInputs) != len(rate.baseInputs) {
		return false
	}
	if len(ar.Comparable) != len(rate.Comparable) || len(ar.Uncomparable) != len(rate.Uncomparable) {
		return false
	}
	// Check comparables (less compute intensive)
	for _, val := range ar.Comparable {
		if !sliceutils.Contains(rate.Comparable, val) {
			return false
		}
	}
	// Then do the more compute-intensive less-comparable values
	left := make([]pairs.Pair[int, inefficientComparable.IC[compoundRate]], len(rate.Uncomparable))
	copy(left, rate.Uncomparable)
	for i := 0; i < len(ar.Uncomparable); i++ {
		pairIndex := pairs.IndexForInefficientlyComparable(left, ar.Uncomparable[i].Value)
		if pairIndex == -1 {
			return false
		}
		// ensure numbers are also the same
		if left[pairIndex].Number != ar.Uncomparable[i].Number {
			return false
		}
		// TODO: consider tracking found with an array of bools
		left = append(left[:pairIndex], left[pairIndex+1:]...) // Remove this index from the slice of those left
	}
	return true
}

func (ar AdditiveRate) Resolve(inp ...*Options) (float64, error) {
	baseValues, err := ar.baseInputs.resolve(inp...)
	if err != nil {
		return 0, err // TODO: maybe more stuff here?
	}
	return ar.resolveWithKnownValues(baseValues...)
}

func (ar AdditiveRate) AsMultiplicativeRate() MultiplicativeRate {
	return MultiplicativeRate{
		baseInputs:                ar.baseInputs,
		subRateUncomparableInputs: [][]int{{0}},
		Pairs: power.Pairs[int, Rate, compoundRate]{
			Comparable:   []pairs.Pair[int, Rate]{},
			Uncomparable: []pairs.Pair[int, inefficientComparable.IC[compoundRate]]{{ar, 1}},
		},
	}
}

func (ar AdditiveRate) simplestForm() Rate { // TODO: test
	nRates := len(ar.Pairs.Comparable) + len(ar.Pairs.Uncomparable)
	if nRates == 0 {
		return nil // TODO: what the hell should we do here
	}
	if nRates == 1 {
		if len(ar.Pairs.Comparable) == 1 {
			if ar.Pairs.Comparable[0].Number == 1 {
				return ar.Pairs.Comparable[0].Value // TODO: ensure ok
			}
			return MultiplicativeRate{ // TODO: ensure ok
				baseInputs:                ar.baseInputs,
				subRateUncomparableInputs: ar.subRateUncomparableInputs,
				Pairs: power.Pairs[int, Rate, compoundRate]{
					Comparable: []pairs.Pair[int, Rate]{
						{
							Value:  ar.Pairs.Comparable[0].Value,
							Number: ar.Pairs.Comparable[0].Number,
						},
					},
				},
			}
		}
		if ar.Pairs.Uncomparable[0].Number == 1 {
			if ar.Pairs.Uncomparable[0].Number == 1 {
				return ar.Pairs.Uncomparable[0].Value.InefficientlyComparableValue() // TODO: ensure ok
			}
			return MultiplicativeRate{ // TODO: ensure ok
				baseInputs:                ar.baseInputs,
				subRateUncomparableInputs: ar.subRateUncomparableInputs,
				Pairs: power.Pairs[int, Rate, compoundRate]{
					Uncomparable: []pairs.Pair[int, inefficientComparable.IC[compoundRate]]{
						{
							Value:  ar.Pairs.Uncomparable[0].Value,
							Number: ar.Pairs.Uncomparable[0].Number,
						},
					},
				},
			}
		}
	}
	return ar
}
