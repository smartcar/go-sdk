package smartcar

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/smartcar/go-sdk/helpers/constants"
	"github.com/smartcar/go-sdk/helpers/requests"
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
	Type          string `json:"token_type"`
	ExpiresIn     int    `json:"expires_in"`
	Access        string `json:"access_token"`
	AccessExpiry  time.Time
	Refresh       string `json:"refresh_token"`
	RefreshExpiry time.Time
}

// AuthConnect contains all the fields than can be used to build auth URL.
type AuthConnect struct {
	Auth          AuthClient
	ForceApproval bool
	State         string
	VehicleInfo
	SingleSelect
}

// Error contains error type and message from Smartcar.
type Error struct {
	ErrorType string `json:"error"`
	Message   string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("error: %s, message: %s", e.ErrorType, e.Message)
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

	res, resErr := requests.POST(constants.ExchangeURL, encodedAuth, strings.NewReader(data.Encode()))
	if resErr != nil {
		resErr = errors.New("Auth ClientID missing")
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
	tokens.RefreshExpiry = time.Now().AddDate(0, 0, 60)

	return tokens, nil
}

// RefreshToken renews access token
func RefreshToken(auth AuthClient, refreshToken string) (Tokens, error) {
	authString := auth.ClientID + ":" + auth.ClientSecret
	encodedAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(authString))

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	res, resErr := requests.POST(constants.ExchangeURL, encodedAuth, strings.NewReader(data.Encode()))
	if resErr != nil {
		resErr = errors.New("Auth ClientID missing")
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

// TokenIsExpired checks if the token has expired.
func TokenIsExpired(expiration time.Time) bool {
	return time.Now().After(expiration)
}

// VehicleIsCompatible checks compatibility for a vin with provided scopes.
func VehicleIsCompatible(vin string, auth AuthClient) (bool, error) {
	type CompatibleResponse struct {
		Compatible bool `json:"compatible"`
	}

	compatiblityURL := url.URL{
		Scheme: constants.APIScheme,
		Host:   constants.APIHost,
		Path:   "v1.0/compatibility",
	}

	query := compatiblityURL.Query()
	query.Set("vin", vin)
	query.Set("scope", strings.Join(auth.Scope, " "))
	compatiblityURL.RawQuery = query.Encode()

	fmt.Println(compatiblityURL.String())
	authString := auth.ClientID + ":" + auth.ClientSecret
	encodedAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(authString))

	res, resErr := requests.GET(compatiblityURL.String(), encodedAuth)
	if resErr != nil {
		resErr = errors.New("Yoo")
		return false, resErr
	}
	defer res.Body.Close()
	jsonDecoder := json.NewDecoder(res.Body)

	if res.StatusCode != 200 {
		var err Error
		jsonErr := jsonDecoder.Decode(&err)
		if jsonErr != nil {
			return false, jsonErr
		}

		return false, &Error{err.ErrorType, err.Message}
	}

	var compatibleResponse CompatibleResponse
	jsonErr := jsonDecoder.Decode(&compatibleResponse)
	if jsonErr != nil {
		jsonErr = errors.New("Decoding JSON error")
		return false, jsonErr
	}

	return compatibleResponse.Compatible, nil
}
