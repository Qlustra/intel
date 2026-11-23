package extension

import (
	"github.com/qlustra/intel/capture"
	"github.com/qlustra/intel/rules"
)

type Encounter struct {
	*rules.EncounterContext
	Struct    *capture.StructInterrogator
	Extension *capture.ExtensionInterrogator
}
