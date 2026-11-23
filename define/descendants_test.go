package define

import (
	"testing"

	"github.com/qlustra/intel/capture"
	"github.com/qlustra/intel/logic"
	"github.com/qlustra/intel/options"
	"github.com/qlustra/intel/rules/contract"
	"github.com/qlustra/intel/rules/extension"
	"github.com/qlustra/intel/rules/property"
)

type boundaryDecl struct{}

func (boundaryDecl) declaration() {}

type hiddenContract interface{ HiddenMethod() }
type hiddenExtension struct{}
type ownedBoundary struct {
	hiddenContract
	hiddenExtension
	Hidden      string `role:"authored"`
	hiddenOwned string `role:"authored"`
}

type privateBoundary struct{ private string }
type promotedBoundary struct{ Promoted string }
type privateHelper struct{ privateBoundary }
type promotedHelper struct{ promotedBoundary }
type siblingHost struct {
	ownedBoundary `visible:"no"`
	Sibling       string
	privateHelper
	promotedHelper
}

func TestBoundaryMarkerStillProcessed(t *testing.T) {
	raw, check := capture.Type[boundaryHost]().AsStruct()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	d := &StructDefinition[boundaryDecl]{PropertyComposition: Open, ExtensionComposition: Open, ContractComposition: Open}
	d.IgnoreDescendantsOf[ownedBoundary]()
	called := false
	marker := extension.Extension()
	marker.Expect("marker", func(e extension.Encounter) bool { return e.Extension.FieldName() == "ownedBoundary" })
	d.Extension(options.Required, marker).Encounter(func(e extension.Encounter) logic.Check {
		called = true
		if value, ok := e.Extension.Tags().Lookup("marker"); !ok || value != "yes" {
			t.Errorf("marker tag: %q present=%v", value, ok)
		}
		return logic.Ok()
	})
	if result := d.Declare(raw, nil); !result.Validation.IsOkay() || !called {
		t.Fatalf("marker called=%v: %s", called, result.Validation.Reason())
	}
}

type authoredHelper struct {
	ownedBoundary `marker:"yes"`
	Sibling       string
}
type boundaryHost struct{ authoredHelper }

func TestDefinitionDescendantCandidatesAndView(t *testing.T) {
	raw, check := capture.Type[boundaryHost]().AsStruct()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	d := &StructDefinition[boundaryDecl]{PropertyComposition: Open, ExtensionComposition: Open, ContractComposition: Open}
	d.IgnoreDescendantsOf[ownedBoundary]().IgnoreDescendantsOf[ownedBoundary]()
	properties, extensions, contracts := map[string]bool{}, map[string]bool{}, map[string]bool{}
	p := property.Property()
	p.Expect("observe", func(e property.Encounter) bool {
		properties[e.Property.FieldName()] = true
		if _, c := e.Struct.Field("Hidden"); c.OK() {
			t.Error("callback view exposes hidden property")
		}
		return false
	})
	x := extension.Extension()
	x.Expect("observe", func(e extension.Encounter) bool { extensions[e.Extension.FieldName()] = true; return false })
	c := contract.Contract()
	c.Expect("observe", func(e contract.Encounter) bool { contracts[e.Contract.FieldName()] = true; return false })
	d.Property(options.Optional, p)
	d.Extension(options.Optional, x)
	d.Contract(options.Optional, c)
	if result := d.Declare(raw, nil); !result.Validation.IsOkay() {
		t.Fatal(result.Validation.Reason())
	}
	if !properties["Sibling"] || properties["Hidden"] || !extensions["ownedBoundary"] || extensions["hiddenExtension"] || contracts["hiddenContract"] {
		t.Fatalf("candidates property=%v extension=%v contract=%v", properties, extensions, contracts)
	}
	if _, check := raw.Field("Hidden"); !check.OK() {
		t.Fatal("source interrogator was changed")
	}
	plain := &StructDefinition[boundaryDecl]{PropertyComposition: Open, ExtensionComposition: Open, ContractComposition: Open}
	seen := false
	probe := property.Property()
	probe.Expect("observe", func(e property.Encounter) bool {
		if e.Property.FieldName() == "Hidden" {
			seen = true
		}
		return false
	})
	plain.Property(options.Optional, probe)
	if result := plain.Declare(raw, nil); !result.Validation.IsOkay() || !seen {
		t.Fatalf("unfiltered definition: seen=%v validation=%v", seen, result.Validation.Reason())
	}
}

