package inspect

import (
	"testing"

	"github.com/qlustra/intel/logic"
)

type rejectingTagFormat struct{ Raw string }

func (f *rejectingTagFormat) Parse(raw string)      { f.Raw = raw }
func (f *rejectingTagFormat) Validate() logic.Check { return logic.Nope("rejected", f.Raw) }

func TestTagMissingAndExplicitlyEmpty(t *testing.T) {
	tags := CaptureFieldTag(`empty:"" bad:"value"`)
	for _, tc := range []struct {
		name    string
		present bool
	}{{"missing", false}, {"empty", true}} {
		t.Run(tc.name, func(t *testing.T) {
			raw, present := tags.Lookup(tc.name)
			if raw != "" || present != tc.present {
				t.Fatalf("Lookup = %q, %v", raw, present)
			}
			if got := tags.IsPresent(tc.name); got.OK() != tc.present {
				t.Fatalf("IsPresent = %v", got)
			}
			parsed, check := tags.Parse[TagFormat](tc.name)
			if tc.present {
				if !check.OK() || parsed == nil || parsed.Name != tc.name || parsed.Format.Raw != "" {
					t.Fatalf("Parse = %v, %v", parsed, check)
				}
			} else if parsed != nil || check.Code != TagNotPresentCode {
				t.Fatalf("Parse = %v, %v", parsed, check)
			}
		})
	}
	parsed, check := tags.Parse[rejectingTagFormat]("bad")
	if parsed != nil || check.Code != InvalidTagFormatCode || check.Cause == nil || check.Cause.Code != "rejected" || check.Cause.Reason != "value" {
		t.Fatalf("invalid parse = %v, %v", parsed, check)
	}
}
