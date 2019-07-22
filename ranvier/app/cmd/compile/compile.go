package compile

import (
	"github.com/eddieowens/kaa"
	"github.com/eddieowens/ranvier/commons/validator"
	"github.com/eddieowens/ranvier/lang"
	"github.com/eddieowens/ranvier/lang/compiler"
	"github.com/spf13/cobra"
	"os"
)

const CmdKey = "CompileCmd"

type compileCmd struct {
	Validator validator.Validator `inject:"Validator"`
}

type compilePayload struct {
	Path       string `arg:"0,optional"`
	DryRun     bool   `flag:"dry-run"`
	OutputType string `flag:"output-type"`
	OutputDir  string `flag:"output-dir" validate:"file"`
	Root       string `flag:"root" validate:"file"`
	Force      bool   `flag:"force"`
}

func (c *compileCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Short: "Compile config schemas locally to single configuration file",
		Long:  longDescription,
		Use:   "compile [flags] [schema file]",
		Run:   kaa.Handle(c.Run),
	}

	cmd.Flags().BoolP("dry-run", "d", false, dryRunUsage)
	cmd.Flags().BoolP("force", "f", false, forceUsage)
	cmd.Flags().StringP("output-type", "t", "json", outputTypeUsage)
	cmd.Flags().StringP("output-dir", "o", "", outputDirectoryUsage)
	cmd.Flags().StringP("root", "r", "", rootUsage)

	return cmd
}

func (c *compileCmd) Run(ctx kaa.Context) error {
	p := new(compilePayload)
	if err := ctx.Bind(p); err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if p.Root == "" {
		p.Root = cwd
	}

	if p.OutputDir == "" {
		p.OutputDir = cwd
	}

	if err := c.Validator.Struct(p); err != nil {
		return err
	}

	comp := lang.NewCompiler()

	compileOptions := compiler.CompileOptions{
		Type: p.OutputType,
		ParseOptions: compiler.ParseOptions{
			Root: p.Root,
		},
		OutputDirectory: p.OutputDir,
		DryRun:          p.DryRun,
	}

	if p.Path == "" {
		_, err = comp.CompileAll(p.Root, compiler.CompileAllOptions{
			CompileOptions: compileOptions,
			Force:          p.Force,
		})
	} else {
		_, err = comp.Compile(p.Path, compileOptions)
	}
	return err
}

const longDescription = `Validates and compiles a Ranvier configuration schema file into the target configuration file 
type e.g. json. If no path to a schema is provided, all schema files are validated and compiled.`

const dryRunUsage = `Validates the configuration file(s) but does not output anything.`

const outputTypeUsage = `The output file type of the compiled configuration file(s). Can be either json, toml, yaml, 
or yml.`

const outputDirectoryUsage = `The target directory for your compiled file(s). Defaults to the current working directory`

const rootUsage = `The base directory for the configuration files. This directory determines the name of the 
compiled configuration file(s) as well as where to search recursively when compiling multiple files. The configuration 
file passed into the command is relative to the root directory. Defaults to the current working directory.`

const forceUsage = `When compiling multiple configuration files, continue compiling even when one fails. Only applicable
when compiling more than one configuration file.`
