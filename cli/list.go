package cli

import (
	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo"

	"github.com/spf13/cobra"
)

func NewListCommand() *cobra.Command {
	options := &waldo.ListOptions{}

	cmd := &cobra.Command{
		Use:     "list [-l | --long] [-u | --user]",
		Short:   "List defined recipes.",
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ioStreams := lib.NewIOStreams(
				cmd.InOrStdin(),
				cmd.OutOrStdout(),
				cmd.ErrOrStderr())

			exitOnError(
				cmd,
				waldo.NewListAction(
					options,
					ioStreams).Perform())
		}}

	cmd.Flags().BoolVarP(&options.LongFormat, "long", "l", false, "Show recipes in long format.")
	cmd.Flags().BoolVarP(&options.UserInfo, "user", "u", false, "Show user-specific information for each recipe.")

	cmd.SetUsageTemplate(`
USAGE: waldo list [-l | --long] [-u | --user]

OPTIONS:
  -l, --long              Show recipes in long format.
  -u, --user              Show user-specific information for each recipe.
`)

	return cmd
}
