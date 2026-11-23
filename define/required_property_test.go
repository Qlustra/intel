package define

import (
	"testing"

	"github.com/qlustra/intel/capture"
	"github.com/qlustra/intel/internal/interfaces"
	"github.com/qlustra/intel/logic"
	"github.com/qlustra/intel/options"
	"github.com/qlustra/intel/rules/property"
)

func requiredPropertyDefinition(composition BodyComposition) *StructDefinition[boundaryDecl] {
	return &StructDefinition[boundaryDecl]{PropertyComposition: composition, ExtensionComposition: Open, ContractComposition: Open}
}

func requiredPropertyNamed(name string) property.RuleAPI {
	r := property.Property()
	r.Expect("field name", func(e property.Encounter) bool { return e.Property.FieldName() == name })
	return r
}

func requiredPropertyBody[T any](t *testing.T) *capture.StructInterrogator {
	t.Helper()
	body, check := capture.Type[T]().AsStruct()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	return body
}

func TestRequiredPropertyEmptyAndFiltered(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		d := requiredPropertyDefinition(Open)
		d.Property(options.Required, requiredPropertyNamed("Wanted"))
		if result := d.Declare(requiredPropertyBody[struct{}](t), nil); result.Validation.IsOkay() {
			t.Fatal("empty body satisfied required property")
		}
	})
	t.Run("filtered", func(t *testing.T) {
		d := requiredPropertyDefinition(Open)
		d.IgnoreDescendantsOf[ownedBoundary]()
		d.Property(options.Required, requiredPropertyNamed("Hidden"))
		if result := d.Declare(requiredPropertyBody[struct{ ownedBoundary }](t), nil); result.Validation.IsOkay() {
			t.Fatal("filtered descendant satisfied required property")
		}
	})
	t.Run("optional empty", func(t *testing.T) {
		d := requiredPropertyDefinition(Open)
		d.Property(options.Optional, requiredPropertyNamed("Wanted"))
		if result := d.Declare(requiredPropertyBody[struct{}](t), nil); !result.Validation.IsOkay() {
			t.Fatal(result.Validation.Reason())
		}
	})
}

func TestRequiredPropertyMatchAnywhere(t *testing.T) {
	for _, tc := range []struct {
		name string
		body *capture.StructInterrogator
	}{
		{"unrelated first", requiredPropertyBody[struct{ Other, Wanted string }](t)},
		{"unrelated last", requiredPropertyBody[struct{ Wanted, Other string }](t)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := requiredPropertyDefinition(Open)
			d.Property(options.Required, requiredPropertyNamed("Wanted"))
			if result := d.Declare(tc.body, nil); !result.Validation.IsOkay() {
				t.Fatal(result.Validation.Reason())
			}
		})
	}
}

func TestRequiredPropertyIndependentRegistrations(t *testing.T) {
	for _, tc := range []struct {
		name  string
		body  *capture.StructInterrogator
		valid bool
	}{
		{"both", requiredPropertyBody[struct{ First, Second string }](t), true},
		{"missing second", requiredPropertyBody[struct{ First string }](t), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := requiredPropertyDefinition(Open)
			d.Property(options.Required, requiredPropertyNamed("First"))
			d.Property(options.Required, requiredPropertyNamed("Second"))
			result := d.Declare(tc.body, nil)
			if result.Validation.IsOkay() != tc.valid {
				t.Fatalf("valid=%v, want %v: %s", result.Validation.IsOkay(), tc.valid, result.Validation.Reason())
			}
		})
	}
	t.Run("shared field", func(t *testing.T) {
		d := requiredPropertyDefinition(Open)
		calls := 0
		for range 2 {
			d.Property(options.Required, requiredPropertyNamed("First")).Encounter(func(property.Encounter) logic.Check { calls++; return logic.Ok() })
		}
		result := d.Declare(requiredPropertyBody[struct{ First string }](t), nil)
		if !result.Validation.IsOkay() || calls != 2 {
			t.Fatalf("calls=%d: %s", calls, result.Validation.Reason())
		}
	})
}

