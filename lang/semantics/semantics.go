package semantics

import (
	"github.com/eddieowens/ranvier/lang/domain"
	"github.com/eddieowens/ranvier/lang/services"
)

const AnalyzerKey = "Analyzer"

type Analyzer interface {
	Semantics(manifest *domain.Schema) error
}

type analyzer struct {
	Validator services.Validator `inject:"Validator"`
}

func (s *analyzer) Semantics(manifest *domain.Schema) error {
	return s.Validator.Validate(manifest)
}
