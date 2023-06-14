package xcode

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/waldoapp/waldo-go-cli/lib"
)

func askConfiguration(configurations []string, ios *lib.IOStreams) string {
	pr := ios.PromptReader()

	sort.Strings(configurations)

	idx := pr.ReadChoose(
		"Available Xcode configurations",
		configurations,
		"Choose a configuration")

	return configurations[idx]
}

func askProject(projects []string, ios *lib.IOStreams) string {
	pr := ios.PromptReader()

	sort.Strings(projects)

	idx := pr.ReadChoose(
		"Available Xcode workspaces and projects",
		projects,
		"Choose a workspace or project")

	return projects[idx]
}

func askScheme(schemes []string, ios *lib.IOStreams) string {
	pr := ios.PromptReader()

	sort.Strings(schemes)

	idx := pr.ReadChoose(
		"Available Xcode schemes",
		schemes,
		"Choose a scheme")

	return schemes[idx]
}

func determineConfiguration(project string, xi *XcodeInfo, verbose bool, ios *lib.IOStreams) (string, error) {
	configs := xi.Configurations()

	if len(configs) > 1 {
		if verbose {
			ios.Printf("\nMore than one Xcode configuration found in %q\n", project)
		}

		return askConfiguration(configs, ios), nil
	}

	if len(configs) == 1 {
		ios.Printf("\nOnly one Xcode configuration found in %q: %q\n", project, configs[0])

		return configs[0], nil
	}

	if verbose {
		ios.Printf("\nNo Xcode configurations found in %q\n", project)
	}

	return "", nil
}

func determineProject(projects []string, verbose bool, ios *lib.IOStreams) (string, error) {
	projects = lib.Map(projects, func(project string) string {
		return filepath.Base(project)
	})

	if len(projects) > 1 {
		if verbose {
			ios.Printf("\nMore than one Xcode workspace or project found\n")
		}

		return askProject(projects, ios), nil
	}

	if len(projects) == 1 {
		ios.Printf("\nOnly one Xcode workspace or project found: %q\n", projects[0])

		return projects[0], nil
	}

	return "", errors.New("No Xcode workspaces or projects found")
}

func determineScheme(project string, xi *XcodeInfo, verbose bool, ios *lib.IOStreams) (string, error) {
	schemes := xi.Schemes()

	if len(schemes) > 1 {
		if verbose {
			ios.Printf("\nMore than one Xcode scheme found in %q\n", project)
		}

		return askScheme(schemes, ios), nil
	}

	if len(schemes) == 1 {
		ios.Printf("\nOnly one Xcode scheme found in %q: %q\n", project, schemes[0])

		return schemes[0], nil
	}

	return "", fmt.Errorf("No Xcode schemes found in %q", project)
}
