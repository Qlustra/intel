package field

import (
	"github.com/qlustra/intel/capture"
	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/rules"
)

type Encounter struct {
	*rules.EncounterContext
	Struct *capture.StructInterrogator
	Field  *inspect.FieldDescriptor
}
