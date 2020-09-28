package test

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"net/http/httptest"
	"os"
	"singo/conf"
	"singo/model"
	"singo/server"
	"singo/service"
	"strings"
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

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)

	var response map[string]interface{}
	err = json.Unmarshal([]byte(recorder.Body.String()), &response)
	assert.Nil(t, err)

	expectResponse := map[string]interface{}( gin.H{
		"code":    float64(401),
		"message": "cookie token is empty",
	})
	assert.Equal(t, expectResponse, response)
}
func TestPing(t *testing.T) {
	router := server.NewRouter()
	_, response := httpPostJson(t, router, "/api/v1/ping", nil)

	expectResponse := gin.H{
		"code": float64(0),
		"msg":  "Pong",
	}
	assert.Equal(t, map[string]interface{}(expectResponse), response)
}
func TestRegister(t *testing.T) {
	router := server.NewRouter()

	userRegisterService := service.UserRegisterDto{
		UserName: "user1",
		Password: "pass1",
	}

	code, response := httpPostJson(t, router, "/api/v1/user/register", userRegisterService)

	assert.Equal(t, http.StatusOK, code)

	expectResponse := gin.H{
		"ID":       float64(1),
		"username": "user1",
	}

	for k := range expectResponse {
		assert.Equal(t, expectResponse[k], response[k])
	}

}
func TestLogin(t *testing.T) {
	router := server.NewRouter()

	userRegisterService := service.UserRegisterDto{
		UserName: "user1",
		Password: "pass1",
	}

	httpPostJson(t, router, "/api/v1/user/register", userRegisterService)

	userLoginDto := service.UserLoginDto{
		UserName: "user1",
		Password: "pass1",
	}

	code, response := httpPostJson(t, router, "/api/v1/user/login", userLoginDto)

	assert.Equal(t, http.StatusOK, code)

	expectResponse := gin.H{
		"code": float64(200),
	}

	for k := range expectResponse {
		assert.Equal(t, expectResponse[k], response[k])
	}

}
func TestMain(m *testing.M) {
	conf.Init()
	dbName := viper.GetString("mysql.dbname")
	if !strings.Contains(dbName, "test") {
		panic("本测试会清空数据库,禁止在 数据库名 不包含 test 的 数据库上运行")
	}

	if _, err := model.Exec(fmt.Sprintf("drop database %v", dbName)); err != nil {
		panic("删除数据库失败")
	}

	model.Database() //重新创建空白的数据库

	exitCode := m.Run()
	os.Exit(exitCode)
}

func httpPost(t *testing.T, router *gin.Engine, url string, body string) (responseCode int, responseText string) {
	request, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
	assert.Nil(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	return recorder.Code, recorder.Body.String()
}
func httpPostJson(t *testing.T, router *gin.Engine, url string, body interface{}) (responseCode int, responseMap map[string]interface{}) {
	js, _ := json.Marshal(body)
	code, text := httpPost(t, router, url, string(js))

	var response map[string]interface{}
	err := json.Unmarshal([]byte(text), &response)
	assert.Nil(t, err)

	return code, response
}
