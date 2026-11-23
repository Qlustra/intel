package interfaces

type Validation interface {
	Okay() Validation
	Nope(reason string) Validation
	Note(reason string) Validation
	Branch(branch string) Validation
	IsOkay() bool
	Reason() string
	BranchName() string
	Trace() string
}
