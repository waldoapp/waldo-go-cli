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
		RunE: func(cmd *cobra.Command, args []string) error {
			ioStreams := lib.NewIOStreams(
				cmd.InOrStdin(),
				cmd.OutOrStdout(),
				cmd.ErrOrStderr())

			return waldo.NewVersionAction(
				options,
				ioStreams).Perform()
		}}

	cmd.SetUsageTemplate(`
USAGE: waldo version
`)

	return cmd
}
