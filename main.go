package main

import (
	"fmt"
	"os"

	"github.com/waldoapp/waldo-go-cli/cli"
	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/data"
)

func main() {
	cmd := cli.NewRootCommand()

	defer func() {
		ioStreams := lib.NewIOStreams(
			cmd.InOrStdin(),
			cmd.OutOrStdout(),
			cmd.ErrOrStderr())

		if err := recover(); err != nil {
			ioStreams.EmitError(
				data.CLIPrefix,
				fmt.Errorf("Unhandled panic: %v", err))

			os.Exit(1)
		}
	}()

	cmd.Execute()
}
