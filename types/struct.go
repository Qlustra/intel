package types

import (
	structEngine "github.com/qlustra/intel/engines/struct"
	typesEngine "github.com/qlustra/intel/engines/types"
)

// --- Interrogator

type StructInterrogator struct {
	fieldInterrogatorMember
	*typesEngine.TypeInspector
	Fields     *StructFields
	Properties *StructProperties
	Extensions *EmbeddedStructs
	Contracts  *EmbeddedInterfaces
	Embeddings *StructEmbeddings
}

func captureStruct(inspector *typesEngine.TypeInspector) *StructInterrogator {
	descriptorRegistry := structEngine.NewDescriptorRegistry()

	fieldsRegistry := createFieldRegistry(captureField)
	propertiesRegistry := createFieldRegistry(captureProperty)
	structsRegistry := createFieldRegistry(captureExtensions)
	interfacesRegistry := createFieldRegistry(captureMembership)
	embeddingsRegistry := createFieldRegistry(captureEmbedding)

	interrogator := &StructInterrogator{
		TypeInspector: inspector,
		Fields: &StructFields{
			fieldsInspector: structEngine.InspectFields(inspector, descriptorRegistry, nil, nil),
			fields:          fieldsRegistry,
			properties:      propertiesRegistry,
			structs:         structsRegistry,
			interfaces:      interfacesRegistry,
			embeddings:      embeddingsRegistry,
		},
		Properties: &StructProperties{
			fieldsInspector: structEngine.InspectFields(inspector, descriptorRegistry, structEngine.IsProperty, nil),
			fields:          propertiesRegistry,
		},
		Extensions: &EmbeddedStructs{
			fieldsInspector: structEngine.InspectFields(inspector, descriptorRegistry, structEngine.IsEmbeddedStruct, nil),
			fields:          structsRegistry,
		},
		Contracts: &EmbeddedInterfaces{
			fieldsInspector: structEngine.InspectFields(inspector, descriptorRegistry, structEngine.IsEmbeddedInterface, nil),
			fields:          interfacesRegistry,
		},
		Embeddings: &StructEmbeddings{
			fieldsInspector: structEngine.InspectFields(inspector, descriptorRegistry, structEngine.IsEmbedding, nil),
			fields:          embeddingsRegistry,
			structs:         structsRegistry,
			interfaces:      interfacesRegistry,
		},
	}

	return interrogator
}

func captureExtensions(descriptor *structEngine.FieldDescriptor) *StructInterrogator {
	return captureStruct(descriptor.FieldTypeInspector())
}

// --- Field

func (i *StructInterrogator) Field(name string) (*FieldInterrogator, bool) {
	return i.Fields.Lookup(name)
}

// --- Property

func (i *StructInterrogator) Property(name string) (*PropertyInterrogator, bool) {
	return i.Properties.Lookup(name)
}

// --- Embeddings

func (i *StructInterrogator) NamedEmbedding(name string) (*EmbeddingInterrogator, bool) {
	return i.Embeddings.Named(name)
}

func (i *StructInterrogator) EmbeddingOfType[T any]() (*EmbeddingInterrogator, bool, bool) {
	return i.Embeddings.OfType[T]()
}

func (i *StructInterrogator) EmbedsType[T any]() bool {
	_, found, _ := i.EmbeddingOfType[T]()
	return found
}

func (i *StructInterrogator) ExtensionOfType[T any]() (*StructInterrogator, bool, bool) {
	return i.Extensions.OfType[T]()
}

func (i *StructInterrogator) ExtendsStruct[T any]() bool {
	_, found, _ := i.ExtensionOfType[T]()
	return found
}

func (i *StructInterrogator) ContractOfType[T any]() (*InterfaceInterrogator, bool, bool) {
	return i.Contracts.OfType[T]()
}

func (i *StructInterrogator) EmbedsContract[T any]() bool {
	_, found, _ := i.ContractOfType[T]()
	return found
}
