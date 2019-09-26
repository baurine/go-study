package middleware

import (
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now().UTC()
		path := c.Request.URL.Path

		reg := regexp.MustCompile("(/v1/user|/login|/sd)")
		if !reg.MatchString(path) {
			return
		}

		c.Next()

		end := time.Now().UTC()
		latency := end.Sub(start)

		log.Infof("start: %s, end: %s, latency: %s", start, end, latency)
	}
}
