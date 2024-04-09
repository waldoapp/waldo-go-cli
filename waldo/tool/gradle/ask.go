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
			ios.Printf("\nMore than one module found\n")
		}

		return askModule(modules, ios), nil
	}

	if len(modules) == 1 {
		ios.Printf("\nOnly one module found: %q\n", modules[0])

		return modules[0], nil
	}

	return "", errors.New("No modules found")
}

func DetermineVariant(module string, variants []string, verbose bool, ios *lib.IOStreams) (string, error) {
	if len(variants) > 1 {
		if verbose {
			ios.Printf("\nMore than one build variant found in module %q\n", module)
		}

		return askVariant(variants, ios), nil
	}

	if len(variants) == 1 {
		ios.Printf("\nOnly one build variant found in module %q: %q\n", module, variants[0])

		return variants[0], nil
	}

	return "", fmt.Errorf("No build variants found in module %q", module)
}

//-----------------------------------------------------------------------------

func askModule(modules []string, ios *lib.IOStreams) string {
	sort.Strings(modules)

	idx := ios.PromptReader().ReadChoose(
		"Available modules",
		modules,
		"Choose a module")

	return modules[idx]
}

func askVariant(variants []string, ios *lib.IOStreams) string {
	sort.Strings(variants)

	idx := ios.PromptReader().ReadChoose(
		"Available build variants",
		variants,
		"Choose a build variant")

	return variants[idx]
}
