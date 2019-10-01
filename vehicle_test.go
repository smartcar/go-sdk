package smartcar

import (
	"io/ioutil"
	"testing"

	"github.com/smartcar/go-sdk/helpers/requests"
	"github.com/smartcar/go-sdk/helpers/test"
	"github.com/smartcar/go-sdk/helpers/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type VehicleTestSuite struct {
	suite.Suite
	vehicleID          string
	vehicleMake        string
	vehicleModel       string
	vehicleYear        int
	vehicleAccessToken string
	authorization      string
	vehicle            Vehicle
}

func (suite *VehicleTestSuite) SetupSuite() {
	suite.vehicleID = "id"
	suite.vehicleMake = "Tesla"
	suite.vehicleModel = "S"
	suite.vehicleYear = 2019
	suite.vehicleAccessToken = "accessToken"
	suite.authorization = "Bearer " + suite.vehicleAccessToken
}

func (suite *VehicleTestSuite) SetupTest() {
	suite.vehicle = Vehicle{
		ID:          suite.vehicleID,
		AccessToken: suite.vehicleAccessToken,
	}
}

func (suite *VehicleTestSuite) AfterTest() {
	test.ClearMock()
}

func (suite *VehicleTestSuite) TestSetUnitsError() {
	// Arrange
	vehicle := Vehicle{}

	// Act
	err := vehicle.SetUnits("unit")

	// Assert
	if err == nil {
		suite.T().Error("Should not throw")
	}
	assert.EqualError(suite.T(), err, "unit must either be metric or imperial")
	assert.Equal(suite.T(), vehicle.UnitSystem, "")
}

func (suite *VehicleTestSuite) TestSetUnitsMetric() {
	// Arrange
	units := "metric"

	// Act
	err := suite.vehicle.SetUnits(units)

	// Assert
	if err != nil {
		suite.T().Error("Should not throw")
	}
	assert.Equal(suite.T(), suite.vehicle.UnitSystem, units)
}

func (suite *VehicleTestSuite) TestSetUnitsImperial() {
	// Arrange
	units := "imperial"

	// Act
	err := suite.vehicle.SetUnits(units)

	// Assert
	if err != nil {
		suite.T().Error("Should not throw")
	}
	assert.Equal(suite.T(), suite.vehicle.UnitSystem, units)
}

func (suite *VehicleTestSuite) TestInfo() {
	// Arrange
	expectedResponse := VehicleInfoResponse{
		ID:    suite.vehicleID,
		Make:  suite.vehicleMake,
		Model: suite.vehicleModel,
		Year:  suite.vehicleYear,
	}
	mockResponse := map[string]interface{}{
		"id":    suite.vehicleID,
		"make":  suite.vehicleMake,
		"model": suite.vehicleModel,
		"year":  suite.vehicleYear,
	}
	url := utils.BuildVehicleURL("/", suite.vehicleID)
	test.MockRequest(requests.GET, url, suite.authorization, 200, mockResponse)

	// Act
	res, err := suite.vehicle.Info()

	// Assert
	if err != nil {
		suite.T().Error(err)
	}
	assert.Equal(suite.T(), expectedResponse, res)
}

func (suite *VehicleTestSuite) TestVIN() {
	// Arrange
	vin := "1234567890ABCDEFG"
	expectedResponse := vin
	mockResponse := map[string]interface{}{
		"vin": vin,
	}
	url := utils.BuildVehicleURL("/vin", suite.vehicleID)
	test.MockRequest(requests.GET, url, suite.authorization, 200, mockResponse)

	// Act
	res, err := suite.vehicle.VIN()

	// Assert
	if err != nil {
		suite.T().Error(err)
	}
	assert.Equal(suite.T(), res, expectedResponse)
}

func (suite *VehicleTestSuite) TestLock() {
	// Arrange
	status := "success"
	expectedResponse := VehicleResponse{Status: status}
	mockResponse := map[string]interface{}{
		"status": status,
	}
	url := utils.BuildVehicleURL("/security", suite.vehicleID)
	test.MockRequest(requests.POST, url, suite.authorization, 200, mockResponse)

	// Act
	res, err := suite.vehicle.Lock()

	// Assert
	if err != nil {
		suite.T().Error(err)
	}
	assert.Equal(suite.T(), res, expectedResponse)
}

func (suite *VehicleTestSuite) TestRequest() {
	// Arrange
	path := "/"
	url := utils.BuildVehicleURL(path, suite.vehicle.ID)
	expectedBody := "It worked"

	test.MockRequest(requests.GET, url, suite.authorization, 200, expectedBody)

	// Act
	req, err := suite.vehicle.request(path, requests.GET, nil)

	// Assert
	if err != nil {
		suite.T().Error(err)
	}
	body, _ := ioutil.ReadAll(req.Body)
	assert.Equal(suite.T(), 200, req.StatusCode)
	assert.Equal(suite.T(), expectedBody, string(body))
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestVehicleTestSuite(t *testing.T) {
	suite.Run(t, new(VehicleTestSuite))
}
