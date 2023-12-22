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
	"strings"

	"github.com/spf13/viper"
)

func NewConfiguration(filename string, path string, searchPaths []string) Configuration {
	c := Configuration{
		filename: filename,
		path:     path,
		paths:    searchPaths,
	}

	return c
}

type Configuration struct {
	filename string
	path     string
	paths    []string
}

func (c Configuration) getViperConfig() (string, string) {
	parts := strings.Split(c.filename, ".")
	if len(parts) < 2 {
		parts = append(parts, "yaml")
		slog.Debug("could not split filename, adding yaml", "filename", c.filename, "parts", parts)
	}
	slog.Debug("getting config file", "parts", parts)
	return parts[0], parts[1]
}

func (c Configuration) GetViper() *viper.Viper {
	v := viper.New()

	if c.path != "" {
		fullPath := c.path + "/" + c.filename
		slog.Debug("setting config file", "file", fullPath)
		v.SetConfigFile(fullPath)
	} else {
		configName, configType := c.getViperConfig()
		slog.Debug("setting config name", "name", configName)
		v.SetConfigName(configName)
		slog.Debug("setting config type", "type", configType)
		v.SetConfigType(configType)

		for _, path := range c.paths {
			slog.Debug("adding config search path", "path", path)
			v.AddConfigPath(path)
		}
	}

	return v
}
