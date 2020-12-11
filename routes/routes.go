package routes

import (
	"os"

	"github.com/yswijaya531/simobiplus/controllers"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

//Init ...
func Init() *echo.Echo {
	
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{		
		Format: `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}","host":"${host}",` +
		`"method":"${method}","uri":"${uri}","status":${status},"error":"${error}","latency":${latency},` +
		`"latency_human":"${latency_human}","bytes_in":${bytes_in},` +
		`"bytes_out":${bytes_out}}` + "\n",
		Output: os.Stdout,
	}))
	e.Use(middleware.Recover())
		
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