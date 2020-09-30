package serializer

import "github.com/gin-gonic/gin"

// ErrResponse 基础序列化器
type ErrResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// 三位数错误编码为复用http原本含义
// 五位数错误编码为应用自定义错误
// 五开头的五位数错误编码为服务器端错误，比如数据库操作失败
// 四开头的五位数错误编码为客户端错误，有时候是客户端代码写错了，有时候是用户操作错误
const (
	// StatusDBError 数据库操作失败
	StatusDBError = 50001
	// StatusEncryptError 加密失败
	StatusEncryptError = 50002
	// StatusRegisterError 插入注册信息失败
	StatusRegisterError = 50003
	// StatusModelToDtoError Model 转 Dto 失败
	StatusModelToDtoError = 50004
	// StatusParamError 各种奇奇怪怪的参数错误
	StatusParamError = 40001
	// StatusUsernameRepeat 用户名重复
	StatusUsernameRepeat = 40002
	// StatusDtoToModelError Dto 转 Model 失败
	StatusDtoToModelError = 40003
	// StatusParamNotValid 参数不合法
	StatusParamNotValid = 40004
)

// Err 通用错误处理
func Err(errCode int, msg string, err error) ErrResponse {
	res := ErrResponse{
		Code:    errCode,
		Message: msg,
	}
	// 生产环境隐藏底层报错
	if err != nil && gin.Mode() != gin.ReleaseMode {
		res.Error = err.Error()
	}
	return res
}
