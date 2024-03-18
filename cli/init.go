package cli

import (
	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/lib/xcmd"
	"github.com/waldoapp/waldo-go-cli/waldo"
)

func NewInitCommand() *xcmd.Command {
	options := &waldo.InitOptions{}

	cmd := &xcmd.Command{
		Use:   "init [-v | --verbose]",
		Short: "Create an empty Waldo configuration.",
		Args:  xcmd.ExactArgs(0),
		Run: func(cmd *xcmd.Command, args []string) {
			ioStreams := lib.NewIOStreams(
				cmd.InOrStdin(),
				cmd.OutOrStdout(),
				cmd.ErrOrStderr())

			exitOnError(
				cmd,
				waldo.NewInitAction(
					options,
					ioStreams).Perform())
		}}

	cmd.Flags().BoolVarP(&options.Verbose, "verbose", "v", false, "Show extra verbiage.")

	cmd.SetUsageTemplate(`
USAGE: waldo init [-v | --verbose]

OPTIONS:
  -v, --verbose           Show extra verbiage.
`)

	return cmd
}
