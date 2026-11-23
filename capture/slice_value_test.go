package capture

import (
	"reflect"
	"testing"
)

type sliceItem struct{ Number int }
type namedItems []sliceItem
type SliceEmbedded struct{ Promoted namedItems }
type sliceParent struct {
	SliceEmbedded
	Direct  []sliceItem
	private []sliceItem
}

func TestSliceAppendRetainsParentAndElements(t *testing.T) {
	parent := sliceParent{Direct: []sliceItem{{Number: 99}}}
	v, _ := Value(&parent)
	s, _ := v.AsStruct()
	for _, name := range []string{"Direct", "Promoted"} {
		t.Run(name, func(t *testing.T) {
			field, check := s.Field(descriptor(t, reflect.TypeFor[sliceParent](), name))
			if !check.OK() {
				t.Fatal(check.Why())
			}
			slice, check := field.AsSlice()
			if !check.OK() {
				t.Fatal(check.Why())
			}
			for n := 1; n <= 100; n++ {
				element, check := slice.AppendZero()
				if !check.OK() {
					t.Fatal(check.Why())
				}
				address, check := element.Address()
				if !check.OK() || address.(*sliceItem).Number != 0 {
					t.Fatal("element is not fresh zero")
				}
				structure, check := element.AsStruct()
				if !check.OK() {
					t.Fatal(check.Why())
				}
				number, check := structure.Field(descriptor(t, reflect.TypeFor[sliceItem](), "Number"))
				if !check.OK() {
					t.Fatal(check.Why())
				}
				numberAddress, check := number.Address()
				if !check.OK() {
					t.Fatal(check.Why())
				}
				*numberAddress.(*int) = n
			}
			var got []sliceItem
			if name == "Direct" {
				if parent.Direct[0].Number != 99 {
					t.Fatal("existing element changed")
				}
				got = parent.Direct[1:]
			} else {
				got = parent.Promoted
			}
			if len(got) != 100 {
				t.Fatalf("length %d", len(got))
			}
			for n, item := range got {
				if item.Number != n+1 {
					t.Fatalf("element %d = %d", n, item.Number)
				}
			}
		})
	}
}

func TestSliceAppendGenericAndCaptureSemantics(t *testing.T) {
	var numbers []int
	v, _ := Value(&numbers)
	s, check := v.AsSlice()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	e, check := s.AppendZero()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	a, _ := e.Address()
	if *a.(*int) != 0 {
		t.Fatal("integer not zero")
	}
	*a.(*int) = 42
	if len(numbers) != 1 || numbers[0] != 42 {
		t.Fatal("slice pointer not retained")
	}

	pointers := []*int{new(int)}
	p, _ := Value(&pointers)
	ps, _ := p.AsSlice()
	pe, check := ps.AppendZero()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	pa, _ := pe.Address()
	if *pa.(**int) != nil || pointers[1] != nil {
		t.Fatal("pointer element not nil")
	}
	*pa.(**int) = &numbers[0]
	if pointers[1] != &numbers[0] {
		t.Fatal("pointer element detached")
	}

	original := []int{1}
	copyValue, _ := Value(original)
	copySlice, _ := copyValue.AsSlice()
	if _, check := copySlice.AppendZero(); !check.OK() {
		t.Fatal(check.Why())
	}
	if len(original) != 1 {
		t.Fatal("value capture changed original slice header")
	}
}

func TestSliceValueCaptureCopiesHeaderAndSharesBackingArray(t *testing.T) {
	original := make([]int, 1, 3)
	original[0] = 17
	original[:3][1] = 29

	value, _ := Value(original)
	slice, check := value.AsSlice()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	element, check := slice.AppendZero()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	address, check := element.Address()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	if *address.(*int) != 0 {
		t.Fatal("appended element was not reset to zero")
	}
	*address.(*int) = 41

	if len(original) != 1 {
		t.Fatalf("original length changed to %d", len(original))
	}
	if original[0] != 17 {
		t.Fatalf("original first element changed to %d", original[0])
	}
	if original[:2][1] != 41 {
		t.Fatalf("backing-array mutation not visible in original: %d", original[:2][1])
	}
}

func TestSliceElementsCaptureAddressesWithoutModification(t *testing.T) {
	items := namedItems{{Number: 11}, {Number: 29}}
	value, _ := Value(&items)
	slice, check := value.AsSlice()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	elements, check := slice.Elements()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	if len(elements) != 2 || items[0].Number != 11 || items[1].Number != 29 {
		t.Fatalf("elements or slice changed: %#v", items)
	}
	for index, want := range []int{11, 29} {
		address, check := elements[index].Address()
		if !check.OK() {
			t.Fatal(check.Why())
		}
		if address.(*sliceItem) != &items[index] || address.(*sliceItem).Number != want {
			t.Fatalf("element %d address/value mismatch", index)
		}
	}
	empty, _ := Value(&[]sliceItem{})
	emptySlice, _ := empty.AsSlice()
	got, check := emptySlice.Elements()
	if !check.OK() || len(got) != 0 {
		t.Fatalf("empty elements = %v, %v", got, check)
	}
}

func TestSliceAppendElementAddressDoesNotSurviveGrowth(t *testing.T) {
	retained := make([]int, 0, 1)
	value, _ := Value(&retained)
	slice, check := value.AsSlice()
	if !check.OK() {
		t.Fatal(check.Why())
	}

	first, check := slice.AppendZero()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	firstAddress, check := first.Address()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	*firstAddress.(*int) = 23

	second, check := slice.AppendZero()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	secondAddress, check := second.Address()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	*firstAddress.(*int) = 37

	if len(retained) != 2 {
		t.Fatalf("retained length = %d, want 2", len(retained))
	}
	if retained[0] != 23 {
		t.Fatalf("current first element = %d, want copied pre-growth value 23", retained[0])
	}
	if retained[1] != 0 || *secondAddress.(*int) != 0 {
		t.Fatal("second appended element was not zero")
	}
}

func TestSliceRejectsInvalidInputs(t *testing.T) {
	var nilValue *ValueInterrogator
	var nilSlice *SliceValueInterrogator
	number := 7
	nonslice, _ := Value(&number)
	private := sliceParent{private: []sliceItem{{Number: 5}}}
	privateValue := reflect.ValueOf(&private).Elem().FieldByName("private")
	bad := []*ValueInterrogator{nilValue, {}, nonslice,
		{val: reflect.ValueOf((*[]int)(nil)), typ: Type[[]int]()},
		{val: privateValue.Addr(), typ: Type[[]sliceItem]()},
		{val: reflect.ValueOf(&number), typ: Type[[]int]()},
	}
	for n, value := range bad {
		if _, check := value.AsSlice(); check.OK() {
			t.Fatalf("invalid value %d accepted", n)
		}
	}
	badSlices := []*SliceValueInterrogator{nilSlice, {},
		{val: privateValue.Addr(), typ: &SliceInterrogator{}},
		{val: reflect.ValueOf(&number), typ: &SliceInterrogator{}},
	}
	for n, slice := range badSlices {
		if _, check := slice.AppendZero(); check.OK() {
			t.Fatalf("invalid slice %d accepted", n)
		}
	}
	v, _ := Value(&private)
	st, _ := v.AsStruct()
	if _, check := st.Field(descriptor(t, reflect.TypeFor[sliceParent](), "private")); check.OK() {
		t.Fatal("private field accepted")
	}
	if number != 7 || len(private.private) != 1 || private.private[0].Number != 5 {
		t.Fatal("invalid operation mutated input")
	}
}
