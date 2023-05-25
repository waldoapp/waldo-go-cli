package tool

import (
	"errors"
	"fmt"

	"github.com/waldoapp/waldo-go-cli/lib"
)

type ExpoBuilder struct {
}

//-----------------------------------------------------------------------------

func FindExpoPaths(path string) []string {
	return make([]string, 0)
}

func IsPossibleExpoContainer(path string) bool {
	return false
}

func MakeExpoBuilder(absPath, relPath string, verbose bool, ios *lib.IOStreams) (*ExpoBuilder, string, string, error) {
	return nil, "", "", errors.New("Don’t know how to make an Expo recipe")
}

func NewExpoBuilder() *ExpoBuilder {
	return &ExpoBuilder{}
}

//-----------------------------------------------------------------------------

func (eb *ExpoBuilder) Build(basePath string, clean, verbose bool, ios *lib.IOStreams) (*ArtifactMetadata, error) {
	return nil, fmt.Errorf("Don’t know how to build this app with Expo!")
}

func (eb *ExpoBuilder) Summarize() string {
	return ""
}
