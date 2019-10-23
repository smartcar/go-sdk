package smartcar

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SmartcarTestSuite struct {
	suite.Suite
	client client
}

type fakeSmartcarClient struct{}

func (s *fakeSmartcarClient) Call(p backendClientParams) error {
	return nil
}

func newFakeSmartcarClient() backendClient {
	return &fakeSmartcarClient{}
}

func (s *SmartcarTestSuite) SetupTest() {
	s.client = client{
		sC: newFakeVehicleClient(),
	}
}

func (s *SmartcarTestSuite) TestGetUserID() {
	res, err := s.client.GetUserID(context.TODO(), &UserIDParams{})

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), res)
}

func (s *SmartcarTestSuite) TestGetVehicleIDs() {
	res, err := s.client.GetVehicleIDs(context.TODO(), &VehicleIDsParams{})

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), res)
}

func (s *SmartcarTestSuite) TestIsTokenExpiredFalse() {
	expired := s.client.IsTokenExpired(&TokenExpiredParams{time.Now()})

	assert.False(s.T(), expired)
}

func (s *SmartcarTestSuite) TestIsTokenExpiredTrue() {
	timeBefore := time.Now().Add(time.Second * -10)
	expired := s.client.IsTokenExpired(&TokenExpiredParams{Expiry: timeBefore})

	assert.True(s.T(), expired)
}

func (s *SmartcarTestSuite) TestIsVINCompatible() {
	res, err := s.client.IsVINCompatible(context.TODO(), &VINCompatibleParams{})

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), res)
}

// TODO: implement this tests mocking vehicle.GetPermissions
// func (s *SmartcarTestSuite) TestHasPermissions() {
// 	v := vehicle{}
// 	res, err := s.client.HasPermissions(context.TODO(), &v, &PermissionsParams{})

// 	assert.Nil(s.T(), err)
// 	assert.NotNil(s.T(), res)
// }

func (s *SmartcarTestSuite) TestNewVehicleEmpty() {
	res := s.client.NewVehicle(&VehicleParams{})

	expectedVehicle := &vehicle{
		client:        s.client.sC,
		requestParams: requestParams{UnitSystem: Metric},
	}

	assert.NotNil(s.T(), res)
	assert.Equal(s.T(), expectedVehicle, res)
}

func (s *SmartcarTestSuite) TestNewVehicle() {
	res := s.client.NewVehicle(&VehicleParams{UnitSystem: Imperial})

	expectedVehicle := &vehicle{
		client:        s.client.sC,
		requestParams: requestParams{UnitSystem: Imperial},
	}

	assert.NotNil(s.T(), res)
	assert.Equal(s.T(), expectedVehicle, res)
}

func (s *SmartcarTestSuite) TestNewAuth() {
	res := s.client.NewAuth(&AuthParams{})

	assert.NotNil(s.T(), res)
}

func (s *SmartcarTestSuite) TestNewClient() {
	res := NewClient()

	assert.NotNil(s.T(), res)
}

func TestSmartcarTestSuite(t *testing.T) {
	suite.Run(t, new(SmartcarTestSuite))
}
