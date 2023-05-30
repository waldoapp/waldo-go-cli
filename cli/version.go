package cli

import (
	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo"

	"github.com/spf13/cobra"
)

func NewVersionCommand() *cobra.Command {
	options := &waldo.VersionOptions{}

	cmd := &cobra.Command{
		Use:   "version [options]",
		Short: "Display version information.",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ioStreams := lib.NewIOStreams(
				cmd.InOrStdin(),
				cmd.OutOrStdout(),
				cmd.ErrOrStderr())

			exitOnError(
				cmd,
				waldo.NewVersionAction(
					options,
					ioStreams).Perform())
		}}

	cmd.SetUsageTemplate(`
USAGE: waldo version
`)

	return cmd
}
