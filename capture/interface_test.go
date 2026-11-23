package capture

import (
	"reflect"
	"testing"
)

type methodResult interface{ Result() }

type inspectedInterface interface {
	Run(string, int) (methodResult, error)
	Only() methodResult
}

type inspectedEmptyInterface interface{}

func TestInterfaceMethodInspection(t *testing.T) {
	interrogator, check := Type[inspectedInterface]().AsInterface()
	if !check.OK() {
		t.Fatalf("AsInterface: %v", check)
	}

	method, check := interrogator.Method("Run")
	if !check.OK() {
		t.Fatalf("Method(Run): %v", check)
	}
	if method.Name() != "Run" || method.InputArity() != 2 || method.OutputArity() != 2 {
		t.Fatalf("unexpected method metadata: %s, inputs=%d, outputs=%d", method.Name(), method.InputArity(), method.OutputArity())
	}
	if method.ResultAssignableTo(reflect.TypeFor[methodResult]()) {
		t.Fatal("multiple-result method reported a singular assignable result")
	}

	only, check := interrogator.Method("Only")
	if !check.OK() {
		t.Fatalf("Method(Only): %v", check)
	}
	if only.InputArity() != 0 || only.OutputArity() != 1 || !only.ResultAssignableTo(reflect.TypeFor[methodResult]()) {
		t.Fatal("single result method metadata or assignability was incorrect")
	}
	if only.ResultAssignableTo(reflect.TypeFor[error]()) {
		t.Fatal("unassignable result reported assignable")
	}
}

func TestInterfaceInterrogatorIsEmpty(t *testing.T) {
	for _, tc := range []struct {
		name string
		make func() *InterfaceInterrogator
		want bool
	}{
		{"named empty", func() *InterfaceInterrogator { i, _ := Type[inspectedEmptyInterface]().AsInterface(); return i }, true},
		{"named nonempty", func() *InterfaceInterrogator { i, _ := Type[inspectedInterface]().AsInterface(); return i }, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.make().IsEmpty(); got != tc.want {
				t.Fatalf("IsEmpty() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestInterfaceMethodInspectionRejectsMissingMethod(t *testing.T) {
	interrogator, check := Type[inspectedInterface]().AsInterface()
	if !check.OK() {
		t.Fatalf("AsInterface: %v", check)
	}
	if method, check := interrogator.Method("Missing"); check.OK() || method != nil {
		t.Fatalf("Method(Missing) = (%v, %v), want (nil, failure)", method, check)
	}
}

type inspectedPrivateMethod interface{ hidden() string }

func TestInterfaceMethodInspectionRejectsUnexportedMethod(t *testing.T) {
	interrogator, check := Type[inspectedPrivateMethod]().AsInterface()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	if method, check := interrogator.Method("hidden"); check.OK() || method != nil {
		t.Fatalf("Method(hidden) = (%v, %v), want (nil, failure)", method, check)
	}
}
