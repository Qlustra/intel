package compile

import (
	"github.com/qlustra/intel/capture"
	"github.com/qlustra/intel/define"
)

type TypeCompiler struct {
	definition *define.TypeDefinition
	targetType *capture.TypeInterrogator
}

func Type[T any](definition *define.TypeDefinition) *TypeCompiler {
	compiler := &TypeCompiler{
		definition: definition,
		targetType: capture.Type[T](),
	}
	return compiler
}
