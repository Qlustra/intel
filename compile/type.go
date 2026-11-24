package compile

import (
	"github.com/qlustra/intel/capture"
	"github.com/qlustra/intel/define"
)

// TypeCompiler pairs a type definition with the reflected target type.
// Construct one with [Type].
type TypeCompiler struct {
	definition *define.TypeDefinition
	targetType *capture.TypeInterrogator
}

// Type constructs a compiler for definition and T's reflected type. This
// constructor records the inputs; it does not compile or return a validation
// result.
func Type[T any](definition *define.TypeDefinition) *TypeCompiler {
	compiler := &TypeCompiler{
		definition: definition,
		targetType: capture.Type[T](),
	}
	return compiler
}
