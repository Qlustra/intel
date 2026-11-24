// Package contract defines rules over struct bodies.
package contract

import (
	"fmt"

	"github.com/qlustra/intel/capture"
	"github.com/qlustra/intel/logic"
)

// Rule is a composable rule for a captured struct type encounter.
type Rule struct {
	RuleType
}

// Contract creates an unconstrained struct-body rule.
func Contract() *Rule {
	return &Rule{}
}

func (r *Rule) ContractEmbeddedAs(name string) *Rule {
	r.ExpectCheck(fmt.Sprintf("lookup %v", name), func(encounter Encounter) logic.Check {
		_, check := encounter.Contracts.Lookup(name)
		return check
	})
	return r
}

func (r *Rule) ContractDeclaredAs(name string) *Rule {
	r.ExpectCheck(fmt.Sprintf("declared as %v", name), func(encounter Encounter) logic.Check {
		_, check := encounter.Contracts.Declared(name)
		return check
	})
	return r
}

func (r *Rule) ContractEmbeds[T any]() *Rule {
	interrogator := capture.Type[T]()
	r.ExpectCheck(fmt.Sprintf("is %v an embedded interface", interrogator.String()), func(encounter Encounter) logic.Check {
		_, check := encounter.Contracts.Embeds[T]()
		return check
	})
	return r
}

func (r *Rule) ContractSatisfies[T any]() *Rule {
	interrogator := capture.Type[T]()
	r.ExpectCheck(fmt.Sprintf("satisfies %v", interrogator.String()), func(encounter Encounter) logic.Check {
		_, check := encounter.Contracts.Satisfies[T]()
		return check
	})
	return r
}
