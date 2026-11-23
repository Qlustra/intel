package property

import (
	"github.com/qlustra/intel/internal/interfaces"
	"github.com/qlustra/intel/rules"
)

type (
	Domain       = interfaces.PropertyDomain
	RuleAPI      = interfaces.RuleApi[Encounter, Domain]
	ConditionAPI = rules.ConditionRule[Encounter, Domain]
	RuleType     = rules.Rule[Encounter, Domain]
	SetType      = rules.Rules[RuleAPI, Encounter, Domain]
)
