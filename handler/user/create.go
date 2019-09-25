package user

import (
	"fmt"

	"apiserver/handler"
	"apiserver/pkg/errno"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// Create user
func Create(c *gin.Context) {
	var r CreateRequest

	var err error
	if err := c.Bind(&r); err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
		return
	}

	admin2 := c.Param("username")
	log.Infof("URL username: %s", admin2)

	desc := c.Query("desc")
	log.Infof("URL key param desc: %s", desc)

	contentType := c.GetHeader("Content-Type")
	log.Infof("Header Content-Type: %s", contentType)
	log.Debugf("username is: [%s], password is [%s]", r.Username, r.Password)

	if r.Username == "" {
		err = errno.New(errno.ErrUserNotFound, fmt.Errorf("username doesn't exist")).Add("This is a add message.")
		handler.SendResponse(c, err, nil)
		return
	}
	if r.Password == "" {
		err = fmt.Errorf("password is empty")
		handler.SendResponse(c, err, nil)
		return
	}

	resp := CreateReponse{Username: r.Username}
	handler.SendResponse(c, nil, resp)
}
