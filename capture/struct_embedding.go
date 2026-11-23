package capture

import (
	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/logic"
)

type EmbeddingInterrogator struct {
	fieldInterrogatorMember
	*inspect.FieldDescriptor
	structs   *FieldRegistry[ExtensionInterrogator]
	contracts *FieldRegistry[ContractInterrogator]
}

func captureEmbedding(descriptor *inspect.FieldDescriptor) *EmbeddingInterrogator {
	interrogator := &EmbeddingInterrogator{
		FieldDescriptor: descriptor,
	}
	return interrogator
}

func (i *EmbeddingInterrogator) AsStruct() (*ExtensionInterrogator, logic.Check) {
	descriptor := i.FieldDescriptor
	check := descriptor.IsEmbeddedStruct()
	if !check.OK() {
		return nil, check
	}
	return i.structs.register(i.FieldDescriptor), check
}

func (i *EmbeddingInterrogator) AsInterface() (*ContractInterrogator, logic.Check) {
	descriptor := i.FieldDescriptor
	check := descriptor.IsEmbeddedInterface()
	if !check.OK() {
		return nil, check
	}
	return i.contracts.register(i.FieldDescriptor), check
}
