package semantics

import (
	"fmt"
	"github.com/eddieowens/ranvier/commons/validator"
	"github.com/eddieowens/ranvier/lang/domain"
	"strings"
)

const AnalyzerKey = "Analyzer"

type Analyzer interface {
	Semantics(manifest *domain.Schema) error
}

type analyzer struct {
	Validator validator.Validator `inject:"Validator"`
}

func (s *analyzer) Semantics(manifest *domain.Schema) error {
	err := s.Validator.Struct(manifest)
	if err != nil {
		if vErrs, ok := err.(validator.ValidationErrors); ok {
			for i, v := range vErrs {
				vErrs[i].Msg = fmt.Sprintf(
					"Failed to compile %s due to field %s: %s",
					manifest.Path,
					strings.ToLower(v.OriginalError.StructField()),
					v.Msg,
				)
			}
		}
		return err
	}
	return nil
}
