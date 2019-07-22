package app

import (
	"github.com/eddieowens/axon"
	"github.com/eddieowens/kaa"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const Key = "App"

var Version string

type App interface {
	kaa.Cmd
	Start() error
}

type app struct {
	Commands []axon.Instance `inject:"Commands"`
}

func (a *app) Command() *cobra.Command {
	cmd := &cobra.Command{
		Short:   "A tool for interacting with Ranvier",
		Use:     "ranvier [flags] [command]",
		Args:    cobra.MinimumNArgs(1),
		Version: Version,
	}

	for _, a := range a.Commands {
		cmd.AddCommand(a.GetStructPtr().(kaa.Cmd).Command())
	}

	return cmd
}

func (a *app) Start() error {
	kaa.OnError = func(ctx kaa.Context) int {
		color.Red("An error occurred:\n%s", ctx.Error().Error())

		return 1
	}
	return a.Command().Execute()
}
