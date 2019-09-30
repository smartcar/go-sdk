package smartcar

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/smartcar/go-sdk/helpers/constants"
	"github.com/smartcar/go-sdk/helpers/requests"
	"github.com/smartcar/go-sdk/helpers/test"
	utils "github.com/smartcar/go-sdk/helpers/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthenticationTestSuite struct {
	suite.Suite
	clientID     string
	clientSecret string
	redirectURI  string
	scope        []string
	testMode     bool
	state        string
	make         string
	vin          string
}

func (suite *AuthenticationTestSuite) SetupSuite() {
	suite.clientID = "clientId"
	suite.clientSecret = "clientSecret"
	suite.redirectURI = "redirectUri"
	suite.scope = []string{"scope"}
	suite.testMode = true
	suite.state = "state"
	suite.make = "tesla"
	suite.vin = "1234567890ABCDEFG"
}

func (suite *AuthenticationTestSuite) AfterTest() {
	test.ClearMock()
}

func (suite *AuthenticationTestSuite) TestGetAuthURLMissingClientID() {
	// Arrange
	authClient := AuthClient{}

	// Act
	_, err := authClient.GetAuthURL(AuthURLOptions{})

	// Assert
	assert.EqualError(suite.T(), err, "AuthClient.ClientID missing")
}

func (suite *AuthenticationTestSuite) TestGetAuthURLMissingRedirectURI() {
	// Arrange
	authClient := AuthClient{ClientID: suite.clientID}

	// Act
	_, err := authClient.GetAuthURL(AuthURLOptions{})

	// Assert
	assert.EqualError(suite.T(), err, "AuthClient.RedirectURI missing")
}

func (suite *AuthenticationTestSuite) TestGetAuthURLEmtpyConfig() {
	// Arrange
	authClient := AuthClient{ClientID: suite.clientID, RedirectURI: suite.redirectURI}
	expectedURL := fmt.Sprintf("https://connect.smartcar.com/oauth/authorize?approval_prompt=auto&client_id=%s&redirect_uri=%s&response_type=code",
		suite.clientID,
		suite.redirectURI,
	)

	// Act
	url, err := authClient.GetAuthURL(AuthURLOptions{})

	// Assert
	if err != nil {
		suite.T().Error("Should not have thrown")
	}
	assert.Equal(suite.T(), expectedURL, url)
}

func (suite *AuthenticationTestSuite) TestGetAuthURLApprovalPrompt() {
	// Arrange
	authClient := AuthClient{ClientID: suite.clientID, RedirectURI: suite.redirectURI}
	authURLOptions := AuthURLOptions{ForceApproval: true}
	expectedURL := fmt.Sprintf("https://connect.smartcar.com/oauth/authorize?approval_prompt=force&client_id=%s&redirect_uri=%s&response_type=code",
		suite.clientID,
		suite.redirectURI,
	)

	// Act
	url, err := authClient.GetAuthURL(authURLOptions)

	// Assert
	if err != nil {
		suite.T().Error("Should not have thrown")
	}
	assert.Equal(suite.T(), expectedURL, url)
}

func (suite *AuthenticationTestSuite) TestGetAuthURLScope() {
	// Arrange
	authClient := AuthClient{ClientID: suite.clientID, RedirectURI: suite.redirectURI, Scope: suite.scope}
	authURLOptions := AuthURLOptions{}
	expectedURL := fmt.Sprintf("https://connect.smartcar.com/oauth/authorize?approval_prompt=auto&client_id=%s&redirect_uri=%s&response_type=code&scope=%s",
		suite.clientID,
		suite.redirectURI,
		strings.Join(suite.scope, " "),
	)

	// Act
	url, err := authClient.GetAuthURL(authURLOptions)

	// Assert
	if err != nil {
		suite.T().Error("Should not have thrown")
	}
	assert.Equal(suite.T(), expectedURL, url)
}

