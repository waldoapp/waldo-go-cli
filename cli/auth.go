package cli

import (
	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo"

	"github.com/spf13/cobra"
)

func NewAuthCommand() *cobra.Command {
	options := &waldo.AuthOptions{}

	cmd := &cobra.Command{
		Use:   "auth [-v | --verbose] <api-token>",
		Short: "Authenticate user access to Waldo.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ioStreams := lib.NewIOStreams(
				cmd.InOrStdin(),
				cmd.OutOrStdout(),
				cmd.ErrOrStderr())

			options.APIToken = args[0]

			exitOnError(
				cmd,
				waldo.NewAuthAction(
					options,
					ioStreams).Perform())
		}}

	cmd.Flags().BoolVarP(&options.Verbose, "verbose", "v", false, "Show extra verbiage.")

	cmd.SetUsageTemplate(`
USAGE: waldo auth [-v | --verbose] <api-token>

ARGUMENTS:
  <api-token>             The API token to authenticate with.

OPTIONS:
  -v, --verbose           Show extra verbiage.
`)

	return cmd
}
