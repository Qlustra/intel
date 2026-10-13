package value

type ValueInterrogator struct {
	val any
}

func CaptureValue(value any) *ValueInterrogator {
	interrogator := &ValueInterrogator{
		val: value,
	}
	return interrogator
}
