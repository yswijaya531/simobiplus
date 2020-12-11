package controllers

import (
	"time"

	log "github.com/sirupsen/logrus"
)

func deferCheckout(req interface{}, apiString string) {

	defer func(begin time.Time) {
		elapsed := float64(time.Since(begin).Nanoseconds()) / float64(1e6)
		log.WithFields(log.Fields{
			"api ":     apiString,
			"request ": result,
			"elapsed":  elapsed,
		}).Info("Sending HTTP request")
	}(time.Now())

	return
}
