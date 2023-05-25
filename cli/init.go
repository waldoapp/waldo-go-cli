package cli

import (
	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo"

	"github.com/spf13/cobra"
)

func NewInitCommand() *cobra.Command {
	options := &waldo.InitOptions{}

	cmd := &cobra.Command{
		Use:   "init [options]",
		Short: "Create an empty Waldo configuration.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ioStreams := lib.NewIOStreams(
				cmd.InOrStdin(),
				cmd.OutOrStdout(),
				cmd.ErrOrStderr())

			return waldo.NewInitAction(
				options,
				ioStreams).Perform()
		}}

	cmd.Flags().BoolVarP(&options.Verbose, "verbose", "v", false, "Display extra verbiage.")

	cmd.SetUsageTemplate(`
USAGE: waldo init [options]

OPTIONS:
  -v, --verbose   Display extra verbiage.
`)

	return cmd
}
