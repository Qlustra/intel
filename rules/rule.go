// Package rules provides the generic rule engine and the shared validation
// facilities used by its domain adapters.
package rules

import (
	"errors"

	"github.com/qlustra/intel/internal/interfaces"
	"github.com/qlustra/intel/logic"
	"github.com/qlustra/intel/options"
)

type (
	// Rule describes a named expectation for a domain subject. Configure its
	// predicate with [Rule.Expect], [Rule.ExpectCheck], or [Rule.ExpectMatch], and
	// combine rules with [ConditionRule] or [Rules].
	Rule[T any, D any] struct {
		interfaces.RuleSeal[D]
		name        string
		errors      []error
		interceptor func(T) logic.Check
		optionality options.Optionality
		expectation func(interfaces.Validation, T) (logic.Check, bool)
	}
)

func (r *Rule[T, D]) RuleName() string {
	return r.name
}

func (r *Rule[T, D]) Error(err error) {
	if r.errors == nil {
		r.errors = make([]error, 0)
	}
	r.errors = append(r.errors, err)
}

func (r *Rule[T, D]) Optional() {
	r.optionality = options.Optional
}

func (r *Rule[T, D]) Required() {
	r.optionality = options.Required
}

func (r *Rule[T, D]) encounter(subject T) logic.Check {
	if r.interceptor != nil {
		return r.interceptor(subject)
	}
	return logic.Ok()
}

// Encounter configures a callback run after this rule's expectation matches.
func (r *Rule[T, D]) Encounter(interceptor func(T) logic.Check) *Rule[T, D] {
	r.interceptor = interceptor
	return r
}

func (r *Rule[T, D]) Expect(name string, expectation func(T) bool) {
	if expectation == nil {
		r.ExpectCheck(name, nil)
		return
	}
	r.ExpectCheck(name, func(subject T) logic.Check {
		if expectation(subject) {
			return logic.Check{Code: logic.CheckPassed}
		}
		return logic.Check{Code: "rule_expectation_not_met", Reason: "unmet rule expectation"}
	})
}

// ExpectCheck registers a predicate: a failed check is a valid non-match.
// Use ExpectMatch when evaluation can fail independently of matching.
func (r *Rule[T, D]) ExpectCheck(name string, expectation func(T) logic.Check) {
	if expectation == nil {
		r.ExpectMatch(name, nil)
		return
	}
	r.ExpectMatch(name, func(_ interfaces.Validation, subject T) (logic.Check, bool) {
		return expectation(subject), true
	})
}

// ExpectMatch gives nested rules the current branch; the check describes the
// match and the boolean says whether evaluating it was valid.
func (r *Rule[T, D]) ExpectMatch(name string, expectation func(ValidationBranch, T) (logic.Check, bool)) {
	r.name = name
	r.expectation = expectation
}

// ValidateRule returns whether the rule matched and whether evaluation succeeded.
func (r *Rule[T, D]) ValidateRule(validation ValidationBranch, subject T) (matched bool, valid bool) {
	ruleMatch := validation.Branch(r.name)
	if len(r.errors) != 0 {
		ruleMatch.Nope(errors.Join(r.errors...).Error())
		return false, false
	}
	if r.expectation == nil {
		ruleMatch.Nope("rule expectation is not configured")
		return false, false
	}
	check, valid := r.expectation(ruleMatch, subject)
	if !valid || !check.OK() {
		reason := check.Why()
		if reason == "" {
			reason = "rule expectation evaluation failed"
		}
		ruleMatch.Nope(reason)
		return false, valid
	}
	check = r.encounter(subject)
	if !check.OK() {
		ruleMatch.Nope(check.Why())
		return true, false
	}
	ruleMatch.Okay()
	return true, true
}
