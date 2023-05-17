package tool

import (
	"errors"
	"fmt"

	"github.com/waldoapp/waldo-go-cli/lib"
)

type CustomBuilder struct {
}

//-----------------------------------------------------------------------------

func MakeCustomBuilder(absPath, relPath string, verbose bool, ios *lib.IOStreams) (*CustomBuilder, error) {
	return nil, errors.New("Don’t know how to make a custom recipe")
}

func NewCustomBuilder() *CustomBuilder {
	return &CustomBuilder{}
}

//-----------------------------------------------------------------------------

func (cb *CustomBuilder) Build(basePath string, verbose bool, ios *lib.IOStreams) error {
	return fmt.Errorf("Don’t know how to custom build this app!")
}

func (cb *CustomBuilder) Summarize() string {
	return ""
}
