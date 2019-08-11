package compiler

import (
	"errors"
	"fmt"
	"github.com/eddieowens/ranvier/lang/domain"
)

const SchemaPackerKey = "SchemaPacker"

type SchemaPack interface {
	Schemas() map[string]domain.CompiledSchema
	Path() string
}

type SchemaPacker interface {
	AddSchema(pack SchemaPack, schema *domain.CompiledSchema) error
}

type schemaPackImpl struct {
	schemas map[string]domain.CompiledSchema
	path    string
}

func (p *schemaPackImpl) Path() string {
	return p.path
}

type schemaPackerImpl struct {
}

func (p *schemaPackerImpl) AddSchema(pack SchemaPack, schema *domain.CompiledSchema) error {
	if schema == nil {
		return errors.New("schema cannot be nil")
	}
	schemas := pack.Schemas()
	v, exists := schemas[schema.Name]
	if exists {
		return errors.New(fmt.Sprintf("cannot add schema %s from %s. already exists from %s", v.Name, schema.Path, v.Path))
	}
	schemas[schema.Name] = *schema
	return nil
}

func (p *schemaPackImpl) Schemas() map[string]domain.CompiledSchema {
	return p.schemas
}

func NewSchemaPack(path string) SchemaPack {
	return &schemaPackImpl{
		schemas: make(map[string]domain.CompiledSchema),
		path:    path,
	}
}
