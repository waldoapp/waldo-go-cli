package tool

import (
	"errors"
	"fmt"

	"github.com/waldoapp/waldo-go-cli/lib"
)

type FlutterBuilder struct {
}

//-----------------------------------------------------------------------------

func FindFlutterPaths(path string) []string {
	return make([]string, 0)
}

func IsPossibleFlutterContainer(path string) bool {
	return false
}

func MakeFlutterBuilder(absPath, relPath string, verbose bool, ios *lib.IOStreams) (*FlutterBuilder, string, string, error) {
	return nil, "", "", errors.New("Don’t know how to make a Flutter recipe")
}

func NewFlutterBuilder() *FlutterBuilder {
	return &FlutterBuilder{}
}

//-----------------------------------------------------------------------------

func (fb *FlutterBuilder) Build(basePath string, clean, verbose bool, ios *lib.IOStreams) (*ArtifactMetadata, error) {
	return nil, fmt.Errorf("Don’t know how to build this app with Flutter!")
}

func (fb *FlutterBuilder) Summarize() string {
	return ""
}
