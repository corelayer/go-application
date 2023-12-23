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

package console

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/corelayer/go-application/pkg/base"
)

var Command = base.Command{
	Cobra: &cobra.Command{
		Use:           "console",
		Short:         "console mode",
		Long:          "example 1 - console mode",
		RunE:          execute,
		PreRunE:       executePreRun,
		SilenceErrors: true,
		SilenceUsage:  true,
		// Command annotation "logtarget" can be set on a sub-command to enforce logging to a file
		Annotations: map[string]string{
			"logtarget": "console.log",
		},
	},
	SubCommands: nil,
	Configure:   configure,
}

func executePreRun(cmd *cobra.Command, args []string) error {
	fmt.Println("CONSOLE PRE RUN")
	return nil
}

func execute(cmd *cobra.Command, args []string) error {
	fmt.Println("CONSOLE")
	slog.Info("CONSOLE INFO")
	slog.Error("CONSOLE ERROR")
	return nil
}

func configure(cmd *cobra.Command) {
	fmt.Println("console")
}
