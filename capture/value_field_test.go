package capture

import (
	"reflect"
	"testing"

	"github.com/qlustra/intel/inspect"
)

type fieldInner struct{ Number int }
type fieldOuter struct {
	fieldInner
	Direct  string
	Pointer *int
	private int
}
type FieldEmbedded struct{ Number int }
type fieldPointerOuter struct{ *FieldEmbedded }
type taggedFieldSource struct {
	Value int `json:"source"`
}
type taggedFieldTarget struct {
	Value int `json:"target"`
}

func descriptor(t *testing.T, typ reflect.Type, name string) *inspect.FieldDescriptor {
	t.Helper()
	f, ok := typ.FieldByName(name)
	if !ok {
		t.Fatalf("missing field %s", name)
	}
	return inspect.CaptureFieldDescriptor(f)
}

func TestCapturedFieldAddresses(t *testing.T) {
	parent := fieldOuter{}
	v, check := Value(&parent)
	if !check.OK() {
		t.Fatal(check.Why())
	}
	s, check := v.AsStruct()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	typ := reflect.TypeFor[fieldOuter]()
	for _, name := range []string{"Direct", "Number", "Pointer"} {
		f, check := s.Field(descriptor(t, typ, name))
		if !check.OK() {
			t.Fatalf("%s: %s", name, check.Why())
		}
		addr, check := f.Address()
		if !check.OK() {
			t.Fatal(check.Why())
		}
		switch name {
		case "Direct":
			*addr.(*string) = "changed"
		case "Number":
			*addr.(*int) = 42
		case "Pointer":
			n := 7
			*addr.(**int) = &n
		}
	}
	if parent.Direct != "changed" || parent.Number != 42 || parent.Pointer == nil || *parent.Pointer != 7 {
		t.Fatalf("parent not updated: %+v", parent)
	}
}

func TestFieldPathChecks(t *testing.T) {
	typ := reflect.TypeFor[fieldPointerOuter]()
	d := descriptor(t, typ, "Number")
	typeOnly, check := Type[fieldPointerOuter]().AsStruct()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	if typeOnly.FieldPathIsValue(d).OK() {
		t.Fatal("pointer path accepted as value-only")
	}
	outer := fieldPointerOuter{}
	v, _ := Value(&outer)
	s, _ := v.AsStruct()
	if _, check := s.Field(d); check.OK() {
		t.Fatal("nil embedded pointer accepted")
	}
	if outer.FieldEmbedded != nil {
		t.Fatal("allocated embedded pointer")
	}
	outer.FieldEmbedded = &FieldEmbedded{}
	f, check := s.Field(d)
	if !check.OK() {
		t.Fatal(check.Why())
	}
	address, check := f.Address()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	*address.(*int) = 19
	if outer.Number != 19 {
		t.Fatal("promoted pointer field detached")
	}
	valueType, _ := Type[fieldOuter]().AsStruct()
	for _, name := range []string{"Number", "Pointer"} {
		if check := valueType.FieldPathIsValue(descriptor(t, reflect.TypeFor[fieldOuter](), name)); !check.OK() {
			t.Fatalf("%s: %s", name, check.Why())
		}
	}
}

func TestFieldDescriptorUsesStructuralCompatibility(t *testing.T) {
	sourceType := reflect.TypeFor[taggedFieldSource]()
	d := descriptor(t, sourceType, "Value")

	// Different tags are metadata, not evidence of ownership or common origin:
	// the descriptor contains no owner identity to validate against the target.
	typeOnly, check := Type[taggedFieldTarget]().AsStruct()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	if check := typeOnly.FieldPathIsValue(d); !check.OK() {
		t.Fatalf("structurally compatible descriptor rejected: %s", check.Why())
	}

	source := taggedFieldSource{Value: 11}
	target := taggedFieldTarget{Value: 22}
	value, check := Value(&target)
	if !check.OK() {
		t.Fatal(check.Why())
	}
	structValue, check := value.AsStruct()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	field, check := structValue.Field(d)
	if !check.OK() {
		t.Fatalf("structurally compatible field rejected: %s", check.Why())
	}
	address, check := field.Address()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	*address.(*int) = 33
	if target.Value != 33 || source.Value != 11 {
		t.Fatalf("mutation target=%d, source=%d; want target=33, source=11", target.Value, source.Value)
	}
}

func TestFieldRejectsInvalidInputs(t *testing.T) {
	parent := fieldOuter{}
	v, _ := Value(&parent)
	s, _ := v.AsStruct()
	typ := reflect.TypeFor[fieldOuter]()
	bad := []struct {
		name  string
		field *inspect.FieldDescriptor
	}{
		{"nil", nil},
		{"private", descriptor(t, typ, "private")},
		{"empty", inspect.CaptureFieldDescriptor(reflect.StructField{Name: "Direct", Type: reflect.TypeFor[string]()})},
		{"negative", inspect.CaptureFieldDescriptor(reflect.StructField{Name: "Direct", Type: reflect.TypeFor[string](), Index: []int{-1}})},
		{"range", inspect.CaptureFieldDescriptor(reflect.StructField{Name: "Direct", Type: reflect.TypeFor[string](), Index: []int{99}})},
		{"nonstruct", inspect.CaptureFieldDescriptor(reflect.StructField{Name: "Direct", Type: reflect.TypeFor[string](), Index: []int{1, 0}})},
		{"wrong name", inspect.CaptureFieldDescriptor(reflect.StructField{Name: "Other", Type: reflect.TypeFor[string](), Index: []int{1}})},
		{"wrong type", inspect.CaptureFieldDescriptor(reflect.StructField{Name: "Direct", Type: reflect.TypeFor[int](), Index: []int{1}})},
	}
	for _, tc := range bad {
		t.Run(tc.name, func(t *testing.T) {
			if _, check := s.Field(tc.field); check.OK() {
				t.Fatal("accepted invalid field")
			}
			if tc.name != "private" && s.typ.FieldPathIsValue(tc.field).OK() {
				t.Fatal("accepted invalid path")
			}
		})
	}
	var sv *StructValueInterrogator
	if _, check := sv.Field(nil); check.OK() {
		t.Fatal("nil struct receiver accepted")
	}
	var vv *ValueInterrogator
	if _, check := vv.Address(); check.OK() {
		t.Fatal("nil value receiver accepted")
	}
	if _, check := (&ValueInterrogator{}).Address(); check.OK() {
		t.Fatal("zero value accepted")
	}
	var st *StructInterrogator
	if st.FieldPathIsValue(nil).OK() {
		t.Fatal("nil type receiver accepted")
	}
	if _, check := Value(nil); check.OK() {
		t.Fatal("nil value accepted")
	}
	var nilParent *fieldOuter
	if _, check := Value(nilParent); check.OK() {
		t.Fatal("nil pointer accepted")
	}
}
