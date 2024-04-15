package gradle

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/waldoapp/waldo-go-cli/lib"
)

type Builder struct {
	Module  string `yaml:"module,omitempty"`
	Variant string `yaml:"variant,omitempty"`
}

type BuildInfo struct {
	Variants []string
}

//-----------------------------------------------------------------------------

func DetectBuildInfo(basePath, module string) (*BuildInfo, error) {
	bi := &BuildInfo{}

	tasks := fetchTasks(basePath, module)

	bi.Variants = extractVariants(tasks)

	return bi, nil
}

func IsPossibleContainer(path string) (bool, bool) {
	wrapperPath := filepath.Join(path, wrapperName())
	kotlinPath := filepath.Join(path, "build.gradle.kts")
	groovyPath := filepath.Join(path, "build.gradle")

	if !lib.IsRegularFile(kotlinPath) && !lib.IsRegularFile(groovyPath) {
		return false, false
	}

	return lib.IsRegularFile(wrapperPath), false
}

func MakeBuilder(basePath string, verbose bool, ios *lib.IOStreams) (*Builder, string, lib.Platform, error) {
	ios.Printf("\nFinding all Gradle modules in %q\n", lib.MakeRelativeToCWD(basePath))

	properties := fetchProperties(basePath, "", ios)

	modules := extractModules(properties["subprojects"])

	module, err := DetermineModule(modules, verbose, ios)

	if err != nil {
		return nil, "", lib.PlatformUnknown, err
	}

	ios.Printf("\nFinding all Gradle build variants in %q\n", module)

	bi, err := DetectBuildInfo(basePath, module)

	if err != nil {
		return nil, "", lib.PlatformUnknown, err
	}

	variant, err := DetermineVariant(module, bi.Variants, verbose, ios)

	if err != nil {
		return nil, "", lib.PlatformUnknown, err
	}

	b := &Builder{
		Module:  module,
		Variant: variant}

	return b, properties["name"], lib.PlatformAndroid, nil
}

//-----------------------------------------------------------------------------

func (b *Builder) Build(basePath string, clean, verbose bool, ios *lib.IOStreams) (string, error) {
	target := b.formatTarget()

	ios.Printf("\nDetecting module properties for %v\n", target)

	properties, err := b.detectModuleProperties(basePath, ios)

	if err != nil {
		return "", err
	}

	ios.Printf("\nDetermining build path for %v\n", target)

	buildPath, err := b.determineBuildPath(properties)

	if err != nil {
		return "", err
	}

	ios.Printf("\nBuilding %v\n", target)

	dashes := "\n" + strings.Repeat("-", 79) + "\n"

	ios.Println(dashes)

	if err = b.build(basePath, clean, verbose, ios); err != nil {
		return "", err
	}

	ios.Println(dashes)

	ios.Printf("\nVerifying build path for %v\n", target)

	return b.verifyBuildPath(buildPath)
}

func (b *Builder) Clean(basePath string, verbose bool, ios *lib.IOStreams) error {
	ios.Printf("\nCleaning %v\n", b.formatTarget())

	return b.clean(basePath, verbose, ios)
}

func (b *Builder) DetermineBuildPath(basePath string, ios *lib.IOStreams) (string, error) {
	target := b.formatTarget()

	ios.Printf("\nDetecting module properties for %v\n", target)

	properties, err := b.detectModuleProperties(basePath, ios)

	if err != nil {
		return "", err
	}

	ios.Printf("\nDetermining build path for %v\n", target)

	return b.determineBuildPath(properties)
}

func (b *Builder) Summarize() string {
	summary := ""

	lib.AppendIfNotEmpty(&summary, "module", b.Module, "=", ", ")
	lib.AppendIfNotEmpty(&summary, "variant", b.Variant, "=", ", ")

	return summary
}

func (b *Builder) VerifyBuildPath(basePath string, ios *lib.IOStreams) (string, error) {
	ios.Printf("\nVerifying build path for %v\n", b.formatTarget())

	return b.verifyBuildPath(basePath)
}

//-----------------------------------------------------------------------------

var (
	projectRE = regexp.MustCompile(`^project ':(.+)'$`)
)

//-----------------------------------------------------------------------------

func candidateVariantFromTask(task string) string {
	if !strings.HasPrefix(task, "assemble") || strings.HasSuffix(task, "Test") {
		return ""
	}

	return task[len("assemble"):]
}

func commonArgs(verbose bool) []string {
	args := []string{"--console=plain"}

	if !verbose {
		args = append(args, "--quiet")
	}

	return args
}

func extractModules(text string) []string {
	if !strings.HasPrefix(text, "[") || !strings.HasSuffix(text, "]") {
		return nil
	}

	modules := []string{}

	for _, project := range strings.Split(text[1:len(text)-1], ", ") {
		matches := projectRE.FindStringSubmatch(project)

		if len(matches) == 2 && len(matches[1]) > 0 {
			modules = append(modules, matches[1])
		}
	}

	return modules
}

