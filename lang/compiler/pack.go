package compiler

import (
	"errors"
	"fmt"
	"github.com/eddieowens/ranvier/lang/domain"
)

const PackerKey = "Packer"

type Pack interface {
	Schemas() map[string]domain.Schema
	Path() string
}

type Packer interface {
	AddSchema(pack Pack, schema *domain.Schema) error
}

type packImpl struct {
	schemas map[string]domain.Schema
	path    string
}

func (p *packImpl) Path() string {
	return p.path
}

type packerImpl struct {
}

func (p *packerImpl) AddSchema(pack Pack, schema *domain.Schema) error {
	if schema == nil {
		return errors.New("schema cannot be nil")
	}
	if schema.IsAbstract {
		return errors.New("cannot add abstract schema to a pack")
	}
	schemas := pack.Schemas()
	v, exists := schemas[schema.Name]
	if exists {
		return errors.New(fmt.Sprintf("cannot add schema %s from %s. already exists from %s", v.Name, schema.Path, v.Path))
	}
	schemas[schema.Name] = *schema
	return nil
}

func (p *packImpl) Schemas() map[string]domain.Schema {
	return p.schemas
}

func NewPack(path string) Pack {
	return &packImpl{
		schemas: make(map[string]domain.Schema),
		path:    path,
	}
}
