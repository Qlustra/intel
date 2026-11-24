package define

import (
	"testing"

	"github.com/qlustra/intel/capture"
	"github.com/qlustra/intel/logic"
	"github.com/qlustra/intel/options"
	"github.com/qlustra/intel/rules"
	"github.com/qlustra/intel/rules/extension"
	"github.com/qlustra/intel/rules/property"
)

type sessionDecl struct {
	ready bool
	links []*Carrier[sessionDecl]
}

func (sessionDecl) declaration() {}

type sessionMarker struct{}
type sessionNode struct {
	sessionMarker
	Next  *sessionNode
	Value string
}
type sessionPair struct {
	First  struct{ FirstOnly int }     `label:"first"`
	Second struct{ SecondOnly string } `label:"second"`
}

func sessionDefinition() *StructDefinition[sessionDecl] {
	d := &StructDefinition[sessionDecl]{PropertyComposition: Open, ExtensionComposition: Open, ContractComposition: Open}
	// Compile identity includes the definition instance.
	d.Definition = d
	return d
}

func sessionBody[T any](t *testing.T) *capture.StructInterrogator {
	t.Helper()
	body, check := capture.Type[T]().AsStruct()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	return body
}

func sessionDeclaration(t *testing.T, ctx *rules.EncounterContext) *sessionDecl {
	t.Helper()
	decl, check := ctx.Access[sessionDecl](DeclarationKey)
	if !check.OK() {
		t.Fatal(check.Why())
	}
	return decl
}

type sessionChildDefinition struct {
	StructDefinition[sessionDecl]
	calls int
}

func (child *sessionChildDefinition) Define() {
	child.SetPropertyComposition(Open)
	p := property.Property()
	p.Expect("value", func(e property.Encounter) bool { return e.Property.FieldName() == "Value" })
	child.Property(options.Required, p).Encounter(func(e property.Encounter) logic.Check {
		child.calls++
		decl, check := e.Access[sessionDecl](DeclarationKey)
		if !check.OK() {
			return check
		}
		decl.ready = true
		return logic.Ok()
	})
}

type sessionParentDefinition struct {
	StructDefinition[sessionDecl]
	child *sessionChildDefinition
}

func (parent *sessionParentDefinition) Define() {
	fields := property.Property()
	fields.Expect("fields", func(e property.Encounter) bool {
		return e.Property.FieldName() == "First" || e.Property.FieldName() == "Second"
	})
	// Both fields use the same target, but the scope carriers must retain the
	// original fields rather than metadata from the selected target.
	scope := parent.Scope[sessionDecl, sessionChildDefinition](fields)
	scope.Target(func(property.Encounter) (*capture.StructInterrogator, logic.Check) {
		return capture.Type[sessionNode]().AsStruct()
	})
	scope.Encounter(func(_ property.Encounter, host *sessionDecl, carrier *Carrier[sessionDecl]) logic.Check {
		host.links = append(host.links, carrier)
		return logic.Ok()
	})
	parent.child, _ = scope.Definition.(*sessionChildDefinition)
}

