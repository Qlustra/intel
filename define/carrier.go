package define

import "github.com/qlustra/intel/capture"

// Carrier represents a field's use of a declaration. Matching registrations may
// produce distinct carriers for the same property, while sharing a declaration.
type Carrier[T any] struct {
	Property    *capture.PropertyInterrogator
	Declaration *T
}
