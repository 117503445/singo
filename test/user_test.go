package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"singo/conf"
	"singo/server"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestUserMeUnauthorized(t *testing.T) {
	body := map[string]interface{}( gin.H{
		"code": float64(401),
		"msg":  "未登录",
	})
	conf.Init()
	router := server.NewRouter()

	w := performRequest(router, "GET", "/api/v1/user/me")

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal([]byte(w.Body.String()), &response)

	assert.Nil(t, err)
	assert.Equal(t, body, response)
}
