package units

import (
	"math"
)

const (
	byteSym = "B" // TODO: might overlap with something else?
	Mu      = "μ"
)

func init() {
	// TODO: ADD Aliases
	// TODO: ADD OBSCURE US Customary Units
	// TODO: ADD Obscure British Imperial Units
	// TODO: ADD Obscure Plasmas Units
	// TODO: ADD CURRENCIES $, pound, euro, etc (requires time-based conversions)
	// TODO: Tesla(T) Gauss(G)
	defineBuiltInUnits()
}

func defineBuiltInUnits() {
	allUnits.addNode(unitless, unitless.Signature()) // Define unitless in units graph
	defineMetricUnits()
	// TODO: add other metric units not in standard engineering prefixes
	// TODO: add freedom units
	// TODO: define special units (counts:Mol,dozen. Temperatures)
	defineStandardCountUnits()
	defineByteUnits()
}

func newMetricBaseUnits(baseSymbol, name, measures string) Units {
	newMetricBase, _ := NewBaseUnit(baseSymbol)
	prefixes := map[string]struct {
		exp  int
		name string
	}{
		"q": {exp: -30, name: "quecto"},
		"r": {exp: -27, name: "ronto"},
		"y": {exp: -24, name: "yocto"},
		"z": {exp: -21, name: "zepto"},
		"a": {exp: -18, name: "atto"},
		"f": {exp: -15, name: "femto"},
		"p": {exp: -12, name: "pico"},
		"n": {exp: -9, name: "nano"},
		// Mu: {exp: -6, name: "micro"}, // TODO: possibly add as alias
		"u": {exp: -6, name: "micro"},
		"m": {exp: -3, name: "milli"},
		"k": {exp: 3, name: "kilo"},
		"M": {exp: 6, name: "mega"},
		"G": {exp: 9, name: "giga"},
		"T": {exp: 12, name: "tera"},
		"P": {exp: 15, name: "peta"},
		"E": {exp: 18, name: "exa"},
		"Z": {exp: 21, name: "zetta"},
		"Y": {exp: 24, name: "yotta"},
		"R": {exp: 27, name: "ronna"},
		"Q": {exp: 30, name: "quetta"},
	}
	// Add non-3-power standard metric units
	centi := struct {
		exp  int
		name string
	}{exp: -2, name: "centi"}
	switch baseSymbol {
	// TODO: add any other non-engineering metric cases
	case "m":
		prefixes["c"] = centi
	default:
		// Do not add or remove any prefixes as a default
	}
	for prefix, data := range prefixes {
		newSig := prefix + baseSymbol
		_, err := NewUnitWithEquivalence(1, newSig, math.Pow10(data.exp), newMetricBase)
		if err != nil {
			panic("failed to create " + newSig + " because " + err.Error()) // TODO: this
		}
	}
	return newMetricBase
}

func defineMetricUnits() {
	// Base units and their prefixes
	m := newMetricBaseUnits("m", "meter", "length") // TODO: this
	//KBase := newMetricBaseUnits("K", "kelvin", "absolute temperature") // TODO: this
	s := newMetricBaseUnits("s", "second", "time") // TODO: this
	_ = newMetricBaseUnits("g", "gram", "mass")    // TODO: this
	//cdBase := newMetricBaseUnits("cd", "candela", "luminous intensity") // TODO: this
	A := newMetricBaseUnits("A", "ampere", "current") // TODO: this
	// standard/compound
	// Get necessary base units for standard/compound
	kg, _ := allUnits.Get("kg")
	// Create standard/compound

	C, _ := NewUnitWithEquivalence(1, "C", 1, A.Mul(s))

	N, _ := NewUnitWithEquivalence(1, "N", 1, kg.Mul(m).Div(s.Pow(2)))
	J, _ := NewUnitWithEquivalence(1, "J", 1, N.Mul(m))
	_, _ = NewUnitWithEquivalence(1, "V", 1, J.Div(C))
	_, _ = NewUnitWithEquivalence(1, "ohm", 1, kg.Mul(m.Pow(2)).Mul(s.Pow(-3)).Mul(A.Pow(-2)))
	_, _ = NewUnitWithEquivalence(1, "hz", 1, s.Pow(-1))
	//if ohm == nil || C == nil || J == nil || V == nil || hz == nil { // TODO: deleteme
	//	panic("NIL THING")
	//}

}

func defineByteUnits() {
	byteUnit, _ := NewBaseUnit(byteSym)
	_, _ = NewUnitWithEquivalence(8, "b", 1, byteUnit)
	_, _ = NewUnitWithEquivalence(4, "nibble", 1, byteUnit) // TODO: aliases: nybble, nyble, and/or nybl to match the spelling of byte
	tenPowBytePrefixes := []string{
		"Ki", //kibi
		"Mi", //mebi
		"Gi", //gibi
		"Ti", //tebi
		"Pi", //pebi
		"Ei", //exbi
		"Zi", //zebi
		"Yi", //yobi
	}
	for i, prefix := range tenPowBytePrefixes {
		bytesPerPrefix := math.Pow(1024, float64(i)) // TODO: would math.Pow(2, float64(i*10)) be faster?
		_, _ = NewUnitWithEquivalence(1, prefix+byteSym, bytesPerPrefix, byteUnit)
		// TODO: make these retrievable via public fxns?
	}
}
