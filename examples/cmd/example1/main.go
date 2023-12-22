/*
 * Copyright 2023 CoreLayer BV
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package main

import (
	"os"

	"github.com/corelayer/go-application/examples/cmd/example1/cmd/sub"
	"github.com/corelayer/go-application/pkg/base"
)

const (
	APPLICATION_NAME   = "example1"
	APPLICATION_TITLE  = "example 1 title"
	APPLICATION_BANNER = "example 1 banner"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	// var err error
	app := base.NewApplication(APPLICATION_NAME, APPLICATION_TITLE, APPLICATION_BANNER, "")

	app.RegisterCommands([]base.Commander{
		sub.Command,
	})
	return app.Run()
}
