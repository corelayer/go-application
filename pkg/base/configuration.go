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
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func NewConfiguration(file string, searchPaths []string) Configuration {
	// Clean search paths
	paths := make([]string, 0)
	for _, p := range searchPaths {
		if !strings.Contains(p, "..") {
			paths = append(paths, filepath.Clean(p))
		}
	}

	// Split file into path and filename
	path, filename := filepath.Split(file)

	// Make sure to clean path, this also causes the path to either be "." or a full path
	if path != "" {
		path = filepath.Clean(path)
	}

	c := Configuration{
		filename: filename,
		path:     path,
		paths:    paths,
	}

	return c
}

type Configuration struct {
	filename string
	path     string
	paths    []string
}

func (c Configuration) GetViperConfig() (string, string) {
	var (
		configName string
		configType string
	)

	fileExtension := filepath.Ext(c.filename)
	if fileExtension == "" {
		configType = "yaml"
	} else {
		configType = strings.TrimPrefix(fileExtension, ".")
	}

	configName = strings.TrimSuffix(c.filename, fileExtension)

	return configName, configType
}

func (c Configuration) GetViper() *viper.Viper {
	v := viper.New()

	// If a full path is specified, set the config file to that path
	if c.path != "" {
		fullPath := filepath.Join(c.path, c.filename)
		slog.Debug("setting config file", "file", fullPath)
		v.SetConfigFile(fullPath)
	} else {
		configName, configType := c.GetViperConfig()
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
