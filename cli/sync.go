package cli

import (
	"github.com/waldoapp/waldo-go-cli/waldo"
	"github.com/waldoapp/waldo-go-cli/waldo/data"

	"github.com/spf13/cobra"
)

func NewSyncCommand() *cobra.Command {
	options := &waldo.SyncOptions{}

	cmd := &cobra.Command{
		Use:   "sync [<recipe-name>]",
		Short: "Shorthand for waldo build followed by waldo upload.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return waldo.NewSyncAction(
				options,
				data.Overrides()).Perform()
		}}

	cmd.Flags().BoolVarP(&options.Clean, "clean", "c", false, "Remove cached artifacts before building.")
	cmd.Flags().StringVar(&options.GitBranch, "git_branch", "", "The originating git commit branch name.")
	cmd.Flags().StringVar(&options.GitCommit, "git_commit", "", "The originating git commit hash.")
	cmd.Flags().StringVarP(&options.UploadToken, "upload_token", "u", "", "The upload token associated with your app.")
	cmd.Flags().StringVar(&options.VariantName, "variant_name", "", "An optional variant name.")
	cmd.Flags().BoolVarP(&options.Verbose, "verbose", "v", false, "Display extra verbiage.")

	return cmd
}
