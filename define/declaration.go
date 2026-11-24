// Package define builds and compiles definitions for inspected Go types.
// Struct definitions validate type rules and body rules; property rules are
// applied to every addressable property, while contract and extension rules
// recognize at most one candidate each per rule. Composition settings have
// domain-specific meanings; see [Closed].
package define

import (
	"github.com/qlustra/intel/capture"
	"github.com/qlustra/intel/logic"
)

// Declaration is implemented by values that can be stored as declarations.
type Declaration interface {
	declaration()
}

// ScopeDeclaration marks a declaration usable as a scope declaration.
type ScopeDeclaration interface {
	Declaration
	scope()
}

// LeafDeclaration marks a declaration usable as a leaf declaration.
type LeafDeclaration interface {
	Declaration
	carrier()
}

// ScopeIndicator marks a property value as a scope and selects the struct body
// entered by a scope handler. See [ScopeDefinition.ResolveTarget].
type ScopeIndicator interface {
	scope()
	Enter(current *capture.StructInterrogator) (*capture.StructInterrogator, logic.Check)
}

// Scope is the ordinary scope marker; its Enter method keeps the containing
// struct as the target body.
type Scope struct{}

func (Scope) scope() {}

// Enter keeps ordinary scopes on the containing property, including when this
// method is promoted through an embedded marker.
func (Scope) Enter(current *capture.StructInterrogator) (*capture.StructInterrogator, logic.Check) {
	return current, logic.Ok()
}

// LeafIndicator marks a value as a leaf for leaf handlers.
type LeafIndicator interface {
	carrier()
}

// Leaf supplies the concrete marker for leaf indicators. Unlike [Scope], it
// does not implement [ScopeIndicator].
type Leaf struct{}

func (Leaf) carrier() {}
