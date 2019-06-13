package cmd

import (
	"github.com/spf13/cobra"
)

const CompileCmdKey = "CompileCmd"

type compileCmd struct {
}

func (c *compileCmd) Cmd() *cobra.Command {
}

func (c *compileCmd) Run(cmd *Command, args []string) error {
}
