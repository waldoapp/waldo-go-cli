package cli

import (
	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo"

	"github.com/spf13/cobra"
)

func NewSessionCommand() *cobra.Command {
	options := &waldo.SessionOptions{}

	cmd := &cobra.Command{
		Use:   "session [options]",
		Short: "Launch Waldo Session against the most recent build uploaded.",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ioStreams := lib.NewIOStreams(
				cmd.InOrStdin(),
				cmd.OutOrStdout(),
				cmd.ErrOrStderr())

			exitOnError(
				cmd,
				waldo.NewSessionAction(
					options,
					ioStreams).Perform())
		}}

	cmd.Flags().StringVarP(&options.Language, "language", "l", "", "The device language.")
	cmd.Flags().StringVarP(&options.Model, "model", "m", "", "The device model.")
	cmd.Flags().StringVarP(&options.OSVersion, "os_version", "o", "", "The OS version of the device.")
	cmd.Flags().BoolVarP(&options.Verbose, "verbose", "v", false, "Display extra verbiage.")

	cmd.SetUsageTemplate(`
USAGE: waldo session [options]

OPTIONS:
  -l, --language <l>     The device language.
  -m, --model <m>        The device model.
  -o, --os_version <v>   The OS version of the device.
  -v, --verbose          Display extra verbiage.
`)

	return cmd
}
