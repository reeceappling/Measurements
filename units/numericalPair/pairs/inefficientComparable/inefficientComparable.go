package inefficientComparable

// TODO: Explain how the user must define comparability, this library does not handle it
//
//go:generate mockery --name IC
type IC[T any] interface { // TODO: rename
	InefficientlyComparableValue() T
	EqualTo(IC[T]) bool
}
