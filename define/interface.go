package define

import (
	"github.com/qlustra/intel/capture"
	"github.com/qlustra/intel/internal/interfaces"
	"github.com/qlustra/intel/rules"
)

type InterfaceDefinition struct {
	interfaces.Definition
	interfaces.Compilable[*capture.InterfaceInterrogator, *rules.EncounterContext]
	Name    TypeNameRules
	Methods map[string]InterfaceDefinition
}

func Interface() *InterfaceDefinition {
	return &InterfaceDefinition{}
}

type InterfaceMethodDefinition struct {
}
