package cli

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "waldo",
		Short: "Waldo CLI"}

	cmd.AddCommand(NewAddCommand())
	cmd.AddCommand(NewBuildCommand())
	cmd.AddCommand(NewInitCommand())
	cmd.AddCommand(NewListCommand())
	cmd.AddCommand(NewRemoveCommand())
	cmd.AddCommand(NewRunCommand())
	cmd.AddCommand(NewSessionCommand())
	cmd.AddCommand(NewSyncCommand())
	cmd.AddCommand(NewTriggerCommand())
	cmd.AddCommand(NewUploadCommand())
	cmd.AddCommand(NewVersionCommand())

	return cmd
}
