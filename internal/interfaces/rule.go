package interfaces

import (
	"github.com/qlustra/intel/logic"
)

type StringDomain struct{}
type TypeDomain struct{}
type StructDomain struct{}
type ContractDomain struct{}
type ExtensionDomain struct{}
type PropertyDomain struct{}
type FieldDomain struct{}

type RuleSeal[D any] struct{}

func (RuleSeal[D]) ruleSeal(D) {}

type RuleApi[T any, D any] interface {
	ruleSeal(D)
	RuleName() string
	Optional()
	Required()
	Expect(name string, expectation func(T) bool)
	ExpectCheck(name string, expectation func(T) logic.Check)
	// ValidateRule reports the match and whether evaluation was valid.
	ValidateRule(validation Validation, subject T) (matched bool, valid bool)
}
