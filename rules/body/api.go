// Package contract defines rules for a struct's body, including its embedded
// fields and properties. The encounter is a captured struct type.
package contract

import (
	"github.com/qlustra/intel/internal/interfaces"
	"github.com/qlustra/intel/rules"
)

type (
	// Domain is the struct-body rule domain shared by [RuleAPI] and [RulesType].
	Domain = interfaces.StructDomain
	// RuleAPI is the struct encounter rule interface.
	RuleAPI = interfaces.RuleApi[Encounter, Domain]
	// ConditionAPI combines struct-body rules with logical operators.
	ConditionAPI = rules.ConditionRule[Encounter, Domain]
	// RuleType is the generic rule implementation for a struct encounter.
	RuleType = rules.Rule[Encounter, Domain]
	// RulesType is the configurable set of struct encounter rules.
	RulesType = rules.Rules[RuleAPI, Encounter, Domain]
)
