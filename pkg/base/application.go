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

package base

import (
	"log/slog"

	"github.com/spf13/cobra"
)

func NewApplication(use string, short string, long string, version string) *Application {
	a := &Application{}

	command := &cobra.Command{
		Use:     use,
		Short:   short,
		Long:    long,
		Args:    cobra.MinimumNArgs(1),
		Version: version,
	}

	a.Command = command
	return a
}

type Application struct {
	Command *cobra.Command
}

func (a *Application) RegisterCommands(c []Commander) {
	for _, cmdr := range c {
		a.Command.AddCommand(cmdr.Initialize())
	}
}

func (a *Application) Run() error {
	if err := a.Command.Execute(); err != nil {
		slog.Error("application terminated unexpectedly", "name", a.Command.Name(), "error", err)
		return err
	}
	return nil
}
