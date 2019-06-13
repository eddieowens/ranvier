package cli

import (
	"github.com/eddieowens/ranvier/cli/cli/cmd"
	"github.com/spf13/cobra"
)

const CliKey = "Cli"

type Cli interface {
	Start() error
}

type cliImpl struct {
	CompileCmd cmd.Command `inject:"CompileCmd"`
	rootCmd    *cobra.Command
}

func (c *cliImpl) Start() error {
	c.rootCmd = &cobra.Command{
		Use:   "ranvier",
		Short: "A CLI for managing your configuration files and interacting with your Ranvier application",
	}

	c.rootCmd.AddCommand(
		c.CompileCmd.Cmd(),
	)

	return c.rootCmd.Execute()
}
