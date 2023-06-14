package gradle

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/waldoapp/waldo-go-cli/lib"
)

var (
	projectRE = regexp.MustCompile(`^project ':(.+)'$`)
)

//-----------------------------------------------------------------------------

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

func fetchProperties(basePath, module string) map[string]string {
	wrapperPath := filepath.Join(basePath, wrapperName())

	verb := "properties"

	if len(module) > 0 {
		verb = module + ":" + verb
	}

	args := append([]string{verb}, commonGradleArgs()...)

	task := lib.NewTask(wrapperPath, args...)

	task.Cwd = basePath

	stdout, _, err := task.Run()

	if err != nil {
		return nil
	}

	return parseProperties(stdout)
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
