package gradle

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/waldoapp/waldo-go-cli/lib"
)

func askModule(modules []string, ios *lib.IOStreams) string {
	pr := ios.PromptReader()

	sort.Strings(modules)

	idx := pr.ReadChoose(
		"Available modules",
		modules,
		"Choose a module")

	return modules[idx]
}

func askVariant(variants []string, ios *lib.IOStreams) string {
	pr := ios.PromptReader()

	sort.Strings(variants)

	idx := pr.ReadChoose(
		"Available build variants",
		variants,
		"Choose a build variant")

	return variants[idx]
}

func determineModule(modules []string, verbose bool, ios *lib.IOStreams) (string, error) {
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

func determineVariant(module string, variants []string, verbose bool, ios *lib.IOStreams) (string, error) {
	if len(variants) > 1 {
		if verbose {
			ios.Printf("\nMore than one build variant found in %q\n", module)
		}

		return askVariant(variants, ios), nil
	}

	if len(variants) == 1 {
		ios.Printf("\nOnly one build variant found in %q: %q\n", module, variants[0])

		return variants[0], nil
	}

	return "", fmt.Errorf("No build variants found in %q", module)
}
