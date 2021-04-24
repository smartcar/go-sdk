package smartcar

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildBasicAuthorization(t *testing.T) {
	// Arrange
	clientID := "clientID"
	clientSecret := "clientSecret"
	authString := fmt.Sprintf("%s:%s", clientID, clientSecret)
	expectedAuthorization := "Basic " + base64.StdEncoding.EncodeToString([]byte(authString))

	// Act
	authorization := buildBasicAuthorization(clientID, clientSecret)

	// Assert
	assert.Equal(t, expectedAuthorization, authorization)
}

func TestBuildBearerAuthorization(t *testing.T) {
	// Arrange
	accessToken := "access"
	expectedAuthorization := fmt.Sprintf("Bearer %s", accessToken)

	// Act
	authorization := buildBearerAuthorization(accessToken)

	// Assert
	assert.Equal(t, expectedAuthorization, authorization)
}

func TestBuildCompatibilityURL(t *testing.T) {
	// Arrange
	expectedURL := "https://api.smartcar.com/v1.0/compatibility/?country=US&scope=scope&vin=vin"

	// Act
	url := buildCompatibilityURL("vin", []string{"scope"}, "")

	// Assert
	assert.Equal(t, url, expectedURL)
}

func TestBuildCompatibilityURLCountry(t *testing.T) {
	// Arrange
	expectedURL := "https://api.smartcar.com/v1.0/compatibility/?country=DE&scope=scope&vin=vin"

	// Act
	url := buildCompatibilityURL("vin", []string{"scope"}, "DE")

	// Assert
	assert.Equal(t, url, expectedURL)
}

func TestBuildCompatibilityURLVersion(t *testing.T) {
	// Arrange
	client := NewClient()
	client.SetAPIVersion("2.0")
	expectedURL := "https://api.smartcar.com/v2.0/compatibility/?country=US&scope=scope&vin=vin"

	// Act
	url := buildCompatibilityURL("vin", []string{"scope"}, "")

	// Assert
	assert.Equal(t, url, expectedURL)
}

func TestBuildVehicleURL(t *testing.T) {
	// Arrange
	ID := "vehicleId"
	path := "/path"
	expectedURL := fmt.Sprintf(vehicleURL, APIVersion) + ID + path

	// Act
	url := buildVehicleURL(path, ID)

	// Assert
	assert.Equal(t, expectedURL, url)
}

func TestBuildVehicleURLVersion(t *testing.T) {
	// Arrange
	client := NewClient()
	client.SetAPIVersion("2.0")
	ID := "vehicleId"
	path := "/path"
	expectedURL := fmt.Sprintf(vehicleURL, "2.0") + ID + path

	// Act
	url := buildVehicleURL(path, ID)

	// Assert
	assert.Equal(t, expectedURL, url)
}
