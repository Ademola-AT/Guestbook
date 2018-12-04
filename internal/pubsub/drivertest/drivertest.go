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
	"errors"
	"strconv"
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

	// MakeNonexistentTopic makes a driver.Topic referencing a topic that
	// does not exist.
	MakeNonexistentTopic(ctx context.Context) (driver.Topic, error)

	// MakeSubscription makes a driver.Subscription subscribed to the given
	// driver.Topic.
	MakeSubscription(ctx context.Context, t driver.Topic) (driver.Subscription, error)

	// MakeNonexistentSubscription makes a driver.Subscription referencing a
	// subscription that does not exist.
	MakeNonexistentSubscription(ctx context.Context) (driver.Subscription, error)

	// Close closes resources used by the harness, but does not call Close
	// on the Topics and Subscriptions generated by the Harness.
	Close()
}

// HarnessMaker describes functions that construct a harness for running tests.
// It is called exactly once per test; Harness.Close() will be called when the test is complete.
type HarnessMaker func(ctx context.Context, t *testing.T) (Harness, error)

// AsTest represents a test of As functionality.
// The conformance test:
// 1. Calls TopicCheck.
// 2. Calls SubscriptionCheck.
type AsTest interface {
	// Name should return a descriptive name for the test.
	Name() string
	// TopicCheck will be called to allow verifcation of Topic.As.
	TopicCheck(t *pubsub.Topic) error
	// SubscriptionCheck will be called to allow verification of Subscription.As.
	SubscriptionCheck(s *pubsub.Subscription) error
}

type verifyAsFailsOnNil struct{}

func (verifyAsFailsOnNil) Name() string {
	return "verify As returns false when passed nil"
}

func (verifyAsFailsOnNil) TopicCheck(t *pubsub.Topic) error {
	if t.As(nil) {
		return errors.New("want Topic.As to return false when passed nil")
	}
	return nil
}

func (verifyAsFailsOnNil) SubscriptionCheck(s *pubsub.Subscription) error {
	if s.As(nil) {
		return errors.New("want Subscription.As to return false when passed nil")
	}
	return nil
}

// RunConformanceTests runs conformance tests for provider implementations of pubsub.
func RunConformanceTests(t *testing.T, newHarness HarnessMaker, asTests []AsTest) {
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
	t.Run("TestNonExistentTopicSucceedsOnOpenButFailsOnSend", func(t *testing.T) {
		testNonExistentTopicSucceedsOnOpenButFailsOnSend(t, newHarness)
	})
	t.Run("TestNonExistentSubscriptionSucceedsOnOpenButFailsOnSend", func(t *testing.T) {
		testNonExistentSubscriptionSucceedsOnOpenButFailsOnSend(t, newHarness)
	})
	asTests = append(asTests, verifyAsFailsOnNil{})
	t.Run("TestAs", func(t *testing.T) {
		for _, st := range asTests {
			if st.Name() == "" {
				t.Fatalf("AsTest.Name is required")
			}
			t.Run(st.Name(), func(t *testing.T) {
				testAs(t, newHarness, st)
			})
		}
	})
}

func testNonExistentTopicSucceedsOnOpenButFailsOnSend(t *testing.T, newHarness HarnessMaker) {
	// Set up.
	ctx := context.Background()
	h, err := newHarness(ctx, t)
	if err != nil {
		t.Fatal(err)
	}
	defer h.Close()

	dt, err := h.MakeNonexistentTopic(ctx)
	if err != nil {
		// Failure shouldn't happen for non-existent topics until messages are sent
		// to them.
		t.Fatalf("creating a local topic that doesn't exist on the server: %v", err)
	}
	top := pubsub.NewTopic(dt)
	defer top.Close()

	m := &pubsub.Message{}
	err = top.Send(ctx, m)
	if err == nil {
		t.Errorf("got no error for send to non-existent topic")
	}
}

func testNonExistentSubscriptionSucceedsOnOpenButFailsOnSend(t *testing.T, newHarness HarnessMaker) {
	// Set up.
	ctx := context.Background()
	h, err := newHarness(ctx, t)
	if err != nil {
		t.Fatal(err)
	}
	defer h.Close()

	ds, err := h.MakeNonexistentSubscription(ctx)
	if err != nil {
		t.Skipf("failed to make non-existent subscription: %v", err)
	}
	sub := pubsub.NewSubscription(ds)
	defer sub.Close()

	_, err = sub.Receive(ctx)
	if err == nil {
		t.Errorf("got no error for send to non-existent topic")
	}
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
	var want []*pubsub.Message
	for i := 0; i < 3; i++ {
		m := &pubsub.Message{
			Body:     []byte(strconv.Itoa(i)),
			Metadata: map[string]string{"a": strconv.Itoa(i)},
		}
		if err := top.Send(ctx, m); err != nil {
			t.Fatal(err)
		}
		want = append(want, m)
	}

	// Receive from the subscription.
	var got []*pubsub.Message
	for i := 0; i < len(want); i++ {
		m, err := sub.Receive(ctx)
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, m)
		m.Ack()
	}

	// Check that the received messages match the sent ones.
	less := func(x, y *pubsub.Message) bool { return bytes.Compare(x.Body, y.Body) < 0 }
	if diff := cmp.Diff(got, want, cmpopts.SortSlices(less), cmpopts.IgnoreUnexported(pubsub.Message{})); diff != "" {
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

func makePair(ctx context.Context, h Harness) (*pubsub.Topic, *pubsub.Subscription, func(), error) {
	dt, err := h.MakeTopic(ctx)
	if err != nil {
		return nil, nil, nil, err
	}
	ds, err := h.MakeSubscription(ctx, dt)
	if err != nil {
		return nil, nil, nil, err
	}
	t := pubsub.NewTopic(dt)
	s := pubsub.NewSubscription(ds)
	cleanup := func() {
		t.Close()
		s.Close()
	}
	return t, s, cleanup, nil
}

// testAs tests the various As functions, using AsTest.
func testAs(t *testing.T, newHarness HarnessMaker, st AsTest) {
	ctx := context.Background()
	h, err := newHarness(ctx, t)
	if err != nil {
		t.Fatal(err)
	}
	defer h.Close()
	top, sub, cleanup, err := makePair(ctx, h)
	defer cleanup()
	if err := st.TopicCheck(top); err != nil {
		t.Error(err)
	}
	if err := st.SubscriptionCheck(sub); err != nil {
		t.Error(err)
	}
}
