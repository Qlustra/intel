// Package rules includes validation traces that record outcomes and nested
// branches produced by rule evaluation.
package rules

import (
	"fmt"
	"strings"

	"github.com/qlustra/intel/internal/interfaces"
)

// ValidationBranch is a node in a validation trace. An ExpectMatch callback
// can pass its supplied branch, or a child from Branch, to a nested rule to
// include the nested rule's diagnostics in the same trace.
type ValidationBranch = interfaces.Validation

// Validation records the outcome, reason, and child branches of an evaluation.
// Create a root with [NewValidation]; nested rules receive branches through
// [Validation.Branch].
type Validation struct {
	branch   string
	valid    bool
	reason   string
	children []*Validation
}

// Okay creates a successful validation node.
func Okay() *Validation { return &Validation{valid: true} }

// Nope creates a failed validation node with its reason.
func Nope(reason string) *Validation { return &Validation{valid: false, reason: reason} }

// NewValidation creates a successful root node with the supplied branch name.
func NewValidation(branch string) *Validation { return &Validation{branch: branch, valid: true} }

func (m *Validation) Okay() interfaces.Validation {
	m.valid = true
	m.reason = ""
	return m
}

func (m *Validation) Nope(reason string) interfaces.Validation {
	m.valid = false
	m.reason = reason
	return m
}

func (m *Validation) Note(reason string) interfaces.Validation {
	m.reason = reason
	return m
}

// Branch appends and returns a successful child node for nested diagnostics.
func (m *Validation) Branch(branch string) interfaces.Validation {
	match := Okay()
	match.branch = branch
	if m.children == nil {
		m.children = make([]*Validation, 0)
	}
	m.children = append(m.children, match)
	return match
}

func (m *Validation) IsOkay() bool {
	return m.valid
}

func (m *Validation) BranchName() string {
	return m.branch
}

func (m *Validation) Reason() string {
	return m.reason
}

func (m *Validation) String() string {
	if m.valid {
		if m.reason != "" {
			return fmt.Sprintf("%s (%s)", m.branch, m.reason)
		}
		return fmt.Sprintf("%s", m.branch)
	}
	return fmt.Sprintf("!%s: %s", m.branch, m.reason)
}

func (m *Validation) trace(level int) string {
	indent := strings.Repeat("  ", level)
	var str strings.Builder
	str.WriteString(fmt.Sprintf("%s%s", indent, m.String()))
	for _, child := range m.children {
		str.WriteString("\n" + child.trace(level+1))
	}
	return str.String()
}

// Trace renders this node and its descendants as an indented diagnostic tree.
func (m *Validation) Trace() string {
	return m.trace(0)
}
