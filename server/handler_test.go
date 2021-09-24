package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dpattmann/furby/mocks"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

var (
	mockToken = oauth2.Token{AccessToken: "foo", RefreshToken: "bar"}
)

func setupMock(storeToken *oauth2.Token, storeError error, authReturn bool) (mockStore *mocks.Store, mockAuth *mocks.Authorization ){
	// New mock controller
	mockStore = &mocks.Store{}
	mockStore.On("GetToken").Return(storeToken, storeError)

	mockAuth = &mocks.Authorization{}
	mockAuth.On("IsAuthorized", httptest.NewRequest(http.MethodGet, "/", nil)).Return(authReturn)

	return
}

func Test_StoreHandler(t *testing.T) {
	t.Run("Return valid token from store", func(t *testing.T) {
		// setup Mock
		mockStore, mockAuth := setupMock(&mockToken, nil, true)

		// create StoreHandler with mocked store and auth
		storeHandler := StoreHandler{store: mockStore, auth: mockAuth}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		responseRecorder := httptest.NewRecorder()

		storeHandler.ServeHTTP(responseRecorder, req)

		res := responseRecorder.Result()

		defer res.Body.Close()
		data, _ := ioutil.ReadAll(res.Body)

		var returnToken oauth2.Token
		_ = json.Unmarshal(data, &returnToken)

		assert.Equal(t, http.StatusOK, responseRecorder.Code)
		assert.Equal(t, mockToken.AccessToken, returnToken.AccessToken)
	})

	t.Run("Return internal server error and no valid token", func(t *testing.T) {
		// setup Mock
		mockStore, mockAuth := setupMock(nil, errors.New("no token found"), true)

		// create StoreHandler with mocked store and auth
		storeHandler := StoreHandler{store: mockStore, auth: mockAuth}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		responseRecorder := httptest.NewRecorder()

		storeHandler.ServeHTTP(responseRecorder, req)

		res := responseRecorder.Result()

		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
	})

	t.Run("Return unauthorized if authorizer denies request", func(t *testing.T) {
		// setup Mock
		mockStore, mockAuth := setupMock(&mockToken, nil, false)

		// create StoreHandler with mocked store and auth
		storeHandler := StoreHandler{store: mockStore, auth: mockAuth}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		responseRecorder := httptest.NewRecorder()

		storeHandler.ServeHTTP(responseRecorder, req)

		res := responseRecorder.Result()

		defer res.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
	})

	t.Run("Return 418 for any other request methode then GET", func(t *testing.T) {
		// setup Mock
		mockStore, mockAuth := setupMock(nil, nil, true)

		// create StoreHandler with mocked store and auth
		storeHandler := StoreHandler{store: mockStore, auth: mockAuth}

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		responseRecorder := httptest.NewRecorder()

		storeHandler.ServeHTTP(responseRecorder, req)

		res := responseRecorder.Result()
		defer res.Body.Close()
		assert.Equal(t, http.StatusTeapot, responseRecorder.Code)
	})
}
