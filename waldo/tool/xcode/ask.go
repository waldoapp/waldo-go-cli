package xcode

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/waldoapp/waldo-go-cli/lib"
)

func DetermineConfiguration(pwName string, configurations []string, verbose bool, ios *lib.IOStreams) (string, error) {
	if len(configurations) > 1 {
		if verbose {
			ios.Printf("\nMore than one Xcode configuration found in %q\n", pwName)
		}

		return askConfiguration(configurations, ios), nil
	}

	if len(configurations) == 1 {
		ios.Printf("\nOnly one Xcode configuration found in %q: %q\n", pwName, configurations[0])

		return configurations[0], nil
	}

	return "", fmt.Errorf("No Xcode configurations found in %q", pwName)
}

func DetermineProject(pwNames []string, verbose bool, ios *lib.IOStreams) (string, error) {
	pwNames = lib.Map(pwNames, func(project string) string {
		return filepath.Base(project)
	})

	if len(pwNames) > 1 {
		if verbose {
			ios.Printf("\nMore than one Xcode workspace or project found\n")
		}

		return askProject(pwNames, ios), nil
	}

	if len(pwNames) == 1 {
		ios.Printf("\nOnly one Xcode workspace or project found: %q\n", pwNames[0])

		return pwNames[0], nil
	}

	return "", errors.New("No Xcode workspaces or projects found")
}

func DetermineScheme(pwName string, schemes []string, verbose bool, ios *lib.IOStreams) (string, error) {
	if len(schemes) > 1 {
		if verbose {
			ios.Printf("\nMore than one Xcode scheme found in %q\n", pwName)
		}

		return askScheme(schemes, ios), nil
	}

	if len(schemes) == 1 {
		ios.Printf("\nOnly one Xcode scheme found in %q: %q\n", pwName, schemes[0])

		return schemes[0], nil
	}

	return "", fmt.Errorf("No Xcode schemes found in %q", pwName)
}

//-----------------------------------------------------------------------------

func askConfiguration(configurations []string, ios *lib.IOStreams) string {
	sort.Strings(configurations)

	idx := ios.PromptReader().ReadChoose(
		"Available Xcode configurations",
		configurations,
		"Choose an Xcode configuration")

	return configurations[idx]
}

func askProject(pwNames []string, ios *lib.IOStreams) string {
	sort.Strings(pwNames)

	idx := ios.PromptReader().ReadChoose(
		"Available Xcode workspaces and projects",
		pwNames,
		"Choose an Xcode workspace or project")

	return pwNames[idx]
}

func askScheme(schemes []string, ios *lib.IOStreams) string {
	sort.Strings(schemes)

	idx := ios.PromptReader().ReadChoose(
		"Available Xcode schemes",
		schemes,
		"Choose an Xcode scheme")

	return schemes[idx]
}
