package capture

import (
	"github.com/qlustra/intel/inspect"
)

type ContractInterrogator struct {
	fieldInterrogatorMember
	*inspect.FieldDescriptor
}

func captureContract(descriptor *inspect.FieldDescriptor) *ContractInterrogator {
	interrogator := &ContractInterrogator{
		FieldDescriptor: descriptor,
	}
	return interrogator
}
