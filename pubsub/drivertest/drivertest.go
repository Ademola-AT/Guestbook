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

// Package drivertest provides a conformance test for implementations of
// driver.
package drivertest // import "gocloud.dev/pubsub/drivertest"

import (
	"bytes"
	"context"
	"errors"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gocloud.dev/internal/escape"
	"gocloud.dev/internal/retry"
	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/driver"
	"golang.org/x/sync/errgroup"
)

// Options contains settings for RunConformanceTests.
type Options struct {
	SkipTestsOfNonExistentThings bool
}

// Harness descibes the functionality test harnesses must provide to run
// conformance tests.
type Harness interface {
	// CreateTopic creates a new topic in the provider and returns a driver.Topic for testing.
	// The topic may have to be removed manually if the test is abruptly terminated or the network connection fails.
	CreateTopic(ctx context.Context, testName string) (dt driver.Topic, cleanup func(), err error)

	// MakeNonexistentTopic makes a driver.Topic referencing a topic that
	// does not exist.
	MakeNonexistentTopic(ctx context.Context) (driver.Topic, error)

	// CreateSubscription creates a new subscription in the provider, subscribed to the given topic, and returns
	// a driver.Subscription for testing.
	// The subscription may have to be cleaned up manually if the test is abruptly terminated or the network connection
	// fails.
	CreateSubscription(ctx context.Context, t driver.Topic, testName string) (ds driver.Subscription, cleanup func(), err error)

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
// 3. Calls TopicErrorCheck.
// 4. Calls SubscriptionErrorCheck.
// 5. Calls MessageCheck.
type AsTest interface {
	// Name should return a descriptive name for the test.
	Name() string
	// TopicCheck will be called to allow verifcation of Topic.As.
	TopicCheck(t *pubsub.Topic) error
	// SubscriptionCheck will be called to allow verification of Subscription.As.
	SubscriptionCheck(s *pubsub.Subscription) error
	// TopicErrorCheck will be called to allow verification of Topic.ErrorAs.
	// The error will be the one returned from SendBatch when called with
	// a non-existent topic.
	TopicErrorCheck(t *pubsub.Topic, err error) error
	// SubscriptionErrorCheck will be called to allow verification of
	// Subscription.ErrorAs.
	// The error will be the one returned from ReceiveBatch when called with
	// a non-existent subscription.
	SubscriptionErrorCheck(s *pubsub.Subscription, err error) error
	// MessageCheck will be called to allow verification of Message.As.
	MessageCheck(m *pubsub.Message) error
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

func (verifyAsFailsOnNil) TopicErrorCheck(t *pubsub.Topic, err error) (ret error) {
	defer func() {
		if recover() == nil {
			ret = errors.New("want Topic.ErrorAs to panic when passed nil")
		}
	}()
	t.ErrorAs(err, nil)
	return nil
}

func (verifyAsFailsOnNil) SubscriptionErrorCheck(s *pubsub.Subscription, err error) (ret error) {
	defer func() {
		if recover() == nil {
			ret = errors.New("want Subscription.ErrorAs to panic when passed nil")
		}
	}()
	s.ErrorAs(err, nil)
	return nil
}

func (verifyAsFailsOnNil) MessageCheck(m *pubsub.Message) error {
	if m.As(nil) {
		return errors.New("want Message.As to return false when passed nil")
	}
	return nil
}

// RunConformanceTests runs conformance tests for provider implementations of pubsub.
func RunConformanceTests(t *testing.T, newHarness HarnessMaker, asTests []AsTest, opts *Options) {
	if opts == nil {
		opts = &Options{}
	}
	tests := map[string]func(t *testing.T, newHarness HarnessMaker){
		"TestSendReceive":                                         testSendReceive,
		"TestSendReceiveTwo":                                      testSendReceiveTwo,
		"TestErrorOnSendToClosedTopic":                            testErrorOnSendToClosedTopic,
		"TestErrorOnReceiveFromClosedSubscription":                testErrorOnReceiveFromClosedSubscription,
		"TestCancelSendReceive":                                   testCancelSendReceive,
		"TestNonExistentTopicSucceedsOnOpenButFailsOnSend":        testNonExistentTopicSucceedsOnOpenButFailsOnSend,
		"TestNonExistentSubscriptionSucceedsOnOpenButFailsOnSend": testNonExistentSubscriptionSucceedsOnOpenButFailsOnSend,
		"TestMetadata":                                            testMetadata,
		"TestNonUTF8MessageBody":                                  testNonUTF8MessageBody,
	}
	if opts.SkipTestsOfNonExistentThings {
		for name := range tests {
			if strings.Contains(name, "NonExistent") {
				delete(tests, name)
			}
		}
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) { test(t, newHarness) })
	}

	asTests = append(asTests, verifyAsFailsOnNil{})
	t.Run("TestAs", func(t *testing.T) {
		for _, st := range asTests {
			if st.Name() == "" {
				t.Fatalf("AsTest.Name is required")
			}
			t.Run(st.Name(), func(t *testing.T) { testAs(t, newHarness, st) })
		}
	})
}

