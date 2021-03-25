// Package smartcar is the official Go SDK of the Smartcar API.
// Smartcar is the only vehicle API built for developers, by developers.
// Learn more about Smartcar here, https://smartcar.com/
package smartcar

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	apiURL2           = "https://api.smartcar.com/v1.0/"
	userURL2          = "https://api.smartcar.com/v1.0/user/"
	vehicleURL2       = "https://api.smartcar.com/v1.0/vehicles/"
	compatibilityURL2 = "https://api.smartcar.com/v1.0/compatibility/"
	exchangeURL       = "https://auth.smartcar.com/oauth/token/"
	apiURL            = "https://api.smartcar.com/v%s/"
	userURL           = "https://api.smartcar.com/v%s/user/"
	vehicleURL        = "https://api.smartcar.com/v%s/vehicles/"
	compatibilityURL  = "https://api.smartcar.com/v%s/compatibility/"
)

// APIVersion is the default version of API to use
var APIVersion string = "1.0"

// UserIDParams is a param in client.GetUserID
type UserIDParams struct {
	Access string
}

// VehicleIDsParams is a param in client.GetVehicleIDs
type VehicleIDsParams struct {
	Access string
}

// TokenExpiredParams is a param in client.IsTokenExpired
type TokenExpiredParams struct {
	Expiry time.Time
}

// VehicleParams is a param in client.NewVehicle
type VehicleParams struct {
	ID          string
	AccessToken string
	UnitSystem  UnitSystem
}

// AuthParams is a param in client.NewAuth
type AuthParams struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Scope        []string
	TestMode     bool
}

// VINCompatibleParams is a param in client.IsVINCompatible
type VINCompatibleParams struct {
	VIN          string
	Scope        []string
	Country      string
	ClientID     string
	ClientSecret string
}

// PermissionsParams is a param in client.HasPermissions
type PermissionsParams struct {
	Permissions []string
}

// GetUserID returns the user ID of the vehicle owner associated with an Access token.
func (c *client) GetUserID(ctx context.Context, params *UserIDParams) (*string, error) {
	target := new(struct {
		ID string
	})
	authorization := buildBearerAuthorization(params.Access)
	versionedUserURL := fmt.Sprintf(userURL, APIVersion)

	return &target.ID, c.sC.Call(backendClientParams{
		ctx:           ctx,
		method:        http.MethodGet,
		url:           versionedUserURL,
		authorization: authorization,
		target:        target,
	})
}

// GetVehicleIds returns IDs of the vehicles associated with an Access token.
func (c *client) GetVehicleIDs(ctx context.Context, params *VehicleIDsParams) (*[]string, error) {
	target := new(struct {
		VehicleIDs []string `json:"vehicles"`
	})
	authorization := buildBearerAuthorization(params.Access)
	versionedVehicleURL := fmt.Sprintf(vehicleURL, APIVersion)

	return &target.VehicleIDs, c.sC.Call(backendClientParams{
		ctx:           ctx,
		method:        http.MethodGet,
		url:           versionedVehicleURL,
		authorization: authorization,
		target:        target,
	})
}

// IsTokenExpired checks if Expiry is expired.
// Note: Does not call Smartcar's API nor makes an http.Request.
func (c *client) IsTokenExpired(params *TokenExpiredParams) bool {
	return time.Now().After(params.Expiry.Add(time.Second * 10))
}

// IsVINCompatible checks if a VIN is compatible for a list scopes.
func (c *client) IsVINCompatible(ctx context.Context, params *VINCompatibleParams) (bool, error) {
	url := buildCompatibilityURL(params.VIN, params.Scope, params.Country)

	isCompatible := new(struct {
		Compatible bool
	})
	authorization := buildBasicAuthorization(params.ClientID, params.ClientSecret)

	return isCompatible.Compatible, c.sC.Call(backendClientParams{
		ctx:           ctx,
		method:        http.MethodGet,
		url:           url,
		authorization: authorization,
		target:        isCompatible,
	})
}

// HasPermissions checks if the vehicle has the permissions passed in.
func (c *client) HasPermissions(ctx context.Context, v Vehicle, params *PermissionsParams) (bool, error) {
	vehiclePermissions, err := v.GetPermissions(ctx)
	if err != nil {
		return false, err
	}

	set := make(map[string]bool)
	for _, value := range vehiclePermissions.Permissions {
		set[value] = true
	}

	for _, value := range params.Permissions {
		value = strings.TrimPrefix(value, "required:")
		if hasPermission, found := set[value]; !found {
			return false, nil
		} else if !hasPermission {
			return false, nil
		}
	}

	return true, nil
}

// NewVehicle creates an instance of Vehicle that allows you to call methods (i.e. GetInfo, GetOdometer, etc) on it and
// send requests to Smartcar's API.
func (c *client) NewVehicle(params *VehicleParams) Vehicle {
	unitSystem := Metric
	if params.UnitSystem != "" {
		unitSystem = params.UnitSystem
	}
	return &vehicle{
		id:            params.ID,
		accessToken:   params.AccessToken,
		client:        c.sC,
		requestParams: requestParams{UnitSystem: unitSystem},
	}
}

// NewAuthClient creates an instance of Auth that allows you to call methods that relate to authentication in Smartcar's API.
func (c *client) NewAuth(params *AuthParams) Auth {
	return &auth{
		clientID:     params.ClientID,
		clientSecret: params.ClientSecret,
		redirectURI:  params.RedirectURI,
		scope:        params.Scope,
		testMode:     params.TestMode,
		sC:           c.sC,
	}
}

// SetAPIVersion sets version of Smartcar API to use
func SetAPIVersion(version string) {
	APIVersion = version
}

// Backend exposes methods needed for executing requests to Smartcar's API.
type backendClient interface {
	Call(backendClientParams) error
}

// backend is an internal helper struct that implements Backend.
type backend struct{}

// getBackend returns a newly created backend.
func newBackend() backendClient {
	return &backend{}
}

type client struct {
	requestParams
	sC backendClient
}

// Client exposes methods that allow you to interact with Smartcar's API that are not part of Vehicle or Auth.
type Client interface {
	GetUserID(context.Context, *UserIDParams) (*string, error)
	GetVehicleIDs(context.Context, *VehicleIDsParams) (*[]string, error)
	IsTokenExpired(*TokenExpiredParams) bool
	IsVINCompatible(context.Context, *VINCompatibleParams) (bool, error)
	HasPermissions(context.Context, Vehicle, *PermissionsParams) (bool, error)
	NewAuth(*AuthParams) Auth
	NewVehicle(*VehicleParams) Vehicle
}

// NewClient creates new SmartcarClient. This is the entry point for communicating with Smartcar's API.
// Note: You cannot use any of the methods on this SDK if you don't call this method.
func NewClient() Client {
	return &client{sC: newBackend()}
}
