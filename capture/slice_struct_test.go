package capture

import (
	"reflect"
	"testing"
)

type sliceMetadataEmbedded struct {
	Promoted int `label:"promoted"`
}
type sliceMetadataElement struct {
	sliceMetadataEmbedded
	Name string `label:"name"`
}
type sliceMetadataList []sliceMetadataElement

func TestStructSliceImmediateElementMetadata(t *testing.T) {
	slice, check := Type[sliceMetadataList]().AsSlice()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	structured, check := slice.AsStructSlice()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	if structured.CapturedType() != reflect.TypeFor[sliceMetadataList]() || structured.Element().CapturedType() != reflect.TypeFor[sliceMetadataElement]() {
		t.Fatal("slice or element type lost")
	}
	for _, tc := range []struct {
		name, tag string
		typ       reflect.Type
		index     []int
	}{
		{"Name", "name", reflect.TypeFor[string](), []int{1}},
		{"Promoted", "promoted", reflect.TypeFor[int](), []int{0, 0}},
	} {
		field, check := structured.Element().Field(tc.name)
		if !check.OK() {
			t.Fatal(check.Why())
		}
		tag, present := field.Tags().Lookup("label")
		if field.ReflectedType() != tc.typ || !reflect.DeepEqual(field.Index(), tc.index) || !present || tag != tc.tag {
			t.Errorf("%s: type=%v index=%v tag=%q present=%v", tc.name, field.ReflectedType(), field.Index(), tag, present)
		}
	}
}

func TestStructSliceRequiresImmediateStruct(t *testing.T) {
	for _, tc := range []struct {
		name string
		typ  reflect.Type
		want bool
	}{
		{"struct", reflect.TypeFor[[]sliceMetadataElement](), true},
		{"named", reflect.TypeFor[sliceMetadataList](), true},
		{"pointer", reflect.TypeFor[[]*sliceMetadataElement](), false},
		{"scalar", reflect.TypeFor[[]int](), false},
		{"nested", reflect.TypeFor[[][]sliceMetadataElement](), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			slice, check := TypeReflection(tc.typ).AsSlice()
			if !check.OK() {
				t.Fatal(check.Why())
			}
			got, check := slice.AsStructSlice()
			if check.OK() != tc.want || (got != nil) != tc.want {
				t.Fatalf("AsStructSlice = %v, %v; want success %v", got, check, tc.want)
			}
		})
	}
}
