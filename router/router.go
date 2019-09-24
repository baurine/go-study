package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Load loads, middlewares, routes, handlers
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// middlewares
	g.Use(gin.Recovery())

	g.Use(mw...)

	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route")
	})

	// svcd := g.Group("/sd")

	return g
}
