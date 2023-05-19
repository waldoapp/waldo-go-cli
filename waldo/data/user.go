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

		saveNeeded := false

		if err == nil {
			saveNeeded, err = ud.migrate()
		}

		if err == nil && saveNeeded {
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

	if err := ud.Save(); err != nil {
		return err
	}

	return nil
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

func (ud *UserData) Path() string {
	return ud.userPath
}

func (ud *UserData) RemoveMetadata(r *Recipe) error {
	delete(ud.MetadataMap, r.Name)

	if err := ud.Save(); err != nil {
		return err
	}

	return nil
}

func (ud *UserData) Save() error {
	data, err := json.Marshal(ud)

	if err != nil {
		return nil
	}

	return os.WriteFile(ud.userPath, data, 0600)
}

//-----------------------------------------------------------------------------

func (ud *UserData) load() error {
	data, err := os.ReadFile(ud.userPath)

	if err == nil {
		err = json.Unmarshal(data, ud)
	}

	return err
}

func (ud *UserData) migrate() (bool, error) {
	saveNeeded := false

	if ud.FormatVersion < udFormatVersion {
		// handle any necessary migration here

		ud.FormatVersion = udFormatVersion

		saveNeeded = true
	} else if ud.FormatVersion > udFormatVersion {
		// may not understand newer format
	}

	if ud.MetadataMap == nil {
		ud.MetadataMap = make(map[string]*tool.ArtifactMetadata)

		saveNeeded = true // ???
	}

	return saveNeeded, nil
}
