package routes

import (
	"GinWeb/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitApi(c *gin.Engine) {
	c.GET("/get", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
	c.GET("/wss", service.WebChat)
}
