package smartcar

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/h2non/gock.v1"
)

type VehicleE2ETestSuite struct {
	suite.Suite
	vehicle                vehicle
	mockAge, mockRequestID string
	mockUnitSystem         UnitSystem
	responseHeaders        ResponseHeaders
}

func (s *VehicleE2ETestSuite) SetupTest() {
	s.vehicle = vehicle{
		id:          "client-id",
		accessToken: "access-token",
		client:      newBackend(),
	}
	s.mockAge = "data-age"
	s.mockRequestID = "request-id"
	s.mockUnitSystem = Metric
	s.responseHeaders = ResponseHeaders{
		Age:        s.mockAge,
		DataAge:    s.mockAge,
		RequestID:  s.mockRequestID,
		UnitSystem: s.mockUnitSystem,
	}
}

func (s *VehicleE2ETestSuite) TearDownTestSuite() {
	gock.Off()
}

func mockVehicleAPI(url, accessToken string, headers ResponseHeaders, response interface{}) {
	gock.New(url).
		MatchHeader("Authorization", buildBearerAuthorization(accessToken)).
		Reply(200).
		SetHeader("Sc-Data-Age", headers.Age).
		SetHeader("Sc-Request-Id", headers.RequestID).
		SetHeader("Sc-Unit-System", string(headers.UnitSystem)).
		JSON(response)
}

func (s *VehicleE2ETestSuite) TestBatchE2E() {
	mockLatitude := 37.4292
	mockLongitude := 122.1381
	mockDistance := 37829.0
	expectedResponse := &Data{
		Odometer: &Odometer{
			Distance: mockDistance,
			ResponseHeaders: ResponseHeaders{
				DataAge:    s.responseHeaders.Age,
				UnitSystem: s.responseHeaders.UnitSystem,
			},
		},
		Location: &Location{
			Latitude:  mockLatitude,
			Longitude: mockLongitude,
			ResponseHeaders: ResponseHeaders{
				DataAge: s.responseHeaders.Age,
			},
		},
	}
	mockURL := buildVehicleURL(string(batchPath), s.vehicle.id)
	mockResponse := map[string]interface{}{
		"responses": []interface{}{
			map[string]interface{}{
				"path": "/odometer",
				"body": map[string]interface{}{
					"distance": mockDistance,
				},
				"code": 200,
				"headers": map[string]interface{}{
					"sc-data-age":    s.responseHeaders.Age,
					"sc-unit-system": s.responseHeaders.UnitSystem,
				},
			},
			map[string]interface{}{
				"path": "/location",
				"body": map[string]interface{}{
					"latitude":  mockLatitude,
					"longitude": mockLongitude,
				},
				"code": 200,
				"headers": map[string]interface{}{
					"sc-data-age": s.responseHeaders.Age,
				},
			},
		},
	}
	mockVehicleAPI(mockURL, s.vehicle.accessToken, s.responseHeaders, mockResponse)

	res, err := s.vehicle.Batch(context.TODO(), OdometerPath, LocationPath)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedResponse, res)
}

func (s *VehicleE2ETestSuite) TestDisconnectE2E() {
	mockStatus := "success"
	expectedResponse := &Disconnect{
		Status:          mockStatus,
		ResponseHeaders: s.responseHeaders,
	}
	mockURL := buildVehicleURL(string(applicationPath), s.vehicle.id)
	mockResponse := map[string]interface{}{"status": mockStatus}
	mockVehicleAPI(mockURL, s.vehicle.accessToken, s.responseHeaders, mockResponse)

	res, err := s.vehicle.Disconnect(context.TODO())

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedResponse, res)
}

func (s *VehicleE2ETestSuite) TestGetBatteryE2E() {
	mockPercentRemaining := 0.3
	mockRange := 40.5
	expectedResponse := &Battery{
		PercentRemaining: mockPercentRemaining,
		Range:            mockRange,
		ResponseHeaders:  s.responseHeaders,
	}
	mockURL := buildVehicleURL(string(BatteryPath), s.vehicle.id)
	mockResponse := map[string]interface{}{"percentRemaining": mockPercentRemaining, "range": mockRange}
	mockVehicleAPI(mockURL, s.vehicle.accessToken, s.responseHeaders, mockResponse)

	res, err := s.vehicle.GetBattery(context.TODO())

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedResponse, res)
}

