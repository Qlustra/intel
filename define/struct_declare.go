package define

import (
	"fmt"

	"github.com/qlustra/intel/capture"
	"github.com/qlustra/intel/internal/interfaces"
	"github.com/qlustra/intel/journal"
	"github.com/qlustra/intel/rules"
	"github.com/qlustra/intel/rules/contract"
	"github.com/qlustra/intel/rules/extension"
	"github.com/qlustra/intel/rules/property"
	"github.com/qlustra/intel/rules/typerule"
)

// Declare creates and validates a fresh body declaration. Compile calls it on a
// cache miss. The identity must not already be registered in the supplied session;
// passing nil starts an independent session.
func (d *StructDefinition[HOST_DECL]) Declare(interrogator *capture.StructInterrogator, compilationJournal *CompilationJournal) *Compilation {
	compilation := &Compilation{Validation: rules.NewValidation(interrogator.String())}
	validation := compilation.Validation
	identity := CompilationIdentity{Definition: d.Definition, Target: interrogator.CapturedType()}
	if compilationJournal == nil {
		compilationJournal = journal.New[CompilationIdentity, *Compilation]()
	}
	if err := compilationJournal.Err(); err != nil {
		validation.Nope(err.Error())
		return compilation
	}
	if compilationJournal.Recorded(identity) {
		validation.Nope("body declaration is already registered in this session; use Compile to reuse it")
		compilationJournal.Fail(fmt.Errorf("declare %v: %s", identity.Target, validation.Reason()))
		return compilation
	}

	context := &rules.EncounterContext{}
	context.Declare(DeclarationKey, new(HOST_DECL))
	context.Declare(JournalKey, compilationJournal)
	compilation.Context = context

	// Register before running rules: nested scopes can link to this declaration
	// while its body is still being populated.
	compilationJournal.Visit(identity, compilation)
	finished := false
	defer func() {
		if !finished {
			if validation.IsOkay() {
				validation.Nope("compilation did not complete")
			}
			compilationJournal.Fail(fmt.Errorf("compile %v: %s", identity.Target, validation.Reason()))
		}
	}()

	if len(d.ignoredDescendantTypes) > 0 {
		interrogator = interrogator.WithoutDescendantsOf(d.ignoredDescendantTypes...)
	}
	d.declareType(interrogator, context, validation)
	if !validation.IsOkay() {
		return compilation
	}
	d.declareContracts(interrogator, context, validation)
	if !validation.IsOkay() {
		return compilation
	}
	d.declareExtensions(interrogator, context, validation)
	if !validation.IsOkay() {
		return compilation
	}
	d.declareProperties(interrogator, context, validation)
	if !validation.IsOkay() {
		return compilation
	}
	if err := compilationJournal.Err(); err != nil {
		validation.Nope(err.Error())
		return compilation
	}
	if check := compilationJournal.Finish(identity); !check.OK() {
		validation.Nope(check.Why())
		return compilation
	}
	finished = true
	return compilation
}

func (d *StructDefinition[HOST_DECL]) declareType(interrogator *capture.StructInterrogator, context *rules.EncounterContext, validation interfaces.Validation) {
	satisfied, valid := d.Name.ValidateRules(validation.Branch("type name"), interrogator.TypeName())
	if !valid {
		validation.Nope("type name rule evaluation failed")
		return
	}
	if !satisfied {
		validation.Nope("type name rules were not matched")
		return
	}

	typeEncounter := typerule.Encounter{
		EncounterContext: context,
		Type:             interrogator.CapturedType(),
	}
	satisfied, valid = d.Type.ValidateRules(validation.Branch("type"), typeEncounter)
	if !valid {
		validation.Nope("struct type rule evaluation failed")
		return
	}
	if !satisfied {
		validation.Nope("struct type rules were not matched")
		return
	}
}

