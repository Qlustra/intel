// Package rules provides generic rules, including logical combinations of
// rules evaluated against the same subject.
package rules

import (
	"github.com/qlustra/intel/internal/interfaces"
)

type (
	logicalOperator int
)

const (
	logicalAnd logicalOperator = iota
	logicalOr
	logicalNone
)

// ConditionRule combines rules using conjunction, disjunction, or negation.
// Construct it with [And], [Or], or [Not].
type ConditionRule[T any, D any] struct {
	Rule[T, D]
	operator   logicalOperator
	conditions []interfaces.RuleApi[T, D]
}

// And matches when every condition matches; evaluation stops at the first
// non-match or invalid evaluation.
func And[T any, D any](conditions ...interfaces.RuleApi[T, D]) *ConditionRule[T, D] {
	return &ConditionRule[T, D]{
		operator:   logicalAnd,
		conditions: conditions,
	}
}

// Or matches when at least one condition matches.
func Or[T any, D any](conditions ...interfaces.RuleApi[T, D]) *ConditionRule[T, D] {
	return &ConditionRule[T, D]{
		operator:   logicalOr,
		conditions: conditions,
	}
}

// Not matches when none of its conditions match.
func Not[T any, D any](conditions ...interfaces.RuleApi[T, D]) *ConditionRule[T, D] {
	return &ConditionRule[T, D]{
		operator:   logicalNone,
		conditions: conditions,
	}
}

// ValidateRule reports the logical match and whether its evaluation was valid.
func (r *ConditionRule[T, D]) ValidateRule(validation interfaces.Validation, subject T) (matched bool, valid bool) {
	switch r.operator {
	case logicalAnd:
		branch := validation.Branch("ALL")
		for _, rule := range r.conditions {
			matched, valid := rule.ValidateRule(branch, subject)
			if !valid {
				branch.Nope("condition evaluation failed")
				return false, false
			}
			if !matched {
				branch.Nope("not all condition were satisfied")
				return false, true
			}
		}
		branch.Okay()
		return true, true
	case logicalOr:
		branch := validation.Branch("EITHER")
		for _, rule := range r.conditions {
			matched, valid := rule.ValidateRule(branch, subject)
			if !valid {
				branch.Nope("condition evaluation failed")
				return false, false
			}
			if matched {
				branch.Okay()
				return true, true
			}
		}
		branch.Nope("no OR alternative was satisfied")
		return false, true
	case logicalNone:
		branch := validation.Branch("NONE")
		for _, rule := range r.conditions {
			matched, valid := rule.ValidateRule(branch, subject)
			if !valid {
				branch.Nope("condition evaluation failed")
				return false, false
			}
			if matched {
				branch.Nope("a forbidden alternative was satisfied")
				return false, true
			}
		}
		branch.Okay()
		return true, true
	default:
		validation.Nope("unknown logical operator")
		return false, false
	}
}