func (s *VehicleE2ETestSuite) TestGetChargeE2E() {
	mockIsPluggedIn := true
	mockState := "FULLY_CHARGED"
	expectedResponse := &Charge{
		IsPluggedIn:     mockIsPluggedIn,
		State:           mockState,
		ResponseHeaders: s.responseHeaders,
	}
	mockURL := buildVehicleURL(string(ChargePath), s.vehicle.id)
	mockResponse := map[string]interface{}{"isPluggedIn": mockIsPluggedIn, "state": mockState}
	mockVehicleAPI(mockURL, s.vehicle.accessToken, s.responseHeaders, mockResponse)

	res, err := s.vehicle.GetCharge(context.TODO())

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedResponse, res)
}

func (s *VehicleE2ETestSuite) TestGetFuelE2E() {
	mockAmountRemaining := 53.2
	mockPercentRemaining := 0.3
	mockRange := 40.5
	expectedResponse := &Fuel{
		AmountRemaining:  mockAmountRemaining,
		PercentRemaining: mockPercentRemaining,
		Range:            mockRange,
		ResponseHeaders:  s.responseHeaders,
	}
	mockURL := buildVehicleURL(string(FuelPath), s.vehicle.id)
	mockResponse := map[string]interface{}{"amountRemaining": mockAmountRemaining, "percentRemaining": mockPercentRemaining, "range": mockRange}
	mockVehicleAPI(mockURL, s.vehicle.accessToken, s.responseHeaders, mockResponse)

	res, err := s.vehicle.GetFuel(context.TODO())

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedResponse, res)
}

func (s *VehicleE2ETestSuite) TestGetInfoE2E() {
	mockYear := 2018
	mockID := "1234"
	mockMake := "TESLA"
	mockModel := "S"
	expectedResponse := &Info{
		ID:              mockID,
		Make:            mockMake,
		Model:           mockModel,
		Year:            mockYear,
		ResponseHeaders: s.responseHeaders,
	}
	mockURL := buildVehicleURL(string(InfoPath), s.vehicle.id)
	mockResponse := map[string]interface{}{"id": mockID, "make": mockMake, "model": mockModel, "year": mockYear}
	mockVehicleAPI(mockURL, s.vehicle.accessToken, s.responseHeaders, mockResponse)

	res, err := s.vehicle.GetInfo(context.TODO())

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedResponse, res)
}

func (s *VehicleE2ETestSuite) TestGetLocationE2E() {
	mockLatitude := 37.4292
	mockLongitude := 122.1381
	expectedResponse := &Location{
		Latitude:        mockLatitude,
		Longitude:       mockLongitude,
		ResponseHeaders: s.responseHeaders,
	}
	mockURL := buildVehicleURL(string(LocationPath), s.vehicle.id)
	mockResponse := map[string]interface{}{"latitude": mockLatitude, "longitude": mockLongitude}
	mockVehicleAPI(mockURL, s.vehicle.accessToken, s.responseHeaders, mockResponse)

	res, err := s.vehicle.GetLocation(context.TODO())

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedResponse, res)
}

func (s *VehicleE2ETestSuite) TestGetOdometerE2E() {
	mockValue := 15444.0232
	expectedResponse := &Odometer{
		Distance:        mockValue,
		ResponseHeaders: s.responseHeaders,
	}
	mockURL := buildVehicleURL(string(OdometerPath), s.vehicle.id)
	mockResponse := map[string]float64{"distance": mockValue}
	mockVehicleAPI(mockURL, s.vehicle.accessToken, s.responseHeaders, mockResponse)

	res, err := s.vehicle.GetOdometer(context.TODO())

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedResponse, res)
}

