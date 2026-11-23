package capture

import (
	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/logic"
)

type ExtensionInterrogator struct {
	fieldInterrogatorMember
	*inspect.FieldDescriptor
}

func captureExtension(descriptor *inspect.FieldDescriptor) *ExtensionInterrogator {
	interrogator := &ExtensionInterrogator{
		FieldDescriptor: descriptor,
	}
	return interrogator
}

func (p *ExtensionInterrogator) AsStruct() (*StructInterrogator, logic.Check) {
	return TypeInspector(p.FieldTypeInspector()).AsStruct()
}
