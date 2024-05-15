package data

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	apiTokenRE = regexp.MustCompile(`^u-[0-9a-fA-F]+$`)
	appIDRE    = regexp.MustCompile(`^app-[0-9a-fA-F]+$`)
	ciTokenRE  = regexp.MustCompile(`^[0-9a-fA-F]+$`)
)

func ValidateAPIToken(token string) error {
	if len(token) == 0 {
		return errors.New("No API token specified")
	}

	if !apiTokenRE.MatchString(token) {
		return fmt.Errorf("Invalid API token syntax: %q", token)
	}

	return nil
}

func ValidateAppID(id string) error {
	if len(id) == 0 {
		return errors.New("No app id specified")
	}

	if !appIDRE.MatchString(id) {
		return fmt.Errorf("Invalid app id syntax: %q", id)
	}

	return nil
}

func ValidateCIToken(token string) error {
	if len(token) == 0 {
		return errors.New("No CI token specified")
	}

	if !ciTokenRE.MatchString(token) {
		return fmt.Errorf("Invalid CI token syntax: %q", token)
	}

	return nil
}

func ValidateUploadToken(token string) error {
	if len(token) == 0 {
		return errors.New("No upload token specified")
	}

	if !apiTokenRE.MatchString(token) && !ciTokenRE.MatchString(token) {
		return fmt.Errorf("Invalid upload token syntax: %q", token)
	}

	return nil
}
