package contract

import (
	"github.com/qlustra/intel/rules"
)

// And matches when every contract condition matches.
func And(conditions ...RuleAPI) *ConditionAPI {
	return rules.And[Encounter](conditions...)
}

// Either matches when at least one contract condition matches.
func Either(conditions ...RuleAPI) *ConditionAPI {
	return rules.Or[Encounter](conditions...)
}

// Not matches when none of the contract conditions match.
func Not(conditions ...RuleAPI) *ConditionAPI {
	return rules.Not[Encounter](conditions...)
}
