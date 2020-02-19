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

package mqttpubsub

import (
	"context"
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/hashicorp/go-multierror"
	"net/url"
	"os"
	"path"
	"sync"
	"time"

	"gocloud.dev/gcerrors"
	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/driver"
)

func init() {
	o := new(defaultDialer)
	pubsub.DefaultURLMux().RegisterTopic(Scheme, o)
	pubsub.DefaultURLMux().RegisterSubscription(Scheme, o)
}

// Scheme is the URL scheme mqttpubsub registers its URLOpeners under on pubsub.DefaultMux.
const Scheme = "mqtt"

// URLOpener opens MQTT URLs like "mqtt://myexchange" for
// topics or "mqtt://myqueue" for subscriptions.
//
// For topics, the URL's host+path is used as the exchange name.
//
// For subscriptions, the URL's host+path is used as the queue name.
//
// No query parameters are supported.
type URLOpener struct {
	// Connection to use for communication with the server.
	Connection MQTTMessenger

	// TopicOptions specifies the options to pass to OpenTopic.
	TopicOptions TopicOptions
	// SubscriptionOptions specifies the options to pass to OpenSubscription.
	SubscriptionOptions SubscriptionOptions
}

// defaultDialer dials a default MQTT server based on the environment
// variable "MQTT_SERVER_URL".
type defaultDialer struct {
	init   sync.Once
	opener *URLOpener
	err    error
}

func (o *defaultDialer) defaultConn(ctx context.Context) (*URLOpener, error) {
	o.init.Do(func() {
		serverURL := os.Getenv("MQTT_SERVER_URL")
		if serverURL == "" {
			o.err = errors.New("MQTT_SERVER_URL environment variable not set")
			return
		}
		_, cancelF := context.WithTimeout(ctx, time.Second*3)
		defer cancelF()

		conn, err := defaultConn(serverURL)
		if err != nil {
			o.err = err
			return
		}
		o.opener = &URLOpener{
			Connection: conn,
		}
	})
	return o.opener, o.err
}

func (o *defaultDialer) OpenTopicURL(ctx context.Context, u *url.URL) (*pubsub.Topic, error) {
	opener, err := o.defaultConn(ctx)
	if err != nil {
		return nil, fmt.Errorf("open topic %v: failed to open default connection: %v", u, err)
	}
	return opener.OpenTopicURL(ctx, u)
}

func (o *defaultDialer) OpenSubscriptionURL(ctx context.Context, u *url.URL) (*pubsub.Subscription, error) {
	opener, err := o.defaultConn(ctx)
	if err != nil {
		return nil, fmt.Errorf("open subscription %v: failed to open default connection: %v", u, err)
	}
	return opener.OpenSubscriptionURL(ctx, u)
}

// OpenTopicURL opens a pubsub.Topic based on u.
func (o *URLOpener) OpenTopicURL(ctx context.Context, u *url.URL) (*pubsub.Topic, error) {
	for param := range u.Query() {
		return nil, fmt.Errorf("open topic %v: invalid query parameter %q", u, param)
	}
	exchangeName := path.Join(u.Host, u.Path)
	return OpenTopic(o.Connection.GetPublisher(), exchangeName, &o.TopicOptions)
}

// OpenSubscriptionURL opens a pubsub.Subscription based on u.
func (o *URLOpener) OpenSubscriptionURL(ctx context.Context, u *url.URL) (*pubsub.Subscription, error) {
	for param := range u.Query() {
		return nil, fmt.Errorf("open subscription %v: invalid query parameter %q", u, param)
	}
	queueName := path.Join(u.Host, u.Path)
	return OpenSubscription(o.Connection.GetSubscriber(), queueName, &o.SubscriptionOptions)
}

type topic struct {
	name string
	conn Publisher

	wg *sync.WaitGroup
	mu *sync.Mutex
	errs *multierror.Error
}

// TopicOptions sets options for constructing a *pubsub.Topic backed by
// MQTT.
type TopicOptions struct{}

// SubscriptionOptions sets options for constructing a *pubsub.Subscription
// backed by MQTT.
type SubscriptionOptions struct{}

func OpenTopic(conn Publisher, name string, _ *TopicOptions) (*pubsub.Topic, error) {
	if conn == nil {
		return nil, errConnRequired
	}
	dt := &topic{
		name,
		conn,
		new(sync.WaitGroup),
		new(sync.Mutex),
		new(multierror.Error),
	}
	return pubsub.NewTopic(dt, nil), nil
}

// SendBatch implements driver.Topic.SendBatch.
func (t *topic) SendBatch(ctx context.Context, msgs []*driver.Message) error {
	if t == nil || t.conn == nil {
		return errNotInitialized
	}

	for _, m := range msgs {
		t.wg.Add(1)
		if err := ctx.Err(); err != nil {
			return err
		}

		payload, err := encodeMessage(m)
		if err != nil {
			return err
		}
		if m.BeforeSend != nil {
			asFunc := func(i interface{}) bool { return false }
			if err := m.BeforeSend(asFunc); err != nil {
				return err
			}
		}
		go func() {
			defer t.wg.Done()
			t.mu.Lock()
			err = t.conn.Publish(t.name, payload)
			if err != nil {
				t.errs.Errors = append(t.errs.Errors, err)
			}
			t.mu.Unlock()
		}()
	}
	t.wg.Wait()
	
	return t.errs.ErrorOrNil()
}

