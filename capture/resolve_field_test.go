package capture

import (
	"reflect"
	"strings"
	"testing"

	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/logic"
)

type resolveEmbedded struct{ Promoted int }
type resolveChild struct{ Value int }

func (c *resolveChild) Increase() int { c.Value++; return c.Value }

type resolveRoot struct {
	Ignored string
	resolveEmbedded
	Child   resolveChild
	Ptr     *resolveChild
	private int
}
type resolveOther resolveRoot
type resolveAlias = resolveRoot
type resolveRootPointer *resolveRoot

func TestResolveFieldLiveAndPaths(t *testing.T) {
	a, b := &resolveRoot{}, &resolveRoot{}
	root := Type[resolveRoot]()
	for _, obj := range []*resolveRoot{a, b} {
		field, check := ResolveField(obj, root, []int{1, 0})
		if !check.OK() {
			t.Fatal(check)
		}
		address, check := field.Address()
		if !check.OK() {
			t.Fatal(check)
		}
		*address.(*int) = 7
	}
	a.Promoted = 8
	if b.Promoted != 7 {
		t.Fatal("instances not isolated")
	}
	field, check := ResolveField(a, root, []int{2})
	if !check.OK() {
		t.Fatal(check)
	}
	child, check := field.AsStruct()
	if !check.OK() {
		t.Fatal(check)
	}
	result, check := child.Invoke[interface{ Increase() int }, int]("Increase", func(v interface{ Increase() int }) (int, logic.Check) { return v.Increase(), logic.Ok() })
	if !check.OK() || result != 1 || a.Child.Value != 1 || b.Child.Value != 0 {
		t.Fatalf("invoke: %v %v", result, check)
	}
	p := &resolveChild{Value: 9}
	a.Ptr = p
	field, check = ResolveField(a, root, []int{3})
	if !check.OK() {
		t.Fatal(check)
	}
	address, check := field.Address()
	if !check.OK() || address.(**resolveChild) != &a.Ptr || *address.(**resolveChild) != p {
		t.Fatal("terminal pointer identity")
	}
	field, check = ResolveField(a, root, []int{3, 0})
	if !check.OK() {
		t.Fatal(check)
	}
	address, check = field.Address()
	if !check.OK() || address.(*int) != &p.Value {
		t.Fatal("pointer intermediate identity")
	}
	a.Ptr = nil
	field, check = ResolveField(a, root, []int{3})
	if !check.OK() {
		t.Fatal(check)
	}
	address, check = field.Address()
	if !check.OK() || address.(**resolveChild) != &a.Ptr || *address.(**resolveChild) != nil {
		t.Fatal("nil terminal pointer slot was not preserved")
	}
	if _, check = ResolveField(a, root, []int{3, 0}); check.OK() || a.Ptr != nil {
		t.Fatal("nil pointer followed or allocated")
	}
	var alias resolveAlias
	if _, check = ResolveField(&alias, root, []int{2, 0}); !check.OK() {
		t.Fatal(check)
	}
	var namedPointer resolveRootPointer = a
	if _, check = ResolveField(namedPointer, root, []int{2, 0}); !check.OK() {
		t.Fatal("pointer carrier name must not change root element identity", check)
	}
}

func TestResolveFieldFailures(t *testing.T) {
	root := Type[resolveRoot]()
	var nilRoot *resolveRoot
	obj := &resolveRoot{}
	for _, instance := range []any{nil, nilRoot, resolveRoot{}, &obj, &resolveOther{}, 3, new(int)} {
		if field, check := ResolveField(instance, root, []int{-1}); check.OK() || field != nil || !strings.Contains(check.Why(), "root type") {
			t.Fatalf("root %T: %v", instance, check)
		}
	}
	for _, expected := range []*TypeInterrogator{nil, {}, TypeReflection(nil), Type[*resolveRoot](), Type[int](), {TypeInspector: inspect.CaptureType(reflect.TypeFor[resolveRoot]())}, {TypeInspector: inspect.CaptureType(reflect.TypeFor[int]()), typ: reflect.TypeFor[resolveRoot]()}} {
		if field, check := ResolveField(obj, expected, []int{0}); check.OK() || field != nil {
			t.Fatalf("bad metadata: %v", check)
		}
	}
	for _, path := range [][]int{nil, {}, {-1}, {99}, {0, 0}, {2, -1}, {2, 99}, {4}} {
		before := append([]int(nil), path...)
		if field, check := ResolveField(obj, root, path); check.OK() || field != nil {
			t.Fatalf("path %v: %v", path, check)
		}
		if !reflect.DeepEqual(before, append([]int(nil), path...)) {
			t.Fatal("path modified")
		}
	}
}
