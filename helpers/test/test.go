package test

import (
	"gopkg.in/h2non/gock.v1"
)

// MockRequest is ...
func MockRequest(method, url, authorization string, statusCode int, response interface{}) {
	gock.New(url).
		MatchHeader("Authorization", authorization).
		Reply(statusCode).
		JSON(response)
}

// ClearMock is...
func ClearMock() {
	defer gock.Off()
}
