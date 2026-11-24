// Package intel provides the starting points for checked interrogation of Go
// types and values. Use [CaptureType] for structural questions and
// [CaptureValue] when a particular value or its storage matters.
package intel

import (
	"github.com/qlustra/intel/capture"
	"github.com/qlustra/intel/logic"
)

// CaptureType describes the exact Go type T without creating an instance.
func CaptureType[T any]() *capture.TypeInterrogator {
	typeInterrogator := capture.Type[T]()
	return typeInterrogator
}

// CaptureValue captures existing storage from a non-nil pointer, or a shallow
// copy of a non-pointer value. See [capture.Value] for the storage and nil
// input rules.
func CaptureValue(v any) (*capture.ValueInterrogator, logic.Check) {
	return capture.Value(v)
}
