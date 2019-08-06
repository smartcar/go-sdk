package smartcar

import (
	"errors"
	"net/url"
	"strings"
)

// AuthClient for interacting with Connect and API.
type AuthClient struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Scope        []string
	TestMode     bool
}

// VehicleInfo for Connect Direct
type VehicleInfo struct {
	Make string
}

// SingleSelect for Connect Match
type SingleSelect struct {
	Vin string
}

// GetAuthURL builds an Auth URL for front-end
func GetAuthURL(auth AuthClient, force bool, state string, vehicleInfo VehicleInfo, singleSelect SingleSelect) (string, error) {
	var err error

	if auth.ClientID == "" {
		err = errors.New("Auth ClientID missing")
	}

	if auth.RedirectURI == "" {
		err = errors.New("Auth RedirectURI missing")
	}

	approvalPrompt := "auto"
	if force {
		approvalPrompt = "force"
	}

	// Build Connect URL from constants.go
	connectURL := url.URL{
		Scheme: ConnectScheme,
		Host:   ConnectHost,
		Path:   ConnectPath,
	}

	query := connectURL.Query()
	query.Set("client_id", auth.ClientID)
	query.Set("response_type", "code")
	query.Set("scope", strings.Join(auth.Scope, " "))
	query.Set("redirect_uri", auth.RedirectURI)
	query.Set("approval_prompt", approvalPrompt)

	if auth.TestMode {
		query.Set("mode", "test")
	}

	if state != "" {
		query.Set("state", state)
	}

	if vehicleInfo != (VehicleInfo{}) {
		if vehicleInfo.Make != "" {
			query.Set("make", vehicleInfo.Make)
		}
	}

	if singleSelect != (SingleSelect{}) {
		if singleSelect.Vin != "" {
			query.Set("vin", singleSelect.Vin)
		}
	}

	connectURL.RawQuery = query.Encode()

	return connectURL.String(), err
}
