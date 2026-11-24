package identifier

import (
	"fmt"
	"regexp"
	"slices"
)

// Rule is a composable rule for a string identifier.
type Rule struct {
	RuleType
}

// String creates an unconstrained identifier rule to add expectations to.
func String() *Rule {
	return &Rule{}
}

// Exactly matches the identifier only when it equals expected.
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

// Pattern matches strings accepted by the Go regular expression pattern. A
// malformed pattern configures an invalid rule that fails during evaluation.
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
