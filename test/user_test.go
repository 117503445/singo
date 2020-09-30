package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"singo/conf"
	"singo/dto"
	"singo/model"
	"singo/router"
	"singo/serializer"
	"singo/util"
	"strings"
	"testing"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var r *gin.Engine

func TestUserMeUnauthorized(t *testing.T) {

	request, err := http.NewRequest(http.MethodGet, "/api/v1/user/me", nil)
	assert.Nil(t, err)

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)

	var response map[string]interface{}
	err = json.Unmarshal([]byte(recorder.Body.String()), &response)
	assert.Nil(t, err)

	expectResponse := map[string]interface{}(gin.H{
		"code":    float64(401),
		"message": "cookie token is empty",
	})
	assert.Equal(t, expectResponse, response)
}
func TestPing(t *testing.T) {

	_, response := httpPost(t, r, "/api/v1/ping", nil, "")

	expectResponse := "\"pong\""
	assert.Equal(t, expectResponse, response)
}
func TestUserRegister(t *testing.T) {

	userCreateUpdateIn := dto.UserCreateUpdateIn{
		UserName: "user1",
		Password: "pass1",
		Avatar:   "https://gw.alicdn.com/tps/TB1W_X6OXXXXXcZXVXXXXXXXXXX-400-400.png",
	}

	code, response := httpPostJson(t, r, "/api/v1/user", nil, userCreateUpdateIn)

	assert.Equal(t, http.StatusOK, code)

	expectResponse := gin.H{
		"id":       float64(2),
		"username": "user1",
	}

	for k := range expectResponse {
		assert.Equal(t, expectResponse[k], response[k])
	}

}
func TestUserRegisterParamNotValidError(t *testing.T) {

	userCreateUpdateIn := dto.UserCreateUpdateIn{
		UserName: "u",
		Password: "pass1",
		Avatar:   "https://gw.alicdn.com/tps/TB1W_X6OXXXXXcZXVXXXXXXXXXX-400-400.png",
	}

	code, response := httpPostJson(t, r, "/api/v1/user", nil, userCreateUpdateIn)

	assert.Equal(t, http.StatusBadRequest, code)

	expectResponse := gin.H{
		"code":    float64(serializer.StatusParamNotValid),
		"message": "StatusParamNotValid",
	}

	for k := range expectResponse {
		assert.Equal(t, expectResponse[k], response[k])
	}

}
func TestUserLogin(t *testing.T) {

	userCreateUpdateIn := dto.UserCreateUpdateIn{
		UserName: "user1",
		Password: "pass1",
		Avatar:   "https://gw.alicdn.com/tps/TB1W_X6OXXXXXcZXVXXXXXXXXXX-400-400.png",
	}

	httpPostJson(t, r, "/api/v1/user", nil, userCreateUpdateIn)

	userLoginDto := dto.UserLoginIn{
		UserName: "user1",
		Password: "pass1",
	}

	code, response := httpPostJson(t, r, "/api/v1/user/login", nil, userLoginDto)

	assert.Equal(t, http.StatusOK, code)

	expectResponse := gin.H{
		"code": float64(200),
	}

	for k := range expectResponse {
		assert.Equal(t, expectResponse[k], response[k])
	}

}
func TestUserMe(t *testing.T) {

	userCreateUpdateIn := dto.UserCreateUpdateIn{
		UserName: "user1",
		Password: "pass1",
		Avatar:   "https://gw.alicdn.com/tps/TB1W_X6OXXXXXcZXVXXXXXXXXXX-400-400.png",
	}

	httpPostJson(t, r, "/api/v1/user", nil, userCreateUpdateIn)

	userLoginDto := dto.UserLoginIn{
		UserName: "user1",
		Password: "pass1",
	}

	_, response := httpPostJson(t, r, "/api/v1/user/login", nil, userLoginDto)
	authorization := "Bearer " + response["token"].(string)

	code, response := httpGetJson(t, r, "/api/v1/user/me", map[string]string{"Authorization": authorization})
	assert.Equal(t, http.StatusOK, code)

	expectResponse := gin.H{
		"id":       float64(2),
		"username": "user1",
		"avatar":   "https://gw.alicdn.com/tps/TB1W_X6OXXXXXcZXVXXXXXXXXXX-400-400.png",
	}

	for k := range expectResponse {
		assert.Equal(t, expectResponse[k], response[k])
	}

}
func initCleanDatabase() {
	dbName := viper.GetString("mysql.dbname")
	if _, err := model.Exec(fmt.Sprintf("drop database %v", dbName)); err != nil {
		panic("删除数据库失败")
	}

	model.InitDatabase() //重新创建空白的数据库
}
func TestMain(m *testing.M) {
	conf.Init()
	dbName := viper.GetString("mysql.dbname")
	if !strings.Contains(dbName, "test") {
		panic("本测试会清空数据库,禁止在 数据库名 不包含 test 的 数据库上运行")
	}
	r = router.NewRouter()
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
func TestUserUpdate(t *testing.T) {

	userCreateUpdateIn := dto.UserCreateUpdateIn{
		UserName: "user1",
		Password: "pass1",
		Avatar:   "https://gw.alicdn.com/tps/TB1W_X6OXXXXXcZXVXXXXXXXXXX-400-400.png",
	}

	httpPostJson(t, r, "/api/v1/user", nil, userCreateUpdateIn)

	userLoginDto := dto.UserLoginIn{
		UserName: "user1",
		Password: "pass1",
	}

	_, response := httpPostJson(t, r, "/api/v1/user/login", nil, userLoginDto)
	authorization := "Bearer " + response["token"].(string)

	userCreateUpdateIn = dto.UserCreateUpdateIn{
		UserName: "user1",
		Password: "pass1",
		Avatar:   "newAva",
	}
	code, response := httpPutJson(t, r, "/api/v1/user", map[string]string{"Authorization": authorization},
		userCreateUpdateIn)

	assert.Equal(t, http.StatusOK, code)
	expectResponse := gin.H{
		"id":       float64(2),
		"username": "user1",
		"avatar":   "newAva",
	}
	for k := range expectResponse {
		assert.Equal(t, expectResponse[k], response[k])
	}

}
func TestUserUpdateRepeatUsernameError(t *testing.T) {

	userCreateUpdateIn := dto.UserCreateUpdateIn{
		UserName: "user1",
		Password: "pass1",
		Avatar:   "https://gw.alicdn.com/tps/TB1W_X6OXXXXXcZXVXXXXXXXXXX-400-400.png",
	}

	httpPostJson(t, r, "/api/v1/user", nil, userCreateUpdateIn)

	userCreateUpdateIn = dto.UserCreateUpdateIn{
		UserName: "user2",
		Password: "pass2",
		Avatar:   "https://gw.alicdn.com/tps/TB1W_X6OXXXXXcZXVXXXXXXXXXX-400-400.png",
	}

	httpPostJson(t, r, "/api/v1/user", nil, userCreateUpdateIn)

	userLoginDto := dto.UserLoginIn{
		UserName: "user1",
		Password: "pass1",
	}

	_, response := httpPostJson(t, r, "/api/v1/user/login", nil, userLoginDto)
	authorization := "Bearer " + response["token"].(string)

	userCreateUpdateIn = dto.UserCreateUpdateIn{
		UserName: "user2",
		Password: "pass1",
		Avatar:   "newAva",
	}
	code, response := httpPutJson(t, r, "/api/v1/user", map[string]string{"Authorization": authorization},
		userCreateUpdateIn)

	assert.Equal(t, http.StatusBadRequest, code)
	expectResponse := gin.H{
		"code":    float64(serializer.StatusUsernameRepeat),
		"message": "Username has already exists.",
	}
	for k := range expectResponse {
		assert.Equal(t, expectResponse[k], response[k])
	}
}
func TestUserUpdateParamNotValidError(t *testing.T) {

	userCreateUpdateIn := dto.UserCreateUpdateIn{
		UserName: "user1",
		Password: "pass1",
		Avatar:   "https://gw.alicdn.com/tps/TB1W_X6OXXXXXcZXVXXXXXXXXXX-400-400.png",
	}

	httpPostJson(t, r, "/api/v1/user", nil, userCreateUpdateIn)

	userCreateUpdateIn = dto.UserCreateUpdateIn{
		UserName: "user2",
		Password: "pass2",
		Avatar:   "https://gw.alicdn.com/tps/TB1W_X6OXXXXXcZXVXXXXXXXXXX-400-400.png",
	}

	httpPostJson(t, r, "/api/v1/user", nil, userCreateUpdateIn)

	userLoginDto := dto.UserLoginIn{
		UserName: "user1",
		Password: "pass1",
	}

	_, response := httpPostJson(t, r, "/api/v1/user/login", nil, userLoginDto)
	authorization := "Bearer " + response["token"].(string)

	userCreateUpdateIn = dto.UserCreateUpdateIn{
		UserName: "u",
		Password: "pass1",
		Avatar:   "newAva",
	}
	code, response := httpPutJson(t, r, "/api/v1/user", map[string]string{"Authorization": authorization},
		userCreateUpdateIn)

	assert.Equal(t, http.StatusBadRequest, code)
	expectResponse := gin.H{
		"code":    float64(serializer.StatusParamNotValid),
		"message": "StatusParamNotValid",
	}
	for k := range expectResponse {
		assert.Equal(t, expectResponse[k], response[k])
	}
}
func TestAdminUserRead(t *testing.T) {

	filePath := util.FilePasswordAdmin
	bytes, err := ioutil.ReadFile(filePath)
	assert.Nil(t, err)
	password := string(bytes)
	assert.Equal(t, 12, len(password))

	userLoginDto := dto.UserLoginIn{
		UserName: "admin",
		Password: password,
	}
	_, response := httpPostJson(t, r, "/api/v1/user/login", nil, userLoginDto)
	authorization := "Bearer " + response["token"].(string)
	code, response := httpGetJson(t, r, "/api/v1/admin/user/1", map[string]string{"Authorization": authorization})
	assert.Equal(t, http.StatusOK, code)

	expectResponse := gin.H{
		"id":       float64(1),
		"username": "admin",
		"avatar":   "",
	}

	for k := range expectResponse {
		assert.Equal(t, expectResponse[k], response[k])
	}
}
func TestAdminUserReadNotFoundError(t *testing.T) {
	initCleanDatabase()

	filePath := util.FilePasswordAdmin
	bytes, err := ioutil.ReadFile(filePath)
	assert.Nil(t, err)
	password := string(bytes)
	assert.Equal(t, 12, len(password))

	userLoginDto := dto.UserLoginIn{
		UserName: "admin",
		Password: password,
	}
	_, response := httpPostJson(t, r, "/api/v1/user/login", nil, userLoginDto)
	authorization := "Bearer " + response["token"].(string)
	code, response := httpGetJson(t, r, "/api/v1/admin/user/2", map[string]string{"Authorization": authorization})
	assert.Equal(t, http.StatusNotFound, code)

	expectResponse := gin.H{
		"code":    float64(404),
		"message": "user not found",
		"error":   "record not found",
	}

	for k := range expectResponse {
		assert.Equal(t, expectResponse[k], response[k])
	}
}
func TestAdminUserReadNoRoleError(t *testing.T) {

	filePath := util.FilePasswordAdmin
	bytes, err := ioutil.ReadFile(filePath)
	assert.Nil(t, err)
	password := string(bytes)
	assert.Equal(t, 12, len(password))

	code, response := httpGetJson(t, r, "/api/v1/admin/user/1", nil)
	assert.Equal(t, http.StatusUnauthorized, code)

	expectResponse := gin.H{
		"code":    float64(401),
		"message": "cookie token is empty",
	}

	for k := range expectResponse {
		assert.Equal(t, expectResponse[k], response[k])
	}
}
func httpRequest(t *testing.T, httpMethod string, router *gin.Engine, url string, headers map[string]string, body string) (responseCode int, responseText string) {
	request, err := http.NewRequest(httpMethod, url, strings.NewReader(body))
	assert.Nil(t, err)
	for key, value := range headers {
		request.Header.Add(key, value)
	}
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	return recorder.Code, recorder.Body.String()
}
func httpRequestJson(t *testing.T, httpMethod string, router *gin.Engine, url string, headers map[string]string, body interface{}) (responseCode int, responseMap map[string]interface{}) {
	js, _ := json.Marshal(body)
	code, text := httpRequest(t, httpMethod, router, url, headers, string(js))

	var response map[string]interface{}
	err := json.Unmarshal([]byte(text), &response)
	assert.Nil(t, err)

	return code, response
}

func httpGet(t *testing.T, router *gin.Engine, url string, headers map[string]string) (responseCode int, responseText string) {
	return httpRequest(t, http.MethodGet, router, url, headers, "")
}
func httpGetJson(t *testing.T, router *gin.Engine, url string, headers map[string]string) (responseCode int, responseMap map[string]interface{}) {
	return httpRequestJson(t, http.MethodGet, router, url, headers, "")
}

func httpPost(t *testing.T, router *gin.Engine, url string, headers map[string]string, body string) (responseCode int, responseText string) {
	return httpRequest(t, http.MethodPost, router, url, headers, body)
}
func httpPostJson(t *testing.T, router *gin.Engine, url string, headers map[string]string, body interface{}) (responseCode int, responseMap map[string]interface{}) {
	return httpRequestJson(t, http.MethodPost, router, url, headers, body)
}

func httpPut(t *testing.T, router *gin.Engine, url string, headers map[string]string, body string) (responseCode int, responseText string) {
	return httpRequest(t, http.MethodPut, router, url, headers, body)
}
func httpPutJson(t *testing.T, router *gin.Engine, url string, headers map[string]string, body interface{}) (responseCode int, responseMap map[string]interface{}) {
	return httpRequestJson(t, http.MethodPut, router, url, headers, body)
}

func httpDelete(t *testing.T, router *gin.Engine, url string, headers map[string]string, body string) (responseCode int, responseText string) {
	return httpRequest(t, http.MethodDelete, router, url, headers, body)
}
func httpDeleteJson(t *testing.T, router *gin.Engine, url string, headers map[string]string, body interface{}) (responseCode int, responseMap map[string]interface{}) {
	return httpRequestJson(t, http.MethodDelete, router, url, headers, body)
}
