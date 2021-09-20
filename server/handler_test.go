package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	authMock "github.com/dpattmann/furby/auth/mock"
	storeMock "github.com/dpattmann/furby/store/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

var (
	mockToken = oauth2.Token{AccessToken: "foo", RefreshToken: "bar"}
)

func Test_StoreHandler(t *testing.T) {
	t.Run("Return valid token from store", func(t *testing.T) {
		// New mock controller
		mockController := gomock.NewController(t)
		defer mockController.Finish()

		// mock store and auth
		mockStore := storeMock.NewStore(mockController)
		mockStore.EXPECT().GetToken().Return(&mockToken, nil)
		mockAuth := authMock.NewAuthorization(mockController)
		mockAuth.EXPECT().IsAuthorized(httptest.NewRequest(http.MethodGet, "/", nil)).Return(true)

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
		// New mock controller
		mockController := gomock.NewController(t)
		defer mockController.Finish()

		// mock store and auth
		mockStore := storeMock.NewStore(mockController)
		mockStore.EXPECT().GetToken().Return(nil, errors.New("no Token found"))
		mockAuth := authMock.NewAuthorization(mockController)
		mockAuth.EXPECT().IsAuthorized(httptest.NewRequest(http.MethodGet, "/", nil)).Return(true)

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
		// New mock controller
		mockController := gomock.NewController(t)
		defer mockController.Finish()

		// mock store and auth
		mockStore := storeMock.NewStore(mockController)
		mockStore.EXPECT().GetToken().Return(nil, errors.New("no Token found"))
		mockAuth := authMock.NewAuthorization(mockController)
		mockAuth.EXPECT().IsAuthorized(httptest.NewRequest(http.MethodGet, "/", nil)).Return(false)

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
		// New mock controller
		mockController := gomock.NewController(t)
		defer mockController.Finish()

		// mock store and auth
		mockStore := storeMock.NewStore(mockController)
		mockAuth := authMock.NewAuthorization(mockController)

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
