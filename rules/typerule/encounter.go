package typerule

import (
	"reflect"

	"github.com/qlustra/intel/rules"
)

// Encounter supplies the shared context and reflected type being evaluated.
type Encounter struct {
	*rules.EncounterContext
	Type reflect.Type
}
