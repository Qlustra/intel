package extension

import (
	"fmt"

	"github.com/qlustra/intel/capture"
	"github.com/qlustra/intel/internal/interfaces"
	"github.com/qlustra/intel/logic"
	"github.com/qlustra/intel/rules/field"
)

const (
	NotMatchedCode       logic.CheckCode = "extension_field_not_matched"
	EvaluationFailedCode logic.CheckCode = "extension_field_evaluation_failed"
)

type Rule struct {
	RuleType
}

func Extension() *Rule {
	return &Rule{}
}

func Field(fieldRule ...field.RuleAPI) *Rule {
	r := Extension()
	rules := field.Rules{}
	rules.Expect(fieldRule...)
	r.ExpectMatch(fmt.Sprintf("field matches"), func(validation interfaces.Validation, encounter Encounter) (logic.Check, bool) {
		fieldEncounter := field.Encounter{
			EncounterContext: encounter.EncounterContext,
			Struct:           encounter.Struct,
			Field:            encounter.Extension.FieldDescriptor,
		}
		satisfied, valid := rules.ValidateRules(validation, fieldEncounter)
		if !valid {
			return logic.Nope(EvaluationFailedCode, "extension field rule evaluation failed"), false
		}
		if !satisfied {
			return logic.Nope(NotMatchedCode, "extension field rule was not satisfied"), true
		}
		return logic.Ok(), true
	})
	return r
}

func DeclaredBy(name string) *Rule {
	r := Extension()
	r.ExpectCheck(fmt.Sprintf("declared by name %v", name), func(encounter Encounter) logic.Check {
		return encounter.Extension.DeclaredBy(name)
	})
	return r
}

func DeclaredAs[T any]() *Rule {
	r := Extension()
	interrogator := capture.Type[T]()
	r.ExpectCheck(fmt.Sprintf("declare as %v", interrogator.String()), func(encounter Encounter) logic.Check {
		return encounter.Extension.DeclaredAs[T]()
	})
	return r
}

func Satisfies[T any]() *Rule {
	r := Extension()
	interrogator := capture.Type[T]()
	r.ExpectCheck(fmt.Sprintf("satisfies %s", interrogator.String()), func(encounter Encounter) logic.Check {
		structInterrogator, check := encounter.Extension.AsStruct()
		if !check.OK() {
			return check
		}
		return structInterrogator.SatisfiesType[T]()
	})
	return r
}
