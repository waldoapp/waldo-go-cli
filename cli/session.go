package cli

import (
	"github.com/waldoapp/waldo-go-cli/waldo"
	"github.com/waldoapp/waldo-go-cli/waldo/data"

	"github.com/spf13/cobra"
)

func NewSessionCommand() *cobra.Command {
	options := &waldo.SessionOptions{}

	cmd := &cobra.Command{
		Use:   "session",
		Short: "Launch Waldo Session against the most recent build uploaded.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return waldo.NewSessionAction(
				options,
				data.Overrides()).Perform()
		}}

	cmd.Flags().StringVarP(&options.Language, "language", "l", "", "The device language.")
	cmd.Flags().StringVarP(&options.Model, "model", "m", "", "The device model.")
	cmd.Flags().StringVarP(&options.OSVersion, "os_version", "o", "", "The OS version of the device.")
	cmd.Flags().BoolVarP(&options.Verbose, "verbose", "v", false, "Display extra verbiage.")

	return cmd
}
