package cli

import (
	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/lib/xcmd"
	"github.com/waldoapp/waldo-go-cli/waldo"
)

func NewAddCommand() *xcmd.Command {
	options := &waldo.AddOptions{}

	cmd := &xcmd.Command{
		Use:   "add [-v | --verbose] <recipe-name>",
		Short: "Add a recipe describing how to build and upload a specific variant of the app.",
		Args:  xcmd.ExactArgs(1),
		Run: func(cmd *xcmd.Command, args []string) {
			ioStreams := lib.NewIOStreams(
				cmd.InOrStdin(),
				cmd.OutOrStdout(),
				cmd.ErrOrStderr())

			options.RecipeName = args[0]

			exitOnError(
				cmd,
				waldo.NewAddAction(
					options,
					ioStreams).Perform())
		}}

	cmd.Flags().BoolVarP(&options.Verbose, "verbose", "v", false, "Show extra verbiage.")

	cmd.SetUsageTemplate(`
USAGE: waldo add [-v | --verbose] <recipe-name>

ARGUMENTS:
  <recipe-name>           The name of the recipe to add.

OPTIONS:
  -v, --verbose           Show extra verbiage.
`)

	return cmd
}
