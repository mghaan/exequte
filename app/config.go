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
package app

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/mghaan/exequte/logger"
)

type Broker struct {
	Host     string `json:"host"`     // mqtt server host
	Port     int    `json:"port"`     // mqtt server port
	Ssl      bool   `json:"ssl"`      // use secure connection
	User     string `json:"username"` // server username
	Password string `json:"password"` // server password
	Client   string `json:"client"`   // client id to present to server
}

type Plugin struct {
	Plugin string          `json:"plugin"`
	Config json.RawMessage `json:"config"`
}

type System struct {
	Debug bool   `json:"debug"` // enable debug log
	Log   string `json:"log"`   // log file
}

type Config struct {
	System  System   `json:"system"` // system configuration
	Mqtt    Broker   `json:"mqtt"`   // mqtt configuration
	Plugins []Plugin `json:"plugins"`
}

func Configure() (*Config, *logger.Logger) {
	workdir, _ := os.Getwd()
	if len(workdir) < 1 {
		workdir = "."
	}

	logs := logger.New()

	var conf string
	flag.StringVar(&conf, "config", workdir+string(os.PathSeparator)+"exequte.yaml", "path to config file")
	flag.Parse()

	cfg := &Config{}
	cfg.System.Debug = false
	cfg.System.Log = ""
	cfg.Mqtt.Host = "127.0.0.1"
	cfg.Mqtt.Port = 1883
	cfg.Mqtt.Ssl = false
	cfg.Mqtt.User = "mqtt"
	cfg.Mqtt.Password = "secret"
	cfg.Mqtt.Client = "exequte"

	file, err := os.Open(conf)
	if err != nil {
		logs.Fatal(logger.SYSTEM, "Config file not found", err)
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		logs.Fatal(logger.SYSTEM, "Unable to parse config", err)
	}

	file.Close()

	if len(cfg.System.Log) > 2 {
		if err = logs.Create(cfg.System.Log); err != nil {
			logs.Fatal(logger.SYSTEM, "Could not open log file", err)
		}
	}

	return cfg, logs
}
