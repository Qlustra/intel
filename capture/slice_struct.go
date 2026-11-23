package capture

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

func (i *StructSliceInterrogator) Element() *StructInterrogator {
	return i.element
}
