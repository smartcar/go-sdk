package smartcar

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthenticationTestSuite struct {
	suite.Suite
	auth auth
}

type fakeBackend struct{}

func (c *fakeBackend) Call(params backendClientParams) error {
	token := &Token{
		ExpiresIn: 7200,
	}
	mapstructure.Decode(token, params.target)
	return nil
}

func newFakeBackend() backendClient {
	return &fakeBackend{}
}

func (s *AuthenticationTestSuite) SetupTest() {
	s.auth = auth{
		clientID:     "client-id",
		clientSecret: "client-secret",
		redirectURI:  "https://example.com",
		scope:        []string{"read_odometer", "read_vin"},
		testMode:     false,
		sC:           newFakeBackend(),
	}
}

func (s *AuthenticationTestSuite) TestGetAuthURLEmptyClientID() {
	params := AuthURLParams{}
	auth := auth{
		clientID: "",
	}

	_, err := auth.GetAuthURL(&params)

	assert.EqualError(s.T(), err, "AuthClient.ClientID missing")
}

func (s *AuthenticationTestSuite) TestGetAuthURLEmptyRedirectURI() {
	params := AuthURLParams{}
	auth := auth{
		clientID:     "client-id",
		clientSecret: "client-secret",
		redirectURI:  "",
	}

	_, err := auth.GetAuthURL(&params)

	assert.EqualError(s.T(), err, "AuthClient.RedirectURI missing")
}

func (s *AuthenticationTestSuite) TestGetAuthURLEmptyParams() {
	expectedScope := strings.Join(s.auth.scope[:], "+")
	expectedURL := url.QueryEscape(s.auth.redirectURI)
	params := AuthURLParams{}
	expectedAuthURL := fmt.Sprintf(
		"https://connect.smartcar.com/oauth/authorize?approval_prompt=auto&client_id=%s&redirect_uri=%s&response_type=code&scope=%s",
		s.auth.clientID,
		expectedURL,
		expectedScope,
	)

	authURL, _ := s.auth.GetAuthURL(&params)

	assert.Equal(s.T(), expectedAuthURL, authURL)
}

func (s *AuthenticationTestSuite) TestGetAuthURLAllParams() {
	make := "TESLA"
	state := "state"
	VIN := "123456789901234567"
        Flags := []string{"country:DE"}
	expectedScope := strings.Join(s.auth.scope[:], "+")
	expectedURL := url.QueryEscape(s.auth.redirectURI)
	expectedFlags := url.QueryEscape(Flags[0])
	params := AuthURLParams{
		ForceApproval: true,
		State:         state,
		MakeBypass:    MakeBypass{Make: make},
		SingleSelect:  SingleSelect{VIN: VIN},
		Flags:         Flags,
	}
	expectedAuthURL := fmt.Sprintf(
		"https://connect.smartcar.com/oauth/authorize?approval_prompt=force&client_id=%s&flags=%s&make=%s&redirect_uri=%s&response_type=code&scope=%s&single_select=true&single_select_vin=%s&state=%s",
		s.auth.clientID,
		expectedFlags,
		make,
		expectedURL,
		expectedScope,
		VIN,
		state,
	)

	authURL, _ := s.auth.GetAuthURL(&params)

	assert.Equal(s.T(), expectedAuthURL, authURL)
}

func (s *AuthenticationTestSuite) TestExchangeCode() {
	token, err := s.auth.ExchangeCode(context.TODO(), &ExchangeCodeParams{})

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), token)
	assert.WithinDuration(s.T(), time.Now().Add(2*time.Hour), token.AccessExpiry, 10*time.Second)
	assert.WithinDuration(s.T(), time.Now().AddDate(0, 0, 60), token.RefreshExpiry, 10*time.Second)
}

func (s *AuthenticationTestSuite) TestExchangeRefreshToken() {
	token, err := s.auth.ExchangeRefreshToken(context.TODO(), &ExchangeRefreshTokenParams{})

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), token)
	assert.WithinDuration(s.T(), time.Now().Add(2*time.Hour), token.AccessExpiry, 10*time.Second)
	assert.WithinDuration(s.T(), time.Now().AddDate(0, 0, 60), token.RefreshExpiry, 10*time.Second)
}

func TestAuthenticationTestSuite(t *testing.T) {
	suite.Run(t, new(AuthenticationTestSuite))
}
