package define

import (
	"fmt"
	"reflect"

	"github.com/qlustra/intel/capture"
	"github.com/qlustra/intel/internal/interfaces"
	"github.com/qlustra/intel/journal"
	"github.com/qlustra/intel/logic"
	"github.com/qlustra/intel/options"
	"github.com/qlustra/intel/rules"
	"github.com/qlustra/intel/rules/contract"
	"github.com/qlustra/intel/rules/extension"
	"github.com/qlustra/intel/rules/property"
)

// StructDefinitionApi is the definition and compilation contract used by
// [Struct], [NewScope], and [NewLeaf].
type StructDefinitionApi[DEC Declaration] interface {
	interfaces.Definition
	interfaces.Compilable[*capture.StructInterrogator, *rules.EncounterContext]
	AsStructDefinition() *StructDefinition[DEC]
}

// DefinitionJournal shares constructed definitions by concrete definition type.
type DefinitionJournal = journal.Journal[reflect.Type, interfaces.Definition]

// CompilationIdentity identifies a compiled body by definition and target type.
type CompilationIdentity struct {
	Definition interfaces.Definition
	Target     reflect.Type
}

// Compilation holds a shared result. While its journal entry is visiting,
// callbacks may link its declaration, but must allow it to finish populating.
type Compilation struct {
	Context    *rules.EncounterContext
	Validation interfaces.Validation
}

// CompilationJournal shares compilation results for [CompilationIdentity] values.
type CompilationJournal = journal.Journal[CompilationIdentity, *Compilation]

const DeclarationKey = "declaration"
const JournalKey = "journal"

// BodyComposition controls whether a body category is open or closed.
type BodyComposition int

const (
	// Closed properties reject unrecognized addressable properties. Closed
	// contracts and extensions instead require a match in their domain; they
	// do not require every candidate to be recognized.
	Closed BodyComposition = iota
	// Open permits unrecognized body elements.
	Open
)

// StructDefinition contains type and body rules for a struct declaration.
// Struct constructs definitions with open contract and extension composition,
// closed property composition, and breadth-first candidate ordering by default.
type StructDefinition[HOST_DECL Declaration] struct {
	interfaces.Definition
	interfaces.Compilable[*capture.StructInterrogator, *rules.EncounterContext]
	// DepthOrdering controls candidate traversal for field rules. Struct sets
	// the default to breadth-first before invoking Define.
	DepthOrdering          options.DepthOrdering
	ContractComposition    BodyComposition
	ExtensionComposition   BodyComposition
	PropertyComposition    BodyComposition
	Name                   TypeNameRules
	Type                   TypeTypeRules
	Journal                *DefinitionJournal
	Contracts              []*StructContractDefinition
	Extensions             []*StructExtensionDefinition
	Properties             []*StructPropertyDefinition
	CarrierDefinition      *StructPropertyDefinition
	ignoredDescendantTypes []reflect.Type
	// Body rules
}

// IgnoreDescendantsOf registers an exact struct type as a field ancestry
// boundary for this definition. The boundary field itself is still processed.
// Configure definitions during Define, before compilation.
func (d *StructDefinition[HOST_DECL]) IgnoreDescendantsOf[T any]() *StructDefinition[HOST_DECL] {
	d.ignoredDescendantTypes = append(d.ignoredDescendantTypes, reflect.TypeFor[T]())
	return d
}