func TestCompilationSessionSharedBodyAndFieldCarriers(t *testing.T) {
	parent, err := Struct[sessionParentDefinition, sessionDecl](nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, validation := parent.Compile(sessionBody[sessionPair](t), nil)
	if !validation.IsOkay() {
		t.Fatal(validation.Reason())
	}
	links := sessionDeclaration(t, ctx).links
	if len(links) != 2 {
		t.Fatalf("got %d carriers", len(links))
	}
	if links[0] == links[1] || links[0].Declaration == nil || links[0].Declaration != links[1].Declaration || !links[0].Declaration.ready {
		t.Fatal("carriers did not share one initialized declaration")
	}
	if parent.child == nil || parent.child.calls != 1 {
		t.Fatal("the shared body did not run its value callback exactly once")
	}
	labels := map[string]string{}
	for _, link := range links {
		value, ok := link.Property.Tags().Lookup("label")
		if !ok {
			t.Fatal("missing field label")
		}
		labels[link.Property.FieldName()] = value
	}
	if labels["First"] != "first" || labels["Second"] != "second" {
		t.Fatalf("field metadata: %v", labels)
	}
}

func TestCompilationSessionRecursiveLinkAndInitialization(t *testing.T) {
	d := sessionDefinition()
	x := extension.Extension()
	x.Expect("marker", func(e extension.Encounter) bool { return e.Extension.FieldName() == "sessionMarker" })
	d.Extension(options.Required, x).Encounter(func(e extension.Encounter) logic.Check {
		sessionDeclaration(t, e.EncounterContext).ready = true
		return logic.Ok()
	})
	p := property.Property()
	p.Expect("next", func(e property.Encounter) bool { return e.Property.FieldName() == "Next" })
	propertyCalls := 0
	d.Property(options.Required, p).Encounter(func(e property.Encounter) logic.Check {
		propertyCalls++
		if propertyCalls != 1 {
			t.Fatal("recursive use re-entered the body instead of linking to it")
		}
		host := sessionDeclaration(t, e.EncounterContext)
		if !host.ready {
			t.Error("extension did not initialize declaration before property")
		}
		ctx, validation := d.Compile(sessionBody[sessionNode](t), e.EncounterContext)
		if !validation.IsOkay() {
			t.Fatal(validation.Reason())
		}
		linked := sessionDeclaration(t, ctx)
		if linked != host || len(linked.links) != 0 {
			t.Fatal("recursive use did not link to the declaration still being built")
		}
		host.links = append(host.links, &Carrier[sessionDecl]{Property: e.Property, Declaration: linked})
		return logic.Ok()
	})
	ctx, validation := d.Compile(sessionBody[sessionNode](t), nil)
	if !validation.IsOkay() {
		t.Fatal(validation.Reason())
	}
	decl := sessionDeclaration(t, ctx)
	if !decl.ready || len(decl.links) != 1 || decl.links[0].Declaration != decl {
		t.Fatal("recursive declaration was not linked and initialized")
	}
}

func TestCompilationSessionFailureIsNotReusable(t *testing.T) {
	parent, child := sessionDefinition(), sessionDefinition()
	fail := true
	calls := 0
	childRule := property.Property()
	childRule.Expect("value", func(e property.Encounter) bool { return e.Property.FieldName() == "Value" })
	child.Property(options.Required, childRule).Encounter(func(e property.Encounter) logic.Check {
		calls++
		if fail {
			return logic.Nope("fixture", "nested failure")
		}
		sessionDeclaration(t, e.EncounterContext).ready = true
		return logic.Ok()
	})
	var failedContext *rules.EncounterContext
	field := property.Property()
	field.Expect("first", func(e property.Encounter) bool { return e.Property.FieldName() == "First" })
	parent.Property(options.Required, field).Encounter(func(e property.Encounter) logic.Check {
		failedContext = e.EncounterContext
		_, validation := child.Compile(sessionBody[sessionNode](t), e.EncounterContext)
		if !validation.IsOkay() {
			return logic.Nope("fixture", "child failed")
		}
		return logic.Ok()
	})
	_, validation := parent.Compile(sessionBody[sessionPair](t), nil)
	if validation.IsOkay() || failedContext == nil || calls != 1 {
		t.Fatalf("nested failure not observed: calls=%d", calls)
	}
	fail = false
	_, retry := child.Compile(sessionBody[sessionNode](t), failedContext)
	if retry.IsOkay() || calls != 1 {
		t.Fatalf("failed session was reused: calls=%d", calls)
	}
	ctx, fresh := child.Compile(sessionBody[sessionNode](t), nil)
	if !fresh.IsOkay() || calls != 2 {
		t.Fatalf("fresh session failed: calls=%d, %s", calls, fresh.Reason())
	}
	if !sessionDeclaration(t, ctx).ready {
		t.Fatal("fresh declaration was not initialized")
	}
}
