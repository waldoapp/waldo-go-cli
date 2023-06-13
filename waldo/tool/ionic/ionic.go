package ionic

import (
	"errors"
	"fmt"

	"github.com/waldoapp/waldo-go-cli/lib"
)

type IonicBuilder struct {
}

//-----------------------------------------------------------------------------

func FindIonicPaths(path string) []string {
	return make([]string, 0)
}

func IsPossibleIonicContainer(path string) bool {
	return false
}

func MakeIonicBuilder(absPath, relPath string, verbose bool, ios *lib.IOStreams) (*IonicBuilder, string, string, error) {
	return nil, "", "", errors.New("Don’t know how to make an Ionic recipe")
}

func NewIonicBuilder() *IonicBuilder {
	return &IonicBuilder{}
}

//-----------------------------------------------------------------------------

func (ib *IonicBuilder) Build(basePath string, platform lib.Platform, clean, verbose bool, ios *lib.IOStreams) (string, error) {
	return "", fmt.Errorf("Don’t know how to build this app with Ionic!")
}

func (ib *IonicBuilder) Summarize() string {
	return ""
}
