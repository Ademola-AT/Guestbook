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

package gcppubsub_test

import (
	"context"
	"fmt"
	"log"

	"gocloud.dev/gcp"
	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/gcppubsub"
)

func ExampleOpenTopic() {
	ctx := context.Background()

	// Get GCP credentials.
	// Here we use a fake JSON credentials file, but you could also use
	// gcp.DefaultCredentials(ctx) to use the default GCP credentials from
	// the environment.
	// See https://cloud.google.com/docs/authentication/production
	// for more info on alternatives.
	creds, err := gcp.FakeDefaultCredentials(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Open a gRPC connection to the GCP Pub Sub API.
	conn, cleanup, err := gcppubsub.Dial(ctx, creds.TokenSource)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()
	pubClient, err := gcppubsub.PublisherClient(ctx, conn)
	if err != nil {
		log.Fatal(err)
	}
	defer pubClient.Close()
	proj, err := gcp.DefaultProjectID(creds)
	if err != nil {
		log.Fatal(err)
	}
	t := gcppubsub.OpenTopic(ctx, pubClient, proj, "example-topic", nil)
	defer t.Shutdown(ctx)
	if err := t.Send(ctx, &pubsub.Message{Body: []byte("example message")}); err != nil {
		fmt.Println("Message send failure is expected due to the fake credentials.")
	}

	// Output:
	// Message send failure is expected due to the fake credentials.
}

func ExampleOpenSubscription() {
	ctx := context.Background()

	// Get GCP credentials.
	// Here we use a fake JSON credentials file, but you could also use
	// gcp.DefaultCredentials(ctx) to use the default GCP credentials from
	// the environment.
	// See https://cloud.google.com/docs/authentication/production
	// for more info on alternatives.
	creds, err := gcp.FakeDefaultCredentials(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Open a gRPC connection to the GCP Pub Sub API.
	conn, cleanup, err := gcppubsub.Dial(ctx, creds.TokenSource)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()
	proj, err := gcp.DefaultProjectID(creds)
	if err != nil {
		log.Fatal(err)
	}
	subClient, err := gcppubsub.SubscriberClient(ctx, conn)
	if err != nil {
		log.Fatal(err)
	}
	defer subClient.Close()
	s := gcppubsub.OpenSubscription(ctx, subClient, proj, "example-subscription", nil)
	defer s.Shutdown(ctx)
	m, err := s.Receive(ctx)
	if err != nil {
		fmt.Println("Message receive failure is expected due to the fake credentials.")
		return
	}
	fmt.Printf("%s\n", m.Body)
	m.Ack()

	// Output:
	// Message receive failure is expected due to the fake credentials.
}
