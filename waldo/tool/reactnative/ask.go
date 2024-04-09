package reactnative

import (
	"errors"
	"sort"

	"github.com/waldoapp/waldo-go-cli/lib"
)

func DetermineMode(modes []string, verbose bool, ios *lib.IOStreams) (string, error) {
	if len(modes) > 1 {
		if verbose {
			ios.Printf("\nMore than supported one build mode found\n")
		}

		return askMode(modes, ios), nil
	}

	if len(modes) == 1 {
		ios.Printf("\nOnly one supported build mode found: %q\n", modes[0])

		return modes[0], nil
	}

	if verbose {
		ios.Printf("\nNo supported build modes found\n")
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

func askMode(modes []string, ios *lib.IOStreams) string {
	sort.Strings(modes)

	idx := ios.PromptReader().ReadChoose(
		"Supported build modes",
		modes,
		"Choose a build mode")

	return modes[idx]
}

func askPlatform(platforms []string, ios *lib.IOStreams) string {
	sort.Strings(platforms)

	idx := ios.PromptReader().ReadChoose(
		"Supported build platforms",
		platforms,
		"Choose a build platform")

	return platforms[idx]
}
