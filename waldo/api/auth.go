package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/data"
)

//-----------------------------------------------------------------------------

type AuthenticateUserResponse struct {
	Email     string `json:"email,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	UserID    string `json:"id"`
}

//-----------------------------------------------------------------------------

func AuthenticateUser(apiToken string, verbose bool, ios *lib.IOStreams) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", getAuthenticateUserEndpoint(), nil)

	if err != nil {
		return "", fmt.Errorf("Unable to authenticate user, error: %v", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Token %v", apiToken))
	req.Header.Add("User-Agent", data.FullVersion())

	if verbose {
		lib.DumpRequest(ios, req, true)
	}

	rsp, err := client.Do(req)

	if err != nil {
		return "", fmt.Errorf("Unable to authenticate user, error: %v", err)
	}

	defer rsp.Body.Close()

	if verbose {
		lib.DumpResponse(ios, rsp, true)
	}

	status := rsp.StatusCode

	if status < 200 || status > 299 {
		return "", fmt.Errorf("Unable to authenticate user, error: %v", rsp.Status)
	}

	aur, err := parseAuthenticateUserResponse(rsp)

	if err != nil {
		return "", fmt.Errorf("Unable to authenticate user, error: %v", err)
	}

	return aur.fullName(), nil
}

//-----------------------------------------------------------------------------

func parseAuthenticateUserResponse(rsp *http.Response) (*AuthenticateUserResponse, error) {
	data, err := io.ReadAll(rsp.Body)

	if err != nil {
		return nil, err
	}

	aur := &AuthenticateUserResponse{}

	if err = json.Unmarshal(data, aur); err != nil {
		return nil, err
	}

	return aur, nil
}

//-----------------------------------------------------------------------------

func (aur *AuthenticateUserResponse) fullName() string {
	if len(aur.FirstName) > 0 && len(aur.LastName) > 0 {
		return aur.FirstName + " " + aur.LastName
	}

	if len(aur.FirstName) > 0 {
		return aur.FirstName
	}

	if len(aur.LastName) > 0 {
		return aur.LastName
	}

	return ""
}
