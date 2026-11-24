package identifier

import "github.com/qlustra/intel/rules"

// And matches when every identifier condition matches.
func And(conditions ...RuleAPI) *ConditionAPI {
	return rules.And[string](conditions...)
}

// Either matches when at least one identifier condition matches.
func Either(conditions ...RuleAPI) *ConditionAPI {
	return rules.Or[string](conditions...)
}

// Not matches when none of the identifier conditions match.
func Not(conditions ...RuleAPI) *ConditionAPI {
	return rules.Not[string](conditions...)
}
