package intel

import (
	"github.com/qlustra/intel/types"
	"github.com/qlustra/intel/value"
)

func CaptureType[T any]() *types.TypeInterrogator {
	typeInterrogator := types.CaptureType[T]()
	return typeInterrogator
}

func CaptureValue(v any) *value.ValueInterrogator {
	valueInterrogator := value.CaptureValue(v)
	return valueInterrogator
}
