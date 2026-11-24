// Package logic provides checks for operations that distinguish success from
// a rejected input or an unavailable result.
package logic

import "fmt"

// CheckCode identifies a check outcome independently of its diagnostic text.
type CheckCode string

const (
	CheckPassed CheckCode = "passed"
)

// Check reports whether an operation succeeded. Test [Check.OK] before
// using the result it accompanies. The zero value does not indicate success.
type Check struct {
	Code   CheckCode
	Reason string
	Cause  *Check
}

// OK reports whether the check has code [CheckPassed].
func (c Check) OK() bool {
	return c.Code == CheckPassed
}

// Why renders a diagnostic for a failed check, including any cause. It
// returns an empty string on success. The wording is not a parsing contract.
func (c Check) Why() string {
	if c.Code == CheckPassed {
		return ""
	}
	if c.Cause == nil {
		return fmt.Sprintf("%s (%s)", c.Code, c.Reason)
	}
	if c.Reason == "" {
		return fmt.Sprintf("%s because %s", c.Code, c.Cause.Why())
	}
	return fmt.Sprintf("%s (%s) because %s", c.Code, c.Reason, c.Cause.Why())
}

// Nope builds a check with the supplied code and reason. Use a code other
// than [CheckPassed] to report failure.
func Nope(nopeCode CheckCode, reason string) Check {
	return Check{
		Code:   nopeCode,
		Reason: reason,
	}
}

// NopeBecause builds a check like [Nope] and attaches the supplied cause.
func NopeBecause(nopeCode CheckCode, reason string, cause Check) Check {
	return Check{
		Code:   nopeCode,
		Reason: reason,
		Cause:  &cause,
	}
}

// Ok returns a successful check.
func Ok() Check { return Check{Code: CheckPassed} }
