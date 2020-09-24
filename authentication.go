package smartcar

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	connectURL = "https://connect.smartcar.com/oauth/authorize"
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
	Enabled bool
	VIN     string
}

// Token is returned by auth.ExchangeCode and auth.ExchangeRefreshToken.
type Token struct {
	Access        string    `json:"access_token"`
	AccessExpiry  time.Time `json:"access_expiry"`
	Refresh       string    `json:"refresh_token"`
	RefreshExpiry time.Time `json:"refresh_expiry"`
	ExpiresIn     int       `json:"expires_in"`
}

// AuthURLParams contains the AuthClient, Pro authorization features and all fields that can be used to construct an auth URL.
type AuthURLParams struct {
	ForceApproval bool
	State         string
	Flags         string
	MakeBypass
	SingleSelect
}

// ExchangeCodeParams struct
type ExchangeCodeParams struct {
	Code string
}

// ExchangeRefreshTokenParams struct
type ExchangeRefreshTokenParams struct {
	Token string
}

// Auth interface is a...
type Auth interface {
	GetAuthURL(*AuthURLParams) (string, error)
	ExchangeCode(context.Context, *ExchangeCodeParams) (*Token, error)
	ExchangeRefreshToken(context.Context, *ExchangeRefreshTokenParams) (*Token, error)
}

// auth is used to store your auth credentials when authenticating with Smartcar.
type auth struct {
	clientID     string
	clientSecret string
	redirectURI  string
	scope        []string
	testMode     bool
	sC           backendClient
}

// GetAuthURL generates Smartcar Connect URL.
func (c *auth) GetAuthURL(params *AuthURLParams) (string, error) {
	forceApproval, state, singleSelect, vehicleInfo, flags := params.ForceApproval, params.State, params.SingleSelect, params.MakeBypass, params.Flags

	if c.clientID == "" {
		return "", errors.New("AuthClient.ClientID missing")
	}

	if c.redirectURI == "" {
		return "", errors.New("AuthClient.RedirectURI missing")
	}

	// Build Connect URL from go
	baseURL, _ := url.Parse(connectURL)
	query := baseURL.Query()
	query.Set("response_type", "code")
	query.Set("client_id", c.clientID)
	query.Set("redirect_uri", c.redirectURI)
	if len(flags) > 0 {
		query.Set("flags", flags)
	}

	if c.scope != nil {
		query.Set("scope", strings.Join(c.scope, " "))
	}

	if c.testMode {
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
			query.Set("make", string(vehicleInfo.Make))
		}
	}

	if singleSelect != (SingleSelect{}) {
		singleSelectEnabled := singleSelect.Enabled
		if singleSelect.VIN != "" {
			query.Set("single_select_vin", singleSelect.VIN)
			singleSelectEnabled = true
		}
		query.Set("single_select", strconv.FormatBool(singleSelectEnabled))
	}

	baseURL.RawQuery = query.Encode()

	return baseURL.String(), nil
}

// ExchangeCode exchanges authorization code for a Token.
func (c *auth) ExchangeCode(ctx context.Context, params *ExchangeCodeParams) (*Token, error) {

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", params.Code)
	data.Set("redirect_uri", c.redirectURI)

	token := &Token{}
	if err := c.request(ctx, http.MethodPost, exchangeURL, data.Encode(), token); err != nil {
		return nil, err
	}

	token.AccessExpiry = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
	token.RefreshExpiry = time.Now().AddDate(0, 0, 60)

	return token, nil
}

// ExchangeRefreshToken exchanges refresh token for a new Token.
func (c *auth) ExchangeRefreshToken(ctx context.Context, params *ExchangeRefreshTokenParams) (*Token, error) {

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", params.Token)

	token := &Token{}
	if err := c.request(ctx, http.MethodPost, exchangeURL, data.Encode(), token); err != nil {
		return nil, err
	}

	timeNow := time.Now()
	token.AccessExpiry = timeNow.Add(time.Duration(token.ExpiresIn) * time.Second)
	token.RefreshExpiry = timeNow.AddDate(0, 0, 60)

	return token, nil
}

// request is an internal function for sending requests, accepts an interface to
func (c *auth) request(ctx context.Context, method string, url string, data string, target interface{}) error {
	authorization := buildBasicAuthorization(c.clientID, c.clientSecret)

	return c.sC.Call(backendClientParams{
		ctx:           ctx,
		method:        method,
		url:           url,
		authorization: authorization,
		body:          strings.NewReader(data),
		target:        target,
	})
}
