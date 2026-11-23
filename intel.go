package intel

import (
	"github.com/qlustra/intel/capture"
	"github.com/qlustra/intel/logic"
)

func CaptureType[T any]() *capture.TypeInterrogator {
	typeInterrogator := capture.Type[T]()
	return typeInterrogator
}

func CaptureValue(v any) (*capture.ValueInterrogator, logic.Check) {
	return capture.Value(v)
}
