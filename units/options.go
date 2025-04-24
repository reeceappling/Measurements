package units

var currentServiceOptions = getBuiltInDefaultOptions() // TODO: unsure if we should do this entire page

type ServiceOptions struct {
	// For New Units
	errorOnDuplicateNewBaseUnit bool
	extraInvalidUnitSymbols     []string
	// For conversions
	allowOverWritingConversionRate bool
}

// TODO: reenable and use
//func resolveServiceOptionsForMethodCall(opts ...func(*ServiceOptions)) *ServiceOptions { // TODO: testMe
//	newOptions := copyCurrentServiceOptions()
//	setServiceOptions(newOptions, opts...)
//	return newOptions
//}

func getBuiltInDefaultOptions() *ServiceOptions { // TODO: testMe
	return &ServiceOptions{
		errorOnDuplicateNewBaseUnit:    true,
		allowOverWritingConversionRate: false,
		extraInvalidUnitSymbols:        []string{},
	}
}

func resetCurrentOptionsToBuiltIn() { // TODO: testMe
	currentServiceOptions = getBuiltInDefaultOptions()
}

// TODO: reenable and use
//func copyCurrentServiceOptions() *ServiceOptions { // TODO: testMe
//	return &(*currentServiceOptions)
//}
//func SetCurrentFromDefaultOptions(opts ...func(*ServiceOptions)) { // TODO: testMe
//	newOptions := getBuiltInDefaultOptions()
//	for _, opt := range opts {
//		opt(newOptions)
//	}
//	currentServiceOptions = newOptions
//}
//func modifyCurrentOptions(opts ...func(*ServiceOptions)) { // TODO: testMe
//	setServiceOptions(currentServiceOptions, opts...)
//}
//func setServiceOptions(service *ServiceOptions, opts ...func(*ServiceOptions)) { // TODO: testMe
//	for _, opt := range opts {
//		opt(service)
//	}
//}
