package interfaces

type Compilable[T any, C any] interface {
	Compile(interrogator T, parentContext C) (C, Validation)
}
