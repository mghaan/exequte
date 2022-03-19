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
package alive

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/mghaan/exequte/app"
)

const PLUGIN string = "ALIVE"

type External struct{}

var Plugin External

type Watchdog struct {
	Interval int    `json:"interval"`
	Topic    string `json:"topic"`
	Process  string `json:"process"`
}

func (plugin *External) Register(data json.RawMessage, server *app.Server) bool {
	var tasks []Watchdog

	if err := json.Unmarshal(data, &tasks); err != nil {
		server.Log().Error(PLUGIN, "Unable to parse tasks", err)

		return false
	}

	for _, task := range tasks {
		app.ScheduleTask(task.Interval, func() {
			ret := -1
			switch runtime.GOOS {
			case "linux":
				ret = proc_check_linux(task.Process)
			}
			server.Log().Info(PLUGIN, fmt.Sprintf("Check process '%s'", task.Process))
			server.Publish("system/alive/"+task.Topic, fmt.Sprintf("%d", ret))
		})
	}

	return true
}

func proc_check_linux(proc string) int {
	cmd := exec.Command("/usr/bin/pgrep", proc)
	cmd.Run()

	ret := cmd.ProcessState.ExitCode()

	return ret
}
