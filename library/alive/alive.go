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
