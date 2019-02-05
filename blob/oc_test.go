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

package blob_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"gocloud.dev/blob"
	"gocloud.dev/blob/memblob"
	"gocloud.dev/gcerrors"
	"gocloud.dev/internal/oc"
	"gocloud.dev/internal/testing/octest"
)

func TestOpenCensus(t *testing.T) {
	ctx := context.Background()
	te := octest.NewTestExporter(blob.OpenCensusViews)
	defer te.Unregister()

	bytes := []byte("foo")
	b := memblob.OpenBucket(nil)
	if err := b.WriteAll(ctx, "key", bytes, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := b.ReadAll(ctx, "key"); err != nil {
		t.Fatal(err)
	}
	if _, err := b.Attributes(ctx, "key"); err != nil {
		t.Fatal(err)
	}
	if err := b.Delete(ctx, "key"); err != nil {
		t.Fatal(err)
	}
	if _, err := b.ReadAll(ctx, "noSuchKey"); err == nil {
		t.Fatal("got nil, want error")
	}

	const provider = "gocloud.dev/blob/memblob"

	diff := octest.Diff(te.Spans(), te.Counts(), "gocloud.dev/blob", provider, []octest.Call{
		{"NewWriter", gcerrors.OK},
		{"NewRangeReader", gcerrors.OK},
		{"Attributes", gcerrors.OK},
		{"Delete", gcerrors.OK},
		{"NewRangeReader", gcerrors.NotFound},
	})
	if diff != "" {
		t.Error(diff)
	}

	// Find and verify the bytes read/written metrics.
	var sawRead, sawWritten bool
	tags := []tag.Tag{{Key: oc.ProviderKey, Value: provider}}
	for !sawRead || !sawWritten {
		data := <-te.Stats
		switch data.View.Name {
		case "gocloud.dev/blob/bytes_read":
			if sawRead {
				continue
			}
			sawRead = true
		case "gocloud.dev/blob/bytes_written":
			if sawWritten {
				continue
			}
			sawWritten = true
		default:
			continue
		}
		diff = cmp.Diff(data.Rows[0], &view.Row{Tags: tags, Data: &view.SumData{Value: float64(len(bytes))}})
		if diff != "" {
			t.Errorf("%s: %s", data.View.Name, diff)
		}
	}
}
