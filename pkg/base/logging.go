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
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/spf13/cobra"
)

type logFormat int

func (h logFormat) String() string {
	return [...]string{"none", "text", "json"}[h]
}

const (
	defaultLogFormat logFormat = iota
	textLogFormat
	jsonLogFormat
)

var logFormats = map[string]logFormat{
	"text": textLogFormat,
	"json": jsonLogFormat,
}

func parseLogHandler(s string) (logFormat, bool) {
	c, ok := logFormats[strings.ToLower(s)]
	return c, ok
}

func GetLogger(cmd *cobra.Command, w io.Writer) (*slog.Logger, error) {
	var (
		err           error
		logLevelFlag  string
		logFormatFlag string
	)

	logLevelFlag, err = cmd.Flags().GetString("loglevel")
	if err != nil {
		return nil, err
	}
	logFormatFlag, err = cmd.Flags().GetString("logformat")
	if err != nil {
		return nil, err
	}

	logFormat, found := parseLogHandler(logFormatFlag)
	if !found {
		return nil, fmt.Errorf("invalid logFormat value %s", logFormatFlag)
	}

	var level slog.Leveler
	switch logLevelFlag {
	case "warn":
		level = slog.LevelWarn
	case "info":
		level = slog.LevelInfo
	case "debug":
		level = slog.LevelDebug
	default:
		level = slog.LevelInfo
	}

	switch logFormat {
	case defaultLogFormat:
		return slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: level})), nil
	case textLogFormat:
		return slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: level})), nil
	case jsonLogFormat:
		return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: level})), nil
	}
	return nil, err
}
