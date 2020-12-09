package middleware

import (
	log "github.com/Sirupsen/logrus"
	be "github.com/wolvex/paymentaggregator"
)


func CheckAuth(msg be.Message) (ok bool) {
		
	if err := Authorize(msg); err == nil {
		return  true
	} else {		
		log.WithField("stacktrace", err.Dump()).Error("Failed to authorize request")	
		return  false
	}
	
}