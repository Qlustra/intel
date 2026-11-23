package capture

import (
	"reflect"
	"testing"

	"github.com/qlustra/intel/options"
)

type filterBoundary struct{ BoundaryLeaf string }
type filterOther struct{ OtherLeaf string }
type filterPointerBoundary struct{ PointerLeaf string }

type filterFixture struct {
	filterBoundary
	*filterPointerBoundary
	Sibling filterOther
	Own     string
}

type filterIdentityMarker struct{ MarkerValue string }

func (filterIdentityMarker) Behavior() {}

type filterIdentityAlias = filterIdentityMarker
type filterIdentityDefined filterIdentityMarker

func (filterIdentityDefined) Behavior() {}

type filterBehaviorOther struct{ OtherValue string }

func (filterBehaviorOther) Behavior() {}

type filterIdentityHost struct {
	filterIdentityMarker `marker:"kept"`
	filterIdentityDefined
	filterBehaviorOther
}

func TestStructFilterExactIdentityAndRoot(t *testing.T) {
	type markerHost struct{ filterIdentityMarker }
	exact, check := Type[markerHost]().AsStruct()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	for _, boundary := range []reflect.Type{reflect.TypeFor[filterIdentityMarker](), reflect.TypeFor[filterIdentityAlias]()} {
		if _, c := exact.WithoutDescendantsOf(boundary).Field("MarkerValue"); c.OK() {
			t.Errorf("%v did not filter exact marker descendant", boundary)
		}
	}
	base, check := Type[filterIdentityHost]().AsStruct()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	for _, boundary := range []reflect.Type{reflect.TypeFor[filterIdentityMarker](), reflect.TypeFor[filterIdentityAlias]()} {
		view := base.WithoutDescendantsOf(boundary)
		for _, name := range []string{"filterIdentityMarker", "filterIdentityDefined", "filterBehaviorOther", "OtherValue"} {
			if _, c := view.Field(name); !c.OK() {
				t.Errorf("%v: missing %s", boundary, name)
			}
		}
		// MarkerValue is ambiguous between the marker and the distinct defined type,
		// so test the defined type on its own below instead of resurrecting it.
	}
	type definedHost struct{ filterIdentityDefined }
	d, check := Type[definedHost]().AsStruct()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	if _, c := d.WithoutDescendantsOf(reflect.TypeFor[filterIdentityMarker]()).Field("MarkerValue"); !c.OK() {
		t.Fatal("distinct defined type with identical layout was filtered")
	}
	root, check := Type[filterIdentityMarker]().AsStruct()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	if _, c := root.WithoutDescendantsOf(reflect.TypeFor[filterIdentityMarker]()).Field("MarkerValue"); !c.OK() {
		t.Fatal("filter treated root as a boundary")
	}
}

type filterDeepMarker struct {
	HiddenTagged  string `role:"valid"`
	hiddenPrivate string
}
type filterDeepInner struct {
	filterDeepMarker `marker:"yes"`
	Authored         string `role:"valid"`
}
type filterDeepMiddle struct{ filterDeepInner }
type filterDeepHost struct {
	filterDeepMiddle
	Direct string
	filterOther
	*filterPointerBoundary
}

func TestStructFilterDeepCompositionAndOrdering(t *testing.T) {
	base, check := Type[filterDeepHost]().AsStruct()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	first := base.WithoutDescendantsOf(reflect.TypeFor[filterDeepMarker]())
	view := first.WithoutDescendantsOf(reflect.TypeFor[filterOther](), reflect.TypeFor[filterPointerBoundary]())
	for _, name := range []string{"HiddenTagged", "hiddenPrivate", "OtherLeaf", "PointerLeaf"} {
		if _, c := view.Field(name); c.OK() {
			t.Errorf("unexpected descendant %s", name)
		}
		if _, c := base.Field(name); !c.OK() {
			t.Errorf("raw source lost %s", name)
		}
	}
	if _, c := first.Field("OtherLeaf"); !c.OK() {
		t.Error("composing view changed earlier view")
	}
	marker, check := view.Field("filterDeepMarker")
	if !check.OK() {
		t.Fatal(check.Why())
	}
	if tag, ok := marker.Tags().Lookup("marker"); !ok || tag != "yes" {
		t.Errorf("marker tag = %q, %v", tag, ok)
	}
	sibling, check := view.Field("Authored")
	if !check.OK() || !reflect.DeepEqual(sibling.Index(), []int{0, 0, 1}) {
		t.Fatalf("authored sibling path = %v (%v)", sibling, check)
	}
	pointer, check := view.Field("filterPointerBoundary")
	if !check.OK() || pointer.ReflectedType() != reflect.TypeFor[*filterPointerBoundary]() {
		t.Fatal("pointer boundary not retained")
	}
	for _, order := range []options.DepthOrdering{options.DepthFirst, options.BreadthFirst} {
		var expected, actual []string
		for name := range base.Fields.AllAddressable(order) {
			if name != "HiddenTagged" && name != "hiddenPrivate" && name != "OtherLeaf" && name != "PointerLeaf" {
				expected = append(expected, name)
			}
		}
		for name := range view.Fields.AllAddressable(order) {
			actual = append(actual, name)
		}
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("order %v: got %v want %v", order, actual, expected)
		}
	}
}

type filterShadowChild struct{ Shadow string }
type filterShadowParent struct {
	filterShadowChild
	Shadow int
}

type filterAmbiguousLeft struct{ Clash string }
type filterAmbiguousRight struct{ Clash int }
type filterAmbiguous struct {
	filterAmbiguousLeft
	filterAmbiguousRight
}

func TestStructWithoutDescendantsOfFiltersExactStructsAsView(t *testing.T) {
	base, check := Type[filterFixture]().AsStruct()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	filtered := base.WithoutDescendantsOf(reflect.TypeFor[filterBoundary](), reflect.TypeFor[filterPointerBoundary]())
	if filtered == base {
		t.Fatal("filter returned the original interrogator")
	}

	got := make(map[string][]int)
	for name, field := range filtered.Fields.AllAddressable() {
		got[name] = field.Index()
	}
	for _, name := range []string{"BoundaryLeaf", "PointerLeaf"} {
		if _, ok := got[name]; ok {
			t.Errorf("filtered fields include descendant %q", name)
		}
	}
	for _, name := range []string{"filterBoundary", "filterPointerBoundary", "Sibling", "Own"} {
		if _, ok := got[name]; !ok {
			t.Errorf("filtered fields omit boundary or sibling %q", name)
		}
	}

	boundary, check := filtered.Field("filterBoundary")
	if !check.OK() || boundary.ReflectedType() != reflect.TypeFor[filterBoundary]() {
		t.Fatalf("boundary lookup = (%v, %v)", boundary, check)
	}
	if _, check := filtered.Field("BoundaryLeaf"); check.OK() {
		t.Fatal("lookup crossed excluded struct boundary")
	}

	// Constructing a filtered view must not alter the originating interrogator.
	if _, check := base.Field("BoundaryLeaf"); !check.OK() {
		t.Fatalf("filter mutated source interrogator: %s", check.Why())
	}
}

func TestStructWithoutDescendantsOfPreservesShadowingAndAmbiguity(t *testing.T) {
	shadow, check := Type[filterShadowParent]().AsStruct()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	shadow = shadow.WithoutDescendantsOf(reflect.TypeFor[filterShadowChild]())
	field, check := shadow.Field("Shadow")
	if !check.OK() || field.ReflectedType() != reflect.TypeFor[int]() {
		t.Fatalf("shadowed lookup = (%v, %v), want direct int field", field, check)
	}
	if _, check := shadow.Field("filterShadowChild"); !check.OK() {
		t.Fatalf("boundary itself was not retained: %s", check.Why())
	}

	ambiguous, check := Type[filterAmbiguous]().AsStruct()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	ambiguous = ambiguous.WithoutDescendantsOf(reflect.TypeFor[filterAmbiguousLeft]())
	if field, check := ambiguous.Field("Clash"); check.OK() || field != nil {
		t.Fatalf("ambiguous lookup = (%v, %v), want checked failure", field, check)
	}
}

func TestStructWithoutDescendantsOfLeavesFieldOrderAndPathsUnchanged(t *testing.T) {
	base, check := Type[filterFixture]().AsStruct()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	filtered := base.WithoutDescendantsOf(reflect.TypeFor[filterBoundary](), reflect.TypeFor[filterPointerBoundary]())
	var baseNames, filteredNames []string
	for name := range base.Fields.AllAddressable() {
		if name != "BoundaryLeaf" && name != "PointerLeaf" {
			baseNames = append(baseNames, name)
		}
	}
	for name, field := range filtered.Fields.AllAddressable() {
		filteredNames = append(filteredNames, name)
		original, check := base.Field(name)
		if !check.OK() || !reflect.DeepEqual(field.Index(), original.Index()) {
			t.Fatalf("field %q path = %v, original %v (%v)", name, field.Index(), original, check)
		}
	}
	if !reflect.DeepEqual(baseNames, filteredNames) {
		t.Fatalf("filter changed surviving field order: base %v filtered %v", baseNames, filteredNames)
	}
}
