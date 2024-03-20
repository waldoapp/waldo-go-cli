package data

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	appIDRE      = regexp.MustCompile(`^app-[0-9a-fA-F]+$`)
	appTokenRE   = regexp.MustCompile(`^[0-9a-fA-F]+$`)
	recipeNameRE = regexp.MustCompile(`^[a-zA-Z][0-9a-zA-Z_-]*$`)
	userTokenRE  = regexp.MustCompile(`^u-[0-9a-fA-F]+$`)
)

func ValidateAppID(id string) error {
	if len(id) == 0 {
		return errors.New("No app id specified")
	}

	if !appIDRE.MatchString(id) {
		return fmt.Errorf("Invalid app id syntax: %q", id)
	}

	return nil
}

func ValidateAppToken(token string) error {
	if len(token) == 0 {
		return errors.New("No upload token specified") // refer to it as _upload_ token for legacy purposes
	}

	if !appTokenRE.MatchString(token) {
		return fmt.Errorf("Invalid upload token syntax: %q", token) // refer to it as _upload_ token for legacy purposes
	}

	return nil
}

func ValidateRecipeName(name string) error {
	if len(name) == 0 {
		return errors.New("No recipe name specified")
	}

	if !recipeNameRE.MatchString(name) {
		return fmt.Errorf("Invalid recipe name syntax: %q", name)
	}

	return nil
}

func ValidateUserToken(token string) error {
	if len(token) == 0 {
		return errors.New("No user token specified")
	}

	if !userTokenRE.MatchString(token) {
		return fmt.Errorf("Invalid user token syntax: %q", token)
	}

	return nil
}
