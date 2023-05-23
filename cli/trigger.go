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
		Use:   "trigger [options]",
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

	cmd.Flags().StringVar(&options.GitCommit, "git_commit", "", "The originating git commit hash.")
	cmd.Flags().StringVar(&options.RuleName, "rule_name", "", "The name of a rule.")
	cmd.Flags().StringVarP(&options.UploadToken, "upload_token", "u", "", "The upload token associated with your app.")
	cmd.Flags().BoolVarP(&options.Verbose, "verbose", "v", false, "Display extra verbiage.")

	cmd.SetUsageTemplate(`
USAGE: waldo trigger [options]

OPTIONS:
      --git_commit    The originating git commit hash.
      --rule_name     The name of a rule.
  -u, --upload_token  The upload token associated with your app.
  -v, --verbose       Display extra verbiage.
`)

	return cmd
}
