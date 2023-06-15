package tool

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/waldoapp/waldo-go-cli/lib"
)

func DetermineBuildPath(items []*BuildPath, verbose bool, ios *lib.IOStreams) (*BuildPath, error) {
	if len(items) > 1 {
		if verbose {
			ios.Printf("\nMore than one build path found\n")
		}

		return askBuildPath(items, ios), nil
	}

	if len(items) == 1 {
		item := items[0]

		ios.Printf("\nOnly one build path found: %q (%s)\n", item.RelPath, item.BuildTool.String())

		return item, nil
	}

	return nil, errors.New("No build paths found")
}

//-----------------------------------------------------------------------------

func askBuildPath(items []*BuildPath, ios *lib.IOStreams) *BuildPath {
	pr := ios.PromptReader()

	sort.Slice(items, func(i, j int) bool {
		return strings.ToLower(items[i].RelPath) < strings.ToLower(items[j].RelPath)
	})

	maxLen := 0

	for _, item := range items {
		rpLen := len(item.RelPath)

		if maxLen < rpLen {
			maxLen = rpLen
		}
	}

	choices := lib.Map(items, func(item *BuildPath) string {
		return fmt.Sprintf("%-*s (%s)", maxLen, item.RelPath, item.BuildTool.String())
	})

	idx := pr.ReadChoose(
		"Available build paths",
		choices,
		"Choose a build path")

	return items[idx]
}
