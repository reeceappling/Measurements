package conv

import (
	"github.com/reeceappling/measurements/units/numericalPair/pairs"
	"github.com/reeceappling/measurements/units/numericalPair/pairs/inefficientComparable"
	"github.com/reeceappling/measurements/units/numericalPair/power"
	"github.com/reeceappling/randomGoStuff/utils"
)

type baseRate interface {
	Rate
	resolveBaseRate(inp ...*Options) (float64, error)
}

// StaticRate // TODO: more info
type StaticRate float64

func (ratio *StaticRate) resolveBaseRate(inp ...*Options) (float64, error) {
	if ratio == nil {
		return 0, ErrNilRatePointer // TODO: test, join err
	}
	return float64(*ratio), nil
}

// NewStaticRate creates a static rate meeting the Rate interface
func NewStaticRate(f float64) *StaticRate {
	if f == 0.0 {
		return nil // TODO: test
	}
	return utils.Pointer(StaticRate(f))
}

// Resolve resolves the underlying float for a *StaticRate. This will only return an error for a nil *StaticRate
func (ratio *StaticRate) Resolve(...*Options) (float64, error) {
	return ratio.resolveBaseRate()
}

// asCompoundRate returns a MultiplicativeRate containing the *StaticRate to the first power
func (rate *StaticRate) AsMultiplicativeRate() MultiplicativeRate {
	return MultiplicativeRate{
		baseInputs:                rateInputs{rate},
		subRateUncomparableInputs: subRateUncomparableInputs{}, // TODO: likely unnecessary
		Pairs: power.Pairs[int, Rate, compoundRate]{
			Comparable:   []pairs.Pair[int, Rate]{{rate, 1}},
			Uncomparable: []pairs.Pair[int, inefficientComparable.IC[compoundRate]]{}, // TODO: likely unnecessary
		},
	}
}

// DynamicRate // TODO: more info (Don't create directly)
type DynamicRate struct { // Also encapsulates time-dependent conversions
	resolver            *RatioFunc
	ratioInputArguments []RatioInputArgument
}

func (rate *DynamicRate) resolveBaseRate(inp ...*Options) (float64, error) {
	return rate.Resolve(inp...)
}

// NewDynamicRate creates a dynamic rate with TODO: more info
// ONLY TO BE USED INTERNALLY TODO: more info
func NewDynamicRate(resolver *RatioFunc, args ...RatioInputArgument) *DynamicRate {
	return &DynamicRate{
		resolver:            resolver,
		ratioInputArguments: args,
	}
}

// Resolve resolves the rate of a *DynamicRate, or errors
func (rate *DynamicRate) Resolve(inpt ...*Options) (float64, error) {
	if rate == nil {
		return 0.0, ErrNilRatePointer
	}
	opts := &Options{}
	if len(inpt) != 0 {
		opts = inpt[0]
	}
	for _, arg := range rate.ratioInputArguments {
		if arg.DefaultValue != nil {
			if _, exists := (*opts)[arg.Name]; !exists {
				opts = opts.Add(OptionField(arg.Name, arg.DefaultValue)) // TODO: can we use add here?
			}
		}
	}
	return rate.resolver.resolve(opts)
}

// asCompoundRate returns a MultiplicativeRate containing the *DynamicRate to the first power
func (rate *DynamicRate) AsMultiplicativeRate() MultiplicativeRate {
	return MultiplicativeRate{
		baseInputs:                rateInputs{rate},
		subRateUncomparableInputs: subRateUncomparableInputs{},
		Pairs: power.Pairs[int, Rate, compoundRate]{
			Comparable: []pairs.Pair[int, Rate]{{rate, 1}},
		},
	}
}

// RatioFunc // TODO: more info
type RatioFunc func(inp *Options) (float64, error)

// resolve resolves a *RatioFunc given input Options
func (rf *RatioFunc) resolve(inpt *Options) (float64, error) {
	if rf == nil {
		return 0.0, ErrNilRateResolver
	}
	return (*rf)(inpt)
}

// ErroneousRate // TODO: more info
// TODO: add Unwrap() error or Unwrap() []error
type ErroneousRate struct { // TODO: test
	Err error
}

func (rate ErroneousRate) resolveBaseRate(inp ...*Options) (float64, error) {
	return rate.Resolve(inp...)
}

// Resolve // TODO: this
func (rate ErroneousRate) Resolve(options ...*Options) (float64, error) { // TODO: test, use
	return 0, rate.Err
}

// asCompoundRate returns a MultiplicativeRate containing the *DynamicRate to the first power
func (rate ErroneousRate) AsMultiplicativeRate() MultiplicativeRate { // TODO: test?
	return MultiplicativeRate{
		baseInputs:                rateInputs{rate},
		subRateUncomparableInputs: subRateUncomparableInputs{},
		Pairs: power.Pairs[int, Rate, compoundRate]{Comparable: []pairs.Pair[int, Rate]{{rate, 1}},
			Uncomparable: []pairs.Pair[int, inefficientComparable.IC[compoundRate]]{},
		},
	}
}
