package inspect

import (
	"reflect"
	"testing"
)

type emptyInterfaceTest interface{}
type nonemptyInterfaceTest interface{ Method() }

func TestInterfaceInspectorIsEmpty(t *testing.T) {
	for _, tc := range []struct {
		name string
		typ  reflect.Type
		want bool
	}{
		{"named empty", reflect.TypeFor[emptyInterfaceTest](), true},
		{"named nonempty", reflect.TypeFor[nonemptyInterfaceTest](), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := InspectInterface(tc.typ).IsEmpty(); got != tc.want {
				t.Fatalf("IsEmpty() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestIsExportedIdentifier(t *testing.T) {
	for _, tc := range []struct {
		name string
		want bool
	}{
		{"Name", true},
		{"Δelta", true},
		{"_name", false},
		{"name", false},
		{"", false},
		{"9Lives", false},
		{"Bad-Name", false},
		{"has space", false},
		{"🙂", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsExportedIdentifier(tc.name); got != tc.want {
				t.Fatalf("IsExportedIdentifier(%q) = %v, want %v", tc.name, got, tc.want)
			}
		})
	}
}
