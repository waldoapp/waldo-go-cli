package gradle

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/waldoapp/waldo-go-cli/lib"
)

type GradleInfo struct {
	Variants []string
}

//-----------------------------------------------------------------------------

func DetectGradleInfo(basePath, module string) (*GradleInfo, error) {
	gi := &GradleInfo{}

	tasks := fetchTasks(basePath, module)

	gi.Variants = extractVariants(tasks)

	return gi, nil
}

//-----------------------------------------------------------------------------

var (
	projectRE = regexp.MustCompile(`^project ':(.+)'$`)
)

//-----------------------------------------------------------------------------

func candidateVariantFromTask(task string) string {
	if !strings.HasPrefix(task, "assemble") {
		return ""
	}

	if strings.HasSuffix(task, "Test") {
		return ""
	}

	return task[len("assemble"):]
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

	args := append([]string{verb}, commonGradleArgs()...)

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

	args := append([]string{"tasks", "--all"}, commonGradleArgs()...)

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

		if strings.HasPrefix(variant, affix) {
			return true
		}

		if strings.HasSuffix(variant, affix) {
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
