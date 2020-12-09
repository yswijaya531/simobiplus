package handlers

import "github.com/labstack/echo"

func PingHandler(e echo.Context) (ping int) {
	defer panicRecovery()
	return 200
}