package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dpattmann/furby/internal/mocks"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func setupForwarderMock(storeToken *oauth2.Token, storeError error, authReturn bool) (mockStore *mocks.Store, mockAuth *mocks.Authorization) {
	// New mock controller
	mockStore = &mocks.Store{}
	mockStore.On("GetToken").Return(storeToken, storeError)

	mockAuth = &mocks.Authorization{}
	mockAuth.On("IsAuthorized", httptest.NewRequest(http.MethodGet, "https://tokenbackend", nil)).Return(authReturn)

	return
}

func Test_Forwarder(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Reflect all request headers
	httpmock.RegisterResponder("GET", "https://tokenbackend", func(request *http.Request) (*http.Response, error) {
		res := httpmock.NewStringResponse(http.StatusOK, "Reflect Header")
		for headerName, headerValues := range request.Header {
			for _, headerValue := range headerValues {
				res.Header.Add(headerName, headerValue)
			}
		}

		return res, nil
	})

	t.Run("Add valid token to HTTP reqeust", func(t *testing.T) {
		// setup Mock
		mockToken = oauth2.Token{AccessToken: "OAuth2 Access Token Value"}
		mockStore, mockAuth := setupForwarderMock(&mockToken, nil, true)

		// create StoreHandler with mocked store and auth
		forwarder := NewForwarder(mockStore, mockAuth)

		req := httptest.NewRequest(http.MethodGet, "https://tokenbackend", nil)
		responseRecorder := httptest.NewRecorder()

		forwarder.ServeHTTP(responseRecorder, req)

		assert.Equal(t, "Bearer: OAuth2 Access Token Value", responseRecorder.Header().Get("Authorization"))
		assert.Equal(t, http.StatusOK, responseRecorder.Code)
	})
}
