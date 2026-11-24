package property

import (
	"fmt"

	"github.com/qlustra/intel/capture"
	"github.com/qlustra/intel/internal/interfaces"
	"github.com/qlustra/intel/logic"
	"github.com/qlustra/intel/rules/field"
	"github.com/qlustra/intel/rules/typerule"
)

const (
	NotTaggedAsCode                 logic.CheckCode = "property_not_tagged_as"
	NotMatchedCode                  logic.CheckCode = "property_field_not_matched"
	EvaluationFailedCode            logic.CheckCode = "property_field_evaluation_failed"
	StructSliceNotMatchedCode       logic.CheckCode = "property_struct_slice_not_matched"
	StructSliceEvaluationFailedCode logic.CheckCode = "property_struct_slice_evaluation_failed"
)

// Rule is a composable rule for a property encounter.
type Rule struct {
	RuleType
}

// Property creates an unconstrained property rule to which expectations may be added.
func Property() *Rule {
	return &Rule{}
}

func Field(fieldRule ...field.RuleAPI) *Rule {
	r := Property()
	rules := field.Rules{}
	rules.Expect(fieldRule...)
	r.ExpectMatch(fmt.Sprintf("field matches"), func(validation interfaces.Validation, encounter Encounter) (logic.Check, bool) {
		fieldEncounter := field.Encounter{
			EncounterContext: encounter.EncounterContext,
			Struct:           encounter.Struct,
			Field:            encounter.Property.FieldDescriptor,
		}
		satisfied, valid := rules.ValidateRules(validation, fieldEncounter)
		if !valid {
			return logic.Nope(EvaluationFailedCode, "property field rule evaluation failed"), false
		}
		if !satisfied {
			return logic.Nope(NotMatchedCode, "property field rule was not satisfied"), true
		}
		return logic.Ok(), true
	})
	return r
}

// StructSlice matches a slice with immediate struct elements, evaluating the
// supplied type rules against its element without changing the property encounter.
func StructSlice(elementRules ...typerule.RuleAPI) *Rule {
	r := Property()
	rules := typerule.RulesType{}
	rules.Expect(elementRules...)
	r.ExpectMatch("struct slice element matches", func(validation interfaces.Validation, encounter Encounter) (logic.Check, bool) {
		slice, check := encounter.Property.AsSlice()
		if !check.OK() {
			return check, true
		}
		structSlice, check := slice.AsStructSlice()
		if !check.OK() {
			return check, true
		}
		elementEncounter := typerule.Encounter{
			EncounterContext: encounter.EncounterContext,
			Type:             structSlice.Element().CapturedType(),
		}
		satisfied, valid := rules.ValidateRules(validation, elementEncounter)
		if !valid {
			return logic.Nope(StructSliceEvaluationFailedCode, "struct slice element rule evaluation failed"), false
		}
		if !satisfied {
			return logic.Nope(StructSliceNotMatchedCode, "struct slice element rules were not satisfied"), true
		}
		return logic.Ok(), true
	})
	return r
}

func DeclaredBy(name string) *Rule {
	r := Property()
	r.ExpectCheck(fmt.Sprintf("declared by name %v", name), func(encounter Encounter) logic.Check {
		return encounter.Property.DeclaredBy(name)
	})
	return r
}

func DeclaredAs[T any]() *Rule {
	r := Property()
	interrogator := capture.Type[T]()
	r.ExpectCheck(fmt.Sprintf("declare as %v", interrogator.String()), func(encounter Encounter) logic.Check {
		return encounter.Property.DeclaredAs[T]()
	})
	return r
}

func Extends[T any]() *Rule {
	r := Property()
	interrogator := capture.Type[T]()
	r.ExpectCheck(fmt.Sprintf("extends %v", interrogator.String()), func(encounter Encounter) logic.Check {
		structInterrogator, check := encounter.Property.AsStruct()
		if !check.OK() {
			return check
		}
		return structInterrogator.ExtendsStruct[T]()
	})
	return r
}

func Satisfies[T any]() *Rule {
	r := Property()
	interrogator := capture.Type[T]()
	r.ExpectCheck(fmt.Sprintf("satisfies %s", interrogator.String()), func(encounter Encounter) logic.Check {
		structInterrogator, check := encounter.Property.AsStruct()
		if !check.OK() {
			return check
		}
		return structInterrogator.SatisfiesType[T]()
	})
	return r
}
