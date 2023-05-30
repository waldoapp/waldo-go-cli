package cli

import (
	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo"

	"github.com/spf13/cobra"
)

func NewListCommand() *cobra.Command {
	options := &waldo.ListOptions{}

	cmd := &cobra.Command{
		Use:     "list [options]",
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

	cmd.Flags().BoolVarP(&options.LongFormat, "long", "l", false, "Display recipes in long format.")
	cmd.Flags().BoolVarP(&options.UserInfo, "user", "u", false, "Display user-specific info for each recipe.")

	cmd.SetUsageTemplate(`
USAGE: waldo list [options]

OPTIONS:
  -l, --long   Display recipes in long format.
  -u, --user   Display user-specific info for each recipe.
`)

	return cmd
}
