// Package capture provides checked interrogation of Go types and values.
// Start with [Type] or [TypeOf] for type metadata, and [Value] for captured
// storage. Narrow types with [TypeInterrogator.AsStruct],
// [TypeInterrogator.AsSlice], or [TypeInterrogator.AsInterface]; narrow values
// with [ValueInterrogator.AsStruct] or [ValueInterrogator.AsSlice].
//
// Struct member views describe types, while field access and invocation operate
// on captured storage. Check each returned check before using its handle. A
// successful type lookup does not guarantee that a live field is accessible.
package capture
