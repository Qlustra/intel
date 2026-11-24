package define

import (
	"github.com/qlustra/intel/capture"
	"github.com/qlustra/intel/internal/interfaces"
	"github.com/qlustra/intel/rules"
)

// InterfaceDefinition holds rules and method definitions for an interface.
type InterfaceDefinition struct {
	interfaces.Definition
	interfaces.Compilable[*capture.InterfaceInterrogator, *rules.EncounterContext]
	Name    TypeNameRules
	Methods map[string]InterfaceDefinition
}

// Interface constructs an empty interface definition.
func Interface() *InterfaceDefinition {
	return &InterfaceDefinition{}
}

type InterfaceMethodDefinition struct {
}
