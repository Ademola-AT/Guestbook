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

package etcdvar_test

import (
	"context"
	"log"

	"go.etcd.io/etcd/clientv3"
	"gocloud.dev/runtimevar"
	"gocloud.dev/runtimevar/etcdvar"
)

func ExampleOpenVariable() {
	// This example is used in https://gocloud.dev/howto/runtimevar/runtimevar/#etcd-ctor

	// Connect to the etcd server.
	client, err := clientv3.NewFromURL("http://your.etcd.server:9999")
	if err != nil {
		log.Fatal(err)
	}

	// Construct a *runtimevar.Variable that watches the variable.
	v, err := etcdvar.OpenVariable(client, "cfg-variable-name", runtimevar.StringDecoder, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer v.Close()
}

func Example_openVariableFromURL() {
	// This example is used in https://gocloud.dev/howto/runtimevar/runtimevar/#etcd-url

	// import _ "gocloud.dev/runtimevar/etcdvar"

	// runtimevar.OpenVariable creates a *runtimevar.Variable from a URL.
	// The default opener connects to an etcd server based on the environment
	// variable ETCD_SERVER_URL.

	// Variables set up elsewhere:
	ctx := context.Background()

	v, err := runtimevar.OpenVariable(ctx, "etcd://myvarname?decoder=string")
	if err != nil {
		log.Fatal(err)
	}
	defer v.Close()
}