// Struct constructs a definition, wires its hooks, and defines its rules.
// The supplied journal reuses definitions of the same concrete type; nil starts
// an independent journal. See [StructDefinition] for defaults.
// Child-construction failures are recorded on the shared journal and returned
// after Define. A failed journal cannot be reused for another construction.
func Struct[
	SD any,
	DEC Declaration,
	P interface {
		*SD
		StructDefinitionApi[DEC]
	},
](definitionJournal *DefinitionJournal) (*SD, error) {
	definitionIdentity := reflect.TypeFor[SD]()
	if definitionJournal == nil {
		definitionJournal = journal.New[reflect.Type, interfaces.Definition]()
	}
	if err := definitionJournal.Err(); err != nil {
		return nil, err
	}
	fail := func(err error) (*SD, error) {
		definitionJournal.Fail(err)
		return nil, err
	}
	check := capture.Type[DEC]().IsStructKind()
	if !check.OK() {
		return fail(fmt.Errorf("construct %v: declaration %v: %s", definitionIdentity, reflect.TypeFor[DEC](), check.Why()))
	}
	if existing, exists := definitionJournal.Existing(definitionIdentity); exists {
		// Entries are registered under their concrete definition type below.
		typed, ok := existing.(P)
		if !ok || typed == nil {
			return fail(fmt.Errorf("construct %v: journal contains %T instead of a non-nil %v", definitionIdentity, existing, reflect.TypeFor[P]()))
		}
		return (*SD)(typed), nil
	}

	instance := new(SD)
	implementation := P(instance)
	base := implementation.AsStructDefinition()
	if base == nil {
		return fail(fmt.Errorf("construct %v: AsStructDefinition returned nil", definitionIdentity))
	}

	base.Journal = definitionJournal

	base.Definition = implementation
	base.Compilable = implementation

	base.DepthOrdering = options.BreadthFirst

	base.ContractComposition = Open
	base.ExtensionComposition = Open
	base.PropertyComposition = Closed

	base.Contracts = make([]*StructContractDefinition, 0)
	base.Extensions = make([]*StructExtensionDefinition, 0)
	base.Properties = make([]*StructPropertyDefinition, 0)

	base.Journal.Visit(definitionIdentity, implementation)
	implementation.Define()
	if err := definitionJournal.Err(); err != nil {
		return nil, err
	}
	if check := definitionJournal.Finish(definitionIdentity); !check.OK() {
		return fail(fmt.Errorf("construct %v: %s", definitionIdentity, check.Why()))
	}

	return instance, nil
}

func (d *StructDefinition[HOST_DECL]) AsStructDefinition() *StructDefinition[HOST_DECL] {
	return d
}

// SetContractComposition sets the contract domain's composition; see [Closed].
func (d *StructDefinition[HOST_DECL]) SetContractComposition(composition BodyComposition) {
	d.ContractComposition = composition
}

// SetExtensionComposition sets the extension domain's composition; see [Closed].
func (d *StructDefinition[HOST_DECL]) SetExtensionComposition(composition BodyComposition) {
	d.ExtensionComposition = composition
}

// SetPropertyComposition sets whether unrecognized properties are accepted.
func (d *StructDefinition[HOST_DECL]) SetPropertyComposition(composition BodyComposition) {
	d.PropertyComposition = composition
}

// Contract appends a contract definition. Each definition scans candidates in
// configured depth order and stops at its first recognized candidate.
func (d *StructDefinition[HOST_DECL]) Contract(optionality options.Optionality, match contract.RuleAPI, rules ...contract.RuleAPI) *StructContractDefinition {
	def := &StructContractDefinition{}
	def.Match(optionality, match)
	def.Expect(rules...)
	if d.Contracts == nil {
		d.Contracts = make([]*StructContractDefinition, 0)
	}
	d.Contracts = append(d.Contracts, def)
	return def
}

// Extension appends an extension definition. Each definition scans candidates
// in configured depth order and stops at its first recognized candidate.
func (d *StructDefinition[HOST_DECL]) Extension(optionality options.Optionality, match extension.RuleAPI, rules ...extension.RuleAPI) *StructExtensionDefinition {
	def := &StructExtensionDefinition{}
	def.Match(optionality, match)
	def.Expect(rules...)
	if d.Extensions == nil {
		d.Extensions = make([]*StructExtensionDefinition, 0)
	}
	d.Extensions = append(d.Extensions, def)
	return def
}

// Property appends a property definition. Every addressable property is tested
// against every definition in registration order; matching does not stop at
// the first definition.
func (d *StructDefinition[HOST_DECL]) Property(optionality options.Optionality, recognizer property.RuleAPI, rules ...property.RuleAPI) *StructPropertyDefinition {
	def := &StructPropertyDefinition{}
	def.Match(optionality, recognizer)
	def.Expect(rules...)
	if d.Properties == nil {
		d.Properties = make([]*StructPropertyDefinition, 0)
	}
	d.Properties = append(d.Properties, def)
	return def
}

//func (d *StructDefinition[HOST_DECL]) Carrier(recognizer property.RuleAPI, rules ...property.RuleAPI) *StructPropertyDefinition {
//	def := &StructPropertyDefinition{}
//	def.Match(options.Required, recognizer)
//	def.Expect(rules...)
//	d.CarrierDefinition = def
//	return def
//}

// Compile routes an invocation to the shared body declaration for this
// definition and target type, declaring it on a cache miss. The parent context
// supplies the shared compilation journal; nil starts an independent session.
func (d *StructDefinition[HOST_DECL]) Compile(interrogator *capture.StructInterrogator, parentContext *rules.EncounterContext) (*rules.EncounterContext, interfaces.Validation) {
	compilationJournal, err := compilationJournalFor(parentContext)
	if err != nil {
		return nil, rules.NewValidation(interrogator.String()).Nope(err.Error())
	}

	identity := CompilationIdentity{Definition: d.Definition, Target: interrogator.CapturedType()}
	compilation, exists := compilationJournal.Existing(identity)
	if exists && (compilation == nil || compilation.Context == nil || compilation.Validation == nil) {
		validation := rules.NewValidation(interrogator.String())
		validation.Nope("compilation journal contains an incomplete result record")
		compilationJournal.Fail(fmt.Errorf("compile %v: %s", identity.Target, validation.Reason()))
		return nil, validation
	}
	if !exists {
		compilation = d.Declare(interrogator, compilationJournal)
	}

	// The incoming field's carrier belongs to the caller, not this shared body.
	return compilation.Context, compilation.Validation
}

func compilationJournalFor(parentContext *rules.EncounterContext) (*CompilationJournal, error) {
	if parentContext == nil {
		return journal.New[CompilationIdentity, *Compilation](), nil
	}
	compilationJournal, check := parentContext.Access[CompilationJournal](JournalKey)
	if !check.OK() {
		return nil, fmt.Errorf("%s", check.Why())
	}
	if err := compilationJournal.Err(); err != nil {
		return nil, err
	}
	return compilationJournal, nil
}

type StructContractDefinition struct {
	contract.Rules
}

type StructExtensionDefinition struct {
	extension.Rules
}

type StructPropertyDefinition struct {
	property.Rules
}

// Scope constructs a child definition and registers an optional property rule
// that compiles it. Its default target is selected by the property's
// [ScopeIndicator.Enter] behavior.
func (d *StructDefinition[HOST_DECL]) Scope[
	DECLARATION Declaration,
	DEF any,
	P interface {
		*DEF
		StructDefinitionApi[DECLARATION]
	},
](
	recognizer property.RuleAPI,
	validators ...property.RuleAPI,
) *ScopeDefinition[HOST_DECL, DECLARATION, DEF, P] {
	scope, err := NewScope[
		HOST_DECL,
		DECLARATION,
		DEF,
		P,
	](d.Journal)
	if err != nil {
		// Keep fluent configuration safe; the factory returns this failure after
		// Define, and this incomplete scope is never registered as a property.
		d.Journal.Fail(err)
		return &ScopeDefinition[HOST_DECL, DECLARATION, DEF, P]{}
	}
	scope.Host(func(encounter property.Encounter) (*HOST_DECL, logic.Check) {
		return encounter.Access[HOST_DECL](DeclarationKey)
	})
	scope.Declared(func(encounter *rules.EncounterContext) (*DECLARATION, logic.Check) {
		return encounter.Access[DECLARATION](DeclarationKey)
	})
	handling := property.Property()
	handling.ExpectMatch("compile scope", func(_ interfaces.Validation, encounter property.Encounter) (logic.Check, bool) {
		check := scope.Handle(encounter)
		return check, check.OK()
	})
	d.Property(options.Optional, recognizer, append(validators, handling)...)
	return scope
}

// Leaf constructs a child definition and registers an optional property rule
// that compiles the property itself as a struct, without invoking Enter.
func (d *StructDefinition[HOST_DECL]) Leaf[
	DECLARATION Declaration,
	DEF any,
	P interface {
		*DEF
		StructDefinitionApi[DECLARATION]
	},
](recognizer property.RuleAPI, validators ...property.RuleAPI) *LeafDefinition[HOST_DECL, DECLARATION, DEF, P] {
	leaf, err := NewLeaf[HOST_DECL, DECLARATION, DEF, P](d.Journal)
	if err != nil {
		d.Journal.Fail(err)
		return &LeafDefinition[HOST_DECL, DECLARATION, DEF, P]{}
	}
	leaf.Host(func(encounter property.Encounter) (*HOST_DECL, logic.Check) {
		return encounter.Access[HOST_DECL](DeclarationKey)
	})
	leaf.Declared(func(encounter *rules.EncounterContext) (*DECLARATION, logic.Check) {
		return encounter.Access[DECLARATION](DeclarationKey)
	})
	handling := property.Property()
	handling.ExpectMatch("compile leaf", func(_ interfaces.Validation, encounter property.Encounter) (logic.Check, bool) {
		check := leaf.Handle(encounter)
		return check, check.OK()
	})
	d.Property(options.Optional, recognizer, append(validators, handling)...)
	return leaf
}
