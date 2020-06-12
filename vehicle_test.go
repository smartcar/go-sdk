package smartcar

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type VehicleTestSuite struct {
	suite.Suite
	vehicle vehicle
}

type fakeVehicleClient struct{}

func (c *fakeVehicleClient) Call(params backendClientParams) error {
	return nil
}

func newFakeVehicleClient() backendClient {
	return &fakeVehicleClient{}

}

func (s *VehicleTestSuite) SetupTest() {
	s.vehicle = vehicle{
		id:          "client-id",
		accessToken: "client-secret",
		client:      newFakeVehicleClient(),
	}
}

func (s *VehicleTestSuite) TestBatch() {
	res, err := s.vehicle.Batch(context.TODO(), OdometerPath, BatteryPath)

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), res)
}

func (s *VehicleTestSuite) TestDisconnect() {
	res, err := s.vehicle.Disconnect(context.TODO())

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), res)
}

func (s *VehicleTestSuite) TestGetBattery() {
	res, err := s.vehicle.GetBattery(context.TODO())

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), res)
}

func (s *VehicleTestSuite) TestGetCharge() {
	res, err := s.vehicle.GetCharge(context.TODO())

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), res)
}

func (s *VehicleTestSuite) TestGetFuel() {
	res, err := s.vehicle.GetFuel(context.TODO())

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), res)
}

func (s *VehicleTestSuite) TestGetInfo() {
	res, err := s.vehicle.GetInfo(context.TODO())

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), res)
}

func (s *VehicleTestSuite) TestGetLocation() {
	res, err := s.vehicle.GetLocation(context.TODO())

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), res)
}

func (s *VehicleTestSuite) TestGetOdometer() {
	res, err := s.vehicle.GetOdometer(context.TODO())

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), res)
}

func (s *VehicleTestSuite) TestGetOil() {
	res, err := s.vehicle.GetOil(context.TODO())

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), res)
}

func (s *VehicleTestSuite) TestGetPermissions() {
	res, err := s.vehicle.GetPermissions(context.TODO())

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), res)
}

func (s *VehicleTestSuite) TestGetTiresPressure() {
	res, err := s.vehicle.GetTiresPressure(context.TODO())

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), res)
}

func (s *VehicleTestSuite) TestGetVIN() {
	res, err := s.vehicle.GetVIN(context.TODO())

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), res)
}

func (s *VehicleTestSuite) TestLock() {
	res, err := s.vehicle.Lock(context.TODO())

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), res)
}

func (s *VehicleTestSuite) TestSetUnitSystem() {
	err := s.vehicle.SetUnitSystem(&UnitsParams{Units: Imperial})

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), Imperial, s.vehicle.UnitSystem)

	err = s.vehicle.SetUnitSystem(&UnitsParams{Units: Metric})

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), Metric, s.vehicle.UnitSystem)
}

func (s *VehicleTestSuite) TestSetUnitSystemError() {
	err := s.vehicle.SetUnitSystem(&UnitsParams{Units: ""})

	assert.NotNil(s.T(), err)
}

func (s *VehicleTestSuite) TestUnlock() {
	res, err := s.vehicle.Unlock(context.TODO())

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), res)
}

func (s *VehicleTestSuite) TestStartCharge() {
	res, err := s.vehicle.StartCharge(context.TODO())

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), res)
}

func (s *VehicleTestSuite) TestStopCharge() {
	res, err := s.vehicle.StopCharge(context.TODO())

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), res)
}

func TestVehicleTestSuite(t *testing.T) {
	suite.Run(t, new(VehicleTestSuite))
}
