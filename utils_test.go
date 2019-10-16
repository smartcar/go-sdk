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
	authorization := BuildBasicAuthorization(clientID, clientSecret)

	// Assert
	assert.Equal(t, expectedAuthorization, authorization)
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

func TestBuildCompatibilityURL(t *testing.T) {
	// Arrange
	expectedURL := "https://api.smartcar.com/v1.0/compatibility/?scope=scope&vin=vin"

	// Act
	url := BuildCompatibilityURL("vin", []string{"scope"})

	// Assert
	assert.Equal(t, url, expectedURL)
}

func TestBuildVehicleURL(t *testing.T) {
	// Arrange
	ID := "vehicleId"
	path := "/path"
	expectedURL := vehicleURL + ID + path

	// Act
	url := BuildVehicleURL(path, ID)

	// Assert
	assert.Equal(t, expectedURL, url)
}
