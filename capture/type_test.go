package capture

import (
	"reflect"
	"testing"
)

func TestTypeOf(t *testing.T) {
	type dynamicValue struct{}
	var value any = dynamicValue{}
	captured, check := TypeOf(value)
	if !check.OK() {
		t.Fatal(check.Why())
	}
	if got := captured.CapturedType(); got != reflect.TypeFor[dynamicValue]() {
		t.Fatalf("TypeOf(interface value) = %v, want %v", got, reflect.TypeFor[dynamicValue]())
	}

	var typedNil *dynamicValue
	captured, check = TypeOf(typedNil)
	if !check.OK() {
		t.Fatal(check.Why())
	}
	if got := captured.CapturedType(); got != reflect.TypeFor[*dynamicValue]() {
		t.Fatalf("TypeOf(typed nil) = %v, want %v", got, reflect.TypeFor[*dynamicValue]())
	}

	if got, check := TypeOf(nil); got != nil || check.OK() {
		t.Fatalf("TypeOf(untyped nil) = %v, %v; want failure", got, check)
	}
}
