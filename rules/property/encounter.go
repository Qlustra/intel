package property

import (
	"github.com/qlustra/intel/capture"
	"github.com/qlustra/intel/rules"
)

// Encounter supplies the shared context, containing struct, and current
// property to property rules. [field.Field] adapts this property to its
// underlying field descriptor when nested field rules are evaluated.
type Encounter struct {
	*rules.EncounterContext
	Struct   *capture.StructInterrogator
	Property *capture.PropertyInterrogator
}
