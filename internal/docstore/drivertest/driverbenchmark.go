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

package drivertest

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"
	"testing"

	"gocloud.dev/internal/docstore"
)

// RunBenchmarks runs benchmarks for docstore drivers.
func RunBenchmarks(b *testing.B, coll *docstore.Collection) {
	if err := cleanUpTable(newDocmap, coll); err != nil {
		b.Fatalf("%+v", err)
	}
	b.Run("BenchmarkSingleActionPut", func(b *testing.B) {
		benchmarkSingleActionPut(25, b, coll)
	})
	b.Run("BenchmarkSingleActionGet", func(b *testing.B) {
		benchmarkSingleActionGet(25, b, coll)
	})
	b.Run("BenchmarkActionListPut", func(b *testing.B) {
		benchmarkActionListPut(100, b, coll)
	})
	b.Run("BenchmarkActionListGet", func(b *testing.B) {
		benchmarkActionListGet(100, b, coll)
	})
	if err := cleanUpTable(newDocmap, coll); err != nil {
		b.Fatalf("%+v", err)
	}
}

func benchmarkSingleActionPut(n int, b *testing.B, coll *docstore.Collection) {
	ctx := context.Background()
	const baseKey = "benchmarksingleaction-put-"
	var nextID uint32

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < n; i++ {
				key := fmt.Sprintf("%s%d", baseKey, atomic.AddUint32(&nextID, 1))
				doc := docmap{KeyField: key, "S": key}
				if err := coll.Put(ctx, doc); err != nil {
					b.Error(err)
				}
			}
		}
	})
}

func benchmarkSingleActionGet(n int, b *testing.B, coll *docstore.Collection) {
	ctx := context.Background()
	const baseKey = "benchmarksingleaction-get-"
	docs := make([]docmap, n)
	puts := coll.Actions()
	for i := 0; i < n; i++ {
		docs[i] = docmap{KeyField: baseKey + strconv.Itoa(i), "n": i}
		puts.Put(docs[i])
	}
	if err := puts.Unordered().Do(ctx); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for _, doc := range docs {
				got := docmap{KeyField: doc[KeyField]}
				if err := coll.Get(ctx, got); err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

func benchmarkActionListPut(n int, b *testing.B, coll *docstore.Collection) {
	ctx := context.Background()
	const baseKey = "benchmarkactionlist-put-"
	var nextID uint32

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			actions := coll.Actions()
			for i := 0; i < n; i++ {
				key := fmt.Sprintf("%s%d", baseKey, atomic.AddUint32(&nextID, 1))
				doc := docmap{KeyField: key, "S": key}
				actions.Put(doc)
			}
			if err := actions.Unordered().Do(ctx); err != nil {
				b.Error(err)
			}
		}
	})
}

func benchmarkActionListGet(n int, b *testing.B, coll *docstore.Collection) {
	ctx := context.Background()
	const baseKey = "benchmarkactionlist-get-"
	docs := make([]docmap, n)
	puts := coll.Actions()
	for i := 0; i < n; i++ {
		docs[i] = docmap{KeyField: baseKey + strconv.Itoa(i), "n": i}
		puts.Put(docs[i])
	}
	if err := puts.Unordered().Do(ctx); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			gets := coll.Actions()
			for _, doc := range docs {
				got := docmap{KeyField: doc[KeyField]}
				gets.Get(got)
			}
			if err := gets.Unordered().Do(ctx); err != nil {
				b.Fatal(err)
			}
		}
	})
}
