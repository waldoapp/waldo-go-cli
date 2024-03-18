package cli

import (
	"os"
	"os/exec"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/lib/xcmd"
	"github.com/waldoapp/waldo-go-cli/waldo/data"
)

var (
	helpTemplate = `OVERVIEW: {{with .Short}}{{. | trimTrailingWhitespaces}}{{end}}
{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`

	usageTemplate = `{{if .HasAvailableSubCommands}}
USAGE: {{.CommandPath}} <subcommand>
{{$cmds := .Commands}}
SUBCOMMANDS:{{range $cmds}}
  {{rpad .Name .NamePadding}}  {{.Short}}{{end}}

Use "{{.CommandPath}} help <subcommand>" for detailed help.
{{else}}
USAGE: {{.UseLine}}
{{if .HasAvailableLocalFlags}}
OPTIONS:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}
{{end}}`
)

func NewRootCommand() *xcmd.Command {
	cmd := &xcmd.Command{
		Use:   "waldo <subcommand>",
		Short: "Work with Waldo from the command line."}

	cmd.SetHelpTemplate(helpTemplate)
	cmd.SetUsageTemplate(usageTemplate)

	cmd.AddCommand(NewAddCommand())
	cmd.AddCommand(NewAuthCommand())
	cmd.AddCommand(NewBuildCommand())
	cmd.AddCommand(NewInitCommand())
	cmd.AddCommand(NewListCommand())
	cmd.AddCommand(NewRemoveCommand())
	cmd.AddCommand(NewSyncCommand())
	cmd.AddCommand(NewTriggerCommand())
	cmd.AddCommand(NewUploadCommand())
	cmd.AddCommand(NewVersionCommand())

	return cmd
}

//-----------------------------------------------------------------------------

func exitOnError(cmd *xcmd.Command, err error) {
	if err != nil {
		ioStreams := lib.NewIOStreams(
			cmd.InOrStdin(),
			cmd.OutOrStdout(),
			cmd.ErrOrStderr())

		ioStreams.EmitError(data.CLIPrefix, err)

		if ee, ok := err.(*exec.ExitError); ok {
			os.Exit(ee.ExitCode())
		}

		os.Exit(1)
	}
}
