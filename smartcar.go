package smartcar

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	helpers "github.com/smartcar/go-sdk/helpers"
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

// Tokens returned from exchange auth code.
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// AuthConnect contains all the fields than can be used to build auth URL.
type AuthConnect struct {
	Auth          AuthClient
	ForceApproval bool
	State         string
	VehicleInfo
	SingleSelect
}

// GetAuthURL builds an Auth URL for front-end
func GetAuthURL(authConnect AuthConnect) (string, error) {
	auth := authConnect.Auth
	vehicleInfo := authConnect.VehicleInfo
	singleSelect := authConnect.SingleSelect
	forceApproval, state := authConnect.ForceApproval, authConnect.State
	var err error

	if auth.ClientID == "" {
		err = errors.New("Auth ClientID missing")
		return "", err
	}

	if auth.RedirectURI == "" {
		err = errors.New("Auth RedirectURI missing")
		return "", err
	}

	approvalPrompt := "auto"
	if forceApproval {
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

	return connectURL.String(), nil
}

// ExchangeCode exchanges auth code for access and refresh tokens
func ExchangeCode(auth AuthClient, authCode string) (Tokens, error) {
	authString := auth.ClientID + ":" + auth.ClientSecret
	encodedAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(authString))

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", authCode)
	data.Set("redirect_uri", auth.RedirectURI)

	response, resErr := helpers.POSTRequest(ExchangeURL, encodedAuth, strings.NewReader(data.Encode()))
	if resErr != nil {
		resErr = errors.New("Auth ClientID missing")
		return Tokens{}, resErr
	}

	jsonDecoder := json.NewDecoder(response)
	var tokens Tokens
	jsonErr := jsonDecoder.Decode(&tokens)
	if jsonErr != nil {
		jsonErr = errors.New("Decoding JSON error")
		return Tokens{}, jsonErr
	}

	return tokens, nil
}
