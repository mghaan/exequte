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
package status

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mghaan/exequte/app"
)

const PLUGIN string = "STATUS"

type External struct{}

var Plugin External

type Report struct {
	Interval int    `json:"interval"`
	Topic    string `json:"topic"`
	Process  string `json:"process"`
}

func (plugin *External) Register(data json.RawMessage, server *app.Server) bool {
	var tasks []Report

	if err := json.Unmarshal(data, &tasks); err != nil {
		server.Log().Error(PLUGIN, "Unable to parse reports", err)

		return false
	}

	for _, task := range tasks {
		args := strings.Split(task.Process, " ")

		if _, err := os.Stat(args[0]); err != nil {
			server.Log().Error(PLUGIN, fmt.Sprintf("%s: No such file or directory", args[0]), err)
			continue
		}

		app.ScheduleTask(task.Interval, func() {
			cmd := exec.Command(args[0])
			if len(args) > 1 {
				for i := 1; i < len(args); i++ {
					cmd.Args = append(cmd.Args, args[i])
				}
			}

			res := "ERROR"
			buff, fail := cmd.CombinedOutput()
			if fail == nil {
				res = string(buff)
				if len(res) < 1 {
					res = "NULL"
				}
			}

			res = strings.TrimSpace(res)

			server.Log().Info(PLUGIN, fmt.Sprintf("Report '%s'", task.Process))
			server.Publish("system/status/"+task.Topic, res)
		})
	}

	return true
}
