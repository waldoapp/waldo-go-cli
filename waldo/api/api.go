package api

import (
	"os"
)

const (
	defaultAuthenticateUserEndpoint = "https://api.waldo.com/1.0/users/me"
	defaultFetchAppsEndpoint        = "https://api.waldo.com/1.0/applications"
)

func getAuthenticateUserEndpoint() string {
	if endpoint := os.Getenv("WALDO_API_AUTHENTICATE_USER_ENDPOINT_OVERRIDE"); len(endpoint) > 0 {
		return endpoint
	}

	return defaultAuthenticateUserEndpoint
}

func getFetchAppsEndpoint() string {
	if endpoint := os.Getenv("WALDO_API_FETCH_APPS_ENDPOINT_OVERRIDE"); len(endpoint) > 0 {
		return endpoint
	}

	return defaultFetchAppsEndpoint
}
