package capture

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/logic"
)

// --- Interrogator

type StructInterrogator struct {
	fieldInterrogatorMember
	*inspect.TypeInspector
	Fields       *StructFields
	Properties   *Properties
	Extensions   *Extensions
	Contracts    *Contracts
	Embeddings   *Embeddings
	fieldMatcher inspect.DescriptorMatcher
}

func captureStruct(inspector *inspect.TypeInspector) *StructInterrogator {
	return captureStructMatching(inspector, nil)
}

func captureStructMatching(inspector *inspect.TypeInspector, matcher inspect.DescriptorMatcher) *StructInterrogator {
	descriptorRegistry := inspect.NewDescriptorRegistry()

	fieldsRegistry := createFieldRegistry(captureField)
	propertiesRegistry := createFieldRegistry(captureProperty)
	extensionRegistry := createFieldRegistry(captureExtension)
	contractRegistry := createFieldRegistry(captureContract)
	embeddingsRegistry := createFieldRegistry(captureEmbedding)

	interrogator := &StructInterrogator{
		TypeInspector: inspector,
		fieldMatcher:  matcher,
		Fields: &StructFields{
			fieldsInspector: inspect.InspectFields(inspector, descriptorRegistry, nil, matcher),
			fields:          fieldsRegistry,
			properties:      propertiesRegistry,
			extensions:      extensionRegistry,
			contracts:       contractRegistry,
			embeddings:      embeddingsRegistry,
		},
		Properties: &Properties{
			fieldsInspector: inspect.InspectFields(inspector, descriptorRegistry, inspect.IsProperty, matcher),
			fields:          propertiesRegistry,
		},
		Extensions: &Extensions{
			fieldsInspector: inspect.InspectFields(inspector, descriptorRegistry, inspect.IsEmbeddedStruct, matcher),
			fields:          extensionRegistry,
		},
		Contracts: &Contracts{
			fieldsInspector: inspect.InspectFields(inspector, descriptorRegistry, inspect.IsEmbeddedInterface, matcher),
			fields:          contractRegistry,
		},
		Embeddings: &Embeddings{
			fieldsInspector: inspect.InspectFields(inspector, descriptorRegistry, inspect.IsEmbedding, matcher),
			fields:          embeddingsRegistry,
			extensions:      extensionRegistry,
			contracts:       contractRegistry,
		},
	}

	return interrogator
}

// WithoutDescendantsOf returns an independent view of Go-visible fields, excluding
// descendants whose paths enter any of the exact struct types. Boundary fields
// themselves remain visible. A single embedded pointer to a boundary is followed.
func (i *StructInterrogator) WithoutDescendantsOf(types ...reflect.Type) *StructInterrogator {
	rootType := i.CapturedType()
	previousMatcher := i.fieldMatcher
	boundaries := make(map[reflect.Type]struct{}, len(types))
	for _, typ := range types {
		if typ != nil && typ.Kind() == reflect.Struct {
			boundaries[typ] = struct{}{}
		}
	}
	matcher := func(field *inspect.FieldDescriptor) logic.Check {
		if previousMatcher != nil {
			if check := previousMatcher(field); !check.OK() {
				return check
			}
		}
		typ := rootType
		path := field.Index()
		for _, index := range path[:len(path)-1] {
			typ = typ.Field(index).Type
			if typ.Kind() == reflect.Pointer {
				typ = typ.Elem()
			}
			if _, excluded := boundaries[typ]; excluded {
				return logic.Nope(inspect.FieldNotMatchedCode, "field is inside an excluded type")
			}
		}
		return logic.Ok()
	}
	view := captureStructMatching(i.TypeInspector, matcher)
	view.fieldInterrogatorMember = i.fieldInterrogatorMember
	return view
}

// --- Value

// CreateValue allocates a fresh zero-valued struct for typed callback invocation.
func (i *StructInterrogator) CreateValue() *StructValueInterrogator {
	return &StructValueInterrogator{
		val: reflect.New(i.CapturedType()),
		typ: i,
	}
}

// --- Field

func (i *StructInterrogator) Field(name string) (*FieldInterrogator, logic.Check) {
	return i.Fields.Lookup(name)
}

