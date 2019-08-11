package semantics

import (
	"fmt"
	"github.com/eddieowens/ranvier/commons/validator"
	"github.com/eddieowens/ranvier/lang/domain"
	"strings"
)

const AnalyzerKey = "Analyzer"

type Analyzer interface {
	Semantics(manifest *domain.ParsedSchema) error
}

type analyzer struct {
	Validator validator.Validator `inject:"Validator"`
}

func (s *analyzer) Semantics(schema *domain.ParsedSchema) error {
	err := s.Validator.Struct(schema)
	if err != nil {
		if vErrs, ok := err.(validator.ValidationErrors); ok {
			for i, v := range vErrs {
				vErrs[i].Msg = fmt.Sprintf(
					"Failed to compile %s due to field %s: %s",
					schema.Path,
					strings.ToLower(v.OriginalError.StructField()),
					v.Msg,
				)
			}
		}
		return err
	}
	return nil
}
