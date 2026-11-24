package inspect

import (
	"fmt"
	"iter"
	"reflect"
	"slices"

	"github.com/qlustra/intel/logic"
	"github.com/qlustra/intel/options"
)

type FieldMatcher func(structField reflect.StructField) logic.Check
type DescriptorMatcher func(descriptor *FieldDescriptor) logic.Check

const (
	FieldLookupFailedCode logic.CheckCode = "struct_field_lookup_failed"
	FieldNotMatchedCode   logic.CheckCode = "struct_field_not_matched"
	FieldAmbiguousCode    logic.CheckCode = "struct_field_ambiguous"
	FieldRejectedCode     logic.CheckCode = "struct_field_rejected"
	FieldNotEmbedsCode    logic.CheckCode = "struct_field_embedding_not_found"
	FieldNotSatisfiesCode logic.CheckCode = "struct_field_not_satisfies"
)

type FieldsInspector struct {
	*DescriptorRegistry
	strct             reflect.Type
	fieldMatcher      FieldMatcher
	descriptorMatcher DescriptorMatcher
}

// InspectFields creates a field inspector for structType. The type must be a
// non-nil struct type, and descriptorRegistry must be non-nil. A nil matcher
// accepts all fields at that filtering stage; the registry may be shared to
// reuse descriptors across views. See [FieldsInspector.AllAddressable] and
// [FieldsInspector.LookupNamed] for the available field views.
func InspectFields(structType *TypeInspector, descriptorRegistry *DescriptorRegistry, fieldMatcher FieldMatcher, descriptorMatcher DescriptorMatcher) *FieldsInspector {
	return &FieldsInspector{
		DescriptorRegistry: descriptorRegistry,
		strct:              structType.CapturedType(),
		fieldMatcher:       fieldMatcher,
		descriptorMatcher:  descriptorMatcher,
	}
}

func (fi *FieldsInspector) matchField(field reflect.StructField) logic.Check {
	if fi.fieldMatcher == nil {
		return logic.Ok()
	}
	return fi.fieldMatcher(field)
}

func (fi *FieldsInspector) matchDescriptor(descriptor *FieldDescriptor) logic.Check {
	if descriptor == nil {
		return logic.Nope(FieldNotMatchedCode, "field descriptor is nil")
	}
	if fi.descriptorMatcher == nil {
		return logic.Ok()
	}
	return fi.descriptorMatcher(descriptor)
}

func (fi *FieldsInspector) tryRegister(field reflect.StructField) (*FieldDescriptor, logic.Check) {
	if check := fi.matchField(field); !check.OK() {
		return nil, check
	}
	descriptor := fi.register(field)
	if check := fi.matchDescriptor(descriptor); !check.OK() {
		return nil, check
	}
	return descriptor, logic.Ok()
}

// AllAddressable iterates the visible fields accepted by this inspector.
// Traversal defaults to breadth-first; an optional ordering selects depth-first
// or breadth-first order without changing Go's visible-field selection. The
// yielded paths retain all promotion steps. Fields rejected by either matcher
// are omitted.
func (fi *FieldsInspector) AllAddressable(ordering ...options.DepthOrdering) iter.Seq2[string, *FieldDescriptor] {
	order := options.BreadthFirst
	if len(ordering) > 0 {
		order = ordering[0]
	}
	return func(yield func(string, *FieldDescriptor) bool) {
		fields := reflect.VisibleFields(fi.strct)
		// VisibleFields already provides depth-first order and full index paths.
		if order != options.DepthFirst {
			slices.SortStableFunc(fields, func(a, b reflect.StructField) int {
				return len(a.Index) - len(b.Index)
			})
		}
		for _, field := range fields {
			descriptor, check := fi.tryRegister(field)
			if check.OK() && !yield(field.Name, descriptor) {
				return
			}
		}
	}
}

// AllDeclared iterates fields declared directly on the inspected struct in
// declaration order, omitting fields rejected by either matcher.
func (fi *FieldsInspector) AllDeclared() iter.Seq2[string, *FieldDescriptor] {
	return func(yield func(string, *FieldDescriptor) bool) {
		for field := range fi.strct.Fields() {
			descriptor, check := fi.tryRegister(field)
			if check.OK() && !yield(descriptor.FieldName(), descriptor) {
				return
			}
		}
	}
}

// LookupNamed resolves name using Go's field selection rules, including
// promotion and ambiguity. It returns a descriptor only when resolution and
// both matchers succeed; inspect the check before using the descriptor.
func (fi *FieldsInspector) LookupNamed(name string) (*FieldDescriptor, logic.Check) {
	if check := IsStructKind(fi.strct); !check.OK() {
		return nil, logic.NopeBecause(FieldLookupFailedCode, fmt.Sprintf("field named '%q' not matched conditions", name), check)
	}
	field, found := fi.strct.FieldByName(name)
	if !found {
		return nil, logic.Nope(FieldLookupFailedCode, fmt.Sprintf("no unambiguous field named %q in %v", name, fi.strct))
	}
	return fi.tryRegister(field)
}

