package inspect

import (
	"reflect"
	"strings"
	"testing"

	"github.com/qlustra/intel/logic"
)

type readable interface{ Read() }
type readWriter interface {
	readable
	Write()
}

type embeddedReader struct{ readable }
type namedReader struct{ Reader readable }
type competingReaders struct {
	readable
	readWriter
}

func inspectContracts[T any]() *FieldsInspector {
	return InspectFields(CaptureType(reflect.TypeFor[T]()), NewDescriptorRegistry(), IsEmbeddedInterface, nil)
}

func TestEmbeddingCheckPreservesRejectedFieldCause(t *testing.T) {
	descriptor, check := inspectContracts[namedReader]().Embeds[readable]()
	if descriptor != nil || check.OK() || check.Code != FieldNotEmbedsCode {
		t.Fatalf("unexpected lookup: descriptor=%v check=%s", descriptor, check.Why())
	}
	// The lookup operation wraps the resolution failure; matcher checks remain
	// intact beneath it rather than acquiring a synthetic FieldRejected layer.
	codes := []logic.CheckCode{FieldNotEmbedsCode, FieldNotMatchedCode, FieldIsNotEmbeddedInterfaceCode, FieldIsNotEmbeddingCode}
	for i, code := range codes {
		if check.Code != code {
			t.Fatalf("cause %d: got %q, want %q", i, check.Code, code)
		}
		if i < len(codes)-1 {
			if check.Cause == nil {
				t.Fatalf("missing cause after %q", code)
			}
			check = *check.Cause
		}
	}
	if !strings.Contains(check.Reason, "Reader") {
		t.Fatalf("missing rejected field name: %s", check.Reason)
	}
}

func TestEmbeddingLookupOutcomes(t *testing.T) {
	descriptor, check := inspectContracts[embeddedReader]().Embeds[readable]()
	if !check.OK() || descriptor == nil || descriptor.ReflectedType() != reflect.TypeFor[readable]() {
		t.Fatalf("expected unique embedding: %s", check.Why())
	}
	missing, check := inspectContracts[struct{}]().Embeds[readable]()
	if missing != nil || check.Code != FieldNotEmbedsCode || check.Cause == nil || check.Cause.Code != FieldNotMatchedCode {
		t.Fatalf("expected missing embedding: %s", check.Why())
	}
	ambiguous, check := inspectContracts[competingReaders]().Satisfies[readable]()
	if ambiguous != nil || check.Code != FieldNotSatisfiesCode || check.Cause == nil || check.Cause.Code != FieldAmbiguousCode {
		t.Fatalf("expected ambiguity: %s", check.Why())
	}
	_, check = inspectContracts[embeddedReader]().Satisfies[int]()
	if check.Code != TypeDoesNotSatisfyCode || check.Cause == nil || check.Cause.Code != TypeIsNotInterfaceKindCode {
		t.Fatalf("expected invalid target cause: %s", check.Why())
	}
}

func TestEmbeddingLookupPrefersShallowestMatch(t *testing.T) {
	type directAndPromoted struct {
		readable
		embeddedReader
	}
	descriptor, check := inspectContracts[directAndPromoted]().Embeds[readable]()
	if !check.OK() || descriptor == nil || len(descriptor.Index()) != 1 {
		t.Fatalf("expected direct embedding to win: %s", check.Why())
	}
}

func TestIsTypeReportsUnnamedTypesAndNil(t *testing.T) {
	check := Is[[]string](reflect.TypeFor[[]int]())
	if check.OK() || !strings.Contains(check.Reason, "[]int") || !strings.Contains(check.Reason, "[]string") {
		t.Fatalf("missing type details: %s", check.Why())
	}
	if check := Is[int](nil); check.OK() {
		t.Fatal("nil type unexpectedly matched")
	}
}

func TestEmbeddingLookupHandlesPromotionShadowingAndCycles(t *testing.T) {
	type promotedReader struct{ *embeddedReader }
	descriptor, check := inspectContracts[promotedReader]().Embeds[readable]()
	if !check.OK() || descriptor == nil || len(descriptor.Index()) != 2 {
		t.Fatalf("expected promoted embedding: %s", check.Why())
	}
	type shadowedReader struct {
		readable int
		embeddedReader
	}
	if descriptor, check := inspectContracts[shadowedReader]().Embeds[readable](); descriptor != nil || check.Code != FieldNotEmbedsCode || check.Cause == nil || check.Cause.Code != FieldNotMatchedCode {
		t.Fatalf("expected shadowed embedding to be unavailable: %s", check.Why())
	}
	type recursiveReader struct{ *recursiveReader }
	if descriptor, check := inspectContracts[recursiveReader]().Embeds[readable](); descriptor != nil || check.Code != FieldNotEmbedsCode || check.Cause == nil || check.Cause.Code != FieldNotMatchedCode {
		t.Fatalf("expected missing field in recursive embedding: %s", check.Why())
	}
}

func TestDescriptorMatcherRetainsItsFailure(t *testing.T) {
	fields := InspectFields(
		CaptureType(reflect.TypeFor[embeddedReader]()),
		NewDescriptorRegistry(),
		IsEmbeddedInterface,
		func(field *FieldDescriptor) logic.Check { return field.IsDirect() },
	)
	if descriptor, check := fields.Embeds[readable](); descriptor == nil || !check.OK() {
		t.Fatalf("direct descriptor rejected: %s", check.Why())
	}
	type promotedReader struct{ embeddedReader }
	fields = InspectFields(
		CaptureType(reflect.TypeFor[promotedReader]()),
		NewDescriptorRegistry(),
		IsEmbeddedInterface,
		func(field *FieldDescriptor) logic.Check { return field.IsDirect() },
	)
	if descriptor, check := fields.Embeds[readable](); descriptor != nil || check.Cause == nil || check.Cause.Cause == nil || check.Cause.Cause.Code != FieldIsNotDeclaredCode {
		t.Fatalf("lost descriptor rejection: %s", check.Why())
	}
}
