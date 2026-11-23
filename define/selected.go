package define

import (
	"github.com/qlustra/intel/capture"
	"github.com/qlustra/intel/logic"
	"github.com/qlustra/intel/rules"
)

// compileSelected preserves the ordinary compilation and journal path for a
// definition selected by a property handler.
func compileSelected[DECL Declaration](definition StructDefinitionApi[DECL], target *capture.StructInterrogator, parent *rules.EncounterContext, journal *CompilationJournal, failureCode logic.CheckCode) (*rules.EncounterContext, logic.Check) {
	encounter, validation := definition.Compile(target, parent)
	if validation == nil {
		return nil, logic.Nope(failureCode, "Compile returned no validation")
	}
	if !validation.IsOkay() {
		// TODO: Preserve the nested compilation cause in the returned failure.
		// Reason alone can hide the useful detail from the enclosing Legend trace.
		return nil, logic.Nope(failureCode, validation.Reason())
	}
	if encounter == nil {
		return nil, logic.Nope(failureCode, "Compile returned no encounter context")
	}
	return encounter, logic.Ok()
}