// matchingField delegates resolution to reflect. An accepted candidate with no
// unique result indicates ambiguity; rejected candidates retain their checks.
func (fi *FieldsInspector) matchingField(parentCode logic.CheckCode, parentReason string, predicate FieldMatcher) (*FieldDescriptor, logic.Check) {
	if check := IsStructKind(fi.strct); !check.OK() {
		return nil, logic.NopeBecause(parentCode, "", logic.NopeBecause(FieldNotMatchedCode, parentReason, check))
	}
	var first *FieldDescriptor
	var cause *logic.Check
	var rejection *logic.Check
	_, unique := fi.strct.FieldByNameFunc(func(name string) bool {
		field, found := fi.strct.FieldByName(name)
		if !found {
			return false
		}
		if check := predicate(field); !check.OK() {
			if cause == nil {
				cause = &check
			}
			return false
		}
		descriptor, check := fi.tryRegister(field)
		if !check.OK() {
			if rejection == nil {
				rejection = &check
			}
			return false
		}
		if first == nil {
			first = descriptor
		}
		return true
	})
	if first != nil {
		if !unique {
			return nil, logic.NopeBecause(parentCode, parentReason, logic.Nope(FieldAmbiguousCode, fmt.Sprintf("fields found multiple times in %v", fi.strct)))
		}
		return first, logic.Ok()
	}
	if rejection != nil {
		cause = rejection
	}
	if cause != nil {
		return nil, logic.NopeBecause(parentCode, parentReason, logic.NopeBecause(FieldNotMatchedCode, "", *cause))
	}
	return nil, logic.NopeBecause(parentCode, parentReason, logic.Nope(FieldNotMatchedCode, ""))
}

// Embeds finds a uniquely selected visible field whose exact type is T. This
// does not require that the field be an embedding; inspect the descriptor when
// that distinction matters.
func (fi *FieldsInspector) Embeds[T any]() (*FieldDescriptor, logic.Check) {
	return fi.matchingField(FieldNotEmbedsCode, fmt.Sprintf("field embedding %v", reflect.TypeFor[T]()), func(field reflect.StructField) logic.Check {
		return Is[T](field.Type)
	})
}

// Satisfies finds a uniquely selected visible field whose type implements
// interface T. T must be an interface type; the operation checks the field
// type's method set, not a value stored in a field.
func (fi *FieldsInspector) Satisfies[T any]() (*FieldDescriptor, logic.Check) {
	if check := IsInterfaceKind(reflect.TypeFor[T]()); !check.OK() {
		return nil, logic.NopeBecause(TypeDoesNotSatisfyCode, "requested contract", check)
	}
	return fi.matchingField(FieldNotSatisfiesCode, fmt.Sprintf("field satisfying %v", reflect.TypeFor[T]()), func(field reflect.StructField) logic.Check {
		return Satisfies[T](field.Type)
	})
}

func (fi *FieldsInspector) declaredField(description string, predicate FieldMatcher) (*FieldDescriptor, logic.Check) {
	if check := IsStructKind(fi.strct); !check.OK() {
		return nil, logic.NopeBecause(FieldNotMatchedCode, description, check)
	}
	var cause *logic.Check
	var rejection *logic.Check
	for field := range fi.strct.Fields() {
		if check := predicate(field); !check.OK() {
			if cause == nil {
				cause = &check
			}
			continue
		}
		descriptor, check := fi.tryRegister(field)
		if check.OK() {
			return descriptor, check
		}
		if rejection == nil {
			rejection = &check
		}
	}
	if rejection != nil {
		cause = rejection
	}
	reason := fmt.Sprintf("%s is not declared in %v", description, fi.strct)
	if cause != nil {
		return nil, logic.NopeBecause(FieldNotMatchedCode, reason, *cause)
	}
	return nil, logic.Nope(FieldNotMatchedCode, reason)
}

// Declared finds a directly declared field by name. Unlike LookupNamed, it
// does not search promoted fields.
func (fi *FieldsInspector) Declared(name string) (*FieldDescriptor, logic.Check) {
	return fi.declaredField(fmt.Sprintf("field named %q", name), func(field reflect.StructField) logic.Check {
		if field.Name == name {
			return logic.Ok()
		}
		return logic.Nope(FieldNotMatchedCode, fmt.Sprintf("%s != %s", field.Name, name))
	})
}

// DeclaredType finds the first directly declared field whose type is exactly
// T, in declaration order.
func (fi *FieldsInspector) DeclaredType[T any]() (*FieldDescriptor, logic.Check) {
	return fi.declaredField(fmt.Sprintf("field of type %v", reflect.TypeFor[T]()), func(field reflect.StructField) logic.Check {
		return Is[T](field.Type)
	})
}

// DeclaredContract finds the first directly declared field whose type
// implements interface T, in declaration order. T must be an interface type.
func (fi *FieldsInspector) DeclaredContract[T any]() (*FieldDescriptor, logic.Check) {
	if check := IsInterfaceKind(reflect.TypeFor[T]()); !check.OK() {
		return nil, logic.NopeBecause(TypeDoesNotSatisfyCode, "requested contract", check)
	}
	return fi.declaredField(fmt.Sprintf("field satisfying %v", reflect.TypeFor[T]()), func(field reflect.StructField) logic.Check {
		return Satisfies[T](field.Type)
	})
}
