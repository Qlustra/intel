package capture

import (
	"fmt"

	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/logic"
)

type FieldInterrogator struct {
	fieldInterrogatorMember
	*inspect.FieldDescriptor
	properties *FieldRegistry[PropertyInterrogator]
	extensions *FieldRegistry[ExtensionInterrogator]
	contracts  *FieldRegistry[ContractInterrogator]
	embeddings *Embeddings
}

func captureField(descriptor *inspect.FieldDescriptor) *FieldInterrogator {
	interrogator := &FieldInterrogator{
		FieldDescriptor: descriptor,
	}
	return interrogator
}

// AsProperty narrows the field to a non-embedded, non-blank property.
// Use [StructInterrogator.Property] when looking up a property by name.
func (i *FieldInterrogator) AsProperty() (*PropertyInterrogator, logic.Check) {
	check := i.FieldDescriptor.IsProperty()
	if !check.OK() {
		return nil, check
	}
	return i.properties.register(i.FieldDescriptor), check
}

// AsEmbedding narrows a field declared using Go embedding syntax. Use
// [EmbeddingInterrogator.AsStruct] or [EmbeddingInterrogator.AsInterface]
// for a further check of the embedded field's type.
func (i *FieldInterrogator) AsEmbedding() (*EmbeddingInterrogator, logic.Check) {
	check := i.FieldDescriptor.IsEmbedding()
	if !check.OK() {
		return nil, check
	}
	return i.embeddings.register(i.FieldDescriptor), check
}

// AsEmbeddedStruct narrows an embedded field of struct kind to an extension.
// A pointer type is not dereferenced; see [Extensions].
func (i *FieldInterrogator) AsEmbeddedStruct() (*ExtensionInterrogator, logic.Check) {
	check := i.FieldDescriptor.IsEmbeddedStruct()
	if !check.OK() {
		return nil, check
	}
	return i.extensions.register(i.FieldDescriptor), check
}

// AsEmbeddedInterface narrows an embedded field of interface kind to a
// contract. A non-embedded interface field does not match; see [Contracts].
func (i *FieldInterrogator) AsEmbeddedInterface() (*ContractInterrogator, logic.Check) {
	check := i.FieldDescriptor.IsEmbeddedInterface()
	if !check.OK() {
		return nil, check
	}
	return i.contracts.register(i.FieldDescriptor), check
}

func (i *FieldInterrogator) String() string {
	return fmt.Sprintf("%s %s", i.FieldName(), i.FieldTypeInspector().TypeName())
}

func (i *FieldInterrogator) Render() string {
	typ := i.FieldTypeInspector()

	fieldName := i.FieldName()
	if typ.IsPointerKind() {
		fieldName = "*" + i.FieldName()
	}

	declaration := fieldName
	if !i.IsEmbedding().OK() {
		declaration += " " + typ.TypeName()
	}

	tags := i.Tags().String()
	if len(tags) > 0 {
		tags = " `" + tags + "`"
	}

	return fmt.Sprintf("%s%s", declaration, tags)
}
