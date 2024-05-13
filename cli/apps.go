package cli

import (
	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo"

	"github.com/spf13/cobra"
)

func NewAppsCommand() *cobra.Command {
	options := &waldo.AppsOptions{}

	cmd := &cobra.Command{
		Use:   "apps [-a | --android] [-i | --ios]",
		Short: "Show available apps for the currently authenticated user.",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ioStreams := lib.NewIOStreams(
				cmd.InOrStdin(),
				cmd.OutOrStdout(),
				cmd.ErrOrStderr())

			exitOnError(
				cmd,
				waldo.NewAppsAction(
					options,
					ioStreams).Perform())
		}}

	cmd.Flags().BoolVarP(&options.AndroidOnly, "android", "a", false, "Show Android apps only.")
	cmd.Flags().BoolVarP(&options.IosOnly, "ios", "i", false, "Show iOS apps only.")

	cmd.SetUsageTemplate(`
USAGE: waldo apps [-a | --android] [-ios | --ios]

OPTIONS:
  -a, --android           Show Android apps only.
  -i, --ios               Show iOS apps only.
`)

	return cmd
}
