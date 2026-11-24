// Package identifier defines rules evaluated against string identifiers.
package identifier

import (
	"github.com/qlustra/intel/internal/interfaces"
	"github.com/qlustra/intel/rules"
)

type (
	// Domain is the identifier-rule domain shared by [RuleAPI] and [RulesType].
	Domain = interfaces.StringDomain
	// RuleAPI is the string identifier rule interface.
	RuleAPI = interfaces.RuleApi[string, Domain]
	// ConditionAPI combines identifier rules with logical operators.
	ConditionAPI = rules.ConditionRule[string, Domain]
	// RuleType is the generic rule implementation for string subjects.
	RuleType = rules.Rule[string, Domain]
	// RulesType is the configurable set of identifier rules.
	RulesType = rules.Rules[RuleAPI, string, Domain]
)
