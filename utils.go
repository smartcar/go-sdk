package smartcar

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
)

// buildBasicAuthorization buids a basic access authentication
func buildBasicAuthorization(id, secret string) string {
	authString := id + ":" + secret

	return "Basic " + base64.StdEncoding.EncodeToString([]byte(authString))
}

// buildBearerAuthorization buids a bearer authentication
func buildBearerAuthorization(accessToken string) string {
	return "Bearer " + accessToken
}

// buildCompatibilityURL based on vin and scope
func buildCompatibilityURL(vin string, scope []string, country string) string {
	versionedURL := fmt.Sprintf(compatibilityURL, APIVersion)
	baseURL, _ := url.Parse(versionedURL)
	query := baseURL.Query()
	query.Set("vin", vin)
	query.Set("scope", strings.Join(scope, " "))
	if len(country) > 0 {
		query.Set("country", country)
	} else {
		query.Set("country", "US")
	}
	baseURL.RawQuery = query.Encode()

	return baseURL.String()
}

// buildVehicleURL buids a vehicle URL with a path and ID
func buildVehicleURL(path, ID string) string {
	versionedVehicleURL := fmt.Sprintf(vehicleURL, APIVersion)
	return versionedVehicleURL + ID + path
}
