// Package contract defines rules for embedded interface fields, also called
// contracts. Its encounter retains the containing struct and the contract.
package contract

import (
	"github.com/qlustra/intel/internal/interfaces"
	"github.com/qlustra/intel/rules"
)

type (
	// Domain is the contract-rule domain shared by [RuleAPI] and [RulesType].
	Domain = interfaces.ContractDomain
	// RuleAPI is the contract encounter rule interface.
	RuleAPI = interfaces.RuleApi[Encounter, Domain]
	// ConditionAPI combines contract rules with logical operators.
	ConditionAPI = rules.ConditionRule[Encounter, Domain]
	// RuleType is the generic rule implementation for a contract encounter.
	RuleType = rules.Rule[Encounter, Domain]
	// RulesType is the configurable set of contract encounter rules.
	RulesType = rules.Rules[RuleAPI, Encounter, Domain]
)
