package capture

import (
	"testing"
	"unsafe"
)

type assignSlice []int
type assignMap map[string]int
type assignFunc func() int
type assignStruct struct{ Value int }
type assignString string

func TestAssignGoAssignable(t *testing.T) {
	s, c := Assign[assignSlice]([]int{3})
	if !c.OK() || s[0] != 3 {
		t.Fatal("slice", c)
	}
	m, c := Assign[assignMap](map[string]int{"a": 4})
	if !c.OK() || m["a"] != 4 {
		t.Fatal("map", c)
	}
	f, c := Assign[assignFunc](func() int { return 5 })
	if !c.OK() || f() != 5 {
		t.Fatal("func", c)
	}
	v, c := Assign[assignStruct](struct{ Value int }{6})
	if !c.OK() || v.Value != 6 {
		t.Fatal("struct", c)
	}
	ch := make(chan int)
	r, c := Assign[<-chan int](ch)
	if !c.OK() || r != ch {
		t.Fatal("channel", c)
	}
	if _, c := Assign[assignString]("x"); c.OK() {
		t.Fatal("converted string")
	}
	if _, c := Assign[int](int8(1)); c.OK() {
		t.Fatal("converted number")
	}
	if _, c := Assign[chan int]((<-chan int)(ch)); c.OK() {
		t.Fatal("reversed channel direction")
	}
}

func TestAssignNil(t *testing.T) {
	if v, c := Assign[any](nil); !c.OK() || v != nil {
		t.Fatal("nil interface", c)
	}
	if v, c := Assign[unsafe.Pointer](nil); !c.OK() || v != nil {
		t.Fatal("nil unsafe pointer", c)
	}
	if _, c := Assign[int](nil); c.OK() {
		t.Fatal("nil int")
	}
	var p *int
	if v, c := Assign[any](p); !c.OK() || v != p {
		t.Fatal("typed nil interface", c)
	}
	if v, c := Assign[*int](p); !c.OK() || v != nil {
		t.Fatal("typed nil pointer", c)
	}
	if _, c := Assign[*string](p); c.OK() {
		t.Fatal("incompatible typed nil")
	}
	var s []int
	if v, c := Assign[assignSlice](s); !c.OK() || v != nil {
		t.Fatal("assignable typed nil slice", c)
	}
}
