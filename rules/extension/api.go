// Package extension defines rules for embedded struct fields. Extension rules
// can delegate to [field] rules for the embedded field descriptor.
package extension

import (
	"github.com/qlustra/intel/internal/interfaces"
	"github.com/qlustra/intel/rules"
)

type (
	// Domain is the extension-rule domain shared by [RuleAPI] and [RulesType].
	Domain = interfaces.ExtensionDomain
	// RuleAPI is the extension encounter rule interface.
	RuleAPI = interfaces.RuleApi[Encounter, Domain]
	// ConditionAPI combines extension rules with logical operators.
	ConditionAPI = rules.ConditionRule[Encounter, Domain]
	// RuleType is the generic rule implementation for an extension encounter.
	RuleType = rules.Rule[Encounter, Domain]
	// RulesType is the configurable set of extension encounter rules.
	RulesType = rules.Rules[RuleAPI, Encounter, Domain]
)
