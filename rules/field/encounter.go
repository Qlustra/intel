package field

import (
	"github.com/qlustra/intel/capture"
	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/rules"
)

// Encounter supplies the shared context, containing struct, and current field
// to field rules. Composite adapters such as property and extension rules
// construct field encounters to evaluate their nested field rules.
type Encounter struct {
	*rules.EncounterContext
	Struct *capture.StructInterrogator
	Field  *inspect.FieldDescriptor
}
