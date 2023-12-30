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

import "testing"

// TestNewConfiguration TODO improve test on search paths
func TestNewConfiguration(t *testing.T) {
	tests := []struct {
		file              string
		searchPaths       []string
		wantedFilename    string
		wantedPath        string
		wantedSearchPaths []string
	}{
		{
			file:              "config.yaml",
			searchPaths:       nil,
			wantedFilename:    "config.yaml",
			wantedPath:        "",
			wantedSearchPaths: nil,
		},
		{
			file:              "./config.yaml",
			searchPaths:       nil,
			wantedFilename:    "config.yaml",
			wantedPath:        ".",
			wantedSearchPaths: nil,
		},
		{
			file:              "/etc/configuration/config.yaml",
			searchPaths:       nil,
			wantedFilename:    "config.yaml",
			wantedPath:        "/etc/configuration",
			wantedSearchPaths: nil,
		},
		{
			file:              "/etc//configuration/config.yaml",
			searchPaths:       nil,
			wantedFilename:    "config.yaml",
			wantedPath:        "/etc/configuration",
			wantedSearchPaths: nil,
		},
		// SearchPath tests
		{
			file: "config.yaml",
			searchPaths: []string{
				"$PWD",
				"/normal/path",
				"//doubleslash/in/front",
				"/doubleslash//in/middle",
				"/etc/../invalid",
			},
			wantedFilename: "config.yaml",
			wantedPath:     "",
			wantedSearchPaths: []string{
				"$PWD",
				"/normal/path",
				"/doubleslash/in/front",
				"/doubleslash/in/middle",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			config := NewConfiguration(tt.file, tt.searchPaths)
			if config.filename != tt.wantedFilename {
				t.Errorf("got filename %s, expected %s", config.filename, tt.wantedFilename)
			}
			if config.path != tt.wantedPath {
				t.Errorf("got path %s, expected %s", config.path, tt.wantedPath)
			}
			for i := range config.paths {
				if config.paths[i] != tt.wantedSearchPaths[i] {
					t.Errorf("unexpected path %s", config.paths[i])
				}
			}
		})
	}
}

func TestConfiguration_GetViperConfig(t *testing.T) {
	tests := []struct {
		file       string
		wantedName string
		wantedType string
	}{
		{
			file:       "",
			wantedName: "",
			wantedType: "yaml",
		},
		{
			file:       "config.yaml",
			wantedName: "config",
			wantedType: "yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			config := Configuration{filename: tt.file}
			configName, configType := config.GetViperConfig()
			if configName != tt.wantedName {
				t.Errorf("got: %s, expected: %s", configName, tt.wantedName)
			}

			if configType != tt.wantedType {
				t.Errorf("got: %s, expected: %s", configType, tt.wantedType)
			}
		})
	}
}
