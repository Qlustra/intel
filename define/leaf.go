package define

import (
	"fmt"

	"github.com/qlustra/intel/capture"
	"github.com/qlustra/intel/logic"
	"github.com/qlustra/intel/rules"
	"github.com/qlustra/intel/rules/property"
)

const (
	LeafHandlingFailedCode    logic.CheckCode = "leaf_handling_failed"
	LeafTargetMissingCode     logic.CheckCode = "leaf_target_missing"
	LeafCompilationFailedCode logic.CheckCode = "leaf_compilation_failed"
)

type LeafEncounterHandler[HOST Declaration, DECL Declaration] func(property.Encounter, *HOST, *Carrier[DECL]) logic.Check

type LeafDefinition[HOST Declaration, DECL Declaration, DEF any, P interface {
	*DEF
	StructDefinitionApi[DECL]
}] struct {
	Definition     StructDefinitionApi[DECL]
	lookupHost     ScopeHostProvider[HOST]
	resolveTarget  ScopeTargetProvider
	lookupDeclared ScopeDeclaredProvider[DECL]
	interceptor    LeafEncounterHandler[HOST, DECL]
}

// NewLeaf constructs a child definition and wraps it in a leaf handler.
// Definition reuse and construction errors follow [Struct].
func NewLeaf[HOST Declaration, DECL Declaration, DEF any, P interface {
	*DEF
	StructDefinitionApi[DECL]
}](journal *DefinitionJournal) (*LeafDefinition[HOST, DECL, DEF, P], error) {
	definition, err := Struct[DEF, DECL, P](journal)
	if err != nil {
		return nil, err
	}
	return &LeafDefinition[HOST, DECL, DEF, P]{Definition: P(definition)}, nil
}

// Host configures lookup of the host declaration used by the encounter handler.
func (c *LeafDefinition[HOST, DECL, DEF, P]) Host(provider ScopeHostProvider[HOST]) *LeafDefinition[HOST, DECL, DEF, P] {
	c.lookupHost = provider
	return c
}

// LookupHost invokes the configured host provider, or fails if none is set.
func (c *LeafDefinition[HOST, DECL, DEF, P]) LookupHost(encounter property.Encounter) (*HOST, logic.Check) {
	if c.lookupHost != nil {
		return c.lookupHost(encounter)
	}
	return nil, logic.Nope("unable to look up host declaration", "no lookup function provided")
}

// Target selects a struct body without replacing the original property.
func (c *LeafDefinition[HOST, DECL, DEF, P]) Target(provider ScopeTargetProvider) *LeafDefinition[HOST, DECL, DEF, P] {
	c.resolveTarget = provider
	return c
}

// ResolveTarget defaults to the property itself, without invoking Enter.
func (c *LeafDefinition[HOST, DECL, DEF, P]) ResolveTarget(encounter property.Encounter) (*capture.StructInterrogator, logic.Check) {
	if c.resolveTarget != nil {
		return c.resolveTarget(encounter)
	}
	return encounter.Property.AsStruct()
}

// Declared configures lookup of the child declaration from the compiled context.
func (c *LeafDefinition[HOST, DECL, DEF, P]) Declared(provider ScopeDeclaredProvider[DECL]) *LeafDefinition[HOST, DECL, DEF, P] {
	c.lookupDeclared = provider
	return c
}

// LookupDeclared invokes the configured child provider, or fails if none is set.
func (c *LeafDefinition[HOST, DECL, DEF, P]) LookupDeclared(encounter *rules.EncounterContext) (*DECL, logic.Check) {
	if c.lookupDeclared != nil {
		return c.lookupDeclared(encounter)
	}
	return nil, logic.Nope("unable to look up child declaration", "no lookup function provided")
}

// Encounter configures the callback invoked after successful child compilation
// and declaration lookup. If unset, handling succeeds without a callback.
func (c *LeafDefinition[HOST, DECL, DEF, P]) Encounter(interceptor LeafEncounterHandler[HOST, DECL]) *LeafDefinition[HOST, DECL, DEF, P] {
	c.interceptor = interceptor
	return c
}

// Handle resolves and compiles the selected property as a struct, then
// optionally calls the configured interceptor with its carrier. Failures are
// recorded in the shared compilation journal.
func (c *LeafDefinition[HOST, DECL, DEF, P]) Handle(encounter property.Encounter) logic.Check {
	if encounter.EncounterContext == nil {
		return logic.Nope(LeafHandlingFailedCode, "leaf has no encounter context")
	}
	journal, check := encounter.Access[CompilationJournal](JournalKey)
	if !check.OK() {
		return logic.NopeBecause(LeafHandlingFailedCode, "cannot access compilation journal", check)
	}
	if err := journal.Err(); err != nil {
		return logic.Nope(LeafHandlingFailedCode, err.Error())
	}
	fail := func(stage string, cause logic.Check) logic.Check {
		check := logic.NopeBecause(LeafHandlingFailedCode, stage, cause)
		journal.Fail(fmt.Errorf("leaf %s: %s", encounter.Property.FieldName(), check.Why()))
		return check
	}
	target, check := c.ResolveTarget(encounter)
	if !check.OK() {
		return fail("cannot resolve leaf target", check)
	}
	if target == nil {
		return fail("cannot resolve leaf target", logic.Nope(LeafTargetMissingCode, "resolver returned no target"))
	}
	definitionEncounter, check := compileSelected(c.Definition, target, encounter.EncounterContext, journal, LeafCompilationFailedCode)
	if !check.OK() {
		return fail("cannot compile leaf", check)
	}
	if err := journal.Err(); err != nil {
		return logic.Nope(LeafHandlingFailedCode, err.Error())
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
	carrier := &Carrier[DECL]{Property: encounter.Property, Declaration: declaration}
	if check := c.interceptor(encounter, host, carrier); !check.OK() {
		return fail("carrier encounter failed", check)
	}
	return logic.Ok()
}
