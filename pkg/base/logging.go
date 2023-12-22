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

func GetLogger(cmd *cobra.Command, w io.Writer) (*slog.Logger, error) {
	var (
		err           error
		logFlag       bool
		logLevelFlag  string
		logFormatFlag string

		format logFormat
		found  bool

		level  slog.Leveler
		source bool
	)

	logFlag, err = cmd.Flags().GetBool("log")
	if err != nil {
		return nil, err
	}

	if !logFlag {
		return slog.New(slog.NewTextHandler(io.Discard, nil)), nil
	}

	logLevelFlag, err = cmd.Flags().GetString("loglevel")
	if err != nil {
		return nil, err
	}
	logFormatFlag, err = cmd.Flags().GetString("logformat")
	if err != nil {
		return nil, err
	}

	format, found = parseLogHandler(logFormatFlag)
	if !found {
		return nil, fmt.Errorf("invalid logFormat value %s", logFormatFlag)
	}

	switch logLevelFlag {
	case "error":
		level = slog.LevelError
	case "warn":
		level = slog.LevelWarn
	case "info":
		level = slog.LevelInfo
	case "debug":
		level = slog.LevelDebug
	default:
		level = slog.LevelInfo
	}

	if level == slog.LevelDebug {
		source = true
	}

	switch format {
	case defaultLogFormat:
		return slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: level, AddSource: source})), nil
	case textLogFormat:
		return slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: level, AddSource: source})), nil
	case jsonLogFormat:
		return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: level, AddSource: source})), nil
	}
	return nil, err
}

func parseLogHandler(s string) (logFormat, bool) {
	c, ok := logFormats[strings.ToLower(s)]
	return c, ok
}
