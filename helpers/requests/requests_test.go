package requests

import (
	"io/ioutil"
	"regexp"
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
	expectedMethod := GET
	authorization := "authorization"
	expectedBody := "Body"
	// This needs to be converted to regex expression because if not
	// MatchHeader does not work.
	expectedUserAgent := regexp.QuoteMeta(getUserAgent())
	gock.New(expectedURL).
		MatchHeader("Authorization", authorization).
		MatchHeader("User-Agent", expectedUserAgent).
		Get("/").
		Reply(200).
		BodyString(expectedBody)

	// Act
	res, err := Request(expectedMethod, expectedURL, authorization, "", nil)
	if err != nil {
		suite.T().Error("Should not be called")
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
	expectedURL := "http://example.com"
	authorization := "authorization"
	expectedUserAgent := regexp.QuoteMeta(getUserAgent())
	gock.New(expectedURL).
		MatchHeader("Authorization", authorization).
		MatchHeader("User-Agent", expectedUserAgent).
		Get("/").
		Reply(401)

	// Act
	_, err := Request(GET, expectedURL, authorization, "", nil)

	// Assert
	if err != nil {
		assert.EqualError(suite.T(), err, "Unauthorized")
	} else {
		suite.T().Error("Should not be called")
	}
}

func TestRequestsSuite(t *testing.T) {
	suite.Run(t, new(RequestsSuiteTest))
}
