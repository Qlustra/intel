package typerule

import (
	"fmt"
	"reflect"

	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/logic"
)

type Rule struct {
	RuleType
}

func Type() *Rule {
	return &Rule{}
}

func (r *Rule) Is[T any]() *Rule {
	return r.IsType(reflect.TypeFor[T]())
}

func (r *Rule) IsType(expectedType reflect.Type) *Rule {
	r.ExpectCheck(fmt.Sprintf("is %v", expectedType), func(encounter Encounter) logic.Check {
		return inspect.IsType(encounter.Type, expectedType)
	})
	return r
}

func (r *Rule) Satisfies[T any]() *Rule {
	return r.SatisfiesType(reflect.TypeFor[T]())
}

func (r *Rule) SatisfiesType(expectedType reflect.Type) *Rule {
	r.ExpectCheck(fmt.Sprintf("satisfies %v", expectedType), func(encounter Encounter) logic.Check {
		return inspect.SatisfiesType(encounter.Type, expectedType)
	})
	return r
}