func (suite *AuthenticationTestSuite) TestGetAuthURLTestMode() {
	// Arrange
	authClient := AuthClient{ClientID: suite.clientID, RedirectURI: suite.redirectURI, TestMode: suite.testMode}
	authURLOptions := AuthURLOptions{}
	expectedURL := fmt.Sprintf("https://connect.smartcar.com/oauth/authorize?approval_prompt=auto&client_id=%s&mode=test&redirect_uri=%s&response_type=code",
		suite.clientID,
		suite.redirectURI,
	)

	// Act
	url, err := authClient.GetAuthURL(authURLOptions)

	// Assert
	if err != nil {
		suite.T().Error("Should not have thrown")
	}
	assert.Equal(suite.T(), expectedURL, url)
}

func (suite *AuthenticationTestSuite) TestGetAuthURLState() {
	// Arrange
	authClient := AuthClient{ClientID: suite.clientID, RedirectURI: suite.redirectURI}
	authURLOptions := AuthURLOptions{State: suite.state}
	expectedURL := fmt.Sprintf("https://connect.smartcar.com/oauth/authorize?approval_prompt=auto&client_id=%s&redirect_uri=%s&response_type=code&state=%s",
		suite.clientID,
		suite.redirectURI,
		suite.state,
	)

	// Act
	url, err := authClient.GetAuthURL(authURLOptions)

	// Assert
	if err != nil {
		suite.T().Error("Should not have thrown")
	}
	assert.Equal(suite.T(), expectedURL, url)
}

func (suite *AuthenticationTestSuite) TestGetAuthURLVehicleInfo() {
	// Arrange
	authClient := AuthClient{ClientID: suite.clientID, RedirectURI: suite.redirectURI}
	authURLOptions := AuthURLOptions{State: suite.state}
	expectedURL := fmt.Sprintf("https://connect.smartcar.com/oauth/authorize?approval_prompt=auto&client_id=%s&redirect_uri=%s&response_type=code&state=%s",
		suite.clientID,
		suite.redirectURI,
		suite.state,
	)

	// Act
	url, err := authClient.GetAuthURL(authURLOptions)

	// Assert
	if err != nil {
		suite.T().Error("Should not have thrown")
	}
	assert.Equal(suite.T(), expectedURL, url)
}

//PRO FEATURES
func (suite *AuthenticationTestSuite) TestGetAuthURLMakeBypass() {
	// Arrange
	makeBypass := MakeBypass{Make: suite.make}
	authClient := AuthClient{ClientID: suite.clientID, RedirectURI: suite.redirectURI}
	authURLOptions := AuthURLOptions{MakeBypass: makeBypass}
	expectedURL := fmt.Sprintf(
		"https://connect.smartcar.com/oauth/authorize?approval_prompt=auto&client_id=%s&make=%s&redirect_uri=%s&response_type=code",
		suite.clientID,
		suite.make,
		suite.redirectURI,
	)

	// Act
	url, err := authClient.GetAuthURL(authURLOptions)

	// Assert
	if err != nil {
		suite.T().Error("Should not have thrown")
	}
	assert.Equal(suite.T(), expectedURL, url)
}

func (suite *AuthenticationTestSuite) TestGetAuthURLSingleSelectVin() {
	// Arrange
	singleSelect := SingleSelect{Vin: suite.vin}
	authClient := AuthClient{ClientID: suite.clientID, RedirectURI: suite.redirectURI}
	authURLOptions := AuthURLOptions{SingleSelect: singleSelect}
	expectedURL := fmt.Sprintf(
		"https://connect.smartcar.com/oauth/authorize?approval_prompt=auto&client_id=%s&redirect_uri=%s&response_type=code&single_select=true&single_select_vin=%s",
		suite.clientID,
		suite.redirectURI,
		suite.vin,
	)

	// Act
	url, err := authClient.GetAuthURL(authURLOptions)

	// Assert
	if err != nil {
		suite.T().Error("Should not have thrown")
	}
	assert.Equal(suite.T(), expectedURL, url)
}

