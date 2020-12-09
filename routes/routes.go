package routes

import (
	"simobiplus/controllers"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

//Init ...
func Init() *echo.Echo {
	
	e := echo.New()
	
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {

	}))
	
	SimobiRoutes(e.Group("/paggr-simobi"))

	return e
}

//SimobiRoutes ...
func SimobiRoutes(g *echo.Group)  {

	g.POST("/advise", controllers.AdviseControllers)

	g.POST("/callback", controllers.CallBackControllers)

	g.POST("/notify", controllers.NotifyControllers)

	g.POST("/ping", controllers.PingControllers)

	g.POST("/payment", controllers.PaymentControllers)

	g.POST("/void", controllers.NotifyControllers)

	
	 
}