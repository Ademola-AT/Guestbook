// Copyright 2018 The Go Cloud Development Kit Authors
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

package fileblob

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"gocloud.dev/blob"
	"gocloud.dev/blob/driver"
	"gocloud.dev/blob/drivertest"
)

type harness struct {
	dir    string
	closer func()
}

func newHarness(ctx context.Context, t *testing.T) (drivertest.Harness, error) {
	dir := path.Join(os.TempDir(), "go-cloud-fileblob")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}
	return &harness{
		dir:    dir,
		closer: func() { _ = os.RemoveAll(dir) },
	}, nil
}

func (h *harness) HTTPClient() *http.Client {
	return nil
}

func (h *harness) MakeDriver(ctx context.Context) (driver.Bucket, error) {
	return openBucket(h.dir, nil)
}

func (h *harness) Close() {
	h.closer()
}

func TestConformance(t *testing.T) {
	drivertest.RunConformanceTests(t, newHarness, []drivertest.AsTest{verifyPathError{}})
}

func BenchmarkFileblob(b *testing.B) {
	dir := path.Join(os.TempDir(), "go-cloud-fileblob")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		b.Fatal(err)
	}
	bkt, err := OpenBucket(dir, nil)
	if err != nil {
		b.Fatal(err)
	}
	drivertest.RunBenchmarks(b, bkt)
}

// File-specific unit tests.
func TestNewBucket(t *testing.T) {
	t.Run("BucketDirMissing", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "fileblob")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(dir)
		_, gotErr := OpenBucket(filepath.Join(dir, "notfound"), nil)
		if gotErr == nil {
			t.Errorf("got nil want error")
		}
	})
	t.Run("BucketIsFile", func(t *testing.T) {
		f, err := ioutil.TempFile("", "fileblob")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(f.Name())
		_, gotErr := OpenBucket(f.Name(), nil)
		if gotErr == nil {
			t.Errorf("got nil want error")
		}
	})
}

func TestMungePath(t *testing.T) {
	tests := []struct {
		URL       string
		Separator uint8
		Want      string
	}{
		{
			URL:       "file:///a/directory",
			Separator: '/',
			Want:      "/a/directory",
		},
		{
			URL:       "file://localhost/a/directory",
			Separator: '/',
			Want:      "/a/directory",
		},
		{
			URL:       "file://localhost/mybucket",
			Separator: '/',
			Want:      "/mybucket",
		},
		{
			URL:       "file:///c:/a/directory",
			Separator: '\\',
			Want:      "c:\\a\\directory",
		},
		{
			URL:       "file://localhost/c:/a/directory",
			Separator: '\\',
			Want:      "c:\\a\\directory",
		},
	}
	for _, tc := range tests {
		u, err := url.Parse(tc.URL)
		if err != nil {
			t.Fatalf("%s: failed to parse URL: %v", tc.URL, err)
		}
		if got := mungeURLPath(u.Path, tc.Separator); got != tc.Want {
			t.Errorf("%s: got %q want %q", tc.URL, got, tc.Want)
		}
	}
}

type verifyPathError struct{}

func (verifyPathError) Name() string { return "verify ErrorAs handles os.PathError" }

func (verifyPathError) BucketCheck(b *blob.Bucket) error             { return nil }
func (verifyPathError) BeforeWrite(as func(interface{}) bool) error  { return nil }
func (verifyPathError) BeforeList(as func(interface{}) bool) error   { return nil }
func (verifyPathError) AttributesCheck(attrs *blob.Attributes) error { return nil }
func (verifyPathError) ReaderCheck(r *blob.Reader) error             { return nil }
func (verifyPathError) ListObjectCheck(o *blob.ListObject) error     { return nil }

func (verifyPathError) ErrorCheck(b *blob.Bucket, err error) error {
	var perr *os.PathError
	if !b.ErrorAs(err, &perr) {
		return errors.New("want ErrorAs to succeed for PathError")
	}
	const wantSuffix = "go-cloud-fileblob/key-does-not-exist"
	if got := perr.Path; !strings.HasSuffix(got, wantSuffix) {
		return fmt.Errorf("got path %q, want suffix %q", got, wantSuffix)
	}
	return nil
}
