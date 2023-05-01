/*
 * Copyright (c) 2021-present Sigma-Soft, Ltd.
 * @author: Nikolay Nikitin
 */

package schemas

// Qualified name
//
// <pkg>.<entity>
//
// Ref to qname.go for constants and methods
type QName struct {
	pkg    string
	entity string
}

// Schema kind enumeration
//
// Ref. schema-kind.go for constants and methods
type SchemaKind uint8

// Data kind enumeration
//
// Ref. data-kind.go for constants and methods
type DataKind uint8

// Field Verification kind
//
// Ref. verify,go for constants and methods
type VerificationKind uint8

// Numeric with OccursUnbounded value
//
// Ref. occurs.go for constants and methods
type Occurs uint16

// Application schemas
//
// Ref to cache.go for implementation
type SchemaCache interface {
	// Returns schema by name.
	//
	// If not found empty Schema with SchemeKind_null is returned
	Schema(name QName) Schema

	// Returns schema by name.
	//
	// Returns nil if not found.
	SchemaByName(name QName) Schema

	// Return count of schemas.
	SchemaCount() int

	// Enumerates all schemas from cache.
	Schemas(func(Schema))
}

// Application schemas builder
//
// Ref to cache.go for implementation
type SchemaCacheBuilder interface {
	SchemaCache

	// Adds new schema specified name and kind.
	//
	// # Panics:
	//   - if name is empty (schemas.NullQName),
	//   - if name is invalid,
	//   - if schema with name already exists.
	Add(name QName, kind SchemaKind) SchemaBuilder

	// Adds new schemas for view.
	AddView(QName) ViewBuilder

	// Must be called after all schemas added. Validates schemas and returns builded schemas or error
	Build() (SchemaCache, error)

	// Has changes since last success build
	HasChanges() bool
}

// Schema describes the entity, such as document, record or view. Schema has fields and containers.
//
// Ref to schema.go for implementation
type Schema interface {
	// Parent cache
	Cache() SchemaCache

	// Schema qualified name.
	QName() QName

	// Schema kind.
	Kind() SchemaKind

	// Finds field by name.
	//
	// Returns nil if not found.
	Field(name string) Field

	// Returns fields count
	FieldCount() int

	// Enumerates all fields in add order.
	Fields(func(Field))

	// Finds container by name.
	//
	// Returns nil if not found.
	Container(name string) Container

	// Returns containers count
	ContainerCount() int

	// Enumerates all containers in add order.
	Containers(func(Container))

	// Finds container schema by constinaer name.
	//
	// Returns nil if not found.
	ContainerSchema(name string) Schema

	// Returns is schema CDoc singleton
	Singleton() bool

	// Validates schema entities
	Validate() error
}

// Schema builder
//
// Ref to schema.go for implementation
type SchemaBuilder interface {
	Schema

	// Adds field specified name and kind.
	//
	// # Panics:
	//   - if name is empty,
	//   - if name is invalid,
	//   - if field with name is already exists,
	//   - if schema kind not supports fields,
	//   - if data kind is not allowed by schema kind.
	AddField(name string, kind DataKind, required bool) SchemaBuilder

	// Adds verified field specified name and kind.
	//
	// # Panics:
	//   - if field name is empty,
	//   - if field name is invalid,
	//   - if field with name is already exists,
	//   - if schema kind not supports fields,
	//   - if data kind is not allowed by schema kind,
	//   - if no verification kinds are specified
	AddVerifiedField(name string, kind DataKind, required bool, vk ...VerificationKind) SchemaBuilder

	// Adds container specified name and occurs.
	//
	// # Panics:
	//   - if name is empty,
	//   - if name is invalid,
	//   - if container with name already exists,
	//   - if invalid occurrences,
	//   - if schema kind does not allow containers,
	//   - if container schema kind is not compatable with schema kind.
	AddContainer(name string, schema QName, min, max Occurs) SchemaBuilder

	// Sets the singleton document flag for CDoc schemas.
	//
	// # Panics:
	//   - if not CDoc schema.
	SetSingleton()

	clear()
}

// View builder
//
// Ref to view.go for implementation
type ViewBuilder interface {
	// Returns view name
	Name() QName

	// Schema returns view schema
	Schema() SchemaBuilder

	// PartKeySchema: returns view partition key schema
	PartKeySchema() SchemaBuilder

	// ClustColsSchema returns view clustering columns schema
	ClustColsSchema() SchemaBuilder

	// FullKeySchema returns view full key (partition key + clustering columns) schema
	FullKeySchema() SchemaBuilder

	// ValueSchema returns view value schema
	ValueSchema() SchemaBuilder

	// AddPartField adds specisified field to view partition key schema. Fields is always required
	AddPartField(name string, kind DataKind) ViewBuilder

	// AddClustColumn adds specisified field to view clustering columns schema. Fields is optional
	AddClustColumn(name string, kind DataKind) ViewBuilder

	// AddValueField adds specisified field to view value schema
	AddValueField(name string, kind DataKind, required bool) ViewBuilder
}

// Describe single field.
//
// Ref to field.go for constants and implementation
type Field interface {
	// Returns field name
	Name() string

	// Returns data kind for field
	DataKind() DataKind

	// Returns is field required
	Required() bool

	// Returns is field verifable
	Verifiable() bool

	// Returns is field verifable by specified verification kind
	VerificationKind(VerificationKind) bool

	// Returns is field has fixed width data kind
	IsFixedWidth() bool

	// Returns is field system
	IsSys() bool
}

// Describes single inclusion of child schema in parent schema.
//
// Ref to container.go for constants and implementation
type Container interface {
	// Returns name of container
	Name() string

	// Returns schema name of container
	Schema() QName

	// Returns minimum occurs
	MinOccurs() Occurs

	// Returns maximum occurs
	MaxOccurs() Occurs

	// Returns is container system
	IsSys() bool
}