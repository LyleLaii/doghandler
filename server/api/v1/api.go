package v1

import (
	"doghandler/dogtimer"
	"doghandler/pkg/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterApi(g *gin.Engine) {
	apiv1 := g.Group("/v1")
	{
		apiv1.POST("/ping/:serviceId", receiverAlert)
	}
}

// ReceiverAlert deal with watchdog message
func receiverAlert(c *gin.Context) {
	serviceId := c.Param("serviceId")
	if dog, ok := dogtimer.Dogs[serviceId]; ok {
		dog.TouchDog()
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": fmt.Sprintf("get service %s barking", serviceId),
		})
	} else {
		logger.LogInfo("DogHandler", fmt.Sprintf("service %s did not register", serviceId))
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": fmt.Sprintf("service %s did not register", serviceId),
		})
	}

}
