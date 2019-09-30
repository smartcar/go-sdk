package requests

import (
	"io/ioutil"
	"testing"

	"github.com/smartcar/go-sdk/helpers/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/h2non/gock.v1"
)

type RequestsSuiteTest struct {
	suite.Suite
}

func (suite *RequestsSuiteTest) AfterTest() {
	test.ClearMock()
}

func (suite *RequestsSuiteTest) TestRequest() {
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
	res, err := Request(expectedMethod, expectedURL, accessToken, "", nil)
	if err != nil {
		suite.T().Error("Expected", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		suite.T().Error("Should not be called")
	}

	// Assert
	assert.Equal(suite.T(), expectedBody, string(body))
}

func (suite *RequestsSuiteTest) TestRequestFail() {
	// Arrange
	expectedURL := ""
	// expectedMethod := "GET"
	accessToken := "access-token"
	gock.New(expectedURL).
		MatchHeader("Authorization", "Bearer "+accessToken)
		// Get("/").
		// Reply(401)

	// Act
	_, err := Request("", expectedURL, accessToken, "", nil)
	if err != nil {
		suite.T().Error("Expected", err)
	}
	// defer res.Body.Close()
	// body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		assert.Error(suite.T(), err, "err")
	}

	// Assert
	// assert.Equal(t, expectedBody, string(body))
}

func TestRequestsSuite(t *testing.T) {
	suite.Run(t, new(RequestsSuiteTest))
}
