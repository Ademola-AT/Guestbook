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

package gcppostgres_test

import (
	"context"
	"log"

	"gocloud.dev/postgres"
	_ "gocloud.dev/postgres/gcppostgres"
)

func Example() {
	// PRAGMA(gocloud.dev): Package this example for gocloud.dev.
	// PRAGMA(gocloud.dev): Add a blank import: _ "gocloud.dev/postgres/gcppostgres"
	// PRAGMA(gocloud.dev): Skip until next blank line.
	ctx := context.Background()

	db, err := postgres.Open(ctx,
		"gcppostgres://myrole@example-project/us-central1/my-instance01/testdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Use database in your program.
	db.Exec("CREATE TABLE foo (bar INT);")
}
