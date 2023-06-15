package expo

import (
	"errors"
	"fmt"

	"github.com/waldoapp/waldo-go-cli/lib"
)

type ExpoBuilder struct {
}

//-----------------------------------------------------------------------------

func IsPossibleExpoContainer(path string) bool {
	return false
}

func MakeExpoBuilder(absPath, relPath string, verbose bool, ios *lib.IOStreams) (*ExpoBuilder, string, lib.Platform, error) {
	return nil, "", lib.PlatformUnknown, errors.New("Don’t know how to make an Expo recipe")
}

//-----------------------------------------------------------------------------

func (eb *ExpoBuilder) Build(basePath string, platform lib.Platform, clean, verbose bool, ios *lib.IOStreams) (string, error) {
	return "", fmt.Errorf("Don’t know how to build this app with Expo!")
}

func (eb *ExpoBuilder) Summarize() string {
	return ""
}
