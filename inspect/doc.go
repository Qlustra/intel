// Package inspect describes Go types, struct fields, interfaces, and field
// tags without capturing a live value. Checked lookups return a logic.Check;
// callers should inspect OK before using the associated result. Start with
// [CaptureType] and [InspectFields] for type and field metadata; use
// [FieldTagInspector] for tag lookup and parsing. [FieldsInspector.LookupNamed]
// follows Go's promotion and ambiguity rules, while
// [FieldsInspector.AllAddressable] lets callers select traversal order. The
// Intel usage manual's inspection chapter describes routes between these
// operations.
package inspect
