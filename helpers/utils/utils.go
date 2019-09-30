package utils

import (
	"encoding/base64"
	"net/url"
	"strings"

	"github.com/smartcar/go-sdk/helpers/constants"
)

// BuildVehicleURL buids a vehicle URL with a path and ID
func BuildVehicleURL(path, ID string) string {
	return GetVehicleURL() + ID + path
}

// GetUserURL returns user URL
func GetUserURL() string {
	URL := url.URL{
		Scheme: constants.APIScheme,
		Host:   constants.APIHost,
		Path:   constants.UserPath,
	}
	return URL.String()
}

// GetVehicleURL returns vehicles URL
func GetVehicleURL() string {
	URL := url.URL{
		Scheme: constants.APIScheme,
		Host:   constants.APIHost,
		Path:   constants.VehiclePath,
	}
	return URL.String()
}

// GetCompatibilityURL returns compatibility URL
func GetCompatibilityURL() url.URL {
	return url.URL{
		Scheme: constants.APIScheme,
		Host:   constants.APIHost,
		Path:   constants.CompatibilityPath,
	}
}

// GetConnectURL returns connect URL
func GetConnectURL() url.URL {
	return url.URL{
		Scheme: constants.ConnectScheme,
		Host:   constants.ConnectHost,
		Path:   constants.ConnectPath,
	}
}

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
	compatibilityURL := GetCompatibilityURL()
	query := compatibilityURL.Query()
	query.Set("vin", vin)
	query.Set("scope", strings.Join(scope, " "))
	compatibilityURL.RawQuery = query.Encode()

	return compatibilityURL.String()
}
