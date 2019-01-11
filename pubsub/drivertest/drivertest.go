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
package drivertest // import "gocloud.dev/pubsub/drivertest"

import (
	"bytes"
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gocloud.dev/internal/retry"
	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/driver"
)

// Harness descibes the functionality test harnesses must provide to run
// conformance tests.
type Harness interface {
	// CreateTopic creates a new topic in the provider and returns a driver.Topic for testing.
	// The topic may have to be cleaned up manually if the test is abruptly terminated or the network connection fails.
	CreateTopic(ctx context.Context) (dt driver.Topic, cleanup func(), err error)

	// MakeNonexistentTopic makes a driver.Topic referencing a topic that
	// does not exist.
	MakeNonexistentTopic(ctx context.Context) (driver.Topic, error)

	// CreateSubscription creates a new subscription in the provider, subscribed to the given topic, and returns
	// a driver.Subscription for testing.
	// The subscription may have to be cleaned up manually if the test is abruptly terminated or the network connection
	// fails.
	CreateSubscription(ctx context.Context, t driver.Topic) (ds driver.Subscription, cleanup func(), err error)

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
	t.Run("TestSendReceiveTwo", func(t *testing.T) {
		testSendReceiveTwo(t, newHarness)
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
	defer top.Shutdown(ctx)

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
	sub := pubsub.NewSubscription(ds, nil)
	defer sub.Shutdown(ctx)

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

	want := publishN(t, ctx, top, 3)
	got := receiveN(t, ctx, sub, len(want))

	// Check that the received messages match the sent ones.
	if diff := diffMessageSets(got, want); diff != "" {
		t.Error(diff)
	}
}

// Receive from two subscriptions to the same topic.
// Verify both get all the messages.
func testSendReceiveTwo(t *testing.T, newHarness HarnessMaker) {
	// Set up.
	ctx := context.Background()
	h, err := newHarness(ctx, t)
	if err != nil {
		t.Fatal(err)
	}
	defer h.Close()

	dt, cleanup, err := h.CreateTopic(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	top := pubsub.NewTopic(dt)
	defer top.Shutdown(ctx)

	var ss []*pubsub.Subscription
	for i := 0; i < 2; i++ {
		ds, cleanup, err := h.CreateSubscription(ctx, dt)
		if err != nil {
			t.Fatal(err)
		}
		defer cleanup()
		s := pubsub.NewSubscription(ds, nil)
		defer s.Shutdown(ctx)
		ss = append(ss, s)
	}
	want := publishN(t, ctx, top, 3)
	for i, s := range ss {
		got := receiveN(t, ctx, s, len(want))
		if diff := diffMessageSets(got, want); diff != "" {
			t.Errorf("sub #%d: %s", i, diff)
		}
	}
}

// Publish n different messages to the topic. Return the messages.
func publishN(t *testing.T, ctx context.Context, top *pubsub.Topic, n int) []*pubsub.Message {
	var ms []*pubsub.Message
	for i := 0; i < n; i++ {
		m := &pubsub.Message{
			Body:     []byte(strconv.Itoa(i)),
			Metadata: map[string]string{"a": strconv.Itoa(i)},
		}
		if err := top.Send(ctx, m); err != nil {
			t.Fatal(err)
		}
		ms = append(ms, m)
	}
	return ms
}

// Receive and ack n messages from sub.
func receiveN(t *testing.T, ctx context.Context, sub *pubsub.Subscription, n int) []*pubsub.Message {
	var ms []*pubsub.Message
	for i := 0; i < n; i++ {
		m, err := sub.Receive(ctx)
		if err != nil {
			t.Fatal(err)
		}
		ms = append(ms, m)
		m.Ack()
	}
	return ms
}

// Find the differences between two sets of messages.
func diffMessageSets(got, want []*pubsub.Message) string {
	less := func(x, y *pubsub.Message) bool { return bytes.Compare(x.Body, y.Body) < 0 }
	return cmp.Diff(got, want, cmpopts.SortSlices(less), cmpopts.IgnoreUnexported(pubsub.Message{}))
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

	top.Shutdown(ctx)

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
	sub.Shutdown(ctx)
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
	if err := top.Send(ctx, m); !isCanceled(err) {
		t.Errorf("top.Send returned %v (%T), want context.Canceled", err, err)
	}
	if _, err := sub.Receive(ctx); !isCanceled(err) {
		t.Errorf("sub.Receive returned %v (%T), want context.Canceled", err, err)
	}
}

func isCanceled(err error) bool {
	if err == context.Canceled {
		return true
	}
	if cerr, ok := err.(*retry.ContextError); ok {
		return cerr.CtxErr == context.Canceled
	}
	return false
}

func makePair(ctx context.Context, h Harness) (*pubsub.Topic, *pubsub.Subscription, func(), error) {
	dt, topicCleanup, err := h.CreateTopic(ctx)
	if err != nil {
		return nil, nil, nil, err
	}
	ds, subCleanup, err := h.CreateSubscription(ctx, dt)
	if err != nil {
		return nil, nil, nil, err
	}
	t := pubsub.NewTopic(dt)
	s := pubsub.NewSubscription(ds, nil)
	cleanup := func() {
		topicCleanup()
		subCleanup()
		t.Shutdown(ctx)
		s.Shutdown(ctx)
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
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	if err := st.TopicCheck(top); err != nil {
		t.Error(err)
	}
	if err := st.SubscriptionCheck(sub); err != nil {
		t.Error(err)
	}
}
