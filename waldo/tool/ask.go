package tool

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/waldoapp/waldo-go-cli/lib"
)

func DetermineApp(platform lib.Platform, items []*AppInfo, verbose bool, ios *lib.IOStreams) (*AppInfo, error) {
	if len(items) > 1 {
		if verbose {
			ios.Printf("\nMore than one %s app found\n", platform)
		}

		return askApp(platform, items, ios), nil
	}

	if len(items) == 1 {
		item := items[0]

		ios.Printf("\nOnly one %s app found: %q (%s)\n", platform, item.AppName, item.AppID)

		return item, nil
	}

	return nil, fmt.Errorf("No %s apps found", platform)
}

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

func askApp(platform lib.Platform, items []*AppInfo, ios *lib.IOStreams) *AppInfo {
	maxLen := maxAppNameLength(items)

	slices.SortStableFunc(items, func(a, b *AppInfo) int {
		return cmp.Compare(strings.ToLower(a.AppName), strings.ToLower(b.AppName))
	})

	choices := lib.Map(items, func(item *AppInfo) string {
		return fmt.Sprintf("%-*s (%s)", maxLen, item.AppName, item.AppID)
	})

	idx := ios.PromptReader().ReadChoose(
		fmt.Sprintf("Available %s apps", platform),
		choices,
		"Choose an app")

	return items[idx]
}

func askBuildPath(items []*BuildPath, ios *lib.IOStreams) *BuildPath {
	maxLen := maxRelPathLength(items)

	slices.SortStableFunc(items, func(a, b *BuildPath) int {
		return cmp.Compare(strings.ToLower(a.RelPath), strings.ToLower(b.RelPath))
	})

	choices := lib.Map(items, func(item *BuildPath) string {
		return fmt.Sprintf("%-*s (%s)", maxLen, item.RelPath, item.BuildTool.String())
	})

	idx := ios.PromptReader().ReadChoose(
		"Available build paths",
		choices,
		"Choose a build path")

	return items[idx]
}

func maxAppNameLength(items []*AppInfo) int {
	maxLen := 0

	for _, item := range items {
		anLen := len(item.AppName)

		if maxLen < anLen {
			maxLen = anLen
		}
	}

	return maxLen
}

func maxRelPathLength(items []*BuildPath) int {
	maxLen := 0

	for _, item := range items {
		rpLen := len(item.RelPath)

		if maxLen < rpLen {
			maxLen = rpLen
		}
	}

	return maxLen
}
