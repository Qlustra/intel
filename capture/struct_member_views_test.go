package capture

import (
	"testing"

	"github.com/qlustra/intel/logic"
)

func TestStructMemberViews(t *testing.T) {
	type Counters struct{ Hits int }
	type Limits struct{ Max int }
	type Source interface{ Read() string }
	type Report struct {
		Counters
		*Limits
		Source
		Backup Source
		_      int
	}

	shape, check := Type[Report]().AsStruct()
	if !check.OK() {
		t.Fatal(check.Why())
	}
	for _, tc := range []struct {
		name                                     string
		property, embedding, extension, contract bool
	}{
		{"Counters", false, true, true, false},
		{"Hits", true, false, false, false},
		{"Limits", false, true, false, false},
		{"Max", true, false, false, false},
		{"Source", false, true, false, true},
		{"Backup", true, false, false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, field := shape.Field(tc.name)
			_, property := shape.Property(tc.name)
			_, embedding := shape.NamedEmbedding(tc.name)
			_, extension := shape.Extensions.Lookup(tc.name)
			_, contract := shape.Contracts.Lookup(tc.name)
			for _, view := range []struct {
				name  string
				check logic.Check
				want  bool
			}{
				{"field", field, true}, {"property", property, tc.property},
				{"embedding", embedding, tc.embedding}, {"extension", extension, tc.extension},
				{"contract", contract, tc.contract},
			} {
				if view.check.OK() != view.want {
					t.Errorf("%s: OK() = %t, want %t (%s)", view.name, view.check.OK(), view.want, view.check.Why())
				}
			}
		})
	}
	if _, check := shape.Properties.Declared("Hits"); check.OK() {
		t.Error("promoted property was reported as directly declared")
	}
	if _, check := shape.Properties.Declared("Backup"); !check.OK() {
		t.Errorf("direct property was not found: %s", check.Why())
	}
	if _, check := shape.Fields.Declared("_"); !check.OK() {
		t.Errorf("blank declaration was not found in the general field view: %s", check.Why())
	}
	if _, check := shape.Properties.Declared("_"); check.OK() {
		t.Error("blank declaration was reported as a property")
	}
}
