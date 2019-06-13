package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

type CobraCmd func(cmd *Command, args []string)

type Runner func(cmd *Command, args []string) error

type Command interface {
	Cmd() *cobra.Command
}

type SubCommand interface {
	Command
	Run(cmd *Command, args []string) error
}

type exitError struct {
	code int
	msg  string
}

func (e *exitError) Error() string {
	return e.msg
}

func NewExitError(code int, msg string, format ...string) error {
	return &exitError{
		code: code,
		msg:  fmt.Sprintf(msg, format),
	}
}

func HandleError(runner Runner) CobraCmd {
	return func(cmd *Command, args []string) {
		err := runner(cmd, args)
		if err != nil {
			if v, ok := err.(*exitError); ok {
				os.Exit(v.code)
			} else {
				os.Exit(1)
			}
		}
	}
}
