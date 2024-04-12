package expo

import (
	"errors"
	"sort"

	"github.com/waldoapp/waldo-go-cli/lib"
)

func DeterminePlatform(verbose bool, ios *lib.IOStreams) (lib.Platform, error) {
	platforms := []string{
		string(lib.PlatformAndroid),
		string(lib.PlatformIos)}

	if len(platforms) > 1 {
		if verbose {
			ios.Printf("\nMore than one Expo build platform supported\n")
		}

		return lib.ParsePlatform(askPlatform(platforms, ios)), nil
	}

	if len(platforms) == 1 {
		ios.Printf("\nOnly one Expo build platform supported: %v\n", platforms[0])

		return lib.ParsePlatform(platforms[0]), nil
	}

	return "", errors.New("No Expo build platforms supported")
}

//-----------------------------------------------------------------------------

func askPlatform(platforms []string, ios *lib.IOStreams) string {
	sort.Strings(platforms)

	idx := ios.PromptReader().ReadChoose(
		"Supported Expo build platforms",
		platforms,
		"Choose an Expo build platform")

	return platforms[idx]
}
