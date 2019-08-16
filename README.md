# Smartcar Go SDK
[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/smartcar/go-sdk)

## Overview

The [Smartcar API](https://smartcar.com/docs) lets you read vehicle data (location, odometer) and send commands to vehicles (lock, unlock) to connected vehicles using HTTP requests.

To make requests to a vehicle a web or mobile application, the end user must connect their vehicle using [Smartcar Connect](https://smartcar.com/docs/api#authorization).
This flow follows the OAuth spec and will return a `code` which can be used to
obtain an access token from Smartcar.

The Smartcar Go SDK provides methods to:

1. Generate the link to redirect to Connect.
2. Make a request to Smartcar with the `code` obtained from Connect to obtain
   access and refresh tokens.
3. Make requests to the Smartcar API to read vehicle data and send commands to
   vehicles using the access token obtained in step 2.

Before integrating with the Go SDK, you'll need to register an application in the [Smartcar Developer portal](https://dashboard.smartcar.com). Once you have registered an application, you will have a Client ID and Client Secret, which will allow you to authorize users.


## Installation
The Smartcar Go SDK is built using Go modules and is easy to install.

Install the smartcar package if you are not using Go modules:
```
go get -u github.com/smartcar/go-sdk
```
Import it using:
```
import smartcar "github.com/smartcar/go-sdk"
```

## Getting Started
After obtaining your client credentials and redirect URI from the Smartcar Dashboard, here's a small guide on getting setup to make requests with the Smartcar API.

1. First initialize an `AuthClient` struct with your client credentials, redirect URI, the scopes required and test mode.
```go
// AuthClient is used to store your auth credentials when authenticating with Smartcar.
auth := smartcar.AuthClient{
	ClientID:     "CLIENT_ID",
	ClientSecret: "CLIENT_SECRET",
	RedirectURI:  "REDIRECT_URI",
	Scope:        []string{"read_vehicle_info"},
	TestMode:     false,
}
```
2. The initialized `AuthClient` must be wrapped in a `AuthConnect` which is used to build an `auth URL`. `AuthConnect` can contain state and structs for Smartcar Pro features such as `ConnectDirect` and `ConnectMatch`, learn more about these features [here](https://smartcar.com/connect/).
```go
connect := smartcar.AuthConnect{
	Auth:     auth,
	ForceApproval: false,
	State:  "",
	ConnectMatch:   smartcar.ConnectMatch{},
	ConnectDirect:     smartcar.ConnectDirect{},
}
```
3. Redirect the user to Smartcar Connect using the URL from `GetAuthURL()`.
```go
authURL, err := smartcar.GetAuthURL(connect)
if err != nil {
  //Handle the err.
}
```
The user will then login, and accept or deny the permissions defined in your `scope`.
  - If the user is already connected to your application, they will not be shown the accept or deny dialog. However the application can force this dialog to be shown by setting `ForceApproval` to `true` in `AuthConnect`.
  - If the user accepts, they will be redirected to your redirect_uri. The query field code will contain an authorization code. This is very important, so save it for later.
  - If the user denies, the query field code will equal "access_denied", so you should handle this somehow.

4. After obtaining an authorization code, exchange it for a `Tokens` struct which will contain access and refresh tokens along with their expiry using `ExchangeCode` which requies your auth credentials and an authorization code.
```go
tokens, err := smartcar.ExchangeCode(auth, "AUTHORIZATION_CODE")
if err != nil {
  //Handle the err.
}
```
5. To make vehicle requests to the Smartcar API, the SDK requires valid access tokens for each request. Access Tokens expire every 2 hours and the current time of expirty will be `Tokens.AccessExpiry`. You can check if an access token is expired using `TokenIsExpired()`.
```go
expired := smartcar.TokenIsExpired(tokens.AccessExpiry)
```
6. To refresh an access token, call `RefreshToken()` with your `AuthClient` and `RefreshToken` to get a new `Token` struct back.
```go
newTokens, err := smartcar.RefreshToken(auth, "REFRESH_TOKEN")
if err != nil {
  //Handle the err.
}
```
You are now to ready to make vehicle requests to the Smartcar API!