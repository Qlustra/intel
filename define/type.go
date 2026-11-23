package define

import (
	"github.com/qlustra/intel/capture"
	"github.com/qlustra/intel/internal/interfaces"
	"github.com/qlustra/intel/rules"
	"github.com/qlustra/intel/rules/identifier"
	"github.com/qlustra/intel/rules/typerule"
)

type TypeDefinition struct {
	interfaces.Definition
	interfaces.Compilable[*capture.TypeInterrogator, *rules.EncounterContext]
	Name TypeNameRules
	Type TypeTypeRules
}

func Type() *TypeDefinition {
	return &TypeDefinition{}
}

type TypeNameRules struct {
	identifier.Rules
}

type TypeTypeRules struct {
	typerule.Rules
}
