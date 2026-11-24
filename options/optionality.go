// Package options supplies choices for matching and traversal. An operation's
// documentation defines where a choice applies and which value it defaults to.
package options

// ======================================== Enums

type (
	// Optionality controls whether a rule registration requires a match.
	// Optionality does not turn failed evaluation into a valid non-match.
	Optionality bool
)

// ======================================== Optionality

const (
	// Optional permits a valid absence of a match.
	Optional Optionality = false
	// Required makes a match an obligation of the registration.
	Required Optionality = true
)
