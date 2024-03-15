package cli

import (
	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo"

	"github.com/spf13/cobra"
)

func NewInitCommand() *cobra.Command {
	options := &waldo.InitOptions{}

	cmd := &cobra.Command{
		Use:   "init [-v | --verbose]",
		Short: "Create an empty Waldo configuration.",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
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
