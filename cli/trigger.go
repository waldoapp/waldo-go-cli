package cli

import (
	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo"

	"github.com/spf13/cobra"
)

func NewTriggerCommand() *cobra.Command {
	options := &waldo.TriggerOptions{}

	cmd := &cobra.Command{
		Use:   "trigger [--git_commit <c>] [--rule_name <r>] [--upload_token <t>] [-v | --verbose]",
		Short: "Trigger a run on Waldo.",
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

	cmd.Flags().StringVar(&options.GitCommit, "git_commit", "", "The originating git commit hash.")
	cmd.Flags().BoolVar(&options.LegacyHelp, "help", false, "Show available options and exit.")
	cmd.Flags().StringVar(&options.RuleName, "rule_name", "", "An optional rule name.")
	cmd.Flags().StringVar(&options.UploadToken, "upload_token", "", "The upload token (overrides WALDO_UPLOAD_TOKEN).")
	cmd.Flags().BoolVarP(&options.Verbose, "verbose", "v", false, "Show extra verbiage.")
	cmd.Flags().BoolVar(&options.LegacyVersion, "version", false, "Show version information and exit.")

	cmd.SetUsageTemplate(`
USAGE: waldo trigger [--git_commit <c>] [--rule_name <r>] [--upload_token <t>] [-v | --verbose]

OPTIONS:
      --git_commit <c>    The originating git commit hash.
      --rule_name <r>     An optional rule name.
      --upload_token <t>  The upload token (overrides WALDO_UPLOAD_TOKEN).
  -v, --verbose           Show extra verbiage.
`)

	return cmd
}
