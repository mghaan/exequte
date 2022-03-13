package run

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/mghaan/exequte/app"
)

const PLUGIN string = "RUN"

type External struct{}

var Plugin External

type Command struct {
	Alias  string
	Script string
	Params bool
}

var commands []Command

func (plugin *External) Register(data json.RawMessage, server *app.Server) bool {
	if err := json.Unmarshal(data, &commands); err != nil {
		server.Log().Error(PLUGIN, "Unable to parse commands", err)

		return false
	}

	server.Subscribe("system/run", func(client paho.Client, message paho.Message) {
		if len(message.Payload()) > 0 {
			payload := string(message.Payload())

			server.Log().Info(PLUGIN, "Received command: "+payload)

			// parse payload
			datas := strings.Split(payload, " ")
			for c := 0; c < len(commands); c++ {
				// check the alias
				if commands[c].Alias == datas[0] {
					cmds := strings.Split(commands[c].Script, " ")
					if finfo, err := os.Stat(cmds[0]); err != nil {
						server.Log().Info(PLUGIN, fmt.Sprintf("%s: No such file or directory", commands[c].Alias))
					} else {
						if finfo.IsDir() {
							server.Log().Info(PLUGIN, fmt.Sprintf("%s: Not an executable file", commands[c].Alias))
						} else {
							// build the command and arguments
							cmd := exec.Command(cmds[0])
							if len(cmds) > 1 {
								for i := 1; i < len(cmds); i++ {
									cmd.Args = append(cmd.Args, cmds[i])
								}
							}

							// append parameters if allowed
							if commands[c].Params {
								if len(datas) > 1 {
									for i := 1; i < len(datas); i++ {
										cmd.Args = append(cmd.Args, datas[i])
									}
								}
							}

							server.Log().Info(PLUGIN, "Run alias '"+commands[c].Alias+"': "+strings.Join(cmd.Args, " "))
							cmd.Start()
						}
					}
				}
			}
		}
	})

	return true
}
