package smartcar

import (
	"net/url"
	"testing"

	"github.com/smartcar/go-sdk/helpers/constants"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestGetUserID(t *testing.T) {
	// Arrange
	userURL := url.URL{
		Scheme: constants.APIScheme,
		Host:   constants.APIHost,
		Path:   constants.UserPath,
	}
	accessToken := "access-token"
	expectedUserID := "userId"
	res := map[string]string{"id": expectedUserID}

	defer gock.Off()
	gock.New(userURL.String()).
		MatchHeader("Authorization", accessToken).
		Get("/").
		Reply(200).
		JSON(res)

	// Act
	userID, _ := GetUserID(accessToken)

	// Assert
	assert.Equal(t, userID, expectedUserID)
}

func TestGetUserIDEmptyAccess(t *testing.T) {
	// Arrange
	userURL := url.URL{
		Scheme: constants.APIScheme,
		Host:   constants.APIHost,
		Path:   constants.UserPath,
	}
	accessToken := ""

	defer gock.Off()
	gock.New(userURL.String()).
		MatchHeader("Authorization", accessToken).
		Get("/").
		Reply(401)

	// Act
	_, err := GetUserID(accessToken)

	// Assert
	if err != nil {
		assert.EqualError(t, err, "Unauthorized")
	} else {
		t.Error("Should not throw here")
	}
}

func TestGetVehicleIds(t *testing.T) {
	// Arrange
	vehiclesURL := url.URL{
		Scheme: constants.APIScheme,
		Host:   constants.APIHost,
		Path:   constants.VehiclePath,
	}
	accessToken := "access-token"
	expectedVehicleIds := []string{"vehicleId", "vehicleId2"}
	res := map[string][]string{"vehicles": expectedVehicleIds}

	defer gock.Off()
	gock.New(vehiclesURL.String()).
		MatchHeader("Authorization", accessToken).
		Get("/").
		Reply(200).
		JSON(res)

	// Act
	vehicleIds, _ := GetVehicleIds(accessToken)

	// Assert
	assert.ElementsMatch(t, vehicleIds, expectedVehicleIds)
}

func TestGetVehicleIdsEmptyAccess(t *testing.T) {
	// Arrange
	vehiclesURL := url.URL{
		Scheme: constants.APIScheme,
		Host:   constants.APIHost,
		Path:   constants.VehiclePath,
	}
	accessToken := ""

	defer gock.Off()
	gock.New(vehiclesURL.String()).
		MatchHeader("Authorization", accessToken).
		Get("/").
		Reply(401)

	// Act
	_, err := GetVehicleIds(accessToken)

	// Assert
	if err != nil {
		assert.EqualError(t, err, "Unauthorized")
	} else {
		t.Error("Should not throw here")
	}
}
