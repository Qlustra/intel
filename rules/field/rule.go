package field

import (
	"fmt"

	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/inspect/tag/formats"
	"github.com/qlustra/intel/logic"
)

const NotTaggedWithCode logic.CheckCode = "field_not_tagged_with_tag"
const NotTaggedAsCode logic.CheckCode = "field_not_tagged_as_format"

// Rule is a composable rule for a field encounter.
type Rule struct {
	RuleType
}

// Field creates an unconstrained field rule to which expectations may be added.
func Field() *Rule {
	return &Rule{}
}

func TaggedWith(name string) *Rule {
	r := Field()
	r.ExpectCheck(fmt.Sprintf("tagged as %s", name), func(encounter Encounter) logic.Check {
		tags := encounter.Field.Tags()
		check := tags.IsPresent(name)
		if !check.OK() {
			return logic.NopeBecause(NotTaggedWithCode, name, check)
		}
		return logic.Ok()
	})
	return r
}

// Tagged parses and validates the tag's format, then applies expectation if supplied.
func Tagged[
	F any,
	FORMAT interface {
		*F
		inspect.TagFormatLike
	},
](name string, valueExpectation string, expectation func(tag *inspect.Tag[F]) bool) *Rule {
	r := Field()
	tagExpectation := name
	if valueExpectation != "" {
		tagExpectation = fmt.Sprintf("%s=\"%s\"", name, valueExpectation)
	}
	r.ExpectCheck(fmt.Sprintf("field tagged as %s", tagExpectation), func(encounter Encounter) logic.Check {
		tags := encounter.Field.Tags()
		tag, check := tags.Parse[F, FORMAT](name)
		if !check.OK() {
			return logic.NopeBecause(NotTaggedAsCode, name, check)
		}
		if expectation != nil && !expectation(tag) {
			return logic.Nope(NotTaggedAsCode, tagExpectation)
		}
		return logic.Ok()
	})
	return r
}

func TaggedAs[
	F any,
	FORMAT interface {
		*F
		inspect.TagFormatLike
	},
](name string, expectation func(tag *inspect.Tag[F]) bool) *Rule {
	return Tagged[F, FORMAT](name, "", expectation)
}

func TaggedAsSingle(name string) *Rule {
	return Tagged[formats.Single](name, "", nil)
}

func TaggedAsSingleValue(name string, expectation string) *Rule {
	return Tagged[formats.Single](name, expectation, func(tag *inspect.Tag[formats.Single]) bool {
		return tag.Format.Raw == expectation
	})
}
