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
	embeddings *FieldRegistry[EmbeddingInterrogator]
}

func captureField(descriptor *inspect.FieldDescriptor) *FieldInterrogator {
	interrogator := &FieldInterrogator{
		FieldDescriptor: descriptor,
	}
	return interrogator
}

func (i *FieldInterrogator) AsProperty() (*PropertyInterrogator, logic.Check) {
	check := i.FieldDescriptor.IsProperty()
	if !check.OK() {
		return nil, check
	}
	return i.properties.register(i.FieldDescriptor), check
}

func (i *FieldInterrogator) AsEmbedding() (*EmbeddingInterrogator, logic.Check) {
	check := i.FieldDescriptor.IsEmbedding()
	if !check.OK() {
		return nil, check
	}
	return i.embeddings.register(i.FieldDescriptor), check
}

func (i *FieldInterrogator) AsEmbeddedStruct() (*ExtensionInterrogator, logic.Check) {
	check := i.FieldDescriptor.IsEmbeddedStruct()
	if !check.OK() {
		return nil, check
	}
	return i.extensions.register(i.FieldDescriptor), check
}

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
