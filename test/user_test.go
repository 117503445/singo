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

	request, err := http.NewRequest(http.MethodPost, "/api/v1/ping", nil)
	assert.Nil(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	err = json.Unmarshal([]byte(recorder.Body.String()), &response)
	assert.Nil(t, err)

	expectResponse := map[string]interface{}( gin.H{
		"code": float64(0),
		"msg":  "Pong",
	})
	assert.Equal(t, expectResponse, response)
}
func TestRegister(t *testing.T) {
	router := server.NewRouter()

	userRegisterService := service.UserRegisterDto{
		UserName: "user1",
		Password: "pass1",
	}

	body, _ := json.Marshal(userRegisterService)
	request, err := http.NewRequest(http.MethodPost, "/api/v1/user/register", strings.NewReader(string(body)))
	assert.Nil(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	err = json.Unmarshal([]byte(recorder.Body.String()), &response)
	assert.Nil(t, err)

	expectResponse := map[string]interface{}( gin.H{
		"ID":       float64(1),
		"username": "user1",
	})

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
