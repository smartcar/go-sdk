package smartcar

import (
	"net/url"

	"github.com/smartcar/go-sdk/helpers/constants"
)

// GetURL exports compatibility URL
func GetURL() string {
	compatiblityURL := url.URL{
		Scheme: constants.APIScheme,
		Host:   constants.APIHost,
		Path:   constants.UserPath,
	}

	return compatiblityURL.String()
}
