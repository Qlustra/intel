package logic

import "fmt"

type CheckCode string

const (
	CheckPassed CheckCode = "passed"
)

type Check struct {
	Code   CheckCode
	Reason string
	Cause  *Check
}

func (c Check) OK() bool {
	return c.Code == CheckPassed
}

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

func Nope(nopeCode CheckCode, reason string) Check {
	return Check{
		Code:   nopeCode,
		Reason: reason,
	}
}

func NopeBecause(nopeCode CheckCode, reason string, cause Check) Check {
	return Check{
		Code:   nopeCode,
		Reason: reason,
		Cause:  &cause,
	}
}

func Ok() Check { return Check{Code: CheckPassed} }
