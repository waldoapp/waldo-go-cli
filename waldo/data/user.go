package data

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/data/tool"
)

const (
	udFormatVersion = 1
)

//-----------------------------------------------------------------------------

type UserData struct {
	FormatVersion int                               `json:"formatVersion"`
	MetadataMap   map[string]*tool.ArtifactMetadata `json:"metadata"` // key is <recipe-name>

	basePath string // absolute
	dirty    bool
	userPath string // absolute
}

//-----------------------------------------------------------------------------

func SetupUserData(cfg *Configuration) *UserData {
	dataPath, err := makeUserSpecificDataPath()

	if err != nil {
		return nil
	}

	fileName := cfg.UniqueID + ".json"

	ud := &UserData{
		basePath: filepath.Dir(dataPath),
		userPath: filepath.Join(dataPath, fileName)}

	if lib.IsRegularFile(ud.userPath) {
		if err := ud.load(); err != nil {
			return nil
		}

		if err := ud.migrate(); err != nil {
			return nil
		}
	} else {
		if err := os.MkdirAll(dataPath, 0700); err != nil {
			return nil
		}

		ud.FormatVersion = udFormatVersion
		ud.MetadataMap = make(map[string]*tool.ArtifactMetadata)

		ud.MarkDirty()
	}

	if err := ud.Save(); err != nil {
		return nil
	}

	return ud
}

//-----------------------------------------------------------------------------

func (ud *UserData) AddMetadata(r *Recipe, am *tool.ArtifactMetadata) error {
	if am == nil {
		return nil
	}

	ud.MetadataMap[r.Name] = am

	ud.MarkDirty()

	return ud.Save()
}

func (ud *UserData) BasePath() string {
	return ud.basePath
}

func (ud *UserData) FindMetadata(r *Recipe) (*tool.ArtifactMetadata, error) {
	if am := ud.MetadataMap[r.Name]; am != nil {
		return am, nil
	}

	return nil, fmt.Errorf("Unable to find metadata for recipe: %q", r.Name)
}

func (ud *UserData) IsDirty() bool {
	return ud.dirty
}

func (ud *UserData) MarkDirty() {
	ud.dirty = true
}

func (ud *UserData) Path() string {
	return ud.userPath
}

func (ud *UserData) RemoveMetadata(r *Recipe) error {
	delete(ud.MetadataMap, r.Name)

	ud.MarkDirty()

	return ud.Save()
}

func (ud *UserData) Save() error {
	if !ud.IsDirty() {
		return nil
	}

	data, err := json.Marshal(ud)

	if err != nil {
		return err
	}

	if err := os.WriteFile(ud.userPath, data, 0600); err != nil {
		return err
	}

	ud.dirty = false

	return nil
}

//-----------------------------------------------------------------------------

func makeUserSpecificDataPath() (string, error) {
	path, err := os.UserConfigDir()

	if err != nil {
		return "", err
	}

	return filepath.Join(path, "waldo"), nil
}

//-----------------------------------------------------------------------------

func (ud *UserData) load() error {
	data, err := os.ReadFile(ud.userPath)

	if err != nil {
		return err
	}

	return json.Unmarshal(data, ud)
}

func (ud *UserData) migrate() error {
	if ud.FormatVersion < udFormatVersion {
		// handle any necessary migration here

		ud.FormatVersion = udFormatVersion

		ud.MarkDirty()
	} else if ud.FormatVersion > udFormatVersion {
		// may not understand newer format
	}

	if ud.MetadataMap == nil {
		ud.MetadataMap = make(map[string]*tool.ArtifactMetadata)

		ud.MarkDirty()
	}

	return nil
}
