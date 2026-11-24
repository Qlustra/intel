// Package identifier defines rules for string identifiers.
package identifier

// Rules embeds the reusable rule-set API specialized to string subjects.
type Rules struct {
	RulesType
}

func (r *Rules) Exactly(expectation string) *Rule {
	rule := Exactly(expectation)
	r.SetExpectations(rule)
	return rule
}
