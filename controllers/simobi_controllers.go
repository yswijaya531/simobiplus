package controllers

import (
	"net/http"

	"simobiplus/handlers"

	"github.com/labstack/echo"
)


func AdviseControllers(c echo.Context) (err error) {
	
	result, err := handlers.NotifyHandler(c)

	return c.JSON(http.StatusOK, result)

}

func CallBackControllers(c echo.Context) (err error) {
	
	result, err := handlers.CallBackHandler(c)

	return c.JSON(http.StatusOK, result)

}

func NotifyControllers(c echo.Context) (err error) {
	
	result, err := handlers.NotifyHandler(c)

	return c.JSON(http.StatusOK, result)

}

func PingControllers(c echo.Context) (err error) {
	
	result := handlers.PingHandler(c)

	return c.JSON(http.StatusOK, result)

}

func PaymentControllers(c echo.Context) (err error) {
	
	result, err := handlers.PaymentHandler(c)

	return c.JSON(http.StatusOK, result)

}

func VoidControllers(c echo.Context) (err error) {
	
	result, err := handlers.VoidHandler(c)

	return c.JSON(http.StatusOK, result)

}