// RunBenchmarks runs benchmarks for provider implementations of pubsub.
func RunBenchmarks(b *testing.B, topic *pubsub.Topic, sub *pubsub.Subscription) {
	b.Run("BenchmarkReceive", func(b *testing.B) {
		benchmark(b, topic, sub, false)
	})
	b.Run("BenchmarkSend", func(b *testing.B) {
		benchmark(b, topic, sub, true)
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
		t.Fatalf("failed to make non-existent subscription: %v", err)
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
	top, sub, cleanup, err := makePair(ctx, h, t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	want := publishN(ctx, t, top, 3)
	got := receiveN(ctx, t, sub, len(want))

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

	dt, cleanup, err := h.CreateTopic(ctx, t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	top := pubsub.NewTopic(dt)
	defer top.Shutdown(ctx)

	var ss []*pubsub.Subscription
	for i := 0; i < 2; i++ {
		ds, cleanup, err := h.CreateSubscription(ctx, dt, t.Name())
		if err != nil {
			t.Fatal(err)
		}
		defer cleanup()
		s := pubsub.NewSubscription(ds, nil)
		defer s.Shutdown(ctx)
		ss = append(ss, s)
	}
	want := publishN(ctx, t, top, 3)
	for i, s := range ss {
		got := receiveN(ctx, t, s, len(want))
		if diff := diffMessageSets(got, want); diff != "" {
			t.Errorf("sub #%d: %s", i, diff)
		}
	}
}

// Publish n different messages to the topic. Return the messages.
func publishN(ctx context.Context, t *testing.T, top *pubsub.Topic, n int) []*pubsub.Message {
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
func receiveN(ctx context.Context, t *testing.T, sub *pubsub.Subscription, n int) []*pubsub.Message {
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
	top, _, cleanup, err := makePair(ctx, h, t.Name())
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
	_, sub, cleanup, err := makePair(ctx, h, t.Name())
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
	top, sub, cleanup, err := makePair(ctx, h, t.Name())
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

func testMetadata(t *testing.T, newHarness HarnessMaker) {
	// Set up.
	ctx := context.Background()
	h, err := newHarness(ctx, t)
	if err != nil {
		t.Fatal(err)
	}
	defer h.Close()

	weirdMetadata := map[string]string{}
	for _, k := range escape.WeirdStrings {
		weirdMetadata[k] = k
	}

	top, sub, cleanup, err := makePair(ctx, h, t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	m := &pubsub.Message{
		Body:     []byte("hello world"),
		Metadata: weirdMetadata,
	}
	if err := top.Send(ctx, m); err != nil {
		t.Fatal(err)
	}

	m, err = sub.Receive(ctx)
	if err != nil {
		t.Fatal(err)
	}
	m.Ack()

	if diff := cmp.Diff(m.Metadata, weirdMetadata); diff != "" {
		t.Fatalf("got\n%v\nwant\n%v\ndiff\n%s", m.Metadata, weirdMetadata, diff)
	}

	// Verify that non-UTF8 strings in metadata key or value fail.
	m = &pubsub.Message{
		Body:     []byte("hello world"),
		Metadata: map[string]string{escape.NonUTF8String: "bar"},
	}
	if err := top.Send(ctx, m); err == nil {
		t.Error("got nil error, expected error for using non-UTF8 string as metadata key")
	}
	m.Metadata = map[string]string{"foo": escape.NonUTF8String}
	if err := top.Send(ctx, m); err == nil {
		t.Error("got nil error, expected error for using non-UTF8 string as metadata value")
	}
}

func testNonUTF8MessageBody(t *testing.T, newHarness HarnessMaker) {
	// Set up.
	ctx := context.Background()
	h, err := newHarness(ctx, t)
	if err != nil {
		t.Fatal(err)
	}
	defer h.Close()

	top, sub, cleanup, err := makePair(ctx, h, t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	// Sort the WeirdStrings map for record/replay consistency.
	var weirdStrings [][]string // [0] = key, [1] = value
	for k, v := range escape.WeirdStrings {
		weirdStrings = append(weirdStrings, []string{k, v})
	}
	sort.Slice(weirdStrings, func(i, j int) bool { return weirdStrings[i][0] < weirdStrings[j][0] })

	// Construct a message body with the weird strings and some non-UTF-8 bytes.
	var body []byte
	for _, v := range weirdStrings {
		body = append(body, []byte(v[1])...)
	}
	body = append(body, []byte(escape.NonUTF8String)...)
	m := &pubsub.Message{Body: body}

	if err := top.Send(ctx, m); err != nil {
		t.Fatal(err)
	}
	m, err = sub.Receive(ctx)
	if err != nil {
		t.Fatal(err)
	}
	m.Ack()

	if diff := cmp.Diff(m.Body, body); diff != "" {
		t.Fatalf("got\n%v\nwant\n%v\ndiff\n%s", m.Body, body, diff)
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

func makePair(ctx context.Context, h Harness, testName string) (*pubsub.Topic, *pubsub.Subscription, func(), error) {
	dt, topicCleanup, err := h.CreateTopic(ctx, testName)
	if err != nil {
		return nil, nil, nil, err
	}
	ds, subCleanup, err := h.CreateSubscription(ctx, dt, testName)
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
	top, sub, cleanup, err := makePair(ctx, h, t.Name())
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
	dt, err := h.MakeNonexistentTopic(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if err := top.Send(ctx, &pubsub.Message{Body: []byte("x")}); err != nil {
		t.Fatal(err)
	}
	m, err := sub.Receive(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.MessageCheck(m); err != nil {
		t.Error(err)
	}

	top = pubsub.NewTopic(dt)
	defer top.Shutdown(ctx)
	topicErr := top.Send(ctx, &pubsub.Message{})
	if topicErr == nil {
		t.Error("got nil expected error sending to nonexistent topic")
	} else if err := st.TopicErrorCheck(top, topicErr); err != nil {
		t.Error(err)
	}

	ds, err := h.MakeNonexistentSubscription(ctx)
	if err != nil {
		t.Fatal(err)
	}
	sub = pubsub.NewSubscription(ds, nil)
	defer sub.Shutdown(ctx)
	_, subErr := sub.Receive(ctx)
	if subErr == nil {
		t.Error("got nil expected error sending to nonexistent subscription")
	} else if err := st.SubscriptionErrorCheck(sub, subErr); err != nil {
		t.Error(err)
	}
}

// Publishes a large number of messages to topic concurrently, and then times
// how long it takes to send (if timeSend is true) or receive (if timeSend
// is false) them all.
func benchmark(b *testing.B, topic *pubsub.Topic, sub *pubsub.Subscription, timeSend bool) {
	attrs := map[string]string{"label": "value"}
	body := []byte("hello, world")
	const (
		nMessages          = 1000
		concurrencySend    = 10
		concurrencyReceive = 10
	)
	if nMessages%concurrencySend != 0 || nMessages%concurrencyReceive != 0 {
		b.Fatal("nMessages must be divisible by # of sending/receiving goroutines")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !timeSend {
			b.StopTimer()
		}
		if err := publishNConcurrently(topic, nMessages, concurrencySend, attrs, body); err != nil {
			b.Fatalf("publishing: %v", err)
		}
		b.Logf("published %d messages", nMessages)
		if timeSend {
			b.StopTimer()
		} else {
			b.StartTimer()
		}
		if err := receiveNConcurrently(sub, nMessages, concurrencyReceive); err != nil {
			b.Fatalf("receiving: %v", err)
		}
		b.SetBytes(nMessages * 1e6)
		b.Log("MB/s is actually number of messages received per second")
		if timeSend {
			b.StartTimer()
		}
	}
}

func publishNConcurrently(topic *pubsub.Topic, nMessages, nGoroutines int, attrs map[string]string, body []byte) error {
	return runConcurrently(nMessages, nGoroutines, func(ctx context.Context) error {
		return topic.Send(ctx, &pubsub.Message{Metadata: attrs, Body: body})
	})
}

func receiveNConcurrently(sub *pubsub.Subscription, nMessages, nGoroutines int) error {
	return runConcurrently(nMessages, nGoroutines, func(ctx context.Context) error {
		m, err := sub.Receive(ctx)
		if err != nil {
			return err
		}
		m.Ack()
		return nil
	})
}

// Call function f n times concurrently, using g goroutines. g must divide n.
// Wait until all calls complete. If any fail, cancel the remaining ones.
func runConcurrently(n, g int, f func(context.Context) error) error {
	gr, ctx := errgroup.WithContext(context.Background())
	ng := n / g
	for i := 0; i < g; i++ {
		gr.Go(func() error {
			for j := 0; j < ng; j++ {
				if err := f(ctx); err != nil {
					return err
				}
			}
			return nil
		})
	}
	return gr.Wait()

}
