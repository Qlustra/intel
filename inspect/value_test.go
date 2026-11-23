package inspect

import (
	"testing"
	"unsafe"
)

type nilTestFunc func()
type nilTestMap map[string]int
type nilTestSlice []int
type nilTestChan chan int

func TestIsNil(t *testing.T) {
	var pointer *int
	var function nilTestFunc
	var mapping nilTestMap
	var slice nilTestSlice
	var channel nilTestChan
	for _, tc := range []struct {
		name  string
		value any
		want  bool
	}{
		{"nil interface", nil, true},
		{"nil pointer", pointer, true},
		{"nil unsafe pointer", unsafe.Pointer(nil), true},
		{"unsafe pointer", unsafe.Pointer(new(int)), false},
		{"pointer", new(int), false},
		{"nil func", function, true},
		{"func", nilTestFunc(func() {}), false},
		{"nil map", mapping, true},
		{"map", make(nilTestMap), false},
		{"nil slice", slice, true},
		{"slice", make(nilTestSlice, 0), false},
		{"nil channel", channel, true},
		{"channel", make(nilTestChan), false},
		{"scalar", 0, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsNil(tc.value); got != tc.want {
				t.Fatalf("IsNil(%s) = %v, want %v", tc.name, got, tc.want)
			}
		})
	}
}
