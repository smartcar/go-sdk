package smartcar

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/h2non/gock.v1"
)

type SmartcarE2ETestSuite struct {
	suite.Suite
	client client
}

func (s *SmartcarE2ETestSuite) SetupTest() {
	s.client = client{
		requestParams: requestParams{},
		sC:            newBackend(),
	}
}

func (s *SmartcarE2ETestSuite) TearDownTestSuite() {
	gock.Off()
}

func mockSmartcarAPI(url, authorization string, response interface{}) {
	gock.New(url).
		MatchHeader("Authorization", authorization).
		Reply(200).
		JSON(response)
}

func (s *SmartcarE2ETestSuite) TestGetUserIDE2E() {
	mockUserID := "user-id"
	expectedResponse := mockUserID
	mockAccess := "mock-access"
	mockResponse := map[string]interface{}{
		"id": mockUserID,
	}
	mockSmartcarAPI(userURL, buildBearerAuthorization(mockAccess), mockResponse)

	res, err := s.client.GetUserID(context.TODO(), &UserIDParams{
		Access: mockAccess,
	})

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedResponse, *res)
}

func (s *SmartcarE2ETestSuite) TestGetVehicleIDsE2E() {
	mockVehicleIDs := []string{"vehicle-id1", "vehicle-id2"}
	mockAccess := "mock-access"
	expectedResponse := mockVehicleIDs
	mockResponse := map[string]interface{}{
		"vehicles": mockVehicleIDs,
	}
	mockSmartcarAPI(vehicleURL, buildBearerAuthorization(mockAccess), mockResponse)

	res, err := s.client.GetVehicleIDs(context.TODO(), &VehicleIDsParams{
		Access: mockAccess,
	})

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedResponse, *res)
}

func (s *SmartcarE2ETestSuite) TestHasPermissionsE2E() {
	expectedResponse := true
	mockAccess := "mock-access"
	mockVehiclePermissions := []string{"read_odometer", "read_location"}
	mockResponse := map[string]interface{}{
		"permissions": mockVehiclePermissions,
	}
	mockID := "id"
	mockVehicle := s.client.NewVehicle(&VehicleParams{
		ID:          mockID,
		AccessToken: mockAccess,
	})
	mockSmartcarAPI(buildVehicleURL(string(PermissionsPath), mockID), buildBearerAuthorization(mockAccess), mockResponse)

	res, err := s.client.HasPermissions(context.TODO(), mockVehicle, &PermissionsParams{
		Permissions: mockVehiclePermissions,
	})

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedResponse, res)
}

func (s *SmartcarE2ETestSuite) TestHasPermissionsE2EMissingPermissions() {
	expectedResponse := false
	mockAccess := "mock-access"
	mockVehiclePermissions := []string{"read_odometer", "read_location"}
	mockResponse := map[string]interface{}{
		"permissions": mockVehiclePermissions,
	}
	mockID := "id"
	mockVehicle := s.client.NewVehicle(&VehicleParams{
		ID:          mockID,
		AccessToken: mockAccess,
	})
	mockSmartcarAPI(buildVehicleURL(string(PermissionsPath), mockID), buildBearerAuthorization(mockAccess), mockResponse)

	res, err := s.client.HasPermissions(context.TODO(), mockVehicle, &PermissionsParams{
		Permissions: []string{"read_odometer", "read_location", "read_battery"},
	})

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedResponse, res)
}

func (s *SmartcarE2ETestSuite) TestIsVINCompatibleE2E() {
	mockCompatibility := true
	expectedResponse := mockCompatibility
	mockID := "mock-id"
	mockSecret := "mock-secret"
	mockVIN := "mock-vin"
	mockScope := []string{"read_odometer", "read_location"}
	mockResponse := map[string]interface{}{
		"compatible": mockCompatibility,
	}
	mockURL := buildCompatibilityURL(mockVIN, mockScope)
	mockSmartcarAPI(mockURL, buildBasicAuthorization(mockID, mockSecret), mockResponse)

	res, err := s.client.IsVINCompatible(context.TODO(), &VINCompatibleParams{
		VIN:          mockVIN,
		Scope:        mockScope,
		ClientID:     mockID,
		ClientSecret: mockSecret,
	})

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedResponse, res)
}

func TestSmartcarE2ETestSuite(t *testing.T) {
	suite.Run(t, new(SmartcarE2ETestSuite))
}