func extractVariants(tasks []string) []string {
	candidates := make([]string, 0)

	for _, task := range tasks {
		if candidate := candidateVariantFromTask(task); len(candidate) > 0 {
			candidates = append(candidates, candidate)
		}
	}

	return lib.CompactMap(candidates, func(candidate string) (string, bool) {
		return strings.ToLower(candidate[0:1]) + candidate[1:], !isAffix(candidate, candidates)
	})
}

func fetchProperties(basePath, module string, ios *lib.IOStreams) map[string]string {
	wrapperPath := filepath.Join(basePath, wrapperName())

	verb := "properties"

	if len(module) > 0 {
		verb = module + ":" + verb
	}

	args := append([]string{verb}, commonArgs(true)...)

	task := lib.NewTask(wrapperPath, args...)

	task.Cwd = basePath
	task.IOStreams = ios

	stdout, _, err := task.Run()

	if err != nil {
		return nil
	}

	return parseProperties(stdout)
}

func fetchTasks(basePath, module string) []string {
	wrapperPath := filepath.Join(basePath, wrapperName())

	args := append([]string{"tasks", "--all"}, commonArgs(true)...)

	task := lib.NewTask(wrapperPath, args...)

	task.Cwd = basePath

	stdout, _, err := task.Run()

	if err != nil {
		return nil
	}

	return parseTasks(stdout, module)
}

func isAffix(affix string, variants []string) bool {
	for _, variant := range variants {
		if variant == affix {
			continue
		}

		if strings.HasPrefix(variant, affix) || strings.HasSuffix(variant, affix) {
			return true
		}
	}

	return false
}

func parseProperties(text string) map[string]string {
	properties := make(map[string]string)

	for _, line := range strings.Split(text, "\n") {
		pair := strings.SplitN(line, ": ", 2)

		if len(pair) != 2 {
			continue
		}

		key := strings.TrimSpace(pair[0])
		value := strings.TrimSpace(pair[1])

		if len(key) > 0 && len(value) > 0 {
			properties[key] = value
		}
	}

	return properties
}

func parseTasks(text, module string) []string {
	tasks := make([]string, 0)

	prefix := module + ":"
	skip := len(prefix)

	for _, line := range strings.Split(text, "\n") {
		if !strings.HasPrefix(line, prefix) {
			continue
		}

		if idx := strings.Index(line, " "); idx >= 0 {
			line = line[skip:idx]
		} else {
			line = line[skip:]
		}

		if len(line) > 0 {
			tasks = append(tasks, line)
		}
	}

	return tasks
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

func (b *Builder) build(basePath string, clean, verbose bool, ios *lib.IOStreams) error {
	wrapperPath := filepath.Join(basePath, wrapperName())

	args := []string{}

	if clean {
		args = append(args, "clean")
	}

	taskName := fmt.Sprintf("%v:assemble%v", b.Module, strings.Title(b.Variant))

	args = append(args, taskName)

	args = append(args, commonArgs(verbose)...)

	task := lib.NewTask(wrapperPath, args...)

	task.Cwd = basePath
	task.IOStreams = ios

	return task.Execute()
}

func (b *Builder) clean(basePath string, verbose bool, ios *lib.IOStreams) error {
	wrapperPath := filepath.Join(basePath, wrapperName())

	args := []string{"clean"}

	args = append(args, commonArgs(verbose)...)

	task := lib.NewTask(wrapperPath, args...)

	task.Cwd = basePath
	task.IOStreams = ios

	return task.Execute()
}

func (b *Builder) detectModuleProperties(basePath string, ios *lib.IOStreams) (map[string]string, error) {
	return fetchProperties(basePath, b.Module, ios), nil
}

func (b *Builder) determineBuildPath(properties map[string]string) (string, error) {
	buildDir := properties["buildDir"]

	if len(buildDir) == 0 {
		return "", errors.New("Unable to determine build path")
	}

	return filepath.Join(buildDir, "outputs", "apk"), nil
}

func (b *Builder) formatTarget() string {
	result := "Gradle"

	lib.AppendIfNotEmpty(&result, "module", b.Module, ": ", ", ")
	lib.AppendIfNotEmpty(&result, "variant", b.Variant, ": ", ", ")

	return result
}

func (b *Builder) isPossibleBuildArtifact(path, basePath string) bool {
	reldir := filepath.Dir(path)[len(basePath):]
	variant := strings.ReplaceAll(reldir, "/", "")

	return strings.EqualFold(variant, b.Variant)
}

func (b *Builder) parseBuildSettings(text string) map[string]string {
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

func (b *Builder) verifyBuildPath(basePath string) (string, error) {
	var paths []string

	err := filepath.WalkDir(basePath, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !lib.IsRegularFile(path) || !strings.HasSuffix(path, ".apk") {
			return nil
		}

		if b.isPossibleBuildArtifact(path, basePath) {
			paths = append(paths, path)
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if len(paths) == 0 {
		return "", fmt.Errorf("Unable to locate build path in %q", basePath)
	}

	return paths[0], nil // for now, arbitrarily always take first path
}
