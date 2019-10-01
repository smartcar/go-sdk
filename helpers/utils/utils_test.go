package utils

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper functions
func TestBuildVehicleURL(t *testing.T) {
	// Arrange
	ID := "vehicleId"
	path := "/path"
	expectedURL := GetVehicleURL() + ID + path

	// Act
	url := BuildVehicleURL(path, ID)

	// Assert
	assert.Equal(t, expectedURL, url)
}

func TestGetUserURL(t *testing.T) {
	// Arrange
	expectedURL := "https://api.smartcar.com/v1.0/user/"

	// Act
	url := GetUserURL()

	// Assert
	assert.Equal(t, expectedURL, url)
}

func TestGetVehicleURL(t *testing.T) {
	// Arrange
	expectedURL := "https://api.smartcar.com/v1.0/vehicles/"

	// Act
	url := GetVehicleURL()

	// Assert
	assert.Equal(t, expectedURL, url)
}

func TestBuildBearerAuthorization(t *testing.T) {
	// Arrange
	accessToken := "access"
	expectedAuthorization := fmt.Sprintf("Bearer %s", accessToken)

	// Act
	authorization := BuildBearerAuthorization(accessToken)

	// Assert
	assert.Equal(t, expectedAuthorization, authorization)
}

func TestConstructBasicAuthorization(t *testing.T) {
	// Arrange
	clientID := "clientID"
	clientSecret := "clientSecret"
	authString := fmt.Sprintf("%s:%s", clientID, clientSecret)
	expectedAuthorization := "Basic " + base64.StdEncoding.EncodeToString([]byte(authString))

	// Act
	authorization := BuildBasicAuthorization(clientID, clientSecret)

	// Assert
	assert.Equal(t, expectedAuthorization, authorization)
}

func TestBuildCompatibilityURL(t *testing.T) {
	// Arrange
	expectedURL := "https://api.smartcar.com/v1.0/compatibility/?scope=scope&vin=vin"

	// Act
	url := BuildCompatibilityURL("vin", []string{"scope"})

	// Assert
	assert.Equal(t, url, expectedURL)
}
