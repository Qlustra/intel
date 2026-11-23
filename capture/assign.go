package capture

import (
	"fmt"
	"reflect"

	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/logic"
)

// Assign checks that value is assignable to T and returns its typed copy.
func Assign[T any](value any) (T, logic.Check) {
	var result T
	if value == nil {
		if inspect.IsNil(result) {
			return result, logic.Ok()
		}
		return result, logic.Nope(ValueIsNilCode, "nil is not assignable to destination")
	}
	want := reflect.TypeFor[T]()
	got := reflect.TypeOf(value)
	if !got.AssignableTo(want) {
		return result, logic.Nope(ValueTypeMismatchCode, fmt.Sprintf("%v is not assignable to %v", got, want))
	}
	reflect.ValueOf(&result).Elem().Set(reflect.ValueOf(value))
	return result, logic.Ok()
}
