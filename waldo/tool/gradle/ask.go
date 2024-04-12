package gradle

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/waldoapp/waldo-go-cli/lib"
)

func DetermineModule(modules []string, verbose bool, ios *lib.IOStreams) (string, error) {
	modules = lib.Map(modules, func(module string) string {
		return filepath.Base(module)
	})

	if len(modules) > 1 {
		if verbose {
			ios.Printf("\nMore than one Gradle module found\n")
		}

		return askModule(modules, ios), nil
	}

	if len(modules) == 1 {
		ios.Printf("\nOnly one Gradle module found: %q\n", modules[0])

		return modules[0], nil
	}

	return "", errors.New("No Gradle modules found")
}

func DetermineVariant(module string, variants []string, verbose bool, ios *lib.IOStreams) (string, error) {
	if len(variants) > 1 {
		if verbose {
			ios.Printf("\nMore than one Gradle build variant found in module %q\n", module)
		}

		return askVariant(variants, ios), nil
	}

	if len(variants) == 1 {
		ios.Printf("\nOnly one Gradle build variant found in module %q: %q\n", module, variants[0])

		return variants[0], nil
	}

	return "", fmt.Errorf("No Gradle build variants found in module %q", module)
}

//-----------------------------------------------------------------------------

func askModule(modules []string, ios *lib.IOStreams) string {
	sort.Strings(modules)

	idx := ios.PromptReader().ReadChoose(
		"Available Gradle modules",
		modules,
		"Choose a Gradle module")

	return modules[idx]
}

func askVariant(variants []string, ios *lib.IOStreams) string {
	sort.Strings(variants)

	idx := ios.PromptReader().ReadChoose(
		"Available Gradle build variants",
		variants,
		"Choose a Gradle build variant")

	return variants[idx]
}
