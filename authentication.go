// Package smartcar is the official Go SDK of the Smartcar API
package smartcar

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/smartcar/go-sdk/helpers/constants"
	"github.com/smartcar/go-sdk/helpers/requests"
)

// ConnectDirect uses a make to bypass the Smartcar Connect brand selector.
// Smartcar Pro feature.
type ConnectDirect struct {
	Make string
}

// ConnectMatch will only authorize vehicles that match the given properties.
// Smartcar Pro feature.
type ConnectMatch struct {
	Vin string
}

// Tokens contains the tokens and their expiry that are returned from exchanging an authorization code.
type Tokens struct {
	Type          string `json:"token_type"`
	ExpiresIn     int    `json:"expires_in"`
	Access        string `json:"access_token"`
	AccessExpiry  time.Time
	Refresh       string `json:"refresh_token"`
	RefreshExpiry time.Time
}

// AuthConnect contains the AuthClient, Pro authorization features and all fields that can be used to construct an auth URL.
type AuthConnect struct {
	Auth          AuthClient
	ForceApproval bool
	State         string
	ConnectDirect
	ConnectMatch
}

// GetAuthURL uses an AuthConnect to return a Smartcar Connect URL that can be displayed to users.
func GetAuthURL(authConnect AuthConnect) (string, error) {
	auth := authConnect.Auth
	vehicleInfo := authConnect.ConnectDirect
	singleSelect := authConnect.ConnectMatch
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
		Scheme: constants.ConnectScheme,
		Host:   constants.ConnectHost,
		Path:   constants.ConnectPath,
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

	if vehicleInfo != (ConnectDirect{}) {
		if vehicleInfo.Make != "" {
			query.Set("make", vehicleInfo.Make)
		}
	}

	if singleSelect != (ConnectMatch{}) {
		if singleSelect.Vin != "" {
			query.Set("vin", singleSelect.Vin)
		}
	}

	connectURL.RawQuery = query.Encode()

	return connectURL.String(), nil
}

// ExchangeCode takes an AuthClient and an authorization code from Connect to return access and refresh tokens.
func ExchangeCode(auth AuthClient, authCode string) (Tokens, error) {
	authString := auth.ClientID + ":" + auth.ClientSecret
	encodedAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(authString))

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", authCode)
	data.Set("redirect_uri", auth.RedirectURI)

	res, resErr := requests.POST(constants.ExchangeURL, encodedAuth, strings.NewReader(data.Encode()))
	if resErr != nil {
		return Tokens{}, resErr
	}
	defer res.Body.Close()

	var tokens Tokens
	jsonDecoder := json.NewDecoder(res.Body)

	if res.StatusCode != 200 {
		var err Error
		jsonErr := jsonDecoder.Decode(&err)
		if jsonErr != nil {
			jsonErr = errors.New("Decoding JSON error")
			return Tokens{}, jsonErr
		}
		err.Type = requests.HandleStatusCode(res.StatusCode)
		return Tokens{}, &Error{err.Type, err.Name, err.Message, err.Code}
	}

	jsonErr := jsonDecoder.Decode(&tokens)
	if jsonErr != nil {
		jsonErr = errors.New("Decoding JSON error")
		return Tokens{}, jsonErr
	}

	tokens.AccessExpiry = time.Now().Add(time.Duration(tokens.ExpiresIn) * time.Second)
	tokens.RefreshExpiry = time.Now().AddDate(0, 0, 60)

	return tokens, nil
}

// RefreshToken uses a basic AuthClient containing your client ID and a refresh token to return new access tokens.
func RefreshToken(auth AuthClient, refreshToken string) (Tokens, error) {
	authString := auth.ClientID + ":" + auth.ClientSecret
	encodedAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(authString))

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	res, resErr := requests.POST(constants.ExchangeURL, encodedAuth, strings.NewReader(data.Encode()))
	if resErr != nil {
		return Tokens{}, resErr
	}
	defer res.Body.Close()

	var tokens Tokens
	jsonDecoder := json.NewDecoder(res.Body)
	jsonErr := jsonDecoder.Decode(&tokens)
	if jsonErr != nil {
		jsonErr = errors.New("Decoding JSON error")
		return Tokens{}, jsonErr
	}

	tokens.AccessExpiry = time.Now().Add(time.Duration(tokens.ExpiresIn) * time.Second)

	return tokens, nil
}

// TokenIsExpired checks if a token is expired by passing in the token expiry time.
func TokenIsExpired(expiration time.Time) bool {
	return time.Now().After(expiration)
}
