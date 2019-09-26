package user

import (
	"strconv"

	. "apiserver/handler"
	"apiserver/model"
	"apiserver/pkg/errno"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

func Delete(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("id"))
	log.Debugf("userId is %d", userID)
	if err := model.DeleteUser(uint64(userID)); err != nil {
		SendResponse(c, errno.ErrDatabase, nil)
		return
	}
	SendResponse(c, nil, nil)
}
