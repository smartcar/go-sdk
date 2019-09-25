package requests

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestRequest(t *testing.T) {
	// Arrange
	expectedURL := "http://example.com"
	expectedMethod := "GET"
	accessToken := "access-token"
	expectedBody := "Body"
	gock.New(expectedURL).
		MatchHeader("Authorization", "Bearer "+accessToken).
		Get("/").
		Reply(200).
		BodyString(expectedBody)

	// Act
	res, err := Request(expectedMethod, expectedURL, accessToken, nil)
	if err != nil {
		t.Error("Expected", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error("Should not be called")
	}

	// Assert
	assert.Equal(t, string(body), expectedBody)
}

func TestRequestFail(t *testing.T) {
	// Arrange
	expectedURL := ""
	// expectedMethod := "GET"
	accessToken := "access-token"
	gock.New(expectedURL).
		MatchHeader("Authorization", "Bearer "+accessToken)
		// Get("/").
		// Reply(401)

	// Act
	_, err := Request("", expectedURL, accessToken, nil)
	if err != nil {
		t.Error("Expected", err)
	}
	// defer res.Body.Close()
	// body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		assert.Error(t, err, "err")
	}

	// Assert
	// assert.Equal(t, string(body), expectedBody)
}
