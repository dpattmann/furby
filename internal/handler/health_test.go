package handler

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HealthHandler(t *testing.T) {
	t.Run("Should return 200 StatusOk", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		responseRecorder := httptest.NewRecorder()

		HealthHandler().ServeHTTP(responseRecorder, req)

		res := responseRecorder.Result()

		defer res.Body.Close()
		data, _ := ioutil.ReadAll(res.Body)

		assert.Equal(t, http.StatusOK, responseRecorder.Code)
		assert.Equal(t, "Ok", string(data))
	})
}
