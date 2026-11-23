package contract

import (
	"github.com/qlustra/intel/capture"
	"github.com/qlustra/intel/rules"
)

type Encounter struct {
	*rules.EncounterContext
	Struct   *capture.StructInterrogator
	Contract *capture.ContractInterrogator
}
