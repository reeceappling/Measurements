package units

import "sync"

const (
	//According to the National Institute of Standards and Technology (NIST), the current accepted value for NA is:
	//NA = (6.0221415 ± 0.0000010) × 10^23 (https://www.americanscientist.org/article/an-exact-value-for-avogadros-number)
	AvagadrosNumber = 6.02214076e23 // TODO: is this fully correct? Is this the best spot for it?
	// TODO: maybe make public functions that can grab each of these symbols?
	molSym         = "mol"
	dozenSym       = "doz"  // TODO: ensure this is ok, maybe dz instead?
	bakersDozenSym = "bDoz" // TODO: change?
)

var unitless = &baseUnit{
	sym:         "",
	conversions: conversionMap{},
	RWMutex:     sync.RWMutex{},
}

func Unitless() *baseUnit {
	return unitless
}

func defineStandardCountUnits() {
	_, _ = NewUnitWithEquivalence(1, molSym, AvagadrosNumber, unitless)
	_, _ = NewUnitWithEquivalence(1, dozenSym, 12, unitless)
	_, _ = NewUnitWithEquivalence(1, bakersDozenSym, 13, unitless)
	// TODO: define other count units
}
