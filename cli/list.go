package cli

import (
	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo"
	"github.com/waldoapp/waldo-go-cli/waldo/data"

	"github.com/spf13/cobra"
)

func NewListCommand() *cobra.Command {
	options := &waldo.ListOptions{}

	cmd := &cobra.Command{
		Use:     "list [options]",
		Short:   "List defined recipes.",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ioStreams := lib.NewIOStreams(
				cmd.InOrStdin(),
				cmd.OutOrStdout(),
				cmd.ErrOrStderr())

			return waldo.NewListAction(
				options,
				ioStreams,
				data.Overrides()).Perform()
		}}

	cmd.Flags().BoolVarP(&options.LongFormat, "long", "l", false, "List recipes in long format.")

	cmd.SetUsageTemplate(`
USAGE: waldo list [options]

OPTIONS:
  -l, --long   List recipes in long format.
`)

	return cmd
}
