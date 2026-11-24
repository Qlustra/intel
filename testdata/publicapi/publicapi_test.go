package publicapi_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/qlustra/intel"
	"github.com/qlustra/intel/compile"
	"github.com/qlustra/intel/define"
	"github.com/qlustra/intel/logic"
	"github.com/qlustra/intel/options"
	"github.com/qlustra/intel/rules"
	"github.com/qlustra/intel/rules/field"
	"github.com/qlustra/intel/rules/property"
)

type Profile struct {
	Name string `json:"name"`
}

type User struct{ Profile }

func TestREADMEPromotedField(t *testing.T) {
	user, check := intel.CaptureType[User]().AsStruct()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	field, check := user.Field("Name")
	if !check.OK() {
		t.Fatal(check.Why())
	}
	if field.FieldName() != "Name" || field.ReflectedType() != reflect.TypeFor[string]() || field.StructField().Tag.Get("json") != "name" {
		t.Fatalf("field = %s %v %q; want Name string name", field.FieldName(), field.ReflectedType(), field.StructField().Tag.Get("json"))
	}
}

// Embedding the interface supplies declaration membership without calling its
// nil seal methods.
type declaration struct {
	define.ScopeDeclaration
	FieldName string
}

type definition struct {
	define.StructDefinition[declaration]
}

func (d *definition) Define() {
	rule := property.Property()
	rule.Expect("field named Name", func(e property.Encounter) bool {
		return e.Property.FieldName() == "Name"
	})
	d.Property(options.Required, rule).Encounter(func(e property.Encounter) logic.Check {
		declared, check := e.Access[declaration](define.DeclarationKey)
		if !check.OK() {
			return check
		}
		declared.FieldName = e.Property.FieldName()
		return logic.Ok()
	})
}

func TestCompilePublicDefinition(t *testing.T) {
	d, err := define.Struct[definition, declaration](nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("success", func(t *testing.T) {
		compiler := compile.Struct[struct{ Name string }](d)
		if !compiler.IsOkay() {
			t.Fatal(compiler.Trace())
		}
		declared, check := compiler.GetDeclaredValue()
		if !check.OK() {
			t.Fatal(check.Why())
		}
		if declared == nil || declared.FieldName != "Name" {
			t.Fatalf("declared = %#v; want field Name", declared)
		}
	})
	t.Run("missing required field", func(t *testing.T) {
		compiler := compile.Struct[struct{}](d)
		if compiler.IsOkay() {
			t.Fatal("compilation accepted missing required field")
		}
		trace := compiler.Trace()
		for _, fragment := range []string{"required", "propert"} {
			if !strings.Contains(strings.ToLower(trace), fragment) {
				t.Errorf("trace %q does not contain %q", trace, fragment)
			}
		}
		if declared, check := compiler.GetDeclaredValue(); check.OK() || declared != nil {
			t.Fatalf("failed compilation returned declaration %#v, check %v", declared, check)
		}
	})
}

func TestExpectMatchWithNestedValidation(t *testing.T) {
	shape, check := intel.CaptureType[struct {
		Name string `json:"name"`
		Age  int
	}]().AsStruct()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	for _, tc := range []struct {
		name, propertyName string
		nested             field.RuleAPI
		matched, valid     bool
	}{
		{"match", "Name", field.TaggedWith("json"), true, true},
		{"non-match", "Age", field.TaggedWith("json"), false, true},
		{"failed evaluation", "Name", field.Field(), false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			member, check := shape.Property(tc.propertyName)
			if !check.OK() {
				t.Fatal(check.Why())
			}
			rule := property.Property()
			rule.ExpectMatch("property tag", func(branch rules.ValidationBranch, e property.Encounter) (logic.Check, bool) {
				matched, valid := tc.nested.ValidateRule(branch.Branch("field checks"), field.Encounter{
					Struct: e.Struct,
					Field:  e.Property.FieldDescriptor,
				})
				if !valid || !matched {
					return logic.Nope("fixture_field_rule_unsatisfied", "nested field rule did not match"), valid
				}
				return logic.Ok(), true
			})
			validation := rules.NewValidation("property validation")
			matched, valid := rule.ValidateRule(validation, property.Encounter{Struct: shape, Property: member})
			if matched != tc.matched || valid != tc.valid {
				t.Fatalf("match=%v valid=%v; want match=%v valid=%v\n%s", matched, valid, tc.matched, tc.valid, validation.Trace())
			}
			if trace := validation.Trace(); !strings.Contains(trace, "field checks") {
				t.Fatalf("trace did not include the supplied nested branch:\n%s", trace)
			}
		})
	}
}
