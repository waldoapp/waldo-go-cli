package data

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/waldoapp/waldo-go-cli/lib"

	"gopkg.in/yaml.v3"
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

		saveNeeded := false

		if err == nil {
			saveNeeded, err = cfg.migrate()
		}

		if err == nil && saveNeeded {
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

	oldRecipes := cfg.Recipes

	cfg.Recipes = append(cfg.Recipes, recipe)

	if err := cfg.Save(); err != nil {
		cfg.Recipes = oldRecipes

		return err
	}

	return nil
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

func (cfg *Configuration) Path() string {
	return cfg.configPath
}

func (cfg *Configuration) RemoveRecipe(name string) error {
	idx := cfg.findRecipeIndex(name)

	if idx >= len(cfg.Recipes) {
		return fmt.Errorf("Unable to find recipe: %q", name)
	}

	oldRecipes := cfg.Recipes

	cfg.Recipes = append(cfg.Recipes[:idx], cfg.Recipes[idx+1:]...)

	if err := cfg.Save(); err != nil {
		cfg.Recipes = oldRecipes

		return err
	}

	return nil
}

func (cfg *Configuration) Save() error {
	data, err := yaml.Marshal(cfg)

	if err != nil {
		return nil
	}

	return os.WriteFile(cfg.configPath, data, 0644)
}

//-----------------------------------------------------------------------------

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
		err = yaml.Unmarshal(data, cfg)
	}

	return err
}

func (cfg *Configuration) migrate() (bool, error) {
	saveNeeded := false

	//
	// Ensure there is ALWAYS a unique ID associated with this configuration:
	//
	if len(cfg.UniqueID) == 0 {
		uniqueID, err := lib.NewUniqueID()

		if err != nil {
			return false, err
		}

		cfg.UniqueID = uniqueID

		saveNeeded = true
	}

	if cfg.FormatVersion < cfgFormatVersion {
		// handle any necessary migration here

		cfg.FormatVersion = cfgFormatVersion

		saveNeeded = true
	} else if cfg.FormatVersion > cfgFormatVersion {
		// may not understand newer format
	}

	return saveNeeded, nil
}

func (cfg *Configuration) populate() error {
	uniqueID, err := lib.NewUniqueID()

	if err != nil {
		return err
	}

	cfg.UniqueID = uniqueID
	cfg.FormatVersion = cfgFormatVersion

	return nil
}
