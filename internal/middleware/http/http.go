package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func httpNotFound(c *gin.Context) {
	msg := c.Request.Method + " " + c.Request.URL.String()
	c.JSON(http.StatusNotFound, "not found "+msg)
}
