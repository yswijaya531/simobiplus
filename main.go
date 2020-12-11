package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	cm "github.com/yswijaya531/simobiplus/common"
	"github.com/yswijaya531/simobiplus/routes"
)

var logger *log.Entry

func initLogger() {
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.999",
	})

	//log.SetReportCaller(true)
}

func panicRecovery() {
	if r := recover(); r != nil {
		fmt.Printf("Recovering from panic: %v \n", r)
	}
}

func main() {

	defer panicRecovery()
	initLogger()

	cm.LoadConfig() //FromFile(configFile)

	if cm.Config.BackendKey != "" || cm.Config.CertificateFile != "" {
		e := routes.Init()
		//e.AutoTLSManager.Cache = autocert.DirCache(".cache")
		//e.Logger.Fatal(e.StartTLS(cm.Config.ListenPort, cm.Config.BackendKey, cm.Config.CertificateFile))
		e.Logger.Fatal(e.Start(cm.Config.ListenPort))
	} else {

		fmt.Println("Unable to start the server")
		os.Exit(1)

	}

}
