package cli

import (
	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo"
	"github.com/waldoapp/waldo-go-cli/waldo/data"

	"github.com/spf13/cobra"
)

func NewBuildCommand() *cobra.Command {
	options := &waldo.BuildOptions{}

	cmd := &cobra.Command{
		Use:   "build [<recipe-name>]",
		Short: "Build app from recipe.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ioStreams := lib.NewIOStreams(
				cmd.InOrStdin(),
				cmd.OutOrStdout(),
				cmd.ErrOrStderr())

			if len(args) > 0 {
				options.RecipeName = args[0]
			}

			return waldo.NewBuildAction(
				options,
				ioStreams,
				data.Overrides()).Perform()
		}}

	cmd.Flags().BoolVarP(&options.Clean, "clean", "c", false, "Remove cached artifacts before building.")
	cmd.Flags().BoolVarP(&options.Verbose, "verbose", "v", false, "Display extra verbiage.")

	return cmd
}
