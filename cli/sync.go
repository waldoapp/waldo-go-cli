package cli

import (
	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo"

	"github.com/spf13/cobra"
)

func NewSyncCommand() *cobra.Command {
	options := &waldo.SyncOptions{}

	cmd := &cobra.Command{
		Use:   "sync [-c | --clean] [--git_branch <b>] [--git_commit <c>] [--variant_name <n>] [-v | --verbose] [<recipe-name>]",
		Short: "Build an app from the recipe and then upload it to Waldo.",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ioStreams := lib.NewIOStreams(
				cmd.InOrStdin(),
				cmd.OutOrStdout(),
				cmd.ErrOrStderr())

			if len(args) > 0 {
				options.RecipeName = args[0]
			}

			exitOnError(
				cmd,
				waldo.NewSyncAction(
					options,
					ioStreams).Perform())
		}}

	cmd.Flags().BoolVarP(&options.Clean, "clean", "c", false, "Remove cached artifacts before building.")
	cmd.Flags().StringVar(&options.GitBranch, "git_branch", "", "The originating git commit branch name.")
	cmd.Flags().StringVar(&options.GitCommit, "git_commit", "", "The originating git commit hash.")
	cmd.Flags().StringVar(&options.VariantName, "variant_name", "", "An optional variant name.")
	cmd.Flags().BoolVarP(&options.Verbose, "verbose", "v", false, "Show extra verbiage.")

	cmd.SetUsageTemplate(`
USAGE: waldo sync [-c | --clean] [--git_branch <b>] [--git_commit <c>] [--variant_name <n>] [-v | --verbose] [<recipe-name>]

ARGUMENTS:
  <recipe-name>           The name of the recipe to build and upload.

OPTIONS:
  -c, --clean             Remove cached artifacts before building.
      --git_branch <b>    The originating git commit branch name.
      --git_commit <c>    The originating git commit hash.
      --variant_name <n>  An optional variant name.
  -v, --verbose           Show extra verbiage.
`)

	return cmd
}
