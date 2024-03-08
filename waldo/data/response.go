package data

import (
	"encoding/json"
	"io"
	"net/http"
)

//-----------------------------------------------------------------------------

type AuthResponse struct {
	Email     string `json:"email,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	UserID    string `json:"id"`
}

//-----------------------------------------------------------------------------

func ParseAuthResponse(resp *http.Response) (*AuthResponse, error) {
	data, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	ar := &AuthResponse{}

	if err = json.Unmarshal(data, ar); err != nil {
		return nil, err
	}

	return ar, nil
}

//-----------------------------------------------------------------------------

func (ar *AuthResponse) FullName() string {
	if len(ar.FirstName) > 0 && len(ar.LastName) > 0 {
		return ar.FirstName + " " + ar.LastName
	}

	if len(ar.FirstName) > 0 {
		return ar.FirstName
	}

	if len(ar.LastName) > 0 {
		return ar.LastName
	}

	return ""
}
