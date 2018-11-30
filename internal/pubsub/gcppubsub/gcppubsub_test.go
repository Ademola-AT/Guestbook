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

package gcppubsub

import (
	"context"
	"fmt"
	"testing"

	raw "cloud.google.com/go/pubsub/apiv1"
	"github.com/google/go-cloud/internal/pubsub/driver"
	"github.com/google/go-cloud/internal/pubsub/drivertest"
	"github.com/google/go-cloud/internal/testing/setup"
	"google.golang.org/api/option"
)

const (
	// These constants capture values that were used during the last -record.
	//
	// If you want to use --record mode,
	// 1a. Create a topic in your GCP project:
	//    https://console.cloud.google.com/cloudpubsub, then
	//    "Enable API", "Create a topic".
	// 1b. Create a subscription by clicking on the topic, then clicking on
	//    the icon at the top with a "Create subscription" tooltip.
	// 2. Update the topicName constant to your topic name, and the
	//    subscriptionName to your subscription name.
	topicName        = "test-topic"
	subscriptionName = "test-subscription-1"
	projectID        = "go-cloud-test-216917"
)

type harness struct {
	closer    func()
	pubClient *raw.PublisherClient
	subClient *raw.SubscriberClient
}

func newHarness(ctx context.Context, t *testing.T) (drivertest.Harness, error) {
	endPoint := "pubsub.googleapis.com:443"
	conn, done := setup.NewGCPgRPCConn(ctx, t, endPoint)
	pubClient, err := raw.NewPublisherClient(ctx, option.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("making publisher client: %v", err)
	}
	subClient, err := raw.NewSubscriberClient(ctx, option.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("making subscription client: %v", err)
	}
	return &harness{done, pubClient, subClient}, nil
}

func (h *harness) MakeTopic(ctx context.Context) (driver.Topic, error) {
	dt, err := openTopic(ctx, h.pubClient, projectID, topicName)
	if err != nil {
		return nil, fmt.Errorf("opening topic: %v", err)
	}
	return dt, nil
}

func (h *harness) MakeSubscription(ctx context.Context, dt driver.Topic) (driver.Subscription, error) {
	ds, err := openSubscription(ctx, h.subClient, projectID, subscriptionName)
	if err != nil {
		return nil, fmt.Errorf("opening subscription: %v", err)
	}
	return ds, nil
}

func (h *harness) Close() {
	h.pubClient.Close()
	h.subClient.Close()
	h.closer()
}

func TestConformance(t *testing.T) {
	drivertest.RunConformanceTests(t, newHarness)
}
