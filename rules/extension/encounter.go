package extension

import (
	"github.com/qlustra/intel/capture"
	"github.com/qlustra/intel/rules"
)

// Encounter supplies the shared context, containing struct, and embedded
// extension to extension rules.
type Encounter struct {
	*rules.EncounterContext
	Struct    *capture.StructInterrogator
	Extension *capture.ExtensionInterrogator
}
