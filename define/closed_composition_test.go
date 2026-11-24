package define

import (
	"testing"

	"github.com/qlustra/intel/capture"
	"github.com/qlustra/intel/options"
	"github.com/qlustra/intel/rules/contract"
	"github.com/qlustra/intel/rules/extension"
)

type closedKnownContract interface{ Known() }
type closedUnknownContract interface{ Unknown() }
type closedCompositionDecl struct{}

func (closedCompositionDecl) declaration() {}

type closedKnownExtension struct{}
type closedUnknownExtension struct{}
type closedMixedContracts struct {
	closedUnknownContract
	closedKnownContract
}
type closedMixedExtensions struct {
	closedUnknownExtension
	closedKnownExtension
}

func TestClosedCompositionRequiresDomainMatch(t *testing.T) {
	for _, domain := range []string{"contract", "extension"} {
		t.Run(domain, func(t *testing.T) {
			for _, mixed := range []bool{false, true} {
				name := "empty"
				if mixed {
					name = "mixed"
				}
				t.Run(name, func(t *testing.T) {
					d := &StructDefinition[closedCompositionDecl]{PropertyComposition: Open, ContractComposition: Open, ExtensionComposition: Open}
					var raw *capture.StructInterrogator
					if domain == "contract" {
						d.ContractComposition = Closed
						rule := contract.Contract()
						rule.Expect("known", func(e contract.Encounter) bool { return e.Contract.FieldName() == "closedKnownContract" })
						d.Contract(options.Optional, rule)
						if mixed {
							raw, _ = capture.Type[closedMixedContracts]().AsStruct()
						} else {
							raw, _ = capture.Type[struct{}]().AsStruct()
						}
					} else {
						d.ExtensionComposition = Closed
						rule := extension.Extension()
						rule.Expect("known", func(e extension.Encounter) bool { return e.Extension.FieldName() == "closedKnownExtension" })
						d.Extension(options.Optional, rule)
						if mixed {
							raw, _ = capture.Type[closedMixedExtensions]().AsStruct()
						} else {
							raw, _ = capture.Type[struct{}]().AsStruct()
						}
					}
					result := d.Declare(raw, nil)
					if got := result.Validation.IsOkay(); got != mixed {
						t.Fatalf("success = %v, want %v: %s", got, mixed, result.Validation.Reason())
					}
				})
			}
		})
	}
}
