package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"singo/dto"
	"singo/model"
	"singo/serializer"
	"strconv"
)

// AdminUserRead 管理员读取某用户信息
func AdminUserRead(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, serializer.Err(http.StatusBadRequest, "id is not number", err))
		return
	}
	user, err := model.ReadUserById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, serializer.Err(http.StatusNotFound, "user not found", err))
		return
	}
	if userOut, err := dto.UserToUserOut(user); err == nil {
		c.JSON(http.StatusOK, userOut)
	} else {
		c.JSON(http.StatusInternalServerError, serializer.Err(serializer.StatusModelToDtoError, "UserToUserOut failed", err))
	}
}
