package smartcar

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAuthURLMissingClientID(t *testing.T) {
	authClient := AuthClient{}

	_, err := authClient.GetAuthURL(AuthURLOptions{})

	assert.EqualError(t, err, "AuthClient.ClientID missing")
}

func TestGetAuthURLMissingRedirectURI(t *testing.T) {
	authClient := AuthClient{ClientID: "clientId"}

	_, err := authClient.GetAuthURL(AuthURLOptions{})

	assert.EqualError(t, err, "AuthClient.RedirectURI missing")
}

func TestGetAuthURLEmtpyConfig(t *testing.T) {
	clientID := "clientId"
	redirectURI := "redirectUri"
	authClient := AuthClient{ClientID: clientID, RedirectURI: redirectURI}
	expectedURL := fmt.Sprintf("https://connect.smartcar.com/oauth/authorize?approval_prompt=auto&client_id=%s&redirect_uri=%s&response_type=code", clientID, redirectURI)

	url, err := authClient.GetAuthURL(AuthURLOptions{})

	if err != nil {
		t.Error("Should not have thrown")
	}
	assert.Equal(t, url, expectedURL)
}

func TestGetAuthURLApprovalPrompt(t *testing.T) {
	clientID := "clientId"
	redirectURI := "redirectUri"
	authClient := AuthClient{ClientID: clientID, RedirectURI: redirectURI}
	authURLOptions := AuthURLOptions{ForceApproval: true}
	expectedURL := fmt.Sprintf("https://connect.smartcar.com/oauth/authorize?approval_prompt=force&client_id=%s&redirect_uri=%s&response_type=code", clientID, redirectURI)

	url, err := authClient.GetAuthURL(authURLOptions)

	if err != nil {
		t.Error("Should not have thrown")
	}
	assert.Equal(t, url, expectedURL)
}

func TestGetAuthURLScope(t *testing.T) {
	clientID := "clientId"
	redirectURI := "redirectUri"
	scope := []string{"scope"}
	authClient := AuthClient{ClientID: clientID, RedirectURI: redirectURI, Scope: scope}
	authURLOptions := AuthURLOptions{}
	expectedURL := fmt.Sprintf("https://connect.smartcar.com/oauth/authorize?approval_prompt=auto&client_id=%s&redirect_uri=%s&response_type=code&scope=%s", clientID, redirectURI, strings.Join(scope, " "))

	url, err := authClient.GetAuthURL(authURLOptions)

	if err != nil {
		t.Error("Should not have thrown")
	}
	assert.Equal(t, url, expectedURL)
}

func TestGetAuthURLTestMode(t *testing.T) {
	clientID := "clientId"
	redirectURI := "redirectUri"
	testMode := true
	authClient := AuthClient{ClientID: clientID, RedirectURI: redirectURI, TestMode: testMode}
	authURLOptions := AuthURLOptions{}
	expectedURL := fmt.Sprintf("https://connect.smartcar.com/oauth/authorize?approval_prompt=auto&client_id=%s&mode=test&redirect_uri=%s&response_type=code", clientID, redirectURI)

	url, err := authClient.GetAuthURL(authURLOptions)

	if err != nil {
		t.Error("Should not have thrown")
	}
	assert.Equal(t, url, expectedURL)
}

func TestGetAuthURLState(t *testing.T) {
	clientID := "clientId"
	redirectURI := "redirectUri"
	state := "state"
	authClient := AuthClient{ClientID: clientID, RedirectURI: redirectURI}
	authURLOptions := AuthURLOptions{State: state}
	expectedURL := fmt.Sprintf("https://connect.smartcar.com/oauth/authorize?approval_prompt=auto&client_id=%s&redirect_uri=%s&response_type=code&state=%s", clientID, redirectURI, state)

	url, err := authClient.GetAuthURL(authURLOptions)

	if err != nil {
		t.Error("Should not have thrown")
	}
	assert.Equal(t, url, expectedURL)
}

func TestGetAuthURLVehicleInfo(t *testing.T) {
	clientID := "clientId"
	redirectURI := "redirectUri"
	state := "state"
	authClient := AuthClient{ClientID: clientID, RedirectURI: redirectURI}
	authURLOptions := AuthURLOptions{State: state}
	expectedURL := fmt.Sprintf("https://connect.smartcar.com/oauth/authorize?approval_prompt=auto&client_id=%s&redirect_uri=%s&response_type=code&state=%s", clientID, redirectURI, state)

	url, err := authClient.GetAuthURL(authURLOptions)

	if err != nil {
		t.Error("Should not have thrown")
	}
	assert.Equal(t, url, expectedURL)
}

//PRO FEATURES
func TestGetAuthURLMakeBypass(t *testing.T) {
	clientID := "clientId"
	redirectURI := "redirectUri"
	makeBypass := MakeBypass{Make: "make"}
	authClient := AuthClient{ClientID: clientID, RedirectURI: redirectURI}
	authURLOptions := AuthURLOptions{MakeBypass: makeBypass}
	expectedURL := fmt.Sprintf(
		"https://connect.smartcar.com/oauth/authorize?approval_prompt=auto&client_id=%s&redirect_uri=%s&response_type=code",
		clientID,
		redirectURI,
	)

	url, err := authClient.GetAuthURL(authURLOptions)

	if err != nil {
		t.Error("Should not have thrown")
	}
	assert.Equal(t, url, expectedURL)
}
