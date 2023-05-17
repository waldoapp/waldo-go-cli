package cli

import (
	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo"
	"github.com/waldoapp/waldo-go-cli/waldo/data"

	"github.com/spf13/cobra"
)

func NewTriggerCommand() *cobra.Command {
	options := &waldo.TriggerOptions{}

	cmd := &cobra.Command{
		Use:   "trigger",
		Short: "Trigger a test run on Waldo.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ioStreams := lib.NewIOStreams(
				cmd.InOrStdin(),
				cmd.OutOrStdout(),
				cmd.ErrOrStderr())

			return waldo.NewTriggerAction(
				options,
				ioStreams,
				data.Overrides()).Perform()
		}}

	cmd.Flags().StringVarP(&options.GitCommit, "git_commit", "c", "", "The originating git commit hash.")
	cmd.Flags().StringVarP(&options.RuleName, "rule_name", "r", "", "The name of a rule.")
	cmd.Flags().StringVarP(&options.UploadToken, "upload_token", "u", "", "The upload token associated with your app.")
	cmd.Flags().BoolVarP(&options.Verbose, "verbose", "v", false, "Display extra verbiage.")

	return cmd
}
