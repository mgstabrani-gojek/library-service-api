package controller_test

import (
	"gojek/library-service-api/internal/controller"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlePingRequest_GivenNothing_ThenReturnHTTPSuccess(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	controller.HandlePingRequest(w, req)

	res := w.Result()
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)

	expectedResponse := `{"message":"pong"}`

	assert.Equal(t, expectedResponse, string(data))
}
