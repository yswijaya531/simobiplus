package controllers

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/yswijaya531/simobiplus/handlers"

	"github.com/labstack/echo"
)


func AdviseControllers(c echo.Context) (err error) {
	
	result, err := handlers.NotifyHandler(c)

	defer func(begin time.Time) {
		elapsed := float64(time.Since(begin).Nanoseconds()) / float64(1e6)
		log.WithFields(log.Fields{
			"request nya": result,
			"elapsed" : elapsed,
		  }).Info("Sending HTTP request")
	}(time.Now())

	return c.JSON(http.StatusOK, result)

}

func CallBackControllers(c echo.Context) (err error) {
	
	result, err := handlers.CallBackHandler(c)

	defer func(begin time.Time) {
		elapsed := float64(time.Since(begin).Nanoseconds()) / float64(1e6)
		log.WithFields(log.Fields{
			"request nya": result,
			"elapsed" : elapsed,
		  }).Info("Sending HTTP request")
	}(time.Now())

	return c.JSON(http.StatusOK, result)

}

func NotifyControllers(c echo.Context) (err error) {
	
	result, err := handlers.NotifyHandler(c)	
	
	defer func(begin time.Time) {
		elapsed := float64(time.Since(begin).Nanoseconds()) / float64(1e6)
		log.WithFields(log.Fields{
			"request nya": result,
			"elapsed" : elapsed,
		  }).Info("Sending HTTP request")
	}(time.Now())

	return c.JSON(http.StatusOK, result)

}

func PingControllers(c echo.Context) (err error) {
	
	result := handlers.PingHandler(c)

	defer func(begin time.Time) {
		elapsed := float64(time.Since(begin).Nanoseconds()) / float64(1e6)
		log.WithFields(log.Fields{
			"request nya": result,
			"elapsed" : elapsed,
		  }).Info("Sending HTTP request")
	}(time.Now())

	return c.JSON(http.StatusOK, result)

}

func PaymentControllers(c echo.Context) (err error) {
	
	result, err := handlers.PaymentHandler(c)
	
	defer func(begin time.Time) {
		elapsed := float64(time.Since(begin).Nanoseconds()) / float64(1e6)
		log.WithFields(log.Fields{
			"request nya": result,
			"elapsed" : elapsed,
		  }).Info("Sending HTTP request")
	}(time.Now())

	return c.JSON(http.StatusOK, result)

}

func VoidControllers(c echo.Context) (err error) {
	
	result, err := handlers.VoidHandler(c)

	defer func(begin time.Time) {
		elapsed := float64(time.Since(begin).Nanoseconds()) / float64(1e6)
		log.WithFields(log.Fields{
			"request nya": result,
			"elapsed" : elapsed,
		  }).Info("Sending HTTP request")
	}(time.Now())

	return c.JSON(http.StatusOK, result)

}

