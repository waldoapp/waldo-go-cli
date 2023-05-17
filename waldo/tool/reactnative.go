package tool

import (
	"errors"
	"fmt"

	"github.com/waldoapp/waldo-go-cli/lib"
)

type ReactNativeBuilder struct {
}

//-----------------------------------------------------------------------------

func FindReactNativePaths(path string) []string {
	return make([]string, 0)
}

func IsPossibleReactNativeContainer(path string) bool {
	return false
}

func MakeReactNativeBuilder(absPath, relPath string, verbose bool, ios *lib.IOStreams) (*ReactNativeBuilder, error) {
	return nil, errors.New("Don’t know how to make a React Native recipe")
}

func NewReactNativeBuilder() *ReactNativeBuilder {
	return &ReactNativeBuilder{}
}

//-----------------------------------------------------------------------------

func (rnb *ReactNativeBuilder) Build(basePath string, verbose bool, ios *lib.IOStreams) error {
	return fmt.Errorf("Don’t know how to build this app with React Native!")
}

func (rnb *ReactNativeBuilder) Summarize() string {
	return ""
}
