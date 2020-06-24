package routers

import (
	middleware "doghandler/middlewares"
	v1 "doghandler/routers/api/v1"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// InitRouter init router config
func InitRouter() *gin.Engine {
	gin.SetMode(viper.GetString("runmode"))
	r := gin.Default()
	// r.Use(middleware.Cors())
	r.Use(middleware.LoggerToFile())
	// r.Use(middleware.CatchPanic())

	r.GET("/ping", func(c *gin.Context) {
		//输出json结果给调用方
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	apiv1 := r.Group("/v1")
	{
		apiv1.POST("/ping/:serviceid", v1.ReceiverAlert)
	}

	return r
}
