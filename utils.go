package smartcar

import (
	"encoding/base64"
	"net/url"
	"strings"
)

// BuildBasicAuthorization buids a basic access authentication
func BuildBasicAuthorization(id, secret string) string {
	authString := id + ":" + secret

	return "Basic " + base64.StdEncoding.EncodeToString([]byte(authString))
}

// BuildBearerAuthorization buids a bearer authentication
func BuildBearerAuthorization(accessToken string) string {
	return "Bearer " + accessToken
}

// BuildCompatibilityURL based on vin and scope
func BuildCompatibilityURL(vin string, scope []string) string {
	baseURL, _ := url.Parse(compatibilityURL)
	query := baseURL.Query()
	query.Set("vin", vin)
	query.Set("scope", strings.Join(scope, " "))
	baseURL.RawQuery = query.Encode()

	return baseURL.String()
}

// BuildVehicleURL buids a vehicle URL with a path and ID
func BuildVehicleURL(path, ID string) string {
	return vehicleURL + ID + path
}
