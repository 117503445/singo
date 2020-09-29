package test

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"singo/conf"
	"singo/dto"
	"singo/model"
	"singo/server"
	"singo/util"
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
	_, response := httpPost(t, router, "/api/v1/ping", "")

	expectResponse := "\"pong\""
	assert.Equal(t, expectResponse, response)
}
func TestUserRegister(t *testing.T) {
	router := server.NewRouter()

	userRegisterService := dto.UserRegisterIn{
		UserName: "user1",
		Password: "pass1",
		Avatar:   "https://gw.alicdn.com/tps/TB1W_X6OXXXXXcZXVXXXXXXXXXX-400-400.png",
	}

	code, response := httpPostJson(t, router, "/api/v1/user/register", userRegisterService)

	assert.Equal(t, http.StatusOK, code)

	expectResponse := gin.H{
		"ID":       float64(2),
		"username": "user1",
	}

	for k := range expectResponse {
		assert.Equal(t, expectResponse[k], response[k])
	}

}
func TestUserLogin(t *testing.T) {
	router := server.NewRouter()

	userRegisterService := dto.UserRegisterIn{
		UserName: "user1",
		Password: "pass1",
		Avatar:   "https://gw.alicdn.com/tps/TB1W_X6OXXXXXcZXVXXXXXXXXXX-400-400.png",
	}

	httpPostJson(t, router, "/api/v1/user/register", userRegisterService)

	userLoginDto := dto.UserLoginIn{
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
func TestUserMe(t *testing.T) {
	router := server.NewRouter()

	userRegisterService := dto.UserRegisterIn{
		UserName: "user1",
		Password: "pass1",
		Avatar:   "https://gw.alicdn.com/tps/TB1W_X6OXXXXXcZXVXXXXXXXXXX-400-400.png",
	}

	httpPostJson(t, router, "/api/v1/user/register", userRegisterService)

	userLoginDto := dto.UserLoginIn{
		UserName: "user1",
		Password: "pass1",
	}

	_, response := httpPostJson(t, router, "/api/v1/user/login", userLoginDto)
	authorization := "Bearer " + response["token"].(string)

	code, response := httpGetJson(t, router, "/api/v1/user/me", map[string]string{"Authorization": authorization})
	assert.Equal(t, http.StatusOK, code)

	expectResponse := gin.H{
		"ID":       float64(2),
		"username": "user1",
		"avatar":   "https://gw.alicdn.com/tps/TB1W_X6OXXXXXcZXVXXXXXXXXXX-400-400.png",
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

	model.InitDatabase() //重新创建空白的数据库

	exitCode := m.Run()
	os.Exit(exitCode)
}
func TestCreateAdminPasswordTxt(t *testing.T) {
	filePath := util.FilePasswordAdmin
	bytes, err := ioutil.ReadFile(filePath)
	assert.Nil(t, err)
	assert.Equal(t, 12, len(string(bytes)))
}
func TestCreateJwtPasswordTxt(t *testing.T) {
	filePath := util.FilePasswordJWT
	bytes, err := ioutil.ReadFile(filePath)
	assert.Nil(t, err)
	assert.Equal(t, 12, len(string(bytes)))
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

func httpGet(t *testing.T, router *gin.Engine, url string, headers map[string]string) (responseCode int, responseText string) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	assert.Nil(t, err)
	for key, value := range headers {
		request.Header.Add(key, value)
	}
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	return recorder.Code, recorder.Body.String()
}
func httpGetJson(t *testing.T, router *gin.Engine, url string, headers map[string]string) (responseCode int, responseMap map[string]interface{}) {
	code, text := httpGet(t, router, url, headers)

	var response map[string]interface{}
	err := json.Unmarshal([]byte(text), &response)
	assert.Nil(t, err)

	return code, response
}
