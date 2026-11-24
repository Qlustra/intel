package contract

import (
	"github.com/qlustra/intel/capture"
	"github.com/qlustra/intel/rules"
)

// Encounter supplies the shared context, containing struct, and current
// embedded contract to contract rules.
type Encounter struct {
	*rules.EncounterContext
	Struct   *capture.StructInterrogator
	Contract *capture.ContractInterrogator
}
