package cli

import (
	"github.com/waldoapp/waldo-go-cli/waldo"
	"github.com/waldoapp/waldo-go-cli/waldo/data"

	"github.com/spf13/cobra"
)

func NewRunCommand() *cobra.Command {
	options := &waldo.RunOptions{}

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Execute one or more locally-defined scripts against an uploaded build.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return waldo.NewRunAction(
				options,
				data.Overrides()).Perform()
		}}

	cmd.Flags().BoolVarP(&options.Interactive, "interactive", "i", false, "Run interactively.")
	cmd.Flags().BoolVarP(&options.Preview, "preview", "p", false, "")
	cmd.Flags().BoolVarP(&options.Verbose, "verbose", "v", false, "Display extra verbiage.")

	return cmd
}
