package define

import (
	"fmt"

	"github.com/qlustra/intel/capture"
	"github.com/qlustra/intel/logic"
	"github.com/qlustra/intel/rules"
	"github.com/qlustra/intel/rules/property"
)

const (
	ScopeHandlingFailedCode    logic.CheckCode = "scope_handling_failed"
	ScopeTargetMissingCode     logic.CheckCode = "scope_target_missing"
	ScopeCompilationFailedCode logic.CheckCode = "scope_compilation_failed"
)

type ScopeEncounterHandler[HOST_DECL Declaration, SUB_DECL Declaration] func(property.Encounter, *HOST_DECL, *Carrier[SUB_DECL]) logic.Check

type ScopeHostProvider[HOST_DECL Declaration] func(encounter property.Encounter) (*HOST_DECL, logic.Check)

type ScopeTargetProvider func(property.Encounter) (*capture.StructInterrogator, logic.Check)

type ScopeDeclaredProvider[CHILD_DECL Declaration] func(encounter *rules.EncounterContext) (*CHILD_DECL, logic.Check)

type ScopeDefinition[
	HOST_DECL Declaration,
	CHILD_DECL Declaration,
	CHILD_DEF any,
	P interface {
		*CHILD_DEF
		StructDefinitionApi[CHILD_DECL]
	},
] struct {
	Definition     StructDefinitionApi[CHILD_DECL]
	lookupHost     ScopeHostProvider[HOST_DECL]
	resolveTarget  ScopeTargetProvider
	lookupDeclared ScopeDeclaredProvider[CHILD_DECL]
	interceptor    ScopeEncounterHandler[HOST_DECL, CHILD_DECL]
}

// NewScope constructs a child definition and wraps it in a scope handler.
// Definition reuse and construction errors follow [Struct].
func NewScope[
	HOST_DECL Declaration,
	CHILD_DECL Declaration,
	CHILD_DEF any,
	P interface {
		*CHILD_DEF
		StructDefinitionApi[CHILD_DECL]
	},
](definitionJournal *DefinitionJournal) (*ScopeDefinition[HOST_DECL, CHILD_DECL, CHILD_DEF, P], error) {
	definition, err := Struct[CHILD_DEF, CHILD_DECL, P](definitionJournal)
	if err != nil {
		return nil, err
	}
	return &ScopeDefinition[HOST_DECL, CHILD_DECL, CHILD_DEF, P]{
		Definition: P(definition),
	}, nil
}

// Host configures lookup of the host declaration used by the encounter handler.
func (c *ScopeDefinition[HOST_DECL, CHILD_DECL, CHILD_DEF, P]) Host(lookupHost ScopeHostProvider[HOST_DECL]) *ScopeDefinition[HOST_DECL, CHILD_DECL, CHILD_DEF, P] {
	c.lookupHost = lookupHost
	return c
}

// LookupHost invokes the configured host provider, or fails if none is set.
func (c *ScopeDefinition[HOST_DECL, CHILD_DECL, CHILD_DEF, P]) LookupHost(encounter property.Encounter) (*HOST_DECL, logic.Check) {
	if c.lookupHost != nil {
		return c.lookupHost(encounter)
	}
	return nil, logic.Nope("unable to look up host declaration", "no lookup function provided")
}

// Target selects the body to compile without replacing the original property.
func (c *ScopeDefinition[HOST_DECL, CHILD_DECL, CHILD_DEF, P]) Target(provider ScopeTargetProvider) *ScopeDefinition[HOST_DECL, CHILD_DECL, CHILD_DEF, P] {
	c.resolveTarget = provider
	return c
}

// ResolveTarget uses the configured provider, or converts the property to a
// struct and invokes its [ScopeIndicator.Enter] method.
func (c *ScopeDefinition[HOST_DECL, CHILD_DECL, CHILD_DEF, P]) ResolveTarget(encounter property.Encounter) (*capture.StructInterrogator, logic.Check) {
	if c.resolveTarget != nil {
		return c.resolveTarget(encounter)
	}
	interrogator, check := encounter.Property.AsStruct()
	if !check.OK() {
		return nil, check
	}
	return interrogator.CreateValue().Invoke("Enter", func(scope ScopeIndicator) (*capture.StructInterrogator, logic.Check) {
		return scope.Enter(interrogator)
	})
}

