package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/mghaan/exequte/app"
	"github.com/mghaan/exequte/library"
	log "github.com/mghaan/exequte/logger"
)

func main() {
	conf, logs := app.Configure()

	mqtt := app.StartMqtt(logs, conf)
	library.LoadLibrary(mqtt, conf.Plugins)
	library.LoadPlugins(mqtt, conf.Plugins)

	logs.Info(log.SYSTEM, "Initialization complete")

	// main loop
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigc
		mqtt.Disconnect()
		logs.Info(log.SYSTEM, "Terminated")
		os.Exit(0)
	}()

	select {}
}
