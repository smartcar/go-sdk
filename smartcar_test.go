package smartcar

import (
	"net/url"
	"testing"

	"github.com/smartcar/go-sdk/helpers/constants"
	"github.com/smartcar/go-sdk/helpers/requests"
	"github.com/smartcar/go-sdk/helpers/test"
	"github.com/smartcar/go-sdk/helpers/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SmartcarTestSuite struct {
	suite.Suite
	userURL string
}

func (suite *SmartcarTestSuite) SetupSuite() {
	URL := url.URL{
		Scheme: constants.APIScheme,
		Host:   constants.APIHost,
		Path:   constants.UserPath,
	}
	suite.userURL = URL.String()
}

func (suite *SmartcarTestSuite) AfterTest() {
	test.ClearMock()
}

func (suite *SmartcarTestSuite) TestGetUserID() {
	// Arrange
	accessToken := "accessToken"
	expectedUserID := "userId"
	res := map[string]string{"id": expectedUserID}
	authorization := utils.BuildBearerAuthorization(accessToken)

	test.MockRequest(requests.GET, suite.userURL, authorization, 200, res)

	// Act
	userID, _ := GetUserID(accessToken)

	// Assert
	assert.Equal(suite.T(), userID, expectedUserID)
}

func (suite *SmartcarTestSuite) TestGetUserIDEmptyAccess() {
	// Arrange
	userURL := utils.GetUserURL()
	accessToken := ""

	test.MockRequest(requests.GET, userURL, accessToken, 401, nil)

	// Act
	_, err := GetUserID(accessToken)

	// Assert
	if err != nil {
		assert.EqualError(suite.T(), err, "Unauthorized")
	} else {
		suite.T().Error("Should not throw here")
	}
}

func (suite *SmartcarTestSuite) TestGetVehicleIds() {
	// Arrange
	vehicleURL := utils.GetVehicleURL()
	accessToken := "access-token"
	expectedVehicleIds := []string{"vehicleId", "vehicleId2"}
	res := map[string][]string{"vehicles": expectedVehicleIds}

	test.MockRequest(requests.GET, vehicleURL, accessToken, 200, res)

	// Act
	vehicleIds, _ := GetVehicleIds(accessToken)

	// Assert
	assert.ElementsMatch(suite.T(), vehicleIds, expectedVehicleIds)
}

func (suite *SmartcarTestSuite) TestGetVehicleIdsEmptyAccess() {
	// Arrange
	vehicleURL := utils.GetVehicleURL()
	accessToken := ""

	test.MockRequest(requests.GET, vehicleURL, accessToken, 401, nil)

	// Act
	_, err := GetVehicleIds(accessToken)

	// Assert
	if err != nil {
		assert.EqualError(suite.T(), err, "Unauthorized")
	} else {
		suite.T().Error("Should not throw here")
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSmartcarTestSuite(t *testing.T) {
	suite.Run(t, new(SmartcarTestSuite))
}
