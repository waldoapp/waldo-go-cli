package cli

import (
	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo"
	"github.com/waldoapp/waldo-go-cli/waldo/data"

	"github.com/spf13/cobra"
)

func NewUploadCommand() *cobra.Command {
	options := &waldo.UploadOptions{}

	cmd := &cobra.Command{
		Use:   "upload [options] [<recipe-name-or-build-path>]",
		Short: "Upload a build to Waldo.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ioStreams := lib.NewIOStreams(
				cmd.InOrStdin(),
				cmd.OutOrStdout(),
				cmd.ErrOrStderr())

			if len(args) > 0 {
				options.Target = args[0]
			}

			return waldo.NewUploadAction(
				options,
				ioStreams,
				data.Overrides()).Perform()
		}}

	cmd.Flags().StringVar(&options.GitBranch, "git_branch", "", "The originating git commit branch name.")
	cmd.Flags().StringVar(&options.GitCommit, "git_commit", "", "The originating git commit hash.")
	cmd.Flags().StringVarP(&options.UploadToken, "upload_token", "u", "", "The upload token associated with your app.")
	cmd.Flags().StringVar(&options.VariantName, "variant_name", "", "An optional variant name.")
	cmd.Flags().BoolVarP(&options.Verbose, "verbose", "v", false, "Display extra verbiage.")

	cmd.SetUsageTemplate(`
USAGE: waldo upload [options] [<recipe-name-or-build-path>]

ARGUMENTS:
  <recipe-name-or-build-path>  The name of the recipe to upload.

OPTIONS:
      --git_branch             The originating git commit branch name.
      --git_commit             The originating git commit hash.
  -u, --upload_token           The upload token associated with your app.
      --variant_name           An optional variant name.
  -v, --verbose                Display extra verbiage.
`)

	return cmd
}
