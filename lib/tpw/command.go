package tpw

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/waldoapp/waldo-go-cli/lib"
)

type Action interface {
	Perform() error
}

type Argument struct {
	Name     string
	Required bool
	Usage    string
}

type Option struct {
	Name      string
	ShortName string
	ValueName bool
	Usage     string
}

type Command struct {
	IOStreams *lib.IOStreams

	wrappedCommand *cobra.Command
}

// func NewAddCommand() *cobra.Command {
// 	options := &waldo.AddOptions{}

// 	cmd := &cobra.Command{
// 		Use:   "add [-v | --verbose] <recipe-name>",
// 		Short: "Add a recipe describing how to build and upload a specific variant of the app.",
// 		Args:  cobra.ExactArgs(1),
// 		Run: func(cmd *cobra.Command, args []string) {
// 			ioStreams := lib.NewIOStreams(
// 				cmd.InOrStdin(),
// 				cmd.OutOrStdout(),
// 				cmd.ErrOrStderr())

// 			options.RecipeName = args[0]

// 			exitOnError(
// 				cmd,
// 				waldo.NewAddAction(
// 					options,
// 					ioStreams).Perform())
// 		}}

// 	cmd.Flags().BoolVarP(&options.Verbose, "verbose", "v", false, "Show extra verbiage.")

// 	cmd.SetUsageTemplate(`
// USAGE: waldo add [-v | --verbose] <recipe-name>

// ARGUMENTS:
//   <recipe-name>           The name of the recipe to add.

// OPTIONS:
//   -v, --verbose           Show extra verbiage.
// `)

// 	return cmd
// }

//-----------------------------------------------------------------------------

func NewCommand(overview string, options []Option, arguments []Argument, action Action) *Command {
	wc := &cobra.Command{
		Use:   generateUsage(),
		Short: overview,
		Args:  generateArgs(),
		Run: func(cmd *cobra.Command, args []string) {
			ioStreams := lib.NewIOStreams(
				cmd.InOrStdin(),
				cmd.OutOrStdout(),
				cmd.ErrOrStderr())

			options.RecipeName = args[0]

			exitOnError(
				cmd,
				waldo.NewAddAction(
					options,
					ioStreams).Perform())
		}}

	wc.DisableAutoGenTag = true
	wc.DisableFlagsInUseLine = true

	cmd := &Command{
		IOStreams: lib.NewIOStreams(
			wc.InOrStdin(),
			wc.OutOrStdout(),
			wc.ErrOrStderr()),
		wrappedCommand: wc}

	return cmd
}

//-----------------------------------------------------------------------------

func (c *Command) ExitOnError(prefix string, err error) {
	if err != nil {
		c.IOStreams.EmitError(prefix, err)

		if ee, ok := err.(*exec.ExitError); ok {
			os.Exit(ee.ExitCode())
		}

		os.Exit(1)
	}
}

//-----------------------------------------------------------------------------
