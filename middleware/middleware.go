package middleware

import (
	log "github.com/sirupsen/logrus"
	be "gitlab.smartfren.com/paggr/libraries"
)


func CheckAuth(msg be.Message) (ok bool) {
		
	if err := Authorize(msg); err == nil {
		return  true
	} else {		
		log.WithField("stacktrace", err.Dump()).Error("Failed to authorize request")	
		return  false
	}
	
}