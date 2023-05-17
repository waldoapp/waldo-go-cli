package cli

import (
	"github.com/waldoapp/waldo-go-cli/waldo"
	"github.com/waldoapp/waldo-go-cli/waldo/data"

	"github.com/spf13/cobra"
)

func NewSyncCommand() *cobra.Command {
	options := &waldo.SyncOptions{}

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Shorthand for waldo build followed by waldo upload.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return waldo.NewSyncAction(
				options,
				data.Overrides()).Perform()
		}}

	return cmd
}
