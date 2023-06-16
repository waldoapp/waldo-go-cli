package custom

import (
	"errors"
	"fmt"

	"github.com/waldoapp/waldo-go-cli/lib"
)

type CustomBuilder struct {
}

//-----------------------------------------------------------------------------

func MakeCustomBuilder(absPath, relPath string, verbose bool, ios *lib.IOStreams) (*CustomBuilder, string, lib.Platform, error) {
	return nil, "", lib.PlatformUnknown, errors.New("Don’t know how to make a custom recipe")
}

//-----------------------------------------------------------------------------

func (cb *CustomBuilder) Build(basePath string, platform lib.Platform, clean, verbose bool, ios *lib.IOStreams) (string, error) {
	return "", fmt.Errorf("Don’t know how to custom build this app!")
}

func (cb *CustomBuilder) Summarize() string {
	return ""
}
