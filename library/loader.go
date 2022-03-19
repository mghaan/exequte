package library

import (
	"encoding/json"
	"fmt"
	"os"
	"plugin"
	"strings"

	"github.com/mghaan/exequte/app"
	plug_alive "github.com/mghaan/exequte/library/alive"
	plug_run "github.com/mghaan/exequte/library/run"
	plug_status "github.com/mghaan/exequte/library/status"
	"github.com/mghaan/exequte/logger"
)

type Module interface {
	Register(json.RawMessage, *app.Server) bool
}

// Load and activate library plugins.
func LoadLibrary(server *app.Server, plugins []app.Plugin) {
	var mod Module

	for _, plugs := range plugins {
		name := plugs.Plugin
		if strings.Contains(name, ".") || strings.Contains(name, string(os.PathSeparator)) {
			continue
		}

		switch name {
		case "run":
			mod = &plug_run.Plugin
		case "alive":
			mod = &plug_alive.Plugin
		case "status":
			mod = &plug_status.Plugin
		default:
			server.Disconnect()
			server.Log().Fatal(logger.SYSTEM, fmt.Sprintf("Invalid library '%s'", name), nil)
		}

		if !mod.Register(plugs.Config, server) {
			server.Disconnect()
			server.Log().Fatal(logger.SYSTEM, fmt.Sprintf("Unable to initialize library '%s'", name), nil)
		}

		server.Log().Info(logger.SYSTEM, fmt.Sprintf("Activated library '%s'", name))
	}
}

// Load and activate external plugins.
func LoadPlugins(server *app.Server, plugins []app.Plugin) {
	for _, plugs := range plugins {
		name := plugs.Plugin
		if !strings.Contains(name, ".") && !strings.Contains(name, string(os.PathSeparator)) {
			continue
		}

		filepath := name
		plug, err := plugin.Open(filepath)
		if err != nil {
			server.Log().Error(logger.SYSTEM, fmt.Sprintf("Unable to open plugin '%s'", name), err)
			return
		}

		regs, err := plug.Lookup("Plugin")
		if err != nil {
			server.Log().Error(logger.SYSTEM, fmt.Sprintf("Unable to lookup plugin '%s'", name), err)
			return
		}

		dial, ok := regs.(Module)
		if !ok {
			server.Log().Error(logger.SYSTEM, fmt.Sprintf("Unable to analyze plugin '%s'", name), nil)
			return
		}

		if !dial.Register(plugs.Config, server) {
			server.Log().Error(logger.SYSTEM, fmt.Sprintf("Unable to initialize plugin '%s'", name), nil)
			return
		}

		server.Log().Info(logger.SYSTEM, fmt.Sprintf("Activated plugin '%s'", name))
	}
}
