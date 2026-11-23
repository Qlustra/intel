package typerule

import (
	"github.com/qlustra/intel/internal/interfaces"
	"github.com/qlustra/intel/rules"
)

type (
	Domain       = interfaces.TypeDomain
	RuleAPI      = interfaces.RuleApi[Encounter, Domain]
	ConditionAPI = rules.ConditionRule[Encounter, Domain]
	RuleType     = rules.Rule[Encounter, Domain]
	RulesType    = rules.Rules[RuleAPI, Encounter, Domain]
)
