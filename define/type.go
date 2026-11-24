package define

import (
	"github.com/qlustra/intel/capture"
	"github.com/qlustra/intel/internal/interfaces"
	"github.com/qlustra/intel/rules"
	"github.com/qlustra/intel/rules/identifier"
	"github.com/qlustra/intel/rules/typerule"
)

// TypeDefinition combines name and type rules for a type definition.
type TypeDefinition struct {
	interfaces.Definition
	interfaces.Compilable[*capture.TypeInterrogator, *rules.EncounterContext]
	Name TypeNameRules
	Type TypeTypeRules
}

// Type constructs an empty type definition.
func Type() *TypeDefinition {
	return &TypeDefinition{}
}

type TypeNameRules struct {
	identifier.Rules
}

type TypeTypeRules struct {
	typerule.Rules
}
