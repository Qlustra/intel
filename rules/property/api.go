// Package property defines rules for struct properties. Property rules can
// delegate to [field] rules or inspect property types and struct slices.
package property

import (
	"github.com/qlustra/intel/internal/interfaces"
	"github.com/qlustra/intel/rules"
)

type (
	// Domain is the property-rule domain shared by [RuleAPI] and [SetType].
	Domain = interfaces.PropertyDomain
	// RuleAPI is the property encounter rule interface.
	RuleAPI = interfaces.RuleApi[Encounter, Domain]
	// ConditionAPI combines property rules with logical operators.
	ConditionAPI = rules.ConditionRule[Encounter, Domain]
	// RuleType is the generic rule implementation for a property encounter.
	RuleType = rules.Rule[Encounter, Domain]
	// SetType is the configurable set of property encounter rules.
	SetType = rules.Rules[RuleAPI, Encounter, Domain]
)
