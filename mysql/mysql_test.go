// Copyright 2019 The Go Cloud Development Kit Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mysql

import (
	"context"
	"net/url"
	"testing"
)

func TestOpen(t *testing.T) {
	t.Skip("Test not hermetic yet")

	ctx := context.Background()
	dbByURL, err := Open(ctx, "mysql://root@localhost/mysql")
	if err != nil {
		t.Fatal(err)
	}
	if err := dbByURL.Ping(); err != nil {
		t.Error("Ping:", err)
	}
	if err := dbByURL.Close(); err != nil {
		t.Error("Close:", err)
	}
}

func TestConfigFromURL(t *testing.T) {
	const urlstr = "mysql://localhost/db?parseTime=true&interpolateParams=true"
	u, err := url.Parse(urlstr)
	if err != nil {
		t.Fatal(err)
	}
	cfg, err := ConfigFromURL(u)
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.ParseTime {
		t.Error("ParseTime is false, want true")
	}
	if !cfg.InterpolateParams {
		t.Error("InterpolateParams is false, want true")
	}
}
