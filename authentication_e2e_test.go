package smartcar

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/h2non/gock.v1"
)

type AuthE2ETestSuite struct {
	suite.Suite
	auth                   auth
	mockAge, mockRequestID string
	mockUnitSystem         UnitSystem
	responseHeaders        ResponseHeaders
}

func (s *AuthE2ETestSuite) SetupTest() {
	s.auth = auth{
		clientID:     "client-id",
		clientSecret: "client-secret",
		redirectURI:  "redirect-uri",
		scope:        []string{"read_odometer", "read_location"},
		testMode:     false,
		sC:           newBackend(),
	}
}

func (s *AuthE2ETestSuite) TearDownTestSuite() {
	gock.Off()
}

func mockAuthAPI(url, id, secret string, response interface{}) {
	gock.New(url).
		MatchHeader("Authorization", buildBasicAuthorization(id, secret)).
		MatchType("x-www-form-urlencoded").
		Reply(200).
		JSON(response)
}

func (s *AuthE2ETestSuite) TestExchangeCodeE2E() {
	mockAccess := "access"
	mockRefresh := "refresh"
	mockExpiresIn := 7200
	expectedResponse := &Token{
		Access:    mockAccess,
		ExpiresIn: mockExpiresIn,
		Refresh:   mockRefresh,
	}
	mockResponse := map[string]interface{}{
		"access_token":  mockAccess,
		"token_type":    "Bearer",
		"expires_in":    mockExpiresIn,
		"refresh_token": mockRefresh,
	}
	mockAuthAPI(exchangeURL, s.auth.clientID, s.auth.clientSecret, mockResponse)

	res, err := s.auth.ExchangeCode(context.TODO(), &ExchangeCodeParams{})

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedResponse.Access, res.Access)
	assert.Equal(s.T(), expectedResponse.Refresh, res.Refresh)
	assert.Equal(s.T(), expectedResponse.ExpiresIn, res.ExpiresIn)
	assert.NotEmpty(s.T(), res.AccessExpiry)
	assert.NotEmpty(s.T(), res.RefreshExpiry)
}

func (s *AuthE2ETestSuite) TestGetExchangeRefreshTokenE2E() {
	mockAccess := "access"
	mockRefresh := "refresh"
	mockExpiresIn := 7200
	expectedResponse := &Token{
		Access:    mockAccess,
		ExpiresIn: mockExpiresIn,
		Refresh:   mockRefresh,
	}
	mockResponse := map[string]interface{}{
		"access_token":  mockAccess,
		"token_type":    "Bearer",
		"expires_in":    mockExpiresIn,
		"refresh_token": mockRefresh,
	}
	mockAuthAPI(exchangeURL, s.auth.clientID, s.auth.clientSecret, mockResponse)

	res, err := s.auth.ExchangeRefreshToken(context.TODO(), &ExchangeRefreshTokenParams{})

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedResponse.Access, res.Access)
	assert.Equal(s.T(), expectedResponse.Refresh, res.Refresh)
	assert.Equal(s.T(), expectedResponse.ExpiresIn, res.ExpiresIn)
	assert.NotEmpty(s.T(), res.AccessExpiry)
	assert.NotEmpty(s.T(), res.RefreshExpiry)
}

func TestAuthE2ETestSuite(t *testing.T) {
	suite.Run(t, new(AuthE2ETestSuite))
}
