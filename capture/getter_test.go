package capture

import (
	"reflect"
	"strings"
	"testing"
)

type getterValue string
type getterInterface interface{ Get() }
type getterReceiver struct{ calls int }
type getterPrimitive int
type getterPointerCycle *getterPointerCycle

func (r *getterReceiver) Value() getterValue       { r.calls++; return "zero" }
func (r *getterReceiver) Nil() getterInterface     { return nil }
func (r *getterReceiver) NeedsInput(string) string { return "" }
func (r *getterReceiver) Many() (string, error)    { return "", nil }
func (r *getterReceiver) Explode() string          { panic("boom") }
func (p *getterPrimitive) Value() getterValue      { return "primitive" }

func TestZeroReceiverGetterRetainsDeclaredResultType(t *testing.T) {
	receiver, check := Type[**getterReceiver]().ZeroReceiver()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	if got, want := receiver.Type().CapturedType(), reflect.TypeFor[*getterReceiver](); got != want {
		t.Fatalf("Type() = %v, want %v", got, want)
	}
	method, check := receiver.Method("Value")
	if !check.OK() {
		t.Fatal(check.Why())
	}
	if method.ReceiverType().CapturedType() != reflect.TypeFor[*getterReceiver]() || method.InputArity() != 0 || method.OutputArity() != 1 {
		t.Fatal("unexpected method metadata")
	}
	resultType, check := method.Result(0)
	if !check.OK() || resultType.CapturedType() != reflect.TypeFor[getterValue]() {
		t.Fatalf("Result(0) = %v, %v", resultType, check)
	}
	getter, check := receiver.Getter(method)
	if !check.OK() {
		t.Fatal(check.Why())
	}
	for range 2 {
		result, check := getter.Read()
		if !check.OK() {
			t.Fatal(check.Why())
		}
		if result.Type().CapturedType() != reflect.TypeFor[getterValue]() || result.Get() != getterValue("zero") {
			t.Fatalf("unexpected result: type=%v value=%#v", result.Type().CapturedType(), result.Get())
		}
	}
}

func TestZeroReceiverGetterAllowsNilInterfaceResult(t *testing.T) {
	receiver, check := Type[getterReceiver]().ZeroReceiver()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	method, check := receiver.Method("Nil")
	if !check.OK() {
		t.Fatal(check.Why())
	}
	getter, check := receiver.Getter(method)
	if !check.OK() {
		t.Fatal(check.Why())
	}
	result, check := getter.Read()
	if !check.OK() || result.Get() != nil {
		t.Fatalf("Read() = (%#v, %v), want (nil, success)", result, check)
	}
	if result.Type().CapturedType() != reflect.TypeFor[getterInterface]() {
		t.Fatalf("result type = %v", result.Type().CapturedType())
	}
}

func TestZeroReceiverAndGetterRejectInvalidRequests(t *testing.T) {
	for _, typ := range []reflect.Type{nil, reflect.TypeFor[getterInterface](), reflect.TypeFor[*getterInterface]()} {
		if _, check := TypeReflection(typ).ZeroReceiver(); check.OK() {
			t.Fatalf("ZeroReceiver(%v) unexpectedly succeeded", typ)
		}
	}
	if _, check := (&TypeInterrogator{}).ZeroReceiver(); check.OK() {
		t.Fatal("ZeroReceiver on nil TypeInspector unexpectedly succeeded")
	}
	if _, check := Type[getterPointerCycle]().ZeroReceiver(); check.OK() {
		t.Fatal("ZeroReceiver on pointer cycle unexpectedly succeeded")
	}
	primitive, check := Type[***getterPrimitive]().ZeroReceiver()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	if got, want := primitive.Type().CapturedType(), reflect.TypeFor[*getterPrimitive](); got != want {
		t.Fatalf("primitive receiver type = %v, want %v", got, want)
	}
	primitiveMethod, check := primitive.Method("Value")
	if !check.OK() {
		t.Fatal(check.Why())
	}
	primitiveGetter, check := primitive.Getter(primitiveMethod)
	if !check.OK() {
		t.Fatal(check.Why())
	}
	primitiveResult, check := primitiveGetter.Read()
	if !check.OK() || primitiveResult.Get() != getterValue("primitive") {
		t.Fatalf("primitive getter result = (%v, %v)", primitiveResult, check)
	}
	receiver, check := Type[getterReceiver]().ZeroReceiver()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	for _, name := range []string{"missing", "NeedsInput", "Many"} {
		method, check := receiver.Method(name)
		if name == "missing" {
			if check.OK() || method != nil {
				t.Fatalf("Method(%q) = (%v, %v)", name, method, check)
			}
			continue
		}
		if !check.OK() {
			t.Fatal(check.Why())
		}
		if getter, check := receiver.Getter(method); check.OK() || getter != nil {
			t.Fatalf("Getter(%q) unexpectedly succeeded", name)
		}
	}
	method, check := receiver.Method("Explode")
	if !check.OK() {
		t.Fatal(check.Why())
	}
	getter, check := receiver.Getter(method)
	if !check.OK() {
		t.Fatal(check.Why())
	}
	if result, check := getter.Read(); check.OK() || result != nil {
		t.Fatalf("panicking getter Read() = (%v, %v)", result, check)
	} else if !strings.Contains(check.Why(), "boom") {
		t.Fatalf("provider panic lost its cause: %s", check.Why())
	}
}

func TestMethodArgumentsAndIndices(t *testing.T) {
	receiver, check := Type[getterReceiver]().ZeroReceiver()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	method, check := receiver.Method("NeedsInput")
	if !check.OK() {
		t.Fatal(check.Why())
	}
	input, check := method.Input(0)
	if !check.OK() || input.CapturedType() != reflect.TypeFor[string]() {
		t.Fatalf("Input(0) = %v, %v", input, check)
	}
	if _, check := method.Input(-1); check.OK() {
		t.Fatal("negative input index succeeded")
	}
	if _, check := method.Result(1); check.OK() {
		t.Fatal("out-of-range result index succeeded")
	}
}