func TestIgnoreDescendantsOfTaggedAndAuthoredSiblings(t *testing.T) {
	raw, check := capture.Type[siblingHost]().AsStruct()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	d := &StructDefinition[boundaryDecl]{PropertyComposition: Open, ExtensionComposition: Open, ContractComposition: Open}
	d.IgnoreDescendantsOf[ownedBoundary]()
	seen := map[string]bool{}
	p := property.Property()
	p.Expect("observe", func(e property.Encounter) bool { seen[e.Property.FieldName()] = true; return false })
	d.Property(options.Optional, p)
	if result := d.Declare(raw, nil); !result.Validation.IsOkay() {
		t.Fatal(result.Validation.Reason())
	}
	if seen["Hidden"] || seen["hiddenOwned"] || !seen["private"] || !seen["Sibling"] || !seen["Promoted"] {
		t.Fatalf("visible authored candidates: %v", seen)
	}
	strict := &StructDefinition[boundaryDecl]{PropertyComposition: Closed, ExtensionComposition: Open, ContractComposition: Open}
	strict.IgnoreDescendantsOf[ownedBoundary]()
	if result := strict.Declare(raw, nil); result.Validation.IsOkay() {
		t.Fatal("strict definition accepted unmatched authored siblings")
	}
}

func TestIgnoreDescendantsOfIndependentDefinitions(t *testing.T) {
	raw, check := capture.Type[boundaryHost]().AsStruct()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	filtered := &StructDefinition[boundaryDecl]{PropertyComposition: Open, ExtensionComposition: Open, ContractComposition: Open}
	filtered.IgnoreDescendantsOf[ownedBoundary]()
	plain := &StructDefinition[boundaryDecl]{PropertyComposition: Open, ExtensionComposition: Open, ContractComposition: Open}
	seen := false
	p := property.Property()
	p.Expect("observe", func(e property.Encounter) bool { seen = seen || e.Property.FieldName() == "Hidden"; return false })
	plain.Property(options.Optional, p)
	if result := plain.Declare(raw, nil); !result.Validation.IsOkay() || !seen {
		t.Fatalf("unfiltered sibling definition: seen=%v validation=%v", seen, result.Validation.Reason())
	}
	probe := property.Property()
	probe.Expect("observe", func(e property.Encounter) bool {
		if e.Property.FieldName() == "Hidden" {
			t.Error("filtered definition exposed hidden property")
		}
		return false
	})
	filtered.Property(options.Optional, probe)
	if result := filtered.Declare(raw, nil); !result.Validation.IsOkay() {
		t.Fatal(result.Validation.Reason())
	}
}

func TestIgnoredRequiredExtensionAndStrictProperty(t *testing.T) {
	raw, check := capture.Type[boundaryHost]().AsStruct()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	d := &StructDefinition[boundaryDecl]{PropertyComposition: Open, ExtensionComposition: Open, ContractComposition: Open}
	d.IgnoreDescendantsOf[ownedBoundary]()
	hidden := extension.Extension()
	hidden.Expect("hidden", func(e extension.Encounter) bool { return e.Extension.FieldName() == "hiddenExtension" })
	d.Extension(options.Required, hidden)
	if result := d.Declare(raw, nil); result.Validation.IsOkay() {
		t.Fatal("ignored descendant satisfied required extension")
	}
	strict := &StructDefinition[boundaryDecl]{PropertyComposition: Closed, ExtensionComposition: Open, ContractComposition: Open}
	strict.IgnoreDescendantsOf[ownedBoundary]()
	if result := strict.Declare(raw, nil); result.Validation.IsOkay() {
		t.Fatal("unknown authored sibling accepted")
	}
}
