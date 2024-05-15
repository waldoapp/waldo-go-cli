package data

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/lib/tpw"
)

const (
	prfFormatVersion = 1
)

//-----------------------------------------------------------------------------

type Profile struct {
	FormatVersion int    `yaml:"format_version"`
	APIToken      string `yaml:"user_token,omitempty"`

	basePath    string // absolute
	dirty       bool
	profilePath string // absolute
}

//-----------------------------------------------------------------------------

func SetupProfile(ck CreateKind) (*Profile, bool, error) {
	var (
		dataPath string
		create   bool
		err      error
	)

	switch ck {
	case CreateKindAlways:
		dataPath, err = makeHomeDataPath()

		if err != nil {
			return nil, false, err
		}

		create = true

	case CreateKindIfNeeded:
		dataPath = findHomeDataPath()

		if len(dataPath) == 0 {
			dataPath, err = makeHomeDataPath()

			if err != nil {
				return nil, false, err
			}

			create = true
		}

	default: // incl. CreateKindNever
		dataPath = findHomeDataPath()

		if len(dataPath) == 0 {
			return nil, false, errors.New("Waldo profile not found")
		}
	}

	prf := &Profile{
		basePath:    filepath.Dir(dataPath),
		profilePath: filepath.Join(dataPath, "profile.yml")}

	if lib.IsRegularFile(prf.profilePath) {
		if err := prf.load(); err != nil {
			return nil, false, err
		}

		if err := prf.migrate(); err != nil {
			return nil, false, err
		}

		if err := prf.Save(); err != nil {
			return nil, false, err
		}

		return prf, false, nil
	}

	if create {
		if err := os.MkdirAll(dataPath, 0755); err != nil {
			return nil, false, err
		}

		prf.FormatVersion = prfFormatVersion

		prf.MarkDirty()

		if err := prf.Save(); err != nil {
			return nil, false, err
		}

		return prf, true, nil
	}

	return nil, false, errors.New("Waldo profile not found")
}

//-----------------------------------------------------------------------------

func (prf *Profile) BasePath() string {
	return prf.basePath
}

func (prf *Profile) IsDirty() bool {
	return prf.dirty
}

func (prf *Profile) MarkDirty() {
	prf.dirty = true
}

func (prf *Profile) Path() string {
	return prf.profilePath
}

func (prf *Profile) Save() error {
	if !prf.IsDirty() {
		return nil
	}

	data, err := tpw.EncodeToYAML(prf)

	if err != nil {
		return err
	}

	if err := os.WriteFile(prf.profilePath, data, 0644); err != nil {
		return err
	}

	prf.dirty = false

	return nil
}

//-----------------------------------------------------------------------------

func findHomeDataPath() string {
	dirPath, err := os.UserHomeDir()

	if err != nil {
		return ""
	}

	dataPath := filepath.Join(dirPath, ".waldo")
	prfPath := filepath.Join(dirPath, ".waldo", "profile.yml")

	if lib.IsRegularFile(prfPath) {
		return dataPath
	}

	return ""
}

func makeHomeDataPath() (string, error) {
	dirPath, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	return filepath.Join(dirPath, ".waldo"), nil
}

//-----------------------------------------------------------------------------

func (prf *Profile) load() error {
	data, err := os.ReadFile(prf.profilePath)

	if err != nil {
		return err
	}

	return tpw.DecodeFromYAML(data, prf)
}

func (prf *Profile) migrate() error {
	if prf.FormatVersion < prfFormatVersion {
		// handle any necessary migration here

		prf.FormatVersion = prfFormatVersion

		prf.MarkDirty()
	} else if prf.FormatVersion > prfFormatVersion {
		// may not understand newer format
	}

	return nil
}
