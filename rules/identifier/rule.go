package identifier

import (
	"fmt"
	"regexp"
	"slices"
)

type Rule struct {
	RuleType
}

func String() *Rule {
	return &Rule{}
}

func Exactly(expected string) *Rule {
	r := String()
	r.Expect("exactly "+expected, func(subject string) bool {
		return subject == expected
	})
	return r
}

func Among(expected ...string) *Rule {
	r := String()
	r.Expect(fmt.Sprintf("among %v", expected), func(subject string) bool {
		return slices.Contains(expected, subject)
	})
	return r
}

func NotAmong(expected ...string) *Rule {
	r := String()
	r.Expect(fmt.Sprintf("not among %v", expected), func(subject string) bool {
		return !slices.Contains(expected, subject)
	})
	return r
}

func Pattern(pattern string) *Rule {
	r := String()
	matcher, err := regexp.Compile(pattern)
	if err != nil {
		r.Error(err)
		return r
	}
	r.Expect(fmt.Sprintf("match pattern %v", pattern), func(subject string) bool {
		return matcher.MatchString(subject)
	})
	return r
}
