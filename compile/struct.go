package compile

import (
	"github.com/qlustra/intel/capture"
	"github.com/qlustra/intel/define"
	"github.com/qlustra/intel/internal/interfaces"
	"github.com/qlustra/intel/logic"
	"github.com/qlustra/intel/rules"
)

type StructCompiler[D define.Declaration] struct {
	definition   *define.StructDefinition[D]
	targetStruct *capture.StructInterrogator
	context      *rules.EncounterContext
	validation   interfaces.Validation
}

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

func Struct[T any, D define.Declaration](source define.StructDefinitionApi[D]) *StructCompiler[D] {
	return StructReflection(source, capture.Type[T]())
}

func (c *StructCompiler[D]) IsOkay() bool {
	return c.validation.IsOkay()
}

func (c *StructCompiler[D]) Trace() string {
	return c.validation.Trace()
}

func (c *StructCompiler[D]) GetDeclaredValue() (*D, logic.Check) {
	if !c.IsOkay() {
		return nil, logic.Nope("compilation_failed", c.validation.Reason())
	}
	if c.context == nil {
		return nil, logic.Nope("compilation_context_missing", "compiler has no encounter context")
	}
	return c.context.Access[D](define.DeclarationKey)
}
