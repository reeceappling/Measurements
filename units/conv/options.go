package conv

import (
	"github.com/reeceappling/randomGoStuff/utils"
	"maps"
	"slices"
)

// Options should NOT be created directly, please create options with NewOptions().
// Options is only a public alias so that it can be referenced in function and (possibly) type definitions
type Options map[string]any

// NewOptions creates a new, blank, *Options
func NewOptions(fields ...optionField) *Options { // TODO: test that no fields returns nil *Options, with fields returns *Options with all fields in it
	if len(fields) == 0 {
		return nil
	}
	return utils.Pointer[Options](nil).addFields(false, fields)
}

// With returns a new Options map, without modifying the original, adding (or replacing) an option field for each optionField provided
func (opts *Options) set(field optionField) { // TODO: test, use
	(*opts)[field.Key] = field.val
}

// With returns a new Options map, without modifying the original, adding (or replacing) an option field for each optionField provided
func (opts *Options) With(fields ...optionField) *Options { // TODO: test (likely unnecessary)
	return opts.addFields(true, fields)
}

// Add modifies an existing Options map, adding (or replacing) an option field for each optionField provided
func (opts *Options) Add(fields ...optionField) *Options { // TODO: test (likely unnecessary)
	return opts.addFields(false, fields)
}

// Without creates a copy of the options, less any of the named fields.
func (opts *Options) Without(fieldNames ...string) *Options { // TODO: test (likely unnecessary)
	return opts.removeFields(true, fieldNames)
}

// Less modifies an existing Options map, returning its Options less any of the named fields
func (opts *Options) Less(fieldNames ...string) *Options { // TODO: test (likely unnecessary)
	return opts.removeFields(false, fieldNames)
}

// addFields either modifies its receiver, OR clones it and does work on the clone
// The underlying of the returned *Options are equivalent to the underlying of the receiver less any of the named fields, if they exist.
func (opts *Options) addFields(returnNew bool, fields []optionField) *Options {
	if len(fields) == 0 {
		if opts == nil {
			return nil // TODO: test
		}
		return opts.cloneIfNeeded(returnNew) // TODO: test
	}
	out := opts.cloneIfNeededPossiblyNil(returnNew)
	for field := range slices.Values(fields) {
		(*out)[field.Key] = field.val
	}
	return out // TODO: test
}

// removeFields either modifies its receiver, OR clones the receiver and does work on that *Options
// The working *Options are returned, having added (or replaced) the option field for each optionField provided
func (opts *Options) removeFields(returnNew bool, fieldNames []string) *Options {
	if opts == nil {
		return nil // TODO: test
	}
	out := opts.cloneIfNeeded(returnNew)
	for fieldName := range slices.Values(fieldNames) {
		delete(*out, fieldName)
	}
	if len(*out) == 0 {
		return nil // TODO: test
	}
	return out // TODO: test that output is as expected
}

// cloneIfNeeded clones the reciever if doClone==true, otherwise returns the original receiver.
//
// Panics if receiver is nil
func (opts *Options) cloneIfNeeded(doClone bool) *Options {
	if doClone {
		return opts.Clone()
	}
	return opts // TODO: test modifying result does modify original receiver
}

// cloneIfNeededPossiblyNil clones the reciever if doClone==true, otherwise returns the original receiver.
//
// if the input *Options are nil, a new pointer to an empty Options is returned
func (opts *Options) cloneIfNeededPossiblyNil(doClone bool) *Options {
	if opts == nil { // Tested with NewOptions
		return utils.Pointer(Options(map[string]any{})) // TODO: test that nil options returns a pointer to an empty options
	}
	return opts.cloneIfNeeded(doClone) // TODO: test non-nil options returns a pointer to an equivalent underlying
}

// Clone clones an options pointer to another pointer with an equivalent underlying value
func (opts *Options) Clone() *Options {
	if opts == nil {
		return nil // TODO: test
	}
	return utils.Pointer(maps.Clone(*opts)) // TODO: test that modifying doesnt change original
}

// optionField is a key-value genericPair with a Key string and a val that is comparable.
//
// optionField is private so that OptionField() can always ensure that val is comparable.
type optionField struct {
	Key string
	val interface{}
}

// OptionField constructs a new optionField
func OptionField[T comparable](key string, val T) optionField { // TODO: is returning unexported type ok?
	return optionField{
		Key: key,
		val: val,
	}
}

//// ensureNoDuplicates panics if any duplicate field names are present in the inputs
//func ensureNoDuplicates(fields ...optionField) { // TODO: use?
//	keys := sliceutils.Map(fields, func(field optionField) string {
//		return field.Key
//	})
//	for i, key := range keys {
//		for j := i + 1; j < len(keys); j++ {
//			if key == keys[j] {
//				panic("duplicate keys in option fields") // TODO: ensure ok
//			}
//		}
//	}
//}

//func CombineOptions(opts ...*Options) (*Options, error) { // TODO: use this
//	if len(opts) == 0 {
//		return nil, nil // TODO: test
//	}
//	out := Options(map[string]any{})
//	nNonNil := 0
//	var firstNonNil *Options = nil
//	for opt := range slices.Values(opts) {
//		if opt == nil {
//			continue
//		}
//		if nNonNil == 0 {
//			firstNonNil = opt
//		}
//		nNonNil++
//		for name, val := range *opt {
//			if currentVal, exists := out[name]; exists {
//				if val == currentVal {
//					continue
//				}
//				return nil, errors.New("value collision in CombineOptions") // TODO: test
//			}
//			out[name] = val
//		}
//	}
//	if nNonNil <= 1 {
//		return firstNonNil, nil
//	}
//	return &out, nil // TODO: test
//}