const FieldPathInvalidCode logic.CheckCode = "field_path_invalid"

// FieldPathIsValue checks structural path compatibility: indices must resolve
// and the terminal name and type must match. It does not establish descriptor
// provenance or exact metadata/tag identity. A descriptor stores no owner
// identity, and tags cannot prove that it came from this struct. The path must
// have no intermediate pointers; the terminal field may itself be a pointer.
// This inspects types only; it does not inspect field accessibility.
func (i *StructInterrogator) FieldPathIsValue(field *inspect.FieldDescriptor) logic.Check {
	_, pointers, check := i.fieldPath(field)
	if !check.OK() {
		return check
	}
	if pointers {
		return logic.Nope(FieldPathInvalidCode, "field path crosses a pointer")
	}
	return logic.Ok()
}

// fieldPath validates structural compatibility: each index and the terminal
// name and type. A descriptor does not carry its originating struct identity;
// metadata and tags are not provenance evidence.
func (i *StructInterrogator) fieldPath(field *inspect.FieldDescriptor) ([]int, bool, logic.Check) {
	if i == nil || i.TypeInspector == nil || i.CapturedType() == nil || i.CapturedType().Kind() != reflect.Struct || field == nil {
		return nil, false, logic.Nope(FieldPathInvalidCode, "missing struct type or field descriptor")
	}
	path := field.Index()
	if len(path) == 0 {
		return nil, false, logic.Nope(FieldPathInvalidCode, "empty field path")
	}
	typ := i.CapturedType()
	pointers := false
	for n, index := range path {
		if typ.Kind() != reflect.Struct || index < 0 || index >= typ.NumField() {
			return nil, false, logic.Nope(FieldPathInvalidCode, "invalid field index path")
		}
		f := typ.Field(index)
		if n == len(path)-1 {
			if f.Name != field.FieldName() || f.Type != field.ReflectedType() {
				return nil, false, logic.Nope(FieldPathInvalidCode, "field name or type does not match path")
			}
			return path, pointers, logic.Ok()
		}
		typ = f.Type
		if typ.Kind() == reflect.Pointer {
			pointers = true
			typ = typ.Elem()
		}
	}
	return nil, false, logic.Nope(FieldPathInvalidCode, "empty field path")
}

// --- Property

func (i *StructInterrogator) Property(name string) (*PropertyInterrogator, logic.Check) {
	return i.Properties.Lookup(name)
}

// --- Embeddings

func (i *StructInterrogator) NamedEmbedding(name string) (*EmbeddingInterrogator, logic.Check) {
	return i.Embeddings.Named(name)
}

func (i *StructInterrogator) Embedding[T any]() (*EmbeddingInterrogator, logic.Check) {
	return i.Embeddings.Embeds[T]()
}

func (i *StructInterrogator) EmbedsType[T any]() logic.Check {
	_, check := i.Embedding[T]()
	return check
}

func (i *StructInterrogator) Extension[T any]() (*ExtensionInterrogator, logic.Check) {
	return i.Extensions.Embeds[T]()
}

func (i *StructInterrogator) ExtendsStruct[T any]() logic.Check {
	_, check := i.Extension[T]()
	return check
}

func (i *StructInterrogator) Contract[T any]() (*ContractInterrogator, logic.Check) {
	return i.Contracts.Embeds[T]()
}

func (i *StructInterrogator) EmbedsContract[T any]() logic.Check {
	_, check := i.Contract[T]()
	return check
}

func (i *StructInterrogator) String() string {
	return fmt.Sprintf("type %s %s {}", i.TypeName(), i.Kind().String())
}

func (i *StructInterrogator) Render() string {
	var fieldLines strings.Builder
	for _, field := range i.Fields.AllDeclared() {
		fieldLines.WriteString("  ")
		fieldLines.WriteString(field.Render())
		fieldLines.WriteString("\n")
	}
	fields := ""
	if fieldLines.Len() > 0 {
		fields = "\n" + fieldLines.String()
	}
	return fmt.Sprintf("type %s %s {%s}", i.TypeName(), i.Kind().String(), fields)
}
