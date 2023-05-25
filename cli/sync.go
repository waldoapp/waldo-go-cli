package cli

import (
	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo"

	"github.com/spf13/cobra"
)

func NewSyncCommand() *cobra.Command {
	options := &waldo.SyncOptions{}

	cmd := &cobra.Command{
		Use:   "sync [options] [<recipe-name>]",
		Short: "Shorthand for waldo build followed by waldo upload.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ioStreams := lib.NewIOStreams(
				cmd.InOrStdin(),
				cmd.OutOrStdout(),
				cmd.ErrOrStderr())

			return waldo.NewSyncAction(
				options,
				ioStreams).Perform()
		}}

	cmd.Flags().BoolVarP(&options.Clean, "clean", "c", false, "Remove cached artifacts before building.")
	cmd.Flags().StringVar(&options.GitBranch, "git_branch", "", "The originating git commit branch name.")
	cmd.Flags().StringVar(&options.GitCommit, "git_commit", "", "The originating git commit hash.")
	cmd.Flags().StringVarP(&options.UploadToken, "upload_token", "u", "", "The upload token associated with your app.")
	cmd.Flags().StringVar(&options.VariantName, "variant_name", "", "An optional variant name.")
	cmd.Flags().BoolVarP(&options.Verbose, "verbose", "v", false, "Display extra verbiage.")

	cmd.SetUsageTemplate(`
USAGE: waldo sync [options] [<recipe-name>]

ARGUMENTS:
  <recipe-name>            The name of the recipe to build and upload.

OPTIONS:
  -c, --clean              Remove cached artifacts before building.
      --git_branch <b>     The originating git commit branch name.
      --git_commit <c>     The originating git commit hash.
  -u, --upload_token <t>   The upload token associated with your app.
      --variant_name <n>   An optional variant name.
  -v, --verbose            Display extra verbiage.
`)

	return cmd
}
