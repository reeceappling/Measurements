package measurement

import (
	"github.com/reeceappling/measurements/units"
	"github.com/reeceappling/measurements/units/conv"
	sliceutils "github.com/reeceappling/randomGoStuff/utils/slices"
)

type Measurement struct {
	val   conv.Rate   // this is a conv.dynamicRate or conv.MultiplicativeRate // TODO: reword this // TODO: Like when accruing interest over time
	units units.Units // TODO: like for currency conversion rates
}

func (d Measurement) rate() conv.Rate { // Convenience function because val might sound weird in some contexts // TODO: do we really want to keep this?
	return d.val
}

func (d Measurement) Value(opts ...*conv.Options) (float64, error) {
	if len(opts) != 0 {
		return d.val.Resolve(opts[0])
	}
	return d.val.Resolve(conv.NewOptions())
}

func (d Measurement) ConvertTo(units units.Units) Measurement {
	return Measurement{
		val:   conv.MultiplyRates(d.val, d.units.ConversionRateTo(units)),
		units: units,
	}
}

func (d Measurement) Mul(mmts ...Measurement) Measurement {
	allValueRates := []conv.Rate{d.val}
	allValueRates = append(allValueRates,
		sliceutils.Map(mmts, func(mmt Measurement) conv.Rate {
			return mmt.rate()
		})...,
	)
	return Measurement{
		val: conv.MultiplyRates(allValueRates...),
		units: d.units.Mul(
			sliceutils.Map(mmts, func(mmt Measurement) units.Units {
				return mmt.units
			})...,
		),
	}
}

func (d Measurement) Div(mmt Measurement) Measurement {
	return d.Mul(mmt.Pow(-1))
}

func (d Measurement) Plus(mmts ...Measurement) Measurement {
	// ensure measurements are compatible
	compatibleRates := []conv.Rate{d.val}
	for _, mmt := range mmts {
		if mmt.units == d.units {
			compatibleRates = append(compatibleRates, mmt.val)
			continue
		}
		// if measurement units are not the same, try converting
		convRate := mmt.units.ConversionRateTo(d.units)
		if errRate, isErr := convRate.(conv.ErroneousRate); isErr {
			// TODO: what here?
			return Measurement{
				val:   errRate,
				units: d.units,
			}
		}
		compatibleRates = append(compatibleRates, mmt.val, convRate)
	}

	otherValAsCorrectUnits := conv.MultiplyRates(compatibleRates...) // TODO: ensure MultiplyRates works here as this is an addition function
	return Measurement{
		val:   conv.NewAdditionRate(otherValAsCorrectUnits), // TODO: fixMe
		units: d.units,
	}
}

func (d Measurement) Minus(mmts ...Measurement) Measurement {
	negatives := sliceutils.Map(mmts, func(mmt Measurement) Measurement {
		return Measurement{
			val:   conv.MultiplyRates(mmt.rate(), conv.NegativeRate),
			units: mmt.units,
		}
	})
	return d.Plus(negatives...)
}

func (d Measurement) Pow(exp int) Measurement {
	return Measurement{
		val:   conv.RatePow(d.val, exp),
		units: d.units.Pow(exp),
	}
}

//func addLikeMeasurements(mmts ...Measurement) []Measurement { // TODO: fixme! or just don't use!
//	// TODO:
//	sigs := map[units.UnitSignature][]int{}
//	for i, mmt := range mmts {
//		sig := mmt.units.Signature()
//		if currentIndices, exists := sigs[sig]; exists {
//			sigs[sig] = append(currentIndices, i)
//			continue
//		}
//		sigs[sig] = []int{i}
//	}
//	out := make([]Measurement, len(sigs))
//unitsLoop:
//	for _, indices := range sigs {
//		unitsOut := mmts[indices[0]].units
//		// TODO: Get unique rates for each index (except static
//		staticRate := 0.0
//		for _, i := range indices {
//			rate := mmts[i].rate()
//			switch rate.(type) {
//			case *conv.StaticRate:
//				rateVal, _ := rate.(*conv.StaticRate).Resolve()
//				staticRate += rateVal
//			case conv.MultiplicativeRate:
//				// TODO: add compound rates
//			case *conv.DynamicRate:
//				// TODO: add dynamic rates
//			case conv.AdditiveRate:
//				// TODO: add addition rates
//			case conv.ErroneousRate:
//				// TODO: if erroneousRate, err
//			default:
//				// TODO: panic?
//			}
//
//		}
//		// TODO: add all rates
//		out = append(out, Measurement{
//			val:   nil,
//			units: mmts[indices[0]].units,
//		})
//	}
//	// TODO: this
//	uniqueUnits := unique.Make(sliceutils.Map(mmts, func(mmt Measurement) units.UnitSignature {
//		return mmt.units.Signature()
//	}))
//}

func NewStaticMeasurement(value float64, units units.Units) Measurement {
	return Measurement{conv.NewStaticRate(value), units}
}

func NewDynamicMeasurement(value conv.Rate, units units.Units) Measurement {
	return Measurement{value, units}
}
