package identifier

import "github.com/qlustra/intel/rules"

func And(conditions ...RuleAPI) *ConditionAPI {
	return rules.And[string](conditions...)
}

func Either(conditions ...RuleAPI) *ConditionAPI {
	return rules.Or[string](conditions...)
}

func Not(conditions ...RuleAPI) *ConditionAPI {
	return rules.Not[string](conditions...)
}
