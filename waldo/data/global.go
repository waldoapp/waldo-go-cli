package data

import (
	"os"
)

const (
	AgentName    = "Waldo CLI"
	AgentNameOld = "Waldo Agent"
	AgentPrefix  = "waldo"
	AgentVersion = "3.0.0"

	DefaultAPIBuildEndpoint   = "https://api.waldo.com/versions"
	DefaultAPIErrorEndpoint   = "https://api.waldo.com/uploadError"
	DefaultAPIInitEndpoint    = "https://wow.waldo.io/cliSetup" // not wow.waldo.com ???
	DefaultAPITriggerEndpoint = "https://api.waldo.com/suites"

	MaxPostAttempts = 2
)

func Overrides() map[string]string {
	overrides := map[string]string{}

	if override := os.Getenv("WALDO_API_BUILD_ENDPOINT_OVERRIDE"); len(override) > 0 {
		overrides["apiBuildEndpoint"] = override
	}

	if override := os.Getenv("WALDO_API_ERROR_ENDPOINT_OVERRIDE"); len(override) > 0 {
		overrides["apiErrorEndpoint"] = override
	}

	if override := os.Getenv("WALDO_API_INIT_ENDPOINT_OVERRIDE"); len(override) > 0 {
		overrides["apiInitEndpoint"] = override
	}

	if override := os.Getenv("WALDO_API_TRIGGER_ENDPOINT_OVERRIDE"); len(override) > 0 {
		overrides["apiTriggerEndpoint"] = override
	}

	if override := os.Getenv("WALDO_WRAPPER_NAME_OVERRIDE"); len(override) > 0 {
		overrides["wrapperName"] = override
	}

	if override := os.Getenv("WALDO_WRAPPER_VERSION_OVERRIDE"); len(override) > 0 {
		overrides["wrapperVersion"] = override
	}

	return overrides
}
