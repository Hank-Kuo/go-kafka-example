package http

import "github.com/gin-gonic/gin"

func Healthz(engine *gin.Engine) {
	engine.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, "OK")
	})
}
