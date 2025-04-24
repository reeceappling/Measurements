package units

import (
	"errors"
	"fmt"
	"github.com/reeceappling/measurements/units/conv"
	"reflect"
	"slices"
)

var (
	errBadDynamicConversionInputDefaultType = errors.New("type of default value did not match type of function input")
	errSettingConvOnUnits                   = errors.New("failed to set staticConversion on unit")
	errConversionSignatureDNE               = errors.New("staticConversion does not contain that Signature")
	errRatioFunctionMissingInput            = errors.New("missing required input value to conversionRate function")
	errRatioFunctionNil                     = errors.New("conversionRate function must not be nil")
	errRatioFunctionOutputs                 = errors.New("conversionRate function must have return type (float64, error)")
	errRatioFunctionNumOutputs              = errors.New("conversionRate function did not have exactly 2 outputs")
	errRatioFunctionSigNoFloat              = errors.New("conversionRate function output at index 0 must be float64")
	errRatioFunctionSigNoError              = errors.New("conversionRate function output at index 1 must be error")
	errRatioArgumentsNumberMismatch         = errors.New("number of conversionRate function input parameters does not match the number of input fields provided")
	errRatioFunctionAssertingError          = errors.New("failed type assertion to error for conversionRate function output")
	errRatioFunctionAssertingFloat          = errors.New("failed type assertion to float64 for conversionRate function output")
)

type conversion struct {
	conversionPair
	rate conv.Rate
}

// rateFrom resolves the directionally-dependent rate. If going backwards, the conversionRate is inverted
func (convers *conversion) rateFrom(sig UnitSignature) conv.Rate {
	if convers == nil {
		panic("cannot do rateFrom() on nil conversion receiver") // TODO: this ok?
	}
	switch sig {
	case convers.from.Signature():
		return convers.rate
	case convers.to.Signature():
		return conv.RatePow(convers.rate, -1)
	default:
		panic("Signature provided did not match either of the provided units") // TODO: is this ok?
	}
}

// setConversionOnNodePair tries to set each leg of a conversion on the conversionMap(s) of associated Units
func (convers *conversion) setConversionOnNodePair() error { // TODO: ensure we aren't doing deleteConversionFromNodePairBySigs after when used
	if convers == nil {
		return errors.New("cannot setConversionOnNodePair for a nil conversion receiver") // TODO: moveMe
	}
	To, From := convers.To(), convers.From()
	toSig, fromSig := convers.To().Signature(), convers.From().Signature()
	_ = To.setConversion(fromSig, convers)
	_ = From.setConversion(toSig, convers)
	//errOnTo := To.setConversion(fromSig, convers) // TODO: unsure if needed, can only error if conversion already exists
	//errOnFrom := From.setConversion(toSig, convers)
	//if err := errors.Join(errOnTo, errOnFrom); err != nil {
	//	deleteConversionFromNodePairBySigs(To, From, toSig, fromSig)
	//	return errors.Join(err, errSettingConvOnUnits)
	//}
	return nil
}

func (convers *conversion) deleteConversionFromNodePair() {
	if convers == nil {
		return
	}
	To, From := convers.To(), convers.From()
	toSig, fromSig := To.Signature(), From.Signature()
	deleteConversionFromNodePairBySigs(To, From, toSig, fromSig)
}

// conversionMap is an internal map on each Units that stores conversion rates to/from other Units
type conversionMap map[UnitSignature]*conversion

// TODO: time-dependent conversions (dynamicConversion)

type conversionPair struct {
	from, to Units
}

func (conv conversionPair) From() Units {
	return conv.from
}

func (conv conversionPair) To() Units {
	return conv.to
}

func newConversionPair(from, to Units) conversionPair {
	return conversionPair{from, to}
}

type ratioInputField struct {
	name         string
	defaultValue *interface{}
}

// TODO: ratioFunction is of the format func(a,b,c)(float64,error) and inputFields are of the format field("aName",val),field("bName",val),field("cName",val)
func newDynamicConversion(To, From Units, ratioFunction interface{}, inputFields ...ratioInputField) (*conversion, error) {
	if ratioFunction == nil {
		return nil, errRatioFunctionNil
	}
	handler := reflect.ValueOf(ratioFunction)
	handlerType := reflect.TypeOf(ratioFunction)
	if handlerType.Kind() != reflect.Func {
		err := errors.New("invalid ratioFunction handler type kind")                               // TODO: moveMe
		errCustom := fmt.Errorf("handler kind (%s) is not (%s)", handlerType.Kind(), reflect.Func) // TODO: this, rename
		return nil, errors.Join(err, errCustom)
	}

	// check Signature output
	if handlerType.NumOut() != 2 {
		return nil, errors.Join(errRatioFunctionNumOutputs, errRatioFunctionOutputs) // TODO: test
	}
	if handlerType.Out(0) != reflect.TypeFor[float64]() {
		return nil, errors.Join(errRatioFunctionSigNoFloat, errRatioFunctionOutputs) // TODO: test
	}
	if handlerType.Out(1) != reflect.TypeFor[error]() {
		return nil, errors.Join(errRatioFunctionSigNoError, errRatioFunctionOutputs) // TODO: test
	}

	// construct arguments
	var ratioInputArgumentOrder []conv.RatioInputArgument
	numArguments := handlerType.NumIn()
	if numArguments != len(inputFields) {
		return nil, errRatioArgumentsNumberMismatch // TODO: test
	}

	for i, inputField := range inputFields {
		valueType := handlerType.In(i)
		defaultValuePtr := inputField.defaultValue
		newArg := conv.RatioInputArgument{
			Name:         inputField.name,
			ValueType:    valueType,
			DefaultValue: nil,
		}
		required := defaultValuePtr == nil
		if !required {
			defaultValue := *defaultValuePtr
			if reflect.TypeOf(defaultValue) != valueType {
				return nil, errBadDynamicConversionInputDefaultType // TODO: testme
			}
			defaultInputValue := reflect.ValueOf(defaultValue)
			newArg.DefaultValue = &defaultInputValue
		}
		ratioInputArgumentOrder = append(ratioInputArgumentOrder, newArg)
	}

	ratioFunctionWrapper := conv.RatioFunc(
		func(inp *conv.Options) (float64, error) {
			var args []reflect.Value
			for expArg := range slices.Values(ratioInputArgumentOrder) {
				if inp != nil {
					if val, exists := (*inp)[expArg.Name]; exists {
						args = append(args, reflect.ValueOf(val))
						continue
					}
				}
				if expArg.DefaultValue == nil { // If the argument is required
					missingInputFieldErr := fmt.Errorf(`missing input for required field: %s`, expArg.Name)
					return 0, errors.Join(errRatioFunctionMissingInput, missingInputFieldErr) // TODO: test
				}
				args = append(args, *expArg.DefaultValue)
			}

			response := handler.Call(args)
			errVal, ok := response[1].Interface().(error)
			if !ok {
				return 0, errRatioFunctionAssertingError // TODO: test
			}
			if errVal != nil {
				return 0, errVal // TODO: test
			}
			floatVal, ok := response[0].Interface().(float64)
			if !ok {
				return 0, errRatioFunctionAssertingFloat // TODO: test.
			}
			return floatVal, nil // TODO: test
		})

	return &conversion{
		conversionPair: newConversionPair(From, To),
		rate:           conv.NewDynamicRate(&ratioFunctionWrapper, ratioInputArgumentOrder...),
	}, nil
}

func deleteConversionFromNodePairBySigs(To, From Units, toSig, fromSig UnitSignature) {
	To.removeConversion(fromSig)
	From.removeConversion(toSig)
}
