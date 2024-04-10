package flutter

import (
	"errors"
	"sort"

	"github.com/waldoapp/waldo-go-cli/lib"
)

func DetermineFlavor(flavors []string, verbose bool, ios *lib.IOStreams) (string, error) {
	if len(flavors) > 1 {
		if verbose {
			ios.Printf("\nMore than supported one build flavor found\n")
		}

		return askFlavor(flavors, ios), nil
	}

	if len(flavors) == 1 {
		ios.Printf("\nOnly one supported build flavor found: %q\n", flavors[0])

		return flavors[0], nil
	}

	if verbose {
		ios.Printf("\nNo supported build flavors found\n")
	}

	return "", nil
}

func DeterminePlatform(verbose bool, ios *lib.IOStreams) (lib.Platform, error) {
	platforms := []string{
		string(lib.PlatformAndroid),
		string(lib.PlatformIos)}

	if len(platforms) > 1 {
		if verbose {
			ios.Printf("\nMore than one build platform supported\n")
		}

		return lib.ParsePlatform(askPlatform(platforms, ios)), nil
	}

	if len(platforms) == 1 {
		ios.Printf("\nOnly one build platform supported: %q\n", platforms[0])

		return lib.ParsePlatform(platforms[0]), nil
	}

	return "", errors.New("No build platforms supported")
}

//-----------------------------------------------------------------------------

func askFlavor(flavors []string, ios *lib.IOStreams) string {
	sort.Strings(flavors)

	idx := ios.PromptReader().ReadChoose(
		"Supported build flavors",
		flavors,
		"Choose a build flavor")

	return flavors[idx]
}

func askPlatform(platforms []string, ios *lib.IOStreams) string {
	sort.Strings(platforms)

	idx := ios.PromptReader().ReadChoose(
		"Supported build platforms",
		platforms,
		"Choose a build platform")

	return platforms[idx]
}
