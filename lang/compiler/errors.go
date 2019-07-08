package compiler

import "strings"

type SchemaPackError struct {
	errors       []error
	errorStrings []string
}

func NewSchemaPackError(errs ...error) *SchemaPackError {
	p := &SchemaPackError{}
	for _, e := range errs {
		p.AddError(e)
	}
	return p
}

func (p *SchemaPackError) Errors() []error {
	return p.errors
}

func (p *SchemaPackError) AddError(e error) {
	p.errors = append(p.errors, e)
	p.errorStrings = append(p.errorStrings, e.Error())
}

func (p *SchemaPackError) Error() string {
	return strings.Join(p.errorStrings, "\n")
}
