package cli

import (
	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo"

	"github.com/spf13/cobra"
)

func NewRemoveCommand() *cobra.Command {
	options := &waldo.RemoveOptions{}

	cmd := &cobra.Command{
		Use:   "remove [-v | --verbose] <recipe-name>",
		Short: "Remove a recipe.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ioStreams := lib.NewIOStreams(
				cmd.InOrStdin(),
				cmd.OutOrStdout(),
				cmd.ErrOrStderr())

			options.RecipeName = args[0]

			exitOnError(
				cmd,
				waldo.NewRemoveAction(
					options,
					ioStreams).Perform())
		}}

	cmd.Flags().BoolVarP(&options.Verbose, "verbose", "v", false, "Show extra verbiage.")

	cmd.SetUsageTemplate(`
USAGE: waldo remove [-v | --verbose] <recipe-name>

ARGUMENTS:
  <recipe-name>           The name of the recipe to remove.

OPTIONS:
  -v, --verbose           Show extra verbiage.
`)

	return cmd
}
