package middleware

import "github.com/gin-gonic/gin"

func Healthz(engine *gin.Engine) {
	engine.GET("/", func(c *gin.Context) {
		c.JSON(200, "server live!!")
	})
}
