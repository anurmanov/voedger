/*
 * Copyright (c) 2021-present Sigma-Soft, Ltd.
 * @author: Nikolay Nikitin
 */

package schemas

import (
	"errors"
	"fmt"
)

// Implements ISchema and ISchemaBuilder interfaces
type schemasCache struct {
	changes int
	schemas map[QName]*schema
}

func newSchemaCache() *schemasCache {
	cache := schemasCache{
		schemas: make(map[QName]*schema),
	}
	return &cache
}

func (cache *schemasCache) Add(name QName, kind SchemaKind) SchemaBuilder {
	if name == NullQName {
		panic(fmt.Errorf("schema name cannot be empty: %w", ErrNameMissed))
	}
	if ok, err := ValidQName(name); !ok {
		panic(fmt.Errorf("invalid schema name «%v»: %w", name, err))
	}
	if cache.SchemaByName(name) != nil {
		panic(fmt.Errorf("schema name «%s» already used: %w", name, ErrNameUniqueViolation))
	}
	schema := newSchema(cache, name, kind)
	cache.schemas[name] = schema
	cache.changed()
	return schema
}

func (cache *schemasCache) AddView(name QName) ViewBuilder {
	v := newViewBuilder(cache, name)
	cache.changed()
	return &v
}

func (cache *schemasCache) Build() (result SchemaCache, err error) {
	cache.prepare()

	validator := newValidator()
	cache.Schemas(func(schema Schema) {
		err = errors.Join(err, validator.validate(schema))
	})
	if err != nil {
		return nil, err
	}

	cache.changes = 0
	return cache, nil
}

func (cache *schemasCache) HasChanges() bool {
	return cache.changes > 0
}

func (cache *schemasCache) Schema(name QName) Schema {
	if schema := cache.SchemaByName(name); schema != nil {
		return schema
	}
	return NullSchema
}

func (cache *schemasCache) SchemaByName(name QName) Schema {
	if schema, ok := cache.schemas[name]; ok {
		return schema
	}
	return nil
}

func (cache *schemasCache) SchemaCount() int {
	return len(cache.schemas)
}

func (cache *schemasCache) Schemas(enum func(Schema)) {
	for _, schema := range cache.schemas {
		enum(schema)
	}
}

func (cache *schemasCache) changed() {
	cache.changes++
}

func (cache *schemasCache) prepare() {
	cache.Schemas(func(s Schema) {
		if s.Kind() == SchemaKind_ViewRecord {
			cache.prepareViewFullKeySchema(s)
		}
	})
}