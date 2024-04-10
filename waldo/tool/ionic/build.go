package ionic

import (
	"errors"
	"fmt"

	"github.com/waldoapp/waldo-go-cli/lib"
)

type Builder struct {
}

//-----------------------------------------------------------------------------

func IsPossibleContainer(path string) bool {
	return false
}

func MakeBuilder(basePath string, verbose bool, ios *lib.IOStreams) (*Builder, string, lib.Platform, error) {
	return nil, "", lib.PlatformUnknown, errors.New("Don’t know how to make an Ionic recipe")
}

//-----------------------------------------------------------------------------

func (b *Builder) Build(basePath string, platform lib.Platform, clean, verbose bool, ios *lib.IOStreams) (string, error) {
	return "", fmt.Errorf("Don’t know how to build this app with Ionic!")
}

func (b *Builder) Summarize() string {
	return ""
}
