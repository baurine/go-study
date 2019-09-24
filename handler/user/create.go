package user

import (
	"apiserver/pkg/errno"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// Create user
func Create(c *gin.Context) {
	var r struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var err error
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": errno.ErrBind})
		return
	}
	log.Debugf("username is: [%s], password is [%s]", r.Username, r.Password)
	if r.Username == "" {
		err = errno.New(errno.ErrUserNotFound, fmt.Errorf("username doesn't exist")).Add("This is a add message.")
		log.Errorf(err, "Get an error")
	}
	if errno.IsErrUserNotFound(err) {
		log.Debug("err type is ErrUserNotFound")
	}
	if r.Password == "" {
		err = fmt.Errorf("password is empty")
	}

	code, message := errno.DecodeErr(err)
	c.JSON(http.StatusOK, gin.H{"code": code, "message": message})
}