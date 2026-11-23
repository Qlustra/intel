package capture

import (
	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/logic"
)

type PropertyInterrogator struct {
	fieldInterrogatorMember
	*inspect.FieldDescriptor
}

func captureProperty(descriptor *inspect.FieldDescriptor) *PropertyInterrogator {
	interrogator := &PropertyInterrogator{
		FieldDescriptor: descriptor,
	}
	return interrogator
}

func (p *PropertyInterrogator) AsStruct() (*StructInterrogator, logic.Check) {
	return TypeInspector(p.FieldTypeInspector()).AsStruct()
}

func (p *PropertyInterrogator) AsSlice() (*SliceInterrogator, logic.Check) {
	return TypeInspector(p.FieldTypeInspector()).AsSlice()
}
