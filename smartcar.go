package smartcar

import (
	"errors"
	"net/url"
	"strings"
)

// AuthClient for interacting with Connect and API.
type AuthClient struct {
	clientID     string
	clientSecret string
	redirectURI  string
	scope        []string
	testMode     bool
}

// VehicleInfo for Connect Direct
type VehicleInfo struct {
	make string
}

// SingleSelect for Connect Match
type SingleSelect struct {
	vin string
}

// GetAuthURL builds an Auth URL for front-end
func GetAuthURL(auth AuthClient, force bool, state string, vehicleInfo VehicleInfo, singleSelect SingleSelect) (string, error) {
	var err error

	if auth.clientID == "" {
		err = errors.New("Auth ClientID missing")
	}

	if auth.redirectURI == "" {
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
	query.Set("client_id", auth.clientID)
	query.Set("response_type", "code")
	query.Set("scope", strings.Join(auth.scope, " "))
	query.Set("redirect_uri", auth.redirectURI)
	query.Set("approval_prompt", approvalPrompt)

	if auth.testMode {
		query.Set("mode", "test")
	}

	if state != "" {
		query.Set("state", state)
	}

	if vehicleInfo != (VehicleInfo{}) {
		if vehicleInfo.make != "" {
			query.Set("make", vehicleInfo.make)
		}
	}

	if singleSelect != (SingleSelect{}) {
		if singleSelect.vin != "" {
			query.Set("vin", singleSelect.vin)
		}
	}

	connectURL.RawQuery = query.Encode()

	return connectURL.String(), err
}
