// Copyright 2018 The Go Cloud Authors
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

package pubsub_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/go-cloud/internal/pubsub"
	"github.com/google/go-cloud/internal/pubsub/mempubsub"
)

func ExampleSendReceive() {
	// Open a topic and corresponding subscription.
	ctx := context.Background()
	dt := mempubsub.OpenTopic()
	ds := mempubsub.OpenSubscription(dt, time.Second)
	t := pubsub.NewTopic(ctx, dt)
	s := pubsub.NewSubscription(ctx, ds)

	// Send a message to the topic.
	if err := t.Send(ctx, &pubsub.Message{Body: []byte("Hello, world!")}); err != nil {
		log.Fatal(err)
	}

	// Receive a message from the subscription.
	m, err := s.Receive(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Print out the received message.
	fmt.Printf("%s\n", m.Body)

	// Output:
	// Hello, world!
}
