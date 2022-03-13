//go:build ignore

package main

import (
	"encoding/json"

	"github.com/mghaan/exequte/app"
)

const PLUGIN string = "DUMMY"

type External struct{}

var Plugin External

func (plugin *External) Register(data json.RawMessage, server *app.Server) bool {
	return true
}
