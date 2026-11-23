package identifier

import (
	"github.com/qlustra/intel/internal/interfaces"
	"github.com/qlustra/intel/rules"
)

type (
	Domain       = interfaces.StringDomain
	RuleAPI      = interfaces.RuleApi[string, Domain]
	ConditionAPI = rules.ConditionRule[string, Domain]
	RuleType     = rules.Rule[string, Domain]
	RulesType    = rules.Rules[RuleAPI, string, Domain]
)
