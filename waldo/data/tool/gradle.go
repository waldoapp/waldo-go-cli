package tool

import (
	"errors"
	"fmt"

	"github.com/waldoapp/waldo-go-cli/lib"
)

type GradleBuilder struct {
	Module  string `yaml:"module,omitempty"`
	Variant string `yaml:"variant,omitempty"`
}

//-----------------------------------------------------------------------------

func FindGradlePaths(path string) []string {
	return make([]string, 0)
}

func IsPossibleGradleContainer(path string) bool {
	return false
}

func MakeGradleBuilder(absPath, relPath string, verbose bool, ios *lib.IOStreams) (*GradleBuilder, string, string, error) {
	return nil, "", "", errors.New("Don’t know how to make a Gradle recipe")
}

func NewGradleBuilder(module, variant string) *GradleBuilder {
	return &GradleBuilder{
		Module:  module,
		Variant: variant}
}

//-----------------------------------------------------------------------------

func (gb *GradleBuilder) Build(basePath string, clean, verbose bool, ios *lib.IOStreams) (*ArtifactMetadata, error) {
	return nil, fmt.Errorf("Don’t know how to build this app with Gradle!")
}

func (gb *GradleBuilder) Summarize() string {
	summary := ""

	if len(gb.Module) > 0 {
		if len(summary) > 0 {
			summary += ", "
		}

		summary += "module=" + gb.Module

	}

	if len(gb.Variant) > 0 {
		if len(summary) > 0 {
			summary += ", "
		}

		summary += "variant=" + gb.Variant
	}

	return summary
}
