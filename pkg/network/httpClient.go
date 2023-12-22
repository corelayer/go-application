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

package network

import (
	"log/slog"
	"net/http"
	"time"
)

func NewHttpClient(useragent string, timeout int, followRedirects bool) *http.Client {
	return &http.Client{
		Transport:     NewHttpTransport(useragent),
		CheckRedirect: checkRedirect(followRedirects),
		Jar:           nil,
		Timeout:       time.Duration(timeout) * time.Second,
	}
}

func checkRedirect(state bool) func(req *http.Request, via []*http.Request) error {
	switch state {
	case true:
		return nil
	case false:
		slog.Debug("disable http redirects for http client")
		return doNotFollowHttpRedirects
	default:
		return nil
	}
}

// DoNotFollowHttpRedirects information at https://go.dev/src/net/http/client.go - line 72
func doNotFollowHttpRedirects(req *http.Request, via []*http.Request) error {
	slog.Debug("do not follow redirect", "url", req.URL)
	return http.ErrUseLastResponse
}
