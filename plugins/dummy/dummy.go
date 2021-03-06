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

// go build:ignore
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
