// Package typerule defines rules over reflected Go types, including exact type
// identity and interface satisfaction.
package typerule

import (
	"github.com/qlustra/intel/internal/interfaces"
	"github.com/qlustra/intel/rules"
)

type (
	// Domain is the type-rule domain shared by [RuleAPI] and [RulesType].
	Domain = interfaces.TypeDomain
	// RuleAPI is the type encounter rule interface.
	RuleAPI = interfaces.RuleApi[Encounter, Domain]
	// ConditionAPI combines type rules with logical operators.
	ConditionAPI = rules.ConditionRule[Encounter, Domain]
	// RuleType is the generic rule implementation for a type encounter.
	RuleType = rules.Rule[Encounter, Domain]
	// RulesType is the configurable set of type encounter rules.
	RulesType = rules.Rules[RuleAPI, Encounter, Domain]
)
