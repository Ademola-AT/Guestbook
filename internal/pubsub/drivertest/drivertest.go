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

// Package drivertest provides a conformance test for implementations of
// driver.
package drivertest

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/google/go-cloud/internal/pubsub"
	"github.com/google/go-cloud/internal/pubsub/driver"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// Harness descibes the functionality test harnesses must provide to run
// conformance tests.
type Harness interface {
	// MakeTopic makes a driver.Topic for testing.
	MakeTopic(ctx context.Context) (driver.Topic, error)

	// MakeSubscription makes a driver.Subscription subscribed to the given
	// driver.Topic.
	MakeSubscription(ctx context.Context, t driver.Topic) (driver.Subscription, error)

	// Close closes resources used by the harness, but does not call Close
	// on the Topics and Subscriptions generated by the Harness.
	Close()
}

// HarnessMaker describes functions that construct a harness for running tests.
// It is called exactly once per test; Harness.Close() will be called when the test is complete.
type HarnessMaker func(ctx context.Context, t *testing.T) (Harness, error)

// RunConformanceTests runs conformance tests for provider implementations of pubsub.
func RunConformanceTests(t *testing.T, newHarness HarnessMaker) {
	t.Run("TestSendReceive", func(t *testing.T) {
		testSendReceive(t, newHarness)
	})
	t.Run("TestErrorOnSendToClosedTopic", func(t *testing.T) {
		testErrorOnSendToClosedTopic(t, newHarness)
	})
	t.Run("TestErrorOnReceiveFromClosedSubscription", func(t *testing.T) {
		testErrorOnReceiveFromClosedSubscription(t, newHarness)
	})
	t.Run("TestCancelSendReceive", func(t *testing.T) {
		testCancelSendReceive(t, newHarness)
	})
}

func testSendReceive(t *testing.T, newHarness HarnessMaker) {
	// Set up.
	ctx := context.Background()
	h, err := newHarness(ctx, t)
	if err != nil {
		t.Fatal(err)
	}
	defer h.Close()
	top, sub, cleanup, err := makePair(ctx, h)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	// Send to the topic.
	ms := []*pubsub.Message{}
	for i := 0; i < 3; i++ {
		m := &pubsub.Message{
			Body:     []byte(randStr()),
			Metadata: map[string]string{randStr(): randStr()},
		}
		if err := top.Send(ctx, m); err != nil {
			t.Fatal(err)
		}
		ms = append(ms, m)
	}

	// Receive from the subscription.
	ms2 := []*pubsub.Message{}
	for i := 0; i < len(ms); i++ {
		m2, err := sub.Receive(ctx)
		if err != nil {
			t.Fatal(err)
		}
		ms2 = append(ms2, m2)
		m2.Ack()
	}

	// Check that the received messages match the sent ones.
	less := func(x, y *pubsub.Message) bool { return bytes.Compare(x.Body, y.Body) < 0 }
	if diff := cmp.Diff(ms2, ms, cmpopts.SortSlices(less), cmpopts.IgnoreUnexported(pubsub.Message{})); diff != "" {
		t.Error(diff)
	}
}

func testErrorOnSendToClosedTopic(t *testing.T, newHarness HarnessMaker) {
	// Set up.
	ctx := context.Background()
	h, err := newHarness(ctx, t)
	if err != nil {
		t.Fatal(err)
	}
	defer h.Close()
	top, _, cleanup, err := makePair(ctx, h)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	top.Close()

	// Check that sending to the closed topic fails.
	m := &pubsub.Message{}
	if err := top.Send(ctx, m); err == nil {
		t.Error("top.Send returned nil, want error")
	}
}

func testErrorOnReceiveFromClosedSubscription(t *testing.T, newHarness HarnessMaker) {
	ctx := context.Background()
	h, err := newHarness(ctx, t)
	if err != nil {
		t.Fatal(err)
	}
	defer h.Close()
	_, sub, cleanup, err := makePair(ctx, h)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	sub.Close()
	if _, err = sub.Receive(ctx); err == nil {
		t.Error("sub.Receive returned nil, want error")
	}
}

func testCancelSendReceive(t *testing.T, newHarness HarnessMaker) {
	ctx, cancel := context.WithCancel(context.Background())
	h, err := newHarness(ctx, t)
	if err != nil {
		t.Fatal(err)
	}
	defer h.Close()
	top, sub, cleanup, err := makePair(ctx, h)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	cancel()

	m := &pubsub.Message{}
	if err := top.Send(ctx, m); err != context.Canceled {
		t.Errorf("top.Send returned %v, want context.Canceled", err)
	}
	if _, err := sub.Receive(ctx); err != context.Canceled {
		t.Errorf("sub.Receive returned %v, want context.Canceled", err)
	}
}

func randStr() string {
	return fmt.Sprintf("%d", rand.Int())
}

func makePair(ctx context.Context, h Harness) (*pubsub.Topic, *pubsub.Subscription, func(), error) {
	dt, err := h.MakeTopic(ctx)
	if err != nil {
		return nil, nil, nil, err
	}
	ds, err := h.MakeSubscription(ctx, dt)
	if err != nil {
		return nil, nil, nil, err
	}
	t := pubsub.NewTopic(dt, pubsub.NewSendBatcher(dt))
	s := pubsub.NewSubscription(ds, pubsub.NewAckBatcher(ds))
	cleanup := func() {
		t.Close()
		s.Close()
	}
	return t, s, cleanup, nil
}
