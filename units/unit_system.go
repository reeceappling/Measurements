package units

var (
	BritishImperial UnitSystemName = "BritishImperial"
	UsCustomary     UnitSystemName = "UsCustomary"
	Metric          UnitSystemName = "Metric"
	Mixed           UnitSystemName = "Mixed"
	FreedomUnits    UnitSystemName = "FreedomUnits"
	CGS             UnitSystemName = "CGS" // TODO: this, currently unused

	allSystems map[UnitSystemName]unitSystemIdentifier
)

func init() {
	allSystems = map[UnitSystemName]unitSystemIdentifier{}
	// Standard systems
	for _, sysName := range []UnitSystemName{
		BritishImperial,
		UsCustomary,
		Metric,
		Mixed,
	} {
		allSystems[sysName] = &UnitSystem{
			name:  sysName,
			units: map[UnitSymbol]Units{},
		}
	}
	// Alias systems
	for alias, baseName := range map[UnitSystemName]UnitSystemName{
		FreedomUnits: UsCustomary,
	} {
		allSystems[alias] = unitSystemAlias{
			name:    baseName,
			aliasOf: allSystems[baseName].baseSystem(),
		}
	}
}

type unitSystemIdentifier interface {
	/*
		BenchmarkPtrs5_5-10      	100000000	        10.33 ns/op
		BenchmarkPtrs20_5-10     	100000000	        11.99 ns/op
		BenchmarkNames5_5-10     	100000000	        10.94 ns/op
		BenchmarkNames20_5-10    	99360021	        12.25 ns/op
	*/
	baseSystem() *UnitSystem // Pointer benchmarks scored better than name benchmarks

}

func systems() map[UnitSystemName]unitSystemIdentifier {
	return allSystems
}

type UnitSystemName string

type UnitSystem struct {
	name  UnitSystemName
	units unitSystemUnitsMap
}

func (sys *UnitSystem) baseSystem() *UnitSystem {
	return sys
}

type unitSystemUnitsMap map[UnitSymbol]Units // TODO: or map[unitSymbol]&baseUnit?

type unitSystemAlias struct {
	name    UnitSystemName
	aliasOf *UnitSystem
}

func (alias unitSystemAlias) baseSystem() *UnitSystem {
	return alias.aliasOf
}
