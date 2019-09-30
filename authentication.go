package smartcar

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/smartcar/go-sdk/helpers/constants"
	"github.com/smartcar/go-sdk/helpers/requests"
	"github.com/smartcar/go-sdk/helpers/utils"
)

/*

The following 2 are Pro features:
	- MakeByPass
	- Single Select
*/

// MakeBypass uses a make to bypass the Smartcar Connect brand selector.
// Smartcar Pro feature.
type MakeBypass struct {
	Make string
}

// SingleSelect will only authorize vehicles that match the given properties.
// Smartcar Pro feature.
type SingleSelect struct {
	Single bool
	Vin    string
}

// AuthClient is used to store your auth credentials when authenticating with Smartcar.
type AuthClient struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Scope        []string
	TestMode     bool
}

// Token contains the tokens and their expiry that are returned from exchanging an authorization code.
type Token struct {
	ExpiresIn     int    `json:"expires_in"`
	Access        string `json:"access_token"`
	AccessExpiry  time.Time
	Refresh       string `json:"refresh_token"`
	RefreshExpiry time.Time
}

// AuthURLOptions contains the AuthClient, Pro authorization features and all fields that can be used to construct an auth URL.
type AuthURLOptions struct {
	ForceApproval bool
	State         string
	TestMode      bool
	MakeBypass
	SingleSelect
}

// GetAuthURL uses an AuthConnect to return a Smartcar Connect URL that can be displayed to users.
func (authClient *AuthClient) GetAuthURL(options AuthURLOptions) (string, error) {
	forceApproval, state, singleSelect, vehicleInfo := options.ForceApproval, options.State, options.SingleSelect, options.MakeBypass

	if authClient.ClientID == "" {
		return "", errors.New("AuthClient.ClientID missing")
	}

	if authClient.RedirectURI == "" {
		return "", errors.New("AuthClient.RedirectURI missing")
	}

	// Build Connect URL from constants.go
	connectURL := utils.GetConnectURL()
	query := connectURL.Query()
	query.Set("response_type", "code")
	query.Set("client_id", authClient.ClientID)
	query.Set("redirect_uri", authClient.RedirectURI)

	if authClient.Scope != nil {
		query.Set("scope", strings.Join(authClient.Scope, " "))
	}

	if authClient.TestMode {
		query.Set("mode", "test")
	}

	approvalPrompt := "auto"
	if forceApproval {
		approvalPrompt = "force"
	}
	query.Set("approval_prompt", approvalPrompt)

	if state != "" {
		query.Set("state", state)
	}

	if vehicleInfo != (MakeBypass{}) {
		if vehicleInfo.Make != "" {
			query.Set("make", vehicleInfo.Make)
		}
	}

	if singleSelect != (SingleSelect{}) {
		if singleSelect.Vin != "" {
			query.Set("single_select_vin", singleSelect.Vin)
		}
		query.Set("single_select", "true")
	}

	connectURL.RawQuery = query.Encode()

	return connectURL.String(), nil
}

// ExchangeCode takes an authorization code and exchanges it for an access and refresh token.
func (authClient *AuthClient) ExchangeCode(authCode string) (Token, error) {

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", authCode)
	data.Set("redirect_uri", authClient.RedirectURI)

	res, err := authClient.request(requests.POST, constants.ExchangeURL, data.Encode())
	if err != nil {
		return Token{}, err
	}

	formattedResponse := new(Token)
	defer res.Body.Close()
	requests.FormatResponse(res.Body, formattedResponse)

	formattedResponse.AccessExpiry = time.Now().Add(time.Duration(formattedResponse.ExpiresIn) * time.Second)
	formattedResponse.RefreshExpiry = time.Now().AddDate(0, 0, 60)

	return *formattedResponse, nil
}

// RefreshToken uses a basic AuthClient containing your client ID and a refresh token to return new access tokens.
func (authClient *AuthClient) RefreshToken(refreshToken string) (Token, error) {

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	res, resErr := authClient.request(requests.POST, constants.ExchangeURL, data.Encode())
	if resErr != nil {
		return Token{}, resErr
	}

	formattedResponse := new(Token)
	defer res.Body.Close()
	requests.FormatResponse(res.Body, formattedResponse)

	formattedResponse.AccessExpiry = time.Now().Add(time.Duration(formattedResponse.ExpiresIn) * time.Second)
	formattedResponse.RefreshExpiry = time.Now().AddDate(0, 0, 60)

	return *formattedResponse, nil
}

// IsCompatible checks compatibility for a vehicle VIN with Smartcar for the provided scopes.
func (authClient *AuthClient) IsCompatible(vin string) (bool, error) {
	url := utils.BuildCompatibilityURL(vin, authClient.Scope)

	res, resErr := authClient.request(requests.GET, url, "")
	if resErr != nil {
		return false, resErr
	}

	formattedResponse := new(struct {
		Compatible bool `json:"compatible"`
	})
	defer res.Body.Close()
	fmtErr := requests.FormatResponse(res.Body, formattedResponse)
	if fmtErr != nil {
		return false, fmtErr
	}

	return formattedResponse.Compatible, nil
}

func (authClient *AuthClient) request(method string, path string, data string) (http.Response, error) {
	authorization := utils.BuildBasicAuthorization(authClient.ClientID, authClient.ClientSecret)

	return requests.Request(method, path, authorization, "", strings.NewReader(data))
}
