package data

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/lib/tpw"
)

const (
	cfgFormatVersion = 1
)

//-----------------------------------------------------------------------------

type Configuration struct {
	UniqueID      string    `yaml:"unique_id"`
	FormatVersion int       `yaml:"format_version"`
	Recipes       []*Recipe `yaml:"recipes,omitempty"`

	basePath   string // absolute
	configPath string // absolute
	dirty      bool
}

//-----------------------------------------------------------------------------

func SetupConfiguration(create bool) (*Configuration, bool, error) {
	dataPath, err := FindRepoSpecificPath() // /path/to/.waldo

	if err != nil {
		return nil, false, err
	}

	cwd, err := os.Getwd()

	if err != nil {
		return nil, false, err
	}

	if strings.HasSuffix(cwd, "/.waldo") ||
		strings.Contains(cwd, "/.waldo/") {
		return nil, false, errors.New("This operation must be run outside of the .waldo tree")
	}

	cfg := &Configuration{
		basePath:   filepath.Dir(dataPath),
		configPath: filepath.Join(dataPath, "config.yml")}

	if lib.IsRegularFile(cfg.configPath) {
		err = cfg.load()

		if err == nil {
			err = cfg.migrate()
		}

		if err == nil {
			err = cfg.Save()
		}

		if err != nil {
			return nil, false, err
		}

		return cfg, false, nil
	}

	if create {
		err = os.MkdirAll(dataPath, 0755)

		if err != nil {
			return nil, false, err
		}

		err = cfg.populate()

		if err == nil {
			err = cfg.Save()
		}

		if err != nil {
			return nil, false, err
		}

		return cfg, true, nil
	}

	return nil, false, errors.New("Waldo configuration not found")
}

//-----------------------------------------------------------------------------

func (cfg *Configuration) AddRecipe(recipe *Recipe) error {
	idx := cfg.findRecipeIndex(recipe.Name)

	if idx < len(cfg.Recipes) {
		return fmt.Errorf("Recipe already added: %q", recipe.Name)
	}

	cfg.Recipes = append(cfg.Recipes, recipe)

	cfg.MarkDirty()

	return cfg.Save()
}

func (cfg *Configuration) BasePath() string {
	return cfg.basePath
}

func (cfg *Configuration) FindRecipe(name string) (*Recipe, error) {
	if len(name) == 0 {
		cnt := len(cfg.Recipes)

		switch {
		case cnt == 1:
			return cfg.Recipes[0], nil

		case cnt > 1:
			return nil, errors.New("Empty recipe name")

		default:
			return nil, errors.New("No recipes defined")
		}
	}

	idx := cfg.findRecipeIndex(name)

	if idx >= len(cfg.Recipes) {
		return nil, fmt.Errorf("Unable to find recipe: %q", name)
	}

	return cfg.Recipes[idx], nil
}

func (cfg *Configuration) IsDirty() bool {
	return cfg.dirty
}

func (cfg *Configuration) MarkDirty() {
	cfg.dirty = true
}

func (cfg *Configuration) Path() string {
	return cfg.configPath
}

func (cfg *Configuration) RemoveRecipe(name string) error {
	idx := cfg.findRecipeIndex(name)

	if idx >= len(cfg.Recipes) {
		return fmt.Errorf("Unable to find recipe: %q", name)
	}

	cfg.Recipes = append(cfg.Recipes[:idx], cfg.Recipes[idx+1:]...)

	cfg.MarkDirty()

	return cfg.Save()
}

func (cfg *Configuration) Save() error {
	if !cfg.IsDirty() {
		return nil
	}

	data, err := tpw.EncodeToYAML(cfg)

	if err == nil {
		err = os.WriteFile(cfg.configPath, data, 0644)
	}

	if err == nil {
		cfg.dirty = false
	}

	return err
}

//-----------------------------------------------------------------------------

func (cfg *Configuration) ensureUniqueID() error {
	//
	// Ensure there is ALWAYS a unique ID associated with this configuration:
	//
	if len(cfg.UniqueID) == 0 {
		uniqueID, err := tpw.NewUniqueID()

		if err != nil {
			return err
		}

		cfg.UniqueID = uniqueID

		cfg.MarkDirty()
	}

	return nil
}

func (cfg *Configuration) findRecipeIndex(name string) int {
	for idx, recipe := range cfg.Recipes {
		if recipe.Name == name {
			return idx
		}
	}

	return len(cfg.Recipes)
}

func (cfg *Configuration) load() error {
	data, err := os.ReadFile(cfg.configPath)

	if err == nil {
		err = tpw.DecodeFromYAML(data, cfg)
	}

	return err
}

func (cfg *Configuration) migrate() error {
	if cfg.FormatVersion < cfgFormatVersion {
		// handle any necessary migration here

		cfg.FormatVersion = cfgFormatVersion

		cfg.MarkDirty()
	} else if cfg.FormatVersion > cfgFormatVersion {
		// may not understand newer format
	}

	return cfg.ensureUniqueID()
}

func (cfg *Configuration) populate() error {
	cfg.FormatVersion = cfgFormatVersion

	cfg.MarkDirty()

	return cfg.ensureUniqueID()
}