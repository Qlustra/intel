package identifier

type Rules struct {
	RulesType
}

func (r *Rules) Exactly(expectation string) *Rule {
	rule := Exactly(expectation)
	r.SetExpectations(rule)
	return rule
}