// IsRetryable implements driver.Topic.IsRetryable.
func (*topic) IsRetryable(error) bool { return false }

// As implements driver.Topic.As.
func (t *topic) As(i interface{}) bool {
	c, ok := i.(*MQTTMessenger)
	if !ok {
		return false
	}
	*c = MQTTMessenger(&messenger{Publisher: t.conn})
	return true
}

// ErrorAs implements driver.Topic.ErrorAs
func (*topic) ErrorAs(error, interface{}) bool {
	return false
}

// ErrorCode implements driver.Topic.ErrorCode
func (*topic) ErrorCode(err error) gcerrors.ErrorCode {
	return whichError(err)

}

// Close implements driver.Topic.Close.
func (t *topic) Close() error {
	return t.conn.Stop()
}

type subscription struct {
	conn      Subscriber
	topicName string

	mu          *sync.RWMutex
	unackedMsgs []mqtt.Message

	stopChan chan struct{}

	errs *multierror.Error
}

func OpenSubscription(conn Subscriber, topicName string, _ *SubscriptionOptions) (*pubsub.Subscription, error) {
	ds := &subscription{
		conn,
		topicName,
		new(sync.RWMutex),
		make([]mqtt.Message, 0),
		make(chan struct{}),
		new(multierror.Error),
	}
	err := ds.conn.Subscribe(topicName, func(client mqtt.Client, m mqtt.Message) {
		ds.mu.Lock()
		ds.unackedMsgs = append(ds.unackedMsgs, m)
		ds.mu.Unlock()
	})
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			select {
			case <-time.Tick(5 * time.Minute): // unhandled messages cleaning
				ds.unackedMsgs = make([]mqtt.Message, 0)
			case <-ds.stopChan:
				return
			}
		}
	}()

	return pubsub.NewSubscription(ds, nil, nil), nil
}

// ReceiveBatch implements driver.ReceiveBatch.
func (s *subscription) ReceiveBatch(ctx context.Context, maxMessages int) (dms []*driver.Message, err error) {
	if s == nil || s.conn == nil {
		return nil, errConnRequired
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	s.mu.RLock()
	for i, m := range s.unackedMsgs {
		if i >= maxMessages {
			break
		}
		dm, err := decode(m)
		if err != nil {
			s.errs.Errors = append(s.errs.Errors, err)
		}
		dms = append(dms, dm)
	}
	s.mu.RUnlock()

	return dms, s.errs.ErrorOrNil()
}

// SendAcks implements driver.Subscription.SendAcks.
func (s *subscription) SendAcks(ctx context.Context, ids []driver.AckID) error {
	if len(ids) == 0 {
		return nil
	}

	for _, id := range ids {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		msgID, ok := id.(uint16)
		if !ok {
			continue
		}

		if len(s.unackedMsgs) == 0 {
			return nil
		}

		s.mu.RLock()
		for _, m := range s.unackedMsgs {
			if m.MessageID() == msgID {
				m.Ack()
			}
		}
		s.mu.RUnlock()
	}
	return nil
}

// CanNack implements driver.CanNack.
func (s *subscription) CanNack() bool { return true }

// SendNacks implements driver.Subscription.SendNacks. MQTT doesn't have implementation for NACK
// so below method just removes unACKed messages
func (s *subscription) SendNacks(ctx context.Context, ids []driver.AckID) error {
	if len(ids) == 0 {
		return nil
	}
	for _, id := range ids {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		msgID, ok := id.(uint16)
		if !ok {
			continue
		}

		if len(s.unackedMsgs) == 0 {
			return nil
		}

		s.mu.Lock()
		for i, m := range s.unackedMsgs {
			if m.MessageID() == msgID {
				s.unackedMsgs[i] = nil
			}
		}
		s.mu.Unlock()
	}
	return nil
}

// IsRetryable implements driver.Subscription.IsRetryable.
func (s *subscription) IsRetryable(error) bool { return false }

// As implements driver.Subscription.As.
func (s *subscription) As(i interface{}) bool {
	c, ok := i.(*MQTTMessenger)
	if !ok {
		return false
	}
	*c = MQTTMessenger(&messenger{Subscriber: s.conn})
	return true
}

// ErrorAs implements driver.Subscription.ErrorAs
func (*subscription) ErrorAs(error, interface{}) bool {
	return false
}

// ErrorCode implements driver.Subscription.ErrorCode
func (*subscription) ErrorCode(err error) gcerrors.ErrorCode {
	return whichError(err)
}

// Close implements driver.Subscription.Close.
func (s *subscription) Close() error {
	if _, ok := <-s.stopChan; ok {
		s.stopChan <- struct{}{}
		close(s.stopChan)
		return s.Close()
	}
	return nil
}
