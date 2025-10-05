package route

import (
	"github.com/suk-chanthea/ezra/api/controller"

	"github.com/gin-gonic/gin"
)

func EventRoutes(rg *gin.RouterGroup, c *controller.EventController) {
	events := rg.Group("/events")
	{
		events.POST("/", c.Create)
		events.GET("/", c.GetAll)
	}
}
