package middleware

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"singo/dto"
	"singo/model"
	"singo/util"
	"time"
)

var JwtMiddleware *jwt.GinJWTMiddleware

func init() {
	identityKey := "user"
	var err error

	var jwtPassword string
	bytes, err := ioutil.ReadFile(util.FilePasswordJWT)
	if err == nil && len(string(bytes)) == 12 {
		jwtPassword = string(bytes)
	} else {
		jwtPassword = util.RandStringRunes(12)
		_ = ioutil.WriteFile(util.FilePasswordJWT, []byte(jwtPassword), 0777)
	}

	JwtMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "singo zone",
		Key:         []byte(jwtPassword),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model.User); ok {
				m := jwt.MapClaims{
					identityKey: v.Username,
				}
				return m
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			username := claims[identityKey].(string)
			u, err := model.ReadUserByName(username)
			if err != nil {
				return nil
			}
			return &u
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var userLoginDto dto.UserLoginIn
			if err := c.ShouldBindJSON(&userLoginDto); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := userLoginDto.UserName
			password := userLoginDto.Password
			queryUser, err := model.ReadUserByName(username)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			if queryUser.CheckPassword(password) {
				return &queryUser, nil
			}
			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*model.User); ok {
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",
		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		util.Log().Error("jwt.New failed", err)
	}

	if err = JwtMiddleware.MiddlewareInit(); err != nil {
		util.Log().Error("authMiddleware.MiddlewareInit()", err)
	}
}
