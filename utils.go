package smartcar

import (
	"encoding/base64"
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
	baseURL, _ := url.Parse(compatibilityURL)
	query := baseURL.Query()
	query.Set("vin", vin)
	query.Set("scope", strings.Join(scope, " "))
	if len(country) > 0 {
		query.Set("country", country)
	}
	baseURL.RawQuery = query.Encode()

	return baseURL.String()
}

// buildVehicleURL buids a vehicle URL with a path and ID
func buildVehicleURL(path, ID string) string {
	return vehicleURL + ID + path
}
