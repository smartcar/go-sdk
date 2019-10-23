package smartcar

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/h2non/gock.v1"
)

type RequestTestSuite struct {
	suite.Suite
	backend backendClient
}

func (s *RequestTestSuite) SetupTest() {
	s.backend = newBackend()
}

type mockResponse struct {
	SomeKey string
	ResponseHeaders
}
type mockBody struct {
	SomeKey string
}

func (s *RequestTestSuite) TestCall() {
	defer gock.Off()

	mockAge := "data-age"
	mockRequestID := "request-id"
	mockUnitSystem := Metric
	mockValue := "mock value"
	expectedResponse := &mockResponse{
		SomeKey: mockValue,
		ResponseHeaders: ResponseHeaders{
			Age:        mockAge,
			DataAge:    mockAge,
			RequestID:  mockRequestID,
			UnitSystem: mockUnitSystem,
		},
	}
	target := new(mockResponse)
	mockURL := "https://example.com"
	gock.New(mockURL).
		Get("/").
		Reply(200).
		SetHeader("Sc-Data-Age", mockAge).
		SetHeader("Sc-Request-Id", mockRequestID).
		SetHeader("Sc-Unit-System", string(mockUnitSystem)).
		JSON(map[string]string{"someKey": mockValue})

	err := s.backend.Call(backendClientParams{
		ctx:    context.TODO(),
		url:    mockURL,
		method: http.MethodGet,
		target: target,
	})

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedResponse, target)
}

func (s *RequestTestSuite) TestformatHeadersResponse() {
	mockAge := "data-age"
	mockRequestID := "request-id"
	mockUnitSystem := Metric
	mockValue := "mock value"
	expectedResponse := &mockResponse{
		SomeKey: mockValue,
		ResponseHeaders: ResponseHeaders{
			Age:        mockAge,
			DataAge:    mockAge,
			RequestID:  mockRequestID,
			UnitSystem: mockUnitSystem,
		},
	}
	target := new(mockResponse)
	// Make sure the formatter does not override other keys than the ResponseHeaders.
	target.SomeKey = mockValue

	headers := http.Header{}
	headers.Add("Sc-Unit-System", string(mockUnitSystem))
	headers.Add("Sc-Request-Id", mockRequestID)
	headers.Add("Sc-Data-Age", mockAge)

	backend := &backend{}
	err := backend.formatHeadersResponse(headers, target)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedResponse, target)
}

func (s *RequestTestSuite) TestformatBodyResponse() {
	mockValue := "mock value"
	expectedResponse := &mockBody{
		SomeKey: mockValue,
	}
	apiResponse := map[string]string{
		"someKey": mockValue,
	}
	b, err := json.Marshal(apiResponse)
	if err != nil {
		assert.Fail(s.T(), "Marshal went wrong")
	}
	body := ioutil.NopCloser(bytes.NewReader(b))
	target := new(mockBody)

	backend := &backend{}
	err = backend.formatBodyResponse(body, target)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedResponse, target)
}

func (s *RequestTestSuite) TestnewRequestEmtpyParams() {
	mockMethod := http.MethodGet
	mockURL := "https://example.com"
	mockAuthorization := "authorization"
	mockUserAgent := getUserAgent()

	backend := &backend{}
	req, err := backend.newRequest(backendClientParams{
		ctx:           context.TODO(),
		url:           mockURL,
		method:        mockMethod,
		authorization: mockAuthorization,
	})

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), mockAuthorization, req.Header.Get("Authorization"))
	assert.Equal(s.T(), mockUserAgent, req.Header.Get("User-Agent"))
	assert.Equal(s.T(), "", req.Header.Get("Content-Type"))
	assert.Equal(s.T(), "", req.Header.Get("SC-Unit-System"))
	assert.Equal(s.T(), context.TODO(), req.Context())
}

func (s *RequestTestSuite) TestnewRequest() {
	mockUnitSystem := Metric
	mockMethod := http.MethodGet
	mockURL := "https://example.com"
	mockAuthorization := "authorization"
	mockUserAgent := getUserAgent()
	mockBody := strings.NewReader("hello world")
	mockContentType := getBodyType(mockBody)

	backend := &backend{}
	req, err := backend.newRequest(backendClientParams{
		authorization: mockAuthorization,
		body:          mockBody,
		ctx:           context.TODO(),
		method:        mockMethod,
		url:           mockURL,
		requestParams: requestParams{
			UnitSystem: mockUnitSystem,
		},
	})

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), mockAuthorization, req.Header.Get("Authorization"))
	assert.Equal(s.T(), mockUserAgent, req.Header.Get("User-Agent"))
	assert.Equal(s.T(), mockContentType, req.Header.Get("Content-Type"))
	assert.Equal(s.T(), string(mockUnitSystem), req.Header.Get("SC-Unit-System"))
	assert.Equal(s.T(), context.TODO(), req.Context())
}

func (s *RequestTestSuite) TestgetBodyType() {
	expectedBodyType := "application/x-www-form-urlencoded"
	mockBody := strings.NewReader("hello world")

	bodyType := getBodyType(mockBody)

	assert.Equal(s.T(), expectedBodyType, bodyType)
}

func (s *RequestTestSuite) TestgetBodyTypeJson() {
	expectedBodyType := "application/json"
	mockBody := bytes.NewBuffer([]byte("hello"))

	bodyType := getBodyType(mockBody)

	assert.Equal(s.T(), expectedBodyType, bodyType)
}

func TestRequestTestSuite(t *testing.T) {
	suite.Run(t, new(RequestTestSuite))
}