// Declared configures lookup of the child declaration from the compiled context.
func (c *ScopeDefinition[HOST_DECL, CHILD_DECL, CHILD_DEF, P]) Declared(lookupDeclared ScopeDeclaredProvider[CHILD_DECL]) *ScopeDefinition[HOST_DECL, CHILD_DECL, CHILD_DEF, P] {
	c.lookupDeclared = lookupDeclared
	return c
}

// LookupDeclared invokes the configured child provider, or fails if none is set.
func (c *ScopeDefinition[HOST_DECL, CHILD_DECL, CHILD_DEF, P]) LookupDeclared(encounter *rules.EncounterContext) (*CHILD_DECL, logic.Check) {
	if c.lookupDeclared != nil {
		return c.lookupDeclared(encounter)
	}
	return nil, logic.Nope("unable to look up child declaration", "no lookup function provided")
}

// Encounter configures the callback invoked after successful child compilation
// and declaration lookup. If unset, handling succeeds without a callback.
func (c *ScopeDefinition[HOST_DECL, CHILD_DECL, CHILD_DEF, P]) Encounter(interceptor ScopeEncounterHandler[HOST_DECL, CHILD_DECL]) *ScopeDefinition[HOST_DECL, CHILD_DECL, CHILD_DEF, P] {
	c.interceptor = interceptor
	return c
}

// Handle resolves and compiles the selected body, then optionally calls the
// configured interceptor with its carrier. Failures are recorded in the shared
// compilation journal.
func (c *ScopeDefinition[HOST_DECL, CHILD_DECL, CHILD_DEF, P]) Handle(encounter property.Encounter) logic.Check {
	if encounter.EncounterContext == nil {
		return logic.Nope(ScopeHandlingFailedCode, "scope has no encounter context")
	}
	compilationJournal, check := encounter.Access[CompilationJournal](JournalKey)
	if !check.OK() {
		return logic.NopeBecause(ScopeHandlingFailedCode, "cannot access compilation journal", check)
	}
	if err := compilationJournal.Err(); err != nil {
		return logic.Nope(ScopeHandlingFailedCode, err.Error())
	}
	fail := func(stage string, cause logic.Check) logic.Check {
		check := logic.NopeBecause(ScopeHandlingFailedCode, stage, cause)
		compilationJournal.Fail(fmt.Errorf("scope %s: %s", encounter.Property.FieldName(), check.Why()))
		return check
	}

	target, check := c.ResolveTarget(encounter)
	if !check.OK() {
		return fail("cannot resolve scope target", check)
	}
	if target == nil {
		return fail("cannot resolve scope target", logic.Nope(ScopeTargetMissingCode, "resolver returned no target"))
	}

	definitionEncounter, check := compileSelected(c.Definition, target, encounter.EncounterContext, compilationJournal, ScopeCompilationFailedCode)
	if !check.OK() {
		return fail("cannot compile scope", check)
	}
	if err := compilationJournal.Err(); err != nil {
		return logic.Nope(ScopeHandlingFailedCode, err.Error())
	}
	declaration, check := c.LookupDeclared(definitionEncounter)
	if !check.OK() {
		return fail("cannot access child declaration", check)
	}
	if c.interceptor == nil {
		return logic.Ok()
	}
	host, check := c.LookupHost(encounter)
	if !check.OK() {
		return fail("cannot access host declaration", check)
	}
	carrier := &Carrier[CHILD_DECL]{
		Property:    encounter.Property,
		Declaration: declaration,
	}
	if check := c.interceptor(encounter, host, carrier); !check.OK() {
		return fail("carrier encounter failed", check)
	}
	return logic.Ok()
}
