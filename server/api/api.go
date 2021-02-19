package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterApi(g *gin.Engine) {
	g.GET("/ping", testPing)
	//generalApi := g.Group("/-")
	//{
	//	generalApi.POST("/reload", reloadConfig)
	//}
}

func testPing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "pong",
	})
}
