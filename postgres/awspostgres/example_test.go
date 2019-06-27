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

package awspostgres_test

import (
	"context"
	"log"

	"gocloud.dev/postgres"
	_ "gocloud.dev/postgres/awspostgres"
)

func Example() {
	// This example is used in https://gocloud.dev/howto/sql/#aws

	// import _ "gocloud.dev/postgres/awspostgres"

	// Variables set up elsewhere:
	ctx := context.Background()

	// Replace these with your actual settings.
	db, err := postgres.Open(ctx,
		"awspostgres://myrole:swordfish@example01.xyzzy.us-west-1.rds.amazonaws.com/testdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Use database in your program.
	db.ExecContext(ctx, "CREATE TABLE foo (bar INT);")
}
