package v1

import (
	"doghandler/modules"
	"doghandler/pkg/logger"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ReceiverAlert deal with watchdog message
func ReceiverAlert(c *gin.Context) {
	serviceid := c.Param("serviceid")
	if dog, ok := modules.Dogs[serviceid]; ok {
		dog.TouchDog()
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": fmt.Sprintf("get service %s barking", serviceid),
		})
	} else {
		logger.LogInfo("DogHandler", fmt.Sprintf("service %s did not register", serviceid))
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": fmt.Sprintf("service %s did not register", serviceid),
		})
	}

}