// TODO: Test with only single_select=true
func (suite *AuthenticationTestSuite) TestGetAuthURLSingleSelectEmpty() {
	clientID := "clientId"
	redirectURI := "redirectUri"
	authClient := AuthClient{ClientID: clientID, RedirectURI: redirectURI}
	singleSelect := SingleSelect{Enabled: true}
	authURLOptions := AuthURLOptions{SingleSelect: singleSelect}
	expectedURL := fmt.Sprintf(
		"https://connect.smartcar.com/oauth/authorize?approval_prompt=auto&client_id=%s&redirect_uri=%s&response_type=code&single_select=true",
		clientID,
		redirectURI,
	)
	url, err := authClient.GetAuthURL(authURLOptions)
	if err != nil {
		suite.T().Error("Should not have thrown")
	}
	assert.Equal(suite.T(), expectedURL, url)
}

func (suite *AuthenticationTestSuite) TestExchangeCode() {
	// Arrange
	authClient := AuthClient{ClientID: suite.clientID, ClientSecret: suite.clientSecret}
	authCode := "authCode"
	authorization := utils.BuildBasicAuthorization(suite.clientID, suite.clientSecret)
	expectedResponse := map[string]interface{}{
		"access_token":  "expectedAccess",
		"token_type":    "Bearer",
		"refresh_token": "refreshToken",
		"expires_in":    7200,
	}
	// TODO
	// Match the req.Body to make sure it contains the code, redirectURI, and grant_type
	test.MockRequest(requests.POST, constants.ExchangeURL, authorization, 200, expectedResponse)

	// Act
	token, err := authClient.ExchangeCode(authCode)

	// Assert
	if err != nil {
		suite.T().Error(err)
	}
	assert.Equal(suite.T(), expectedResponse["access_token"], token.Access)
	assert.WithinDuration(suite.T(), time.Now().Add(2*time.Hour), token.AccessExpiry, 10*time.Second)
	assert.Equal(suite.T(), expectedResponse["refresh_token"], token.Refresh)
	assert.WithinDuration(suite.T(), time.Now().AddDate(0, 0, 60), token.RefreshExpiry, 10*time.Second)
	assert.Equal(suite.T(), expectedResponse["expires_in"], token.ExpiresIn)
}

func (suite *AuthenticationTestSuite) TestRefreshToken() {
	// Arrange
	authClient := AuthClient{ClientID: suite.clientID, ClientSecret: suite.clientSecret}
	refreshToken := "refreshToken"
	authorization := utils.BuildBasicAuthorization(suite.clientID, suite.clientSecret)
	expectedResponse := map[string]interface{}{
		"access_token":  "expectedAccess",
		"token_type":    "Bearer",
		"refresh_token": "refreshToken",
		"expires_in":    7200,
	}
	// TODO
	// Match the req.Body to make sure it contains the code, redirectURI, and grant_type
	test.MockRequest(requests.POST, constants.ExchangeURL, authorization, 200, expectedResponse)

	// Act
	token, err := authClient.RefreshToken(refreshToken)

	// Assert
	if err != nil {
		suite.T().Error(err)
	}
	assert.Equal(suite.T(), expectedResponse["access_token"], token.Access)
	assert.WithinDuration(suite.T(), time.Now().Add(2*time.Hour), token.AccessExpiry, 10*time.Second)
	assert.Equal(suite.T(), expectedResponse["refresh_token"], token.Refresh)
	assert.WithinDuration(suite.T(), time.Now().AddDate(0, 0, 60), token.RefreshExpiry, 10*time.Second)
	assert.Equal(suite.T(), expectedResponse["expires_in"], token.ExpiresIn)
}

func (suite *AuthenticationTestSuite) TestTokenIsExpiredFalse() {
	// Act
	result := TokenIsExpired(time.Now())

	// Assert
	assert.False(suite.T(), result)
}

func (suite *AuthenticationTestSuite) TestTokenIsExpiredTrue() {
	// Act
	result := TokenIsExpired(time.Now().Add(time.Second * -10))

	// Assert
	assert.True(suite.T(), result)
}

func TestAuthenticationTestSuite(t *testing.T) {
	suite.Run(t, new(AuthenticationTestSuite))
}
