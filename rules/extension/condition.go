package extension

import (
	"github.com/qlustra/intel/rules"
)

func And(conditions ...RuleAPI) *ConditionAPI {
	return rules.And[Encounter](conditions...)
}

func Either(conditions ...RuleAPI) *ConditionAPI {
	return rules.Or[Encounter](conditions...)
}

func Not(conditions ...RuleAPI) *ConditionAPI {
	return rules.Not[Encounter](conditions...)
}
