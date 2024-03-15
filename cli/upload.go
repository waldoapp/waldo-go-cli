package cli

import (
	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo"

	"github.com/spf13/cobra"
)

func NewUploadCommand() *cobra.Command {
	options := &waldo.UploadOptions{}

	cmd := &cobra.Command{
		Use:   "upload [--git_branch <b>] [--git_commit <c>] [--upload_token <t>] [--variant_name <n>] [-v | --verbose ] <recipe-or-build-path>",
		Short: "Upload a build artifact to Waldo.",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ioStreams := lib.NewIOStreams(
				cmd.InOrStdin(),
				cmd.OutOrStdout(),
				cmd.ErrOrStderr())

			if len(args) > 0 {
				options.Target = args[0]
			}

			exitOnError(
				cmd,
				waldo.NewUploadAction(
					options,
					ioStreams).Perform())
		}}

	cmd.Flags().StringVar(&options.GitBranch, "git_branch", "", "The originating git commit branch name.")
	cmd.Flags().StringVar(&options.GitCommit, "git_commit", "", "The originating git commit hash.")
	cmd.Flags().BoolVar(&options.Help, "help", false, "Show available options and exit.")
	cmd.Flags().StringVar(&options.UploadToken, "upload_token", "", "The upload token (overrides WALDO_UPLOAD_TOKEN).")
	cmd.Flags().StringVar(&options.VariantName, "variant_name", "", "An optional variant name.")
	cmd.Flags().BoolVar(&options.Verbose, "verbose", false, "Show extra verbiage.")
	cmd.Flags().BoolVar(&options.Version, "version", false, "Show version information and exit.")

	cmd.SetUsageTemplate(`
USAGE: waldo upload [--git_branch <b>] [--git_commit <c>] [--upload_token <t>] [--variant_name <n>] [-v | --verbose ] <recipe-or-build-path>

ARGUMENTS:
  <recipe-or-build-path>  The name of the recipe associated with, or the direct path to, the build artifact to upload.

OPTIONS:
      --git_branch <b>    The originating git commit branch name.
      --git_commit <c>    The originating git commit hash.
      --upload_token <t>  The upload token (overrides WALDO_UPLOAD_TOKEN).
      --variant_name <n>  An optional variant name.
  -v, --verbose           Show extra verbiage.
`)

	return cmd
}
