package capture

// StructSliceInterrogator describes a slice with struct-kind elements.
type StructSliceInterrogator struct {
	*SliceInterrogator
	element *StructInterrogator
}

func captureStructSlice(slice *SliceInterrogator, element *StructInterrogator) *StructSliceInterrogator {
	interrogator := &StructSliceInterrogator{
		SliceInterrogator: slice,
		element:           element,
	}
	return interrogator
}

// Element returns the struct type view of the slice's elements.
func (i *StructSliceInterrogator) Element() *StructInterrogator {
	return i.element
}
