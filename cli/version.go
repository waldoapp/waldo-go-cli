package cli

import (
	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/lib/xcmd"
	"github.com/waldoapp/waldo-go-cli/waldo"
)

func NewVersionCommand() *xcmd.Command {
	options := &waldo.VersionOptions{}

	cmd := &xcmd.Command{
		Use:   "version [options]",
		Short: "Show version information.",
		Args:  xcmd.ExactArgs(0),
		Run: func(cmd *xcmd.Command, args []string) {
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
