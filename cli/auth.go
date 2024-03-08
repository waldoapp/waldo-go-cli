package cli

import (
	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo"

	"github.com/spf13/cobra"
)

func NewAuthCommand() *cobra.Command {
	options := &waldo.AuthOptions{}

	cmd := &cobra.Command{
		Use:   "auth [options] <user-token>",
		Short: "Authorize user to access Waldo Core API.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ioStreams := lib.NewIOStreams(
				cmd.InOrStdin(),
				cmd.OutOrStdout(),
				cmd.ErrOrStderr())

			options.UserToken = args[0]

			exitOnError(
				cmd,
				waldo.NewAuthAction(
					options,
					ioStreams).Perform())
		}}

	cmd.Flags().BoolVarP(&options.Verbose, "verbose", "v", false, "Display extra verbiage.")

	cmd.SetUsageTemplate(`
USAGE: waldo auth [options] <user-token>

ARGUMENTS:
  <user-token>             The user token to authorize with.

OPTIONS:
  -v, --verbose            Display extra verbiage.
`)

	return cmd
}