func TestRequiredPropertyAllMatchesHandled(t *testing.T) {
	d := requiredPropertyDefinition(Open)
	r := property.Property()
	r.Expect("all", func(property.Encounter) bool { return true })
	var names []string
	d.Property(options.Required, r).Encounter(func(e property.Encounter) logic.Check {
		names = append(names, e.Property.FieldName())
		return logic.Ok()
	})
	result := d.Declare(requiredPropertyBody[struct{ First, Second string }](t), nil)
	if !result.Validation.IsOkay() || len(names) != 2 || names[0] != "First" || names[1] != "Second" {
		t.Fatalf("callbacks=%v: %s", names, result.Validation.Reason())
	}
}

func TestRequiredPropertySessionsIndependent(t *testing.T) {
	d := requiredPropertyDefinition(Open)
	d.Property(options.Required, requiredPropertyNamed("Wanted"))
	matched := requiredPropertyBody[struct{ Wanted string }](t)
	missing := requiredPropertyBody[struct{}](t)
	for _, compile := range []bool{false, true} {
		for _, tc := range []struct {
			body  *capture.StructInterrogator
			valid bool
		}{{matched, true}, {missing, false}, {matched, true}} {
			var valid bool
			if compile {
				_, validation := d.Compile(tc.body, nil)
				valid = validation.IsOkay()
			} else {
				valid = d.Declare(tc.body, nil).Validation.IsOkay()
			}
			if valid != tc.valid {
				t.Errorf("compile=%v body=%s: valid=%v want %v", compile, tc.body, valid, tc.valid)
			}
		}
	}
}

func TestRequiredPropertyClosedRejectsUnknown(t *testing.T) {
	for _, tc := range []struct {
		name string
		body *capture.StructInterrogator
	}{
		{"ordinary", requiredPropertyBody[struct{ Wanted, Other string }](t)},
		{"private", requiredPropertyBody[struct{ Wanted, other string }](t)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := requiredPropertyDefinition(Closed)
			d.Property(options.Required, requiredPropertyNamed("Wanted"))
			if result := d.Declare(tc.body, nil); result.Validation.IsOkay() {
				t.Fatal("closed composition accepted unknown property")
			}
		})
	}
}

func TestRequiredPropertyLaterFailureFatal(t *testing.T) {
	for _, optionality := range []options.Optionality{options.Required, options.Optional} {
		d := requiredPropertyDefinition(Open)
		r := property.Property()
		r.Expect("all", func(property.Encounter) bool { return true })
		calls := 0
		d.Property(optionality, r).Encounter(func(property.Encounter) logic.Check {
			calls++
			if calls == 2 {
				return logic.Nope("test_failure", "second callback failed")
			}
			return logic.Ok()
		})
		result := d.Declare(requiredPropertyBody[struct{ First, Second string }](t), nil)
		if result.Validation.IsOkay() || calls != 2 {
			t.Errorf("optionality=%v calls=%d: %s", optionality, calls, result.Validation.Reason())
		}
	}
}

func TestRequiredPropertyRecognitionInvalid(t *testing.T) {
	for _, optionality := range []options.Optionality{options.Required, options.Optional} {
		d := requiredPropertyDefinition(Open)
		r := property.Property()
		r.ExpectMatch("invalid second recognition", func(_ interfaces.Validation, e property.Encounter) (logic.Check, bool) {
			if e.Property.FieldName() == "Second" {
				return logic.Nope("test_failure", "recognition invalid"), false
			}
			return logic.Ok(), true
		})
		d.Property(optionality, r)
		result := d.Declare(requiredPropertyBody[struct{ First, Second string }](t), nil)
		if result.Validation.IsOkay() {
			t.Errorf("optionality=%v accepted invalid recognition", optionality)
		}
	}
}
