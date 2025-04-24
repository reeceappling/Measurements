package units

import (
	"math/rand"
	"testing"
	"time"
)

type unitSystemIdentifierName interface {
	baseSystem() UnitSystemName
}

type unitSystemIdentifierPtr interface {
	baseSystem() *unitSystemPtr
}

type unitSystemPtr struct {
	name  UnitSystemName
	units unitSystemUnitsMap
}

func (sys *unitSystemPtr) baseSystem() *unitSystemPtr {
	return sys
}

type unitSystemForName struct {
	name  UnitSystemName
	units unitSystemUnitsMap
}

func (u unitSystemForName) baseSystem() UnitSystemName {
	return u.name
}

//type unitSystemUnitsMap map[units.UnitSymbol]units.Units // TODO: or map[unitSymbol]&baseUnit?

type unitSystemAliasPtr struct {
	name    UnitSystemName
	aliasOf *unitSystemPtr
}

func (alias unitSystemAliasPtr) baseSystem() *unitSystemPtr {
	return alias.aliasOf
}

type unitSystemAliasName struct {
	name    UnitSystemName
	aliasOf UnitSystemName
}

func (u unitSystemAliasName) baseSystem() UnitSystemName {
	return u.aliasOf
}

// TODO: TEST searching through various sizes of map[UnitSystemName]unitSystemIdentifier where unitSystemAlias is pointer vs UnitSystemName

func benchmarkPtrs(bases, aliases int, b *testing.B) {
	allSystemsPtrs := map[UnitSystemName]unitSystemIdentifierPtr{}
	baseNames := createNames(bases)
	aliasNames := createNames(aliases)
	// Create bases
	for _, name := range baseNames {
		allSystemsPtrs[name] = &unitSystemPtr{
			name:  name,
			units: nil,
		}
	}
	for _, alias := range aliasNames {
		randBase := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(baseNames))
		allSystemsPtrs[alias] = &unitSystemAliasPtr{
			name:    alias,
			aliasOf: allSystemsPtrs[baseNames[randBase]].baseSystem(),
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res := allSystemsPtrs[aliasNames[i%len(aliasNames)]].baseSystem()
		if res.name == "/////\\" {
		}
	}
}

// from fib_test.go
func benchmarkNames(bases, aliases int, b *testing.B) {
	allSystemsNames := map[UnitSystemName]unitSystemIdentifierName{}
	baseNames := createNames(bases)
	aliasNames := createNames(aliases)
	// Create bases
	for _, name := range baseNames {
		allSystemsNames[name] = &unitSystemForName{
			name:  name,
			units: nil,
		}
	}
	for _, alias := range aliasNames {
		randBase := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(baseNames))
		allSystemsNames[alias] = &unitSystemAliasName{
			name:    alias,
			aliasOf: allSystemsNames[baseNames[randBase]].baseSystem(),
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		aliasName := aliasNames[i%len(aliasNames)]
		res := allSystemsNames[aliasName].baseSystem()
		if _, ok := allSystemsNames[res]; !ok {
			panic("alias name not found: " + aliasName)
		}
	}
}

func BenchmarkPtrs5_5(b *testing.B)   { benchmarkPtrs(5, 5, b) }
func BenchmarkPtrs20_5(b *testing.B)  { benchmarkPtrs(20, 5, b) }
func BenchmarkNames5_5(b *testing.B)  { benchmarkPtrs(5, 5, b) }
func BenchmarkNames20_5(b *testing.B) { benchmarkPtrs(20, 5, b) }

func createNames(bases int) []UnitSystemName {
	out := make([]UnitSystemName, bases)
	for i := 0; i < bases; i++ {
		out[i] = UnitSystemName(StringWithCharset(20))
	}
	return out
}

func StringWithCharset(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()_+-=~[]{};:,.<>/?"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.New(
			rand.NewSource(time.Now().UnixNano())).Intn(len(charset))]
	}
	return string(b)
}
