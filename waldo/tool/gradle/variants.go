package gradle

import (
	"path/filepath"
	"strings"

	"github.com/waldoapp/waldo-go-cli/lib"
)

func DetectBuildVariants(basePath, module string) []string {
	tasks := fetchTasks(basePath, module)

	return extractVariants(tasks)
}

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

func extractVariants(tasks []string) []string {
	candidates := make([]string, 0)

	for _, task := range tasks {
		if candidate := candidateVariantFromTask(task); len(candidate) > 0 {
			candidates = append(candidates, candidate)
		}
	}

	variants := lib.CompactMap(candidates, func(candidate string) bool {
		return !isAffix(candidate, candidates)
	})

	return lib.Map(variants, func(task string) string {
		return strings.ToLower(task[0:1]) + task[1:]
	})
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