func (s *VehicleE2ETestSuite) TestGetOilE2E() {
	mockLifeRemaining := 0.35
	expectedResponse := &Oil{
		LifeRemaining:   mockLifeRemaining,
		ResponseHeaders: s.responseHeaders,
	}
	mockURL := buildVehicleURL(string(OilPath), s.vehicle.id)
	mockResponse := map[string]interface{}{"lifeRemaining": mockLifeRemaining}
	mockVehicleAPI(mockURL, s.vehicle.accessToken, s.responseHeaders, mockResponse)

	res, err := s.vehicle.GetOil(context.TODO())

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedResponse, res)
}

func (s *VehicleE2ETestSuite) TestGetPermissionsE2E() {
	mockPermissions := []string{"read_odometer", "read_batch"}
	expectedResponse := &Permissions{
		Permissions:     mockPermissions,
		ResponseHeaders: s.responseHeaders,
	}
	mockURL := buildVehicleURL(string(PermissionsPath), s.vehicle.id)
	mockResponse := map[string]interface{}{"permissions": mockPermissions}
	mockVehicleAPI(mockURL, s.vehicle.accessToken, s.responseHeaders, mockResponse)

	res, err := s.vehicle.GetPermissions(context.TODO())

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedResponse, res)
}

func (s *VehicleE2ETestSuite) TestGetTiresPressureE2E() {
	mockFrontLeft := 219.0
	mockFrontRight := 219.0
	mockBackLeft := 219.0
	mockBackRight := 219.0
	expectedResponse := &TirePressure{
		FrontLeft:       mockFrontLeft,
		FrontRight:      mockFrontRight,
		BackLeft:        mockBackLeft,
		BackRight:       mockBackRight,
		ResponseHeaders: s.responseHeaders,
	}
	mockURL := buildVehicleURL(string(TirePressurePath), s.vehicle.id)
	mockResponse := map[string]interface{}{
		"backLeft":   mockBackLeft,
		"backRight":  mockBackRight,
		"frontLeft":  mockFrontLeft,
		"frontRight": mockFrontRight,
	}
	mockVehicleAPI(mockURL, s.vehicle.accessToken, s.responseHeaders, mockResponse)

	res, err := s.vehicle.GetTiresPressure(context.TODO())

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedResponse, res)
}

func (s *VehicleE2ETestSuite) TestGetVINE2E() {
	mockVIN := "1234DLFAJ4"
	expectedResponse := &VIN{
		VIN:             mockVIN,
		ResponseHeaders: s.responseHeaders,
	}
	mockURL := buildVehicleURL(string(VINPath), s.vehicle.id)
	mockResponse := map[string]interface{}{"vin": mockVIN}
	mockVehicleAPI(mockURL, s.vehicle.accessToken, s.responseHeaders, mockResponse)

	res, err := s.vehicle.GetVIN(context.TODO())

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedResponse, res)
}

func (s *VehicleE2ETestSuite) TestLockE2E() {
	mockStatus := "1234DLFAJ4"
	expectedResponse := &Security{
		Status:          mockStatus,
		ResponseHeaders: s.responseHeaders,
	}
	mockURL := buildVehicleURL(string(securityPath), s.vehicle.id)
	mockResponse := map[string]interface{}{"status": mockStatus}
	mockVehicleAPI(mockURL, s.vehicle.accessToken, s.responseHeaders, mockResponse)

	res, err := s.vehicle.Lock(context.TODO())

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedResponse, res)
}

func (s *VehicleE2ETestSuite) TestUnlockE2E() {
	mockStatus := "1234DLFAJ4"
	expectedResponse := &Security{
		Status:          mockStatus,
		ResponseHeaders: s.responseHeaders,
	}
	mockURL := buildVehicleURL(string(securityPath), s.vehicle.id)
	mockResponse := map[string]interface{}{"status": mockStatus}
	mockVehicleAPI(mockURL, s.vehicle.accessToken, s.responseHeaders, mockResponse)

	res, err := s.vehicle.Unlock(context.TODO())

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedResponse, res)
}

func TestVehicleE2ETestSuite(t *testing.T) {
	suite.Run(t, new(VehicleE2ETestSuite))
}
