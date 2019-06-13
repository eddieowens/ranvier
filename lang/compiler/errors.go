package compiler

import "strings"

type PackError struct {
	errors       []error
	errorStrings []string
}

func NewPackError(errs ...error) *PackError {
	p := &PackError{}
	for _, e := range errs {
		p.AddError(e)
	}
	return p
}

func (p *PackError) Errors() []error {
	return p.errors
}

func (p *PackError) AddError(e error) {
	p.errors = append(p.errors, e)
	p.errorStrings = append(p.errorStrings, e.Error())
}

func (p *PackError) Error() string {
	return strings.Join(p.errorStrings, "\n")
}
