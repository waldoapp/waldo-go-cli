package gradle

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/waldoapp/waldo-go-cli/lib"
)

type GradleBuilder struct {
	Module  string `yaml:"module,omitempty"`
	Variant string `yaml:"variant,omitempty"`
}

//-----------------------------------------------------------------------------

func IsPossibleGradleContainer(path string) bool {
	wrapperPath := filepath.Join(path, wrapperName())
	kotlinPath := filepath.Join(path, "build.gradle.kts")
	groovyPath := filepath.Join(path, "build.gradle")

	if !lib.IsRegularFile(kotlinPath) && !lib.IsRegularFile(groovyPath) {
		return false
	}

	return lib.IsRegularFile(wrapperPath)
}

func MakeGradleBuilder(absPath, relPath string, verbose bool, ios *lib.IOStreams) (*GradleBuilder, string, lib.Platform, error) {
	ios.Printf("\nSearching for Gradle modules in %q…\n", relPath)

	properties := fetchProperties(absPath, "")

	modules := extractModules(properties["subprojects"])

	module, err := determineModule(modules, verbose, ios)

	if err != nil {
		return nil, "", lib.PlatformUnknown, err
	}

	ios.Printf("\nFinding all build variants in %q…\n", module)

	tasks := fetchTasks(absPath, module)

	variants := extractVariants(tasks)

	variant, err := determineVariant(variants, module, verbose, ios)

	if err != nil {
		return nil, "", lib.PlatformUnknown, err
	}

	gb := newGradleBuilder(module, variant)

	return gb, properties["name"], lib.PlatformAndroid, nil
}

//-----------------------------------------------------------------------------

func (gb *GradleBuilder) Build(basePath string, clean, verbose bool, ios *lib.IOStreams) (string, error) {
	target := gb.formatTarget()

	ios.Printf("\nDetecting module properties for %s…\n", target)

	properties, err := gb.detectModuleProperties(basePath)

	if err != nil {
		return "", err
	}

	ios.Printf("\nDetermining build path for %s…\n", target)

	buildPath, err := gb.determineBuildPath(properties)

	if err != nil {
		return "", err
	}

	ios.Printf("\nBuilding %s…\n", target)

	dashes := "\n" + strings.Repeat("-", 79) + "\n"

	ios.Println(dashes)

	if err = gb.build(basePath, clean, verbose, ios); err != nil {
		return "", err
	}

	ios.Println(dashes)

	ios.Printf("\nVerifying build path for %s…\n", target)

	buildPaths, err := gb.verifyBuildPath(buildPath)

	if err != nil {
		return "", err
	}

	// if len(buildPaths) > 1 {
	// }

	return buildPaths[0], nil // for now, arbitrarily always take first path
}

func (gb *GradleBuilder) Summarize() string {
	summary := ""

	lib.AppendIfNotEmpty(&summary, "module", gb.Module, "=", ", ")
	lib.AppendIfNotEmpty(&summary, "variant", gb.Variant, "=", ", ")

	return summary
}

//-----------------------------------------------------------------------------

func commonGradleArgs() []string {
	return []string{"--console=plain", "--quiet"}
}

func newGradleBuilder(module, variant string) *GradleBuilder {
	return &GradleBuilder{
		Module:  module,
		Variant: variant}
}

func wrapperName() string {
	switch runtime.GOOS {
	case "windows":
		return "gradlew.bat"

	default:
		return "gradlew"
	}
}

//-----------------------------------------------------------------------------

func (gb *GradleBuilder) build(basePath string, clean, verbose bool, ios *lib.IOStreams) error {
	wrapperPath := filepath.Join(basePath, wrapperName())

	args := []string{}

	if clean {
		args = append(args, "clean")
	}

	taskName := fmt.Sprintf("%s:assemble%s", gb.Module, strings.Title(gb.Variant))

	args = append(args, taskName)

	args = append(args, "--console=plain")

	if !verbose {
		args = append(args, "--quiet")
	}

	task := lib.NewTask(wrapperPath, args...)

	task.Cwd = basePath
	task.IOStreams = ios

	return task.Execute()
}

func (gb *GradleBuilder) detectModuleProperties(basePath string) (map[string]string, error) {
	return fetchProperties(basePath, gb.Module), nil
}

func (gb *GradleBuilder) determineBuildPath(properties map[string]string) (string, error) {
	buildDir := properties["buildDir"]

	if len(buildDir) == 0 {
		return "", errors.New("Unable to determine build path")
	}

	return filepath.Join(buildDir, "outputs", "apk"), nil
}

func (gb *GradleBuilder) formatTarget() string {
	result := gb.Module

	lib.AppendIfNotEmpty(&result, "variant", gb.Variant, ": ", ", ")

	return result
}

func (gb *GradleBuilder) isPossibleBuildArtifact(path, basePath string) bool {
	reldir := filepath.Dir(path)[len(basePath):]
	variant := strings.ReplaceAll(reldir, "/", "")

	return strings.EqualFold(variant, gb.Variant)
}

func (gb *GradleBuilder) parseBuildSettings(text string) map[string]string {
	settings := make(map[string]string)

	for _, line := range strings.Split(text, "\n") {
		pair := strings.SplitN(line, "=", 2)

		if len(pair) != 2 {
			continue
		}

		key := strings.TrimSpace(pair[0])
		value := strings.TrimSpace(pair[1])

		if len(key) > 0 && len(value) > 0 {
			settings[key] = value
		}
	}

	return settings
}

func (gb *GradleBuilder) verifyBuildPath(basePath string) ([]string, error) {
	var results []string

	err := filepath.WalkDir(basePath, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !lib.IsRegularFile(path) || !strings.HasSuffix(path, ".apk") {
			return nil
		}

		if gb.isPossibleBuildArtifact(path, basePath) {
			results = append(results, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("Unable to locate build path in %q", basePath)
	}

	return results, nil
}