func (d *StructDefinition[HOST_DECL]) declareContracts(interrogator *capture.StructInterrogator, context *rules.EncounterContext, validation interfaces.Validation) {
	contractsMatched := false
	for _, def := range d.Contracts {
		contractBranch := validation.Branch("contracts")
		matched := false
		for name, contractInterrogator := range interrogator.Contracts.AllAddressable(d.DepthOrdering) {
			fieldBranch := contractBranch.Branch(name)
			encounter := contract.Encounter{
				EncounterContext: context,
				Struct:           interrogator,
				Contract:         contractInterrogator,
			}
			recognized, valid := def.MatchSubject(fieldBranch, encounter)
			if !valid {
				validation.Nope("contract recognition failed")
				return
			}
			if !recognized {
				fieldBranch.Note("unknown")
				continue
			}
			satisfied, valid := def.ValidateExpectations(fieldBranch, encounter)
			if !valid || !satisfied {
				fieldBranch.Nope("invalid")
				validation.Nope("contract expectations were not satisfied")
				return
			}
			matched = true
			contractsMatched = true
			break
		}
		if !matched && def.RequiresMatch() {
			contractBranch.Nope("required")
			validation.Nope("required contract was not matched")
			return
		}
	}
	if !contractsMatched && d.ContractComposition == Closed {
		validation.Nope("unknown contracts are not allowed")
		return
	}
}

func (d *StructDefinition[HOST_DECL]) declareExtensions(interrogator *capture.StructInterrogator, context *rules.EncounterContext, validation interfaces.Validation) {
	extensionsMatched := false
	for _, def := range d.Extensions {
		extensionBranch := validation.Branch("extensions")
		matched := false
		for name, extensionInterrogator := range interrogator.Extensions.AllAddressable(d.DepthOrdering) {
			fieldBranch := extensionBranch.Branch(name)
			encounter := extension.Encounter{
				EncounterContext: context,
				Struct:           interrogator,
				Extension:        extensionInterrogator,
			}
			recognized, valid := def.MatchSubject(fieldBranch, encounter)
			if !valid {
				validation.Nope("extension recognition failed")
				return
			}
			if !recognized {
				fieldBranch.Note("unknown")
				continue
			}
			satisfied, valid := def.ValidateExpectations(fieldBranch, encounter)
			if !valid || !satisfied {
				fieldBranch.Nope("invalid")
				validation.Nope("extension expectations were not satisfied")
				return
			}
			matched = true
			extensionsMatched = true
			break
		}
		if !matched && def.RequiresMatch() {
			extensionBranch.Nope("required")
			validation.Nope("required extension was not matched")
			return
		}
	}
	if !extensionsMatched && d.ExtensionComposition == Closed {
		validation.Nope("unknown extensions are not allowed")
		return
	}
}

func (d *StructDefinition[HOST_DECL]) declareProperties(interrogator *capture.StructInterrogator, context *rules.EncounterContext, validation interfaces.Validation) {
	// Requiredness is a whole-body obligation for each registered definition,
	// independent of whether other fields match it or another definition.
	matchedDefinitions := make(map[*StructPropertyDefinition]bool, len(d.Properties))
	for name, propertyInterrogator := range interrogator.Properties.AllAddressable(d.DepthOrdering) {
		propertyBranch := validation.Branch(name)
		matched := false
		encounter := property.Encounter{
			EncounterContext: context,
			Struct:           interrogator,
			Property:         propertyInterrogator,
		}
		for _, def := range d.Properties {
			ruleBranch := propertyBranch.Branch(name)
			recognized, valid := def.MatchSubject(ruleBranch, encounter)
			if !valid {
				propertyBranch.Nope("invalid")
				validation.Nope("property recognition failed")
				return
			}
			if !recognized {
				continue
			}
			satisfied, valid := def.ValidateExpectations(ruleBranch, encounter)
			if !valid || !satisfied {
				propertyBranch.Nope("invalid")
				validation.Nope("invalid property found")
				return
			}
			matched = true
			matchedDefinitions[def] = true
		}
		if !matched {
			propertyBranch.Note("unknown")
			if d.PropertyComposition == Closed {
				validation.Nope("unknown properties are not allowed")
				return
			}
		}
	}
	missing := false
	for index, def := range d.Properties {
		if def.RequiresMatch() && !matchedDefinitions[def] {
			missing = true
			validation.Branch("required properties").Nope(fmt.Sprintf("definition %d did not match any property", index+1))
		}
	}
	if missing {
		validation.Nope("required properties not matched")
	}
}
