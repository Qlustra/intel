package options

// DepthOrdering controls the order of visible fields, without changing which
// fields are visible or how shadowing and ambiguity are resolved.
type DepthOrdering int

const (
	// DepthFirst visits each embedded field followed by its promoted fields.
	DepthFirst DepthOrdering = iota
	// BreadthFirst visits shallower fields first, retaining declaration order
	// among fields at the same depth.
	BreadthFirst
)
