// Package rules composes recognition and expectations over caller-defined
// subjects. A rule distinguishes a valid non-match from failed evaluation.
// [ValidationBranch] records diagnostics and [EncounterContext] lets related
// callbacks exchange declarations.
package rules

import (
	"fmt"

	"github.com/qlustra/intel/logic"
)

const (
	DeclarationNotFoundCode     logic.CheckCode = "declaration_not_found"
	DeclarationTypeMismatchCode logic.CheckCode = "declaration_type_mismatch"
	DeclarationIsNilCode        logic.CheckCode = "declaration_is_nil"
)

// EncounterContext stores named pointers for sharing declarations among
// related encounters. Retrieve them with the typed [EncounterContext.Access].
type EncounterContext struct {
	declarations map[string]any
}

// Declare stores value under name, replacing any previous declaration there.
func (e *EncounterContext) Declare[T any](name string, value *T) {
	if e.declarations == nil {
		e.declarations = make(map[string]any)
	}
	e.declarations[name] = value
}

// Access returns the stored pointer, or a check explaining why it is unavailable.
func (e *EncounterContext) Access[T any](name string) (*T, logic.Check) {
	value, got := e.declarations[name]
	if !got {
		return nil, logic.Nope(DeclarationNotFoundCode, name)
	}
	declaration, ok := value.(*T)
	if !ok {
		return nil, logic.Nope(DeclarationTypeMismatchCode,
			fmt.Sprintf("%q: stored %T; requested %T", name, value, (*T)(nil)))
	}
	if declaration == nil {
		return nil, logic.Nope(DeclarationIsNilCode, name)
	}
	return declaration, logic.Ok()
}
