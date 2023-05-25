package cli

import (
	"github.com/spf13/cobra"
)

var (
	helpTemplate = `OVERVIEW: {{with .Short}}{{. | trimTrailingWhitespaces}}{{end}}
{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`

	usageTemplate = `{{if .HasAvailableSubCommands}}
USAGE: {{.CommandPath}} <command>
{{$cmds := .Commands}}
COMMANDS:{{range $cmds}}
  {{rpad .Name .NamePadding}}  {{.Short}}{{end}}

Use "{{.CommandPath}} help <command>" for more information about a command.
{{else}}
USAGE: {{.UseLine}}
{{if .HasAvailableLocalFlags}}
OPTIONS:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}
{{end}}`
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "waldo <command>",
		Short: "Work with Waldo from the command line."}

	cmd.CompletionOptions.DisableDefaultCmd = true

	cmd.SetHelpTemplate(helpTemplate)
	cmd.SetUsageTemplate(usageTemplate)

	cmd.AddCommand(fixup(NewAddCommand()))
	cmd.AddCommand(fixup(NewBuildCommand()))
	cmd.AddCommand(fixup(NewInitCommand()))
	cmd.AddCommand(fixup(NewListCommand()))
	cmd.AddCommand(fixup(NewRemoveCommand()))
	// cmd.AddCommand(fixup(NewRunCommand()))
	// cmd.AddCommand(fixup(NewSessionCommand()))
	// cmd.AddCommand(fixup(NewSyncCommand()))
	cmd.AddCommand(fixup(NewTriggerCommand()))
	cmd.AddCommand(fixup(NewUploadCommand()))
	cmd.AddCommand(fixup(NewVersionCommand()))

	return fixup(cmd)
}

//-----------------------------------------------------------------------------

func fixup(cmd *cobra.Command) *cobra.Command {
	cmd.DisableAutoGenTag = true
	cmd.DisableFlagsInUseLine = true

	return cmd
}
