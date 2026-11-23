package rules

import (
	"fmt"
	"strings"

	"github.com/qlustra/intel/internal/interfaces"
)

type Validation struct {
	branch   string
	valid    bool
	reason   string
	children []*Validation
}

func Okay() *Validation { return &Validation{valid: true} }

func Nope(reason string) *Validation { return &Validation{valid: false, reason: reason} }

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

func (m *Validation) Trace() string {
	return m.trace(0)
}
