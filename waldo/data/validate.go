package data

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	recipeNameRE  = regexp.MustCompile(`^[a-zA-Z][0-9a-zA-Z_-]*$`)
	uploadTokenRE = regexp.MustCompile(`^[0-9a-fA-F]{32}$`)
)

func ValidateRecipeName(name string) error {
	if len(name) == 0 {
		return errors.New("Empty recipe name")
	}

	if recipeNameRE.FindString(name) != name {
		return fmt.Errorf("Invalid recipe name syntax: %q", name)
	}

	return nil
}

func ValidateUploadToken(token string) error {
	if len(token) == 0 {
		return errors.New("Empty upload token")
	}

	if uploadTokenRE.FindString(token) != token {
		return fmt.Errorf("Invalid upload token syntax: %q", token)
	}

	return nil
}
