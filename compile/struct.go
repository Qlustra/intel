// Package compile provides typed entry points for compiling definitions
// against reflected Go types. Struct compilation retains its encounter
// context and validation for the lifetime of the returned compiler; see
// [StructCompiler] and [Struct].
package compile

import (
	"github.com/qlustra/intel/capture"
	"github.com/qlustra/intel/define"
	"github.com/qlustra/intel/internal/interfaces"
	"github.com/qlustra/intel/logic"
	"github.com/qlustra/intel/rules"
)

// StructCompiler is the result of compiling a struct definition for a target
// struct. Its validation (including its trace) and encounter context describe
// that compilation; a failed compilation remains failed, and does not expose
// a declared value. The declaration type D keeps successful results typed.
// Scope and Leaf definitions participate through their struct definition
// APIs; see [define.StructDefinition.Scope] and [define.StructDefinition.Leaf]
// for how member uses select bodies.
type StructCompiler[D define.Declaration] struct {
	definition   *define.StructDefinition[D]
	targetStruct *capture.StructInterrogator
	context      *rules.EncounterContext
	validation   interfaces.Validation
}

// StructReflection compiles source against the supplied reflected type.
// A non-struct target produces a failed compiler; otherwise the definition is
// compiled once and its context and validation are retained in the result.
func StructReflection[D define.Declaration](source define.StructDefinitionApi[D], typeInterrogator *capture.TypeInterrogator) *StructCompiler[D] {
	structDefinition := source.AsStructDefinition()
	targetStruct, isStruct := typeInterrogator.AsStruct()

	compiler := &StructCompiler[D]{
		definition:   structDefinition,
		targetStruct: targetStruct,
	}

	if !isStruct.OK() {
		compiler.validation = rules.Nope(isStruct.Reason)
		return compiler
	}

	context, validation := source.Compile(targetStruct, nil)
	compiler.context = context
	compiler.validation = validation

	if !compiler.validation.IsOkay() {
		return compiler
	}

	return compiler
}

// Struct compiles source for T, using T's reflected type. The returned
// [StructCompiler] carries the typed declaration result and compilation
// validation.
func Struct[T any, D define.Declaration](source define.StructDefinitionApi[D]) *StructCompiler[D] {
	return StructReflection(source, capture.Type[T]())
}

// IsOkay reports whether compilation succeeded.
func (c *StructCompiler[D]) IsOkay() bool {
	return c.validation.IsOkay()
}

// Trace returns the validation trace for this compilation. It describes the
// compilation outcome and does not alter the retained result.
func (c *StructCompiler[D]) Trace() string {
	return c.validation.Trace()
}

// GetDeclaredValue accesses the declaration stored in the successful
// compilation context. It returns a failed check if compilation failed or
// the context does not contain a non-nil declaration of type *D.
func (c *StructCompiler[D]) GetDeclaredValue() (*D, logic.Check) {
	if !c.IsOkay() {
		return nil, logic.Nope("compilation_failed", c.validation.Reason())
	}
	if c.context == nil {
		return nil, logic.Nope("compilation_context_missing", "compiler has no encounter context")
	}
	return c.context.Access[D](define.DeclarationKey)
}
