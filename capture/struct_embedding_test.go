package capture

import "testing"

func TestFieldEmbeddingNarrowing(t *testing.T) {
	type Counters struct{ Hits int }
	type Source interface{ Read() string }
	type Limits struct{ Max int }
	type Report struct {
		Counters
		Source
		*Limits
	}

	for _, tc := range []struct {
		name                      string
		wantStruct, wantInterface bool
	}{
		{"Counters", true, false},
		{"Source", false, true},
		{"Limits", false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Start with a fresh view so no direct lookup initializes the embedding.
			shape, check := Type[Report]().AsStruct()
			if !check.OK() {
				t.Fatal(check.Why())
			}
			field, check := shape.Field(tc.name)
			if !check.OK() {
				t.Fatal(check.Why())
			}
			embedding, check := field.AsEmbedding()
			if !check.OK() {
				t.Fatal(check.Why())
			}
			extension, check := embedding.AsStruct()
			if check.OK() != tc.wantStruct || (extension != nil) != tc.wantStruct {
				t.Fatalf("AsStruct: OK() = %t, non-nil = %t, want both %t (%s)", check.OK(), extension != nil, tc.wantStruct, check.Why())
			}
			contract, check := embedding.AsInterface()
			if check.OK() != tc.wantInterface || (contract != nil) != tc.wantInterface {
				t.Fatalf("AsInterface: OK() = %t, non-nil = %t, want both %t (%s)", check.OK(), contract != nil, tc.wantInterface, check.Why())
			}
			if direct, check := shape.NamedEmbedding(tc.name); !check.OK() || direct != embedding {
				t.Fatalf("NamedEmbedding did not return the same embedding (%s)", check.Why())
			}
			if tc.wantStruct {
				if direct, check := shape.Extensions.Lookup(tc.name); !check.OK() || direct != extension {
					t.Fatalf("Extensions.Lookup did not return the same extension (%s)", check.Why())
				}
			}
			if tc.wantInterface {
				if direct, check := shape.Contracts.Lookup(tc.name); !check.OK() || direct != contract {
					t.Fatalf("Contracts.Lookup did not return the same contract (%s)", check.Why())
				}
			}
		})
	}
}
