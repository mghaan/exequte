/*
 * Copyright (C) 2022 Marian Micek
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mghaan/exequte/app"
	"github.com/mghaan/exequte/library"
	log "github.com/mghaan/exequte/logger"
)

const VERSION string = "0.1"

func main() {
	conf, logs := app.Configure()
	logs.Echo("exeQute", fmt.Sprintf("version %s started", VERSION))

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
