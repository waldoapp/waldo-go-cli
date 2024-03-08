package data

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	appTokenRE   = regexp.MustCompile(`^[0-9a-fA-F]+$`)
	recipeNameRE = regexp.MustCompile(`^[a-zA-Z][0-9a-zA-Z_-]*$`)
	userTokenRE  = regexp.MustCompile(`^u-[0-9a-fA-F]+$`)
)

func ValidateAppToken(token string) error {
	if len(token) == 0 {
		return errors.New("Empty app token")
	}

	if !appTokenRE.MatchString(token) {
		return fmt.Errorf("Invalid app token syntax: %q", token)
	}

	return nil
}

func ValidateRecipeName(name string) error {
	if len(name) == 0 {
		return errors.New("Empty recipe name")
	}

	if !recipeNameRE.MatchString(name) {
		return fmt.Errorf("Invalid recipe name syntax: %q", name)
	}

	return nil
}

func ValidateUploadToken(token string) error {
	if len(token) == 0 {
		return errors.New("Empty upload token")
	}

	if !appTokenRE.MatchString(token) && !userTokenRE.MatchString(token) {
		return fmt.Errorf("Invalid upload token syntax: %q", token)
	}

	return nil
}

func ValidateUserToken(token string) error {
	if len(token) == 0 {
		return errors.New("Empty user token")
	}

	if !userTokenRE.MatchString(token) {
		return fmt.Errorf("Invalid user token syntax: %q", token)
	}

	return nil
}
