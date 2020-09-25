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

func TestUserMeUnauthorized(t *testing.T) {
	body := map[string]interface{}( gin.H{
		"code": float64(401),
		"msg":  "未登录",
	})
	conf.Init()
	router := server.NewRouter()

	request, err := http.NewRequest(http.MethodGet, "/api/v1/user/me", nil)
	assert.Nil(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	err = json.Unmarshal([]byte(recorder.Body.String()), &response)

	assert.Nil(t, err)
	assert.Equal(t, body, response)
}
