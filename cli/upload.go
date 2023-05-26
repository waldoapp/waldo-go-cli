package cli

import (
	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo"

	"github.com/spf13/cobra"
)

func NewUploadCommand() *cobra.Command {
	options := &waldo.UploadOptions{}

	cmd := &cobra.Command{
		Use:   "upload [options] <build-path>",
		Short: "Upload build to Waldo.",
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

	cmd.Flags().StringVar(&options.GitBranch, "git_branch", "", "Branch name for originating git commit.")
	cmd.Flags().StringVar(&options.GitCommit, "git_commit", "", "Hash of originating git commit.")
	cmd.Flags().BoolVar(&options.Help, "help", false, "Display available options and exit.")
	cmd.Flags().StringVar(&options.UploadToken, "upload_token", "", "Upload token (overrides WALDO_UPLOAD_TOKEN).")
	cmd.Flags().StringVar(&options.VariantName, "variant_name", "", "Variant name.")
	cmd.Flags().BoolVar(&options.Verbose, "verbose", false, "Display extra verbiage.")
	cmd.Flags().BoolVar(&options.Version, "version", false, "Display version and exit.")

	cmd.SetUsageTemplate(`
USAGE: waldo upload [options] <build-path>

OPTIONS:
  --git_branch <value>    Branch name for originating git commit.
  --git_commit <value>    Hash of originating git commit.
  --help                  Display available options and exit.
  --upload_token <value>  Upload token (overrides WALDO_UPLOAD_TOKEN).
  --variant_name <value>  Variant name.
  --verbose               Display extra verbiage.
  --version               Display version and exit.
`)

	return cmd
}
