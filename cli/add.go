package cli

import (
	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo"
	"github.com/waldoapp/waldo-go-cli/waldo/data"

	"github.com/spf13/cobra"
)

func NewAddCommand() *cobra.Command {
	options := &waldo.AddOptions{}

	cmd := &cobra.Command{
		Use:   "add [options] <recipe-name>",
		Short: "Add a recipe describing how to build and upload a specific variant of the app.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ioStreams := lib.NewIOStreams(
				cmd.InOrStdin(),
				cmd.OutOrStdout(),
				cmd.ErrOrStderr())

			options.RecipeName = args[0]

			return waldo.NewAddAction(
				options,
				ioStreams,
				data.Overrides()).Perform()
		}}

	cmd.Flags().StringVarP(&options.AppName, "app_name", "a", "", "The name associated with your app.")
	cmd.Flags().StringVarP(&options.Platform, "platform", "p", "", "The platform associated with your app.")
	cmd.Flags().StringVarP(&options.UploadToken, "upload_token", "u", "", "The upload token associated with your app.")
	cmd.Flags().BoolVarP(&options.Verbose, "verbose", "v", false, "Display extra verbiage.")

	cmd.SetUsageTemplate(`
USAGE: waldo add [options] <recipe-name>

ARGUMENTS:
  <recipe-name>       The name of the recipe to add.

OPTIONS:
  -a, --app_name      The name associated with your app.
  -p, --platform      The platform associated with your app.
  -u, --upload_token  The upload token associated with your app.
  -v, --verbose       Display extra verbiage.
`)

	return cmd
}