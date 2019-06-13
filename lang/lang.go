package lang

import (
	"github.com/eddieowens/ranvier/lang/compiler"
	"github.com/eddieowens/ranvier/lang/injector"
)

func NewCompiler() compiler.Compiler {
	return injector.CreateInjector().GetStructPtr(compiler.CompilerKey).(compiler.Compiler)
}
