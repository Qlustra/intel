package formats

// TagValidatorLike validates one parsed value without inspecting a declaration or
// consulting its symbol ledger.
type TagValidatorLike interface {
	Validate() error
}

type TagFormatLike interface {
	TagValidatorLike
	Parse(raw string)
}
