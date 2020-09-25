package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"singo/conf"
	"singo/server"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUserMeUnauthorized(t *testing.T) {
	router := server.NewRouter()

	request, err := http.NewRequest(http.MethodGet, "/api/v1/user/me", nil)
	assert.Nil(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	err = json.Unmarshal([]byte(recorder.Body.String()), &response)
	assert.Nil(t, err)

	expectBody := map[string]interface{}( gin.H{
		"code": float64(401),
		"msg":  "未登录",
	})
	assert.Equal(t, expectBody, response)
}

func TestPing(t *testing.T) {
	router := server.NewRouter()

	request, err := http.NewRequest(http.MethodPost, "/api/v1/ping", nil)
	assert.Nil(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	err = json.Unmarshal([]byte(recorder.Body.String()), &response)
	assert.Nil(t, err)

	expectBody := map[string]interface{}( gin.H{
		"code": float64(0),
		"msg":  "Pong",
	})
	assert.Equal(t, expectBody, response)

}

func TestMain(m *testing.M) {
	conf.Init()
	exitCode := m.Run()
	os.Exit(exitCode)
}
