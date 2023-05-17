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
		Use:       "upload [<build-path>]",
		Short:     "Upload a build to Waldo.",
		ValidArgs: []string{"build-path"},
		Args:      cobra.OnlyValidArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ioStreams := lib.NewIOStreams(
				cmd.InOrStdin(),
				cmd.OutOrStdout(),
				cmd.ErrOrStderr())

			options.BuildPath = args[0]

			return waldo.NewUploadAction(
				options,
				ioStreams,
				data.Overrides()).Perform()
		}}

	cmd.Flags().StringVarP(&options.GitBranch, "git_branch", "b", "", "The originating git commit branch name.")
	cmd.Flags().StringVarP(&options.GitCommit, "git_commit", "c", "", "The originating git commit hash.")
	cmd.Flags().StringVarP(&options.UploadToken, "upload_token", "u", "", "The upload token associated with your app.")
	cmd.Flags().StringVarP(&options.VariantName, "variant_name", "n", "", "An optional variant name.")
	cmd.Flags().BoolVarP(&options.Verbose, "verbose", "v", false, "Display extra verbiage.")

	return cmd
}
