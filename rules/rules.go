// Package rules provides composable rules, rule sets, logical conditions, and
// validation traces. The domain adapters under this package specialize the
// generic APIs for identifiers, types, struct bodies, fields, properties,
// extensions, and contracts.
package rules

import (
	"github.com/qlustra/intel/internal/interfaces"
	"github.com/qlustra/intel/logic"
	"github.com/qlustra/intel/options"
)

// ======================================== Enums

type (
	multiplicity int
	matching     int
)

// ======================================== Multiplicity

const (
	/* none means we expect there to be no match */
	none multiplicity = iota
	/* one means we expect there to be only one match */
	one
	/* some means we expect there to be zero or more matches */
	some
	/* many means we expect there to be one or more match */
	many
	/* bounded means we expect there to be a match in some range */
	bounded
)

// ======================================== Matching

const (
	/* exact means we expect the match to be exact */
	exact matching = iota
	/* other means we expect the match to be anything other than the expectation */
	other
)

// ======================================== Struct

type Rules[R interfaces.RuleApi[T, D], T any, D any] struct {
	optionality options.Optionality
	shouldMatch bool
	recognizer  R
	rules       []R
	interceptor func(T) logic.Check
}

// ======================================== Rules

func (r *Rules[R, T, D]) ensureList() {
	if r.rules == nil {
		r.rules = make([]R, 0)
	}
}

func (r *Rules[R, T, D]) addExpectations(rules ...R) {
	r.ensureList()
	r.rules = append(r.rules, rules...)
}

// SetExpectations replaces the set's expectations with rules, in order.
func (r *Rules[R, T, D]) SetExpectations(rules ...R) {
	r.ensureList()
	r.rules = rules
}

// Expect replaces expectations; use [Rules.Also] to append instead.
func (r *Rules[R, T, D]) Expect(rules ...R) {
	r.SetExpectations(rules...)
}

// Also appends expectations after those already configured.
func (r *Rules[R, T, D]) Also(rules ...R) {
	r.addExpectations(rules...)
}

// Match configures a recognizer and its optionality. An unmatched optional
// recognizer skips expectations and the set's encounter callback.
func (r *Rules[R, T, D]) Match(optionality options.Optionality, recognizer R) {
	r.optionality = optionality
	r.recognizer = recognizer
	r.shouldMatch = true
}

// MatchSubject reports the match and whether recognition was valid.
func (r *Rules[R, T, D]) MatchSubject(validation interfaces.Validation, subject T) (matched bool, valid bool) {
	if !r.shouldMatch {
		return false, true
	}
	return r.recognizer.ValidateRule(validation, subject)
}

func (r *Rules[R, T, D]) RequiresMatch() bool {
	return r.shouldMatch && r.optionality == options.Required
}

// ValidateRules reports satisfaction (including optional absence) and validity.
func (r *Rules[R, T, D]) ValidateRules(validation interfaces.Validation, subject T) (satisfied bool, valid bool) {
	if r.shouldMatch {
		matched, valid := r.recognizer.ValidateRule(validation, subject)
		if !valid {
			validation.Nope("rule recognition failed")
			return false, false
		}
		if !matched {
			if r.optionality == options.Required {
				validation.Nope("required rule unmatched")
				return false, true
			}
			// An unmatched optional selector skips expectations and the encounter.
			validation.Okay()
			return true, true
		}
	}
	return r.ValidateExpectations(validation, subject)
}

// ValidateExpectations validates an already recognized subject without running
// the recognizer or its encounter callback a second time. It reports
// satisfaction and whether evaluation was valid.
func (r *Rules[R, T, D]) ValidateExpectations(validation interfaces.Validation, subject T) (satisfied bool, valid bool) {
	for _, rule := range r.rules {
		matched, valid := rule.ValidateRule(validation, subject)
		if !valid {
			validation.Nope("rule expectation evaluation failed")
			return false, false
		}
		if !matched {
			validation.Nope("rule expectation not satisfied")
			return false, true
		}
	}
	check := r.encounter(subject)
	if !check.OK() {
		validation.Nope(check.Why())
		return false, false
	}
	validation.Okay()
	return true, true
}

func (r Rules[R, T, D]) encounter(subject T) logic.Check {
	if r.interceptor != nil {
		return r.interceptor(subject)
	}
	return logic.Ok()
}

// Encounter configures a callback run after expectations are satisfied.
func (r *Rules[R, T, D]) Encounter(interceptor func(T) logic.Check) *Rules[R, T, D] {
	r.interceptor = interceptor
	return r
}

// ======================================== String
