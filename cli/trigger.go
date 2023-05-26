package cli

import (
	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo"

	"github.com/spf13/cobra"
)

func NewTriggerCommand() *cobra.Command {
	options := &waldo.TriggerOptions{}

	cmd := &cobra.Command{
		Use:   "trigger [options]",
		Short: "Trigger run on Waldo.",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ioStreams := lib.NewIOStreams(
				cmd.InOrStdin(),
				cmd.OutOrStdout(),
				cmd.ErrOrStderr())

			exitOnError(
				cmd,
				waldo.NewTriggerAction(
					options,
					ioStreams).Perform())
		}}

	cmd.Flags().StringVar(&options.GitCommit, "git_commit", "", "Hash of originating git commit.")
	cmd.Flags().BoolVar(&options.Help, "help", false, "Display available options and exit.")
	cmd.Flags().StringVar(&options.RuleName, "rule_name", "", "Rule name.")
	cmd.Flags().StringVar(&options.UploadToken, "upload_token", "", "Upload token (overrides WALDO_UPLOAD_TOKEN).")
	cmd.Flags().BoolVar(&options.Verbose, "verbose", false, "Display extra verbiage.")
	cmd.Flags().BoolVar(&options.Version, "version", false, "Display version and exit.")

	cmd.SetUsageTemplate(`
USAGE: waldo trigger [options]

OPTIONS:
  --git_commit <value>    Hash of originating git commit.
  --help                  Display available options and exit.
  --rule_name <value>     Rule name.
  --upload_token <value>  Upload token (overrides WALDO_UPLOAD_TOKEN).
  --verbose               Display extra verbiage.
  --version               Display version and exit.
`)

	return cmd
}
