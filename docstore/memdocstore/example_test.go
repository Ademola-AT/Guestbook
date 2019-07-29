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

package memdocstore_test

import (
	"context"
	"log"

	"gocloud.dev/docstore"
	"gocloud.dev/docstore/memdocstore"
)

func ExampleOpenCollection() {
	// PRAGMA(gocloud.dev): Package this example for gocloud.dev.

	coll, err := memdocstore.OpenCollection("userID", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer coll.Close()
}

func ExampleOpenCollectionWithKeyFunc() {
	// PRAGMA(gocloud.dev): Package this example for gocloud.dev.
	// PRAGMA(gocloud.dev): Skip until next blank line.
	type HighScore struct {
		Game   string
		Player string
	}

	// The name of a document is constructed from the Game and Player fields.
	nameFromDocument := func(doc docstore.Document) interface{} {
		hs := doc.(*HighScore)
		return hs.Game + "|" + hs.Player
	}

	coll, err := memdocstore.OpenCollectionWithKeyFunc(nameFromDocument, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer coll.Close()
}

func Example_openCollectionFromURL() {
	// PRAGMA(gocloud.dev): Package this example for gocloud.dev.
	// PRAGMA(gocloud.dev): Add a blank import: _ "gocloud.dev/docstore/memdocstore"
	// PRAGMA(gocloud.dev): Skip until next blank line.
	ctx := context.Background()

	// docstore.OpenCollection creates a *docstore.Collection from a URL.
	coll, err := docstore.OpenCollection(ctx, "mem://userID")
	if err != nil {
		log.Fatal(err)
	}
	defer coll.Close()
}
