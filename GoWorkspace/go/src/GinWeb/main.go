package main

import (
	"GinWeb/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	routes.InitApi(r)

	r.Run("192.168.10.66:8088")
}
