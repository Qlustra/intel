package define

import (
	"github.com/qlustra/intel/capture"
	"github.com/qlustra/intel/logic"
)

type Declaration interface {
	declaration()
}

type ScopeDeclaration interface {
	Declaration
	scope()
}

type LeafDeclaration interface {
	Declaration
	carrier()
}

type ScopeIndicator interface {
	scope()
	Enter(current *capture.StructInterrogator) (*capture.StructInterrogator, logic.Check)
}

type Scope struct{}

func (Scope) scope() {}

// Enter keeps ordinary scopes on the containing property, including when this
// method is promoted through an embedded marker.
func (Scope) Enter(current *capture.StructInterrogator) (*capture.StructInterrogator, logic.Check) {
	return current, logic.Ok()
}

type LeafIndicator interface {
	carrier()
}

// Leaf supplies the concrete marker for leaf indicators.
type Leaf struct{}

func (Leaf) carrier() {}
