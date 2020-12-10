package controllers

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	cm "github.com/yswijaya531/simobiplus/common"

	"github.com/labstack/echo"
	"github.com/yswijaya531/simobiplus/handlers"
	mw "github.com/yswijaya531/simobiplus/middleware"
	be "gitlab.smartfren.com/paggr/libraries"
)

var result be.Message

func AdviseControllers(c echo.Context) (errs error) {
	
	req := new(be.Message)
		
	if errs = c.Bind(req); errs != nil {
		log.WithField("error", errs).Error("Exception caught")
		return  errs
	}
		
	msg := *req 

	log.WithField("info",msg).Info("Decode Request Simobi Advise API")

	if mw.CheckAuth(msg) {	
		result, errs = handlers.AdviseHandler(msg)
		log.WithField("error", errs).Error("Exception caught")					
	}  else {
	 	result = mw.BuildResponse(result, be.ERR_INVALID_SIGNATURE, "Signature invalid")
	}
		
	defer func(begin time.Time) {
		elapsed := float64(time.Since(begin).Nanoseconds()) / float64(1e6)
		log.WithFields(log.Fields{
			"request nya": result,
			"elapsed" : elapsed,
		  }).Info("Sending HTTP request")
	}(time.Now())

	return c.JSON(http.StatusOK, result)

}

func CallBackControllers(c echo.Context) (errs error) {
	
	req := new(cm.SimobiCallBack)
		
	if errs = c.Bind(req); errs != nil {
		log.WithField("error", errs).Error("Exception caught")
		return  errs
	}
	
	result := *req 

	log.WithField("info",result).Info("Decode Request Simobi Callback API")

	if result, errs = handlers.CallBackHandler(c); errs != nil{
		log.WithField("error", errs).Error("Exception caught")
		return errs
	} 
	
	defer func(begin time.Time) {
		elapsed := float64(time.Since(begin).Nanoseconds()) / float64(1e6)
		log.WithFields(log.Fields{
			"request nya": result,
			"elapsed" : elapsed,
		  }).Info("Sending HTTP request")
	}(time.Now())

	return c.JSON(http.StatusOK, result)

}

func NotifyControllers(c echo.Context)  (errs error) {
	
	req := new(be.Message)
		
	if errs = c.Bind(req); errs != nil {
		log.WithField("error", errs).Error("Exception caught")
		return  errs
	}
	
	msg := *req 

	log.WithField("info",msg).Info("Decode Request Simobi Notify API")

	if mw.CheckAuth(msg) {	
		result, errs = handlers.NotifyHandler(msg)					
	}  else {
	 	result = mw.BuildResponse(result, be.ERR_INVALID_SIGNATURE, "Signature invalid")
	}
		
	defer func(begin time.Time) {
		elapsed := float64(time.Since(begin).Nanoseconds()) / float64(1e6)
		log.WithFields(log.Fields{
			"request nya": result,
			"elapsed" : elapsed,
		  }).Info("Sending HTTP request")
	}(time.Now())

	return c.JSON(http.StatusOK, result)

}

func PingControllers(c echo.Context) (errs error) {
	
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

func PaymentControllers(c echo.Context) (errs error) {
	
	req := new(be.Message)
		
	if errs = c.Bind(req); errs != nil {
		log.WithField("error", errs).Error("Exception caught")
		return  errs
	}
		
	msg := *req 

	log.WithField("info",msg).Info("Decode Request Simobi Payment API")

	if mw.CheckAuth(msg) {	
		result, errs = handlers.PaymentHandler(msg)					
	}  else {
	 	result = mw.BuildResponse(result, be.ERR_INVALID_SIGNATURE, "Signature invalid")
	}
		
	defer func(begin time.Time) {
		elapsed := float64(time.Since(begin).Nanoseconds()) / float64(1e6)
		log.WithFields(log.Fields{
			"request nya": result,
			"elapsed" : elapsed,
		  }).Info("Sending HTTP request")
	}(time.Now())

	return c.JSON(http.StatusOK, result)

}

func VoidControllers(c echo.Context) (errs error) {
	
	req := new(be.Message)
		
	if errs = c.Bind(req); errs != nil {
		log.WithField("error", errs).Error("Exception caught")
		return  errs
	}
		
	msg := *req 

	log.WithField("info",msg).Info("Decode  Request Simobi Void API")

	if mw.CheckAuth(msg) {	
		result, errs = handlers.VoidHandler(msg)					
	}  else {
	 	result = mw.BuildResponse(result, be.ERR_INVALID_SIGNATURE, "Signature invalid")
	}
		
	defer func(begin time.Time) {
		elapsed := float64(time.Since(begin).Nanoseconds()) / float64(1e6)
		log.WithFields(log.Fields{
			"request nya": result,
			"elapsed" : elapsed,
		  }).Info("Sending HTTP request")
	}(time.Now())

	return c.JSON(http.StatusOK, result)

}

