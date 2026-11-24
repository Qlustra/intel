// Package field defines rules for inspecting struct fields and their tags.
// Its encounter includes both the containing struct and the field descriptor.
package field

import (
	"github.com/qlustra/intel/internal/interfaces"
	"github.com/qlustra/intel/rules"
)

type (
	// Domain is the field-rule domain shared by [RuleAPI] and [RulesType].
	Domain = interfaces.FieldDomain
	// RuleAPI is the field encounter rule interface.
	RuleAPI = interfaces.RuleApi[Encounter, Domain]
	// ConditionAPI combines field rules using logical operators.
	ConditionAPI = rules.ConditionRule[Encounter, Domain]
	// RuleType is the generic rule implementation for a field encounter.
	RuleType = rules.Rule[Encounter, Domain]
	// RulesType is the configurable set of field encounter rules.
	RulesType = rules.Rules[RuleAPI, Encounter, Domain]
)
