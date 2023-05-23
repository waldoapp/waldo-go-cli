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

func SetupUserData(cfg *Configuration) (*UserData, error) {
	dataPath, err := FindUserSpecificPath()

	if err != nil {
		return nil, err
	}

	fileName := cfg.UniqueID + ".json"

	ud := &UserData{
		basePath: filepath.Dir(dataPath),
		userPath: filepath.Join(dataPath, fileName)}

	if lib.IsRegularFile(ud.userPath) {
		err = ud.load()

		if err == nil {
			err = ud.migrate()
		}

		if err == nil {
			err = ud.Save()
		}

		if err != nil {
			return nil, err
		}

		return ud, nil
	}

	err = os.MkdirAll(dataPath, 0700)

	if err != nil {
		return nil, err
	}

	ud.FormatVersion = udFormatVersion
	ud.MetadataMap = make(map[string]*tool.ArtifactMetadata)

	ud.MarkDirty()

	err = ud.Save()

	if err != nil {
		return nil, err
	}

	return ud, nil
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
	am := ud.MetadataMap[r.Name]

	if am != nil {
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

	if err == nil {
		err = os.WriteFile(ud.userPath, data, 0600)
	}

	if err == nil {
		ud.dirty = false
	}

	return err
}

//-----------------------------------------------------------------------------

func (ud *UserData) load() error {
	data, err := os.ReadFile(ud.userPath)

	if err == nil {
		err = json.Unmarshal(data, ud)
	}

	return err
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
