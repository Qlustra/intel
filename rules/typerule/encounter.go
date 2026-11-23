package typerule

import (
	"reflect"

	"github.com/qlustra/intel/rules"
)

type Encounter struct {
	*rules.EncounterContext
	Type reflect.Type
}
