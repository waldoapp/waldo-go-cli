package cli

import (
	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo"

	"github.com/spf13/cobra"
)

func NewBuildCommand() *cobra.Command {
	options := &waldo.BuildOptions{}

	cmd := &cobra.Command{
		Use:   "build [-c | --clean] [-v | --verbose] [<recipe-name>]",
		Short: "Build an app from the recipe.",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ioStreams := lib.NewIOStreams(
				cmd.InOrStdin(),
				cmd.OutOrStdout(),
				cmd.ErrOrStderr())

			if len(args) > 0 {
				options.RecipeName = args[0]
			}

			exitOnError(
				cmd,
				waldo.NewBuildAction(
					options,
					ioStreams).Perform())
		}}

	cmd.Flags().BoolVarP(&options.Clean, "clean", "c", false, "Remove cached artifacts before building.")
	cmd.Flags().BoolVarP(&options.Verbose, "verbose", "v", false, "Show extra verbiage.")

	cmd.SetUsageTemplate(`
USAGE: waldo build [-c | --clean] [-v | --verbose] [<recipe-name>]

ARGUMENTS:
  <recipe-name>           The name of the recipe to build.

OPTIONS:
  -c, --clean             Remove cached artifacts before building.
  -v, --verbose           Show extra verbiage.
`)

	return cmd
}
