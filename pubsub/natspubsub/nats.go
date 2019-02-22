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

// Package natspubsub provides a pubsub implementation for NATS.io.
// Use OpenTopic to construct a *pubsub.Topic, and/or OpenSubscription
// to construct a *pubsub.Subscription.
//
// As
//
// natspubsub exposes the following types for As:
//  - Topic: *nats.Conn
//  - Subscription: *nats.Subscription
//  - Message: *nats.Msg

package natspubsub // import "gocloud.dev/pubsub/natspubsub"

import (
	"context"
	"errors"
	"reflect"

	"github.com/nats-io/go-nats"
	"github.com/ugorji/go/codec"
	"gocloud.dev/gcerrors"
	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/driver"
)

var errNotExist = errors.New("natspubsub: topic does not exist")

type topic struct {
	nc   *nats.Conn
	subj string
}

// For encoding we use msgpack from github.com/ugorji/go.
// it was already imported in go-cloud and is reasonably performant.
var mh codec.MsgpackHandle

func init() {
	// driver.Message.Metadata type
	dm := driver.Message{}
	mh.MapType = reflect.TypeOf(dm.Metadata)
	mh.ExplicitRelease = true
}

// We define our own version of message here for encoding that
// only encodes Body and Metadata. Otherwise we would have to
// add codec decorations to driver.Message.
type encMsg struct {
	Body     []byte            `codec:",omitempty"`
	Metadata map[string]string `codec:",omitempty"`
}

// OpenTopic returns a *pubsub.Topic for use with NATS.
// We delay checking for the proper syntax here.
func OpenTopic(nc *nats.Conn, topicName string) *pubsub.Topic {
	return pubsub.NewTopic(openTopic(nc, topicName))
}

// openTopic returns the driver for OpenTopic. This function exists so the test
// harness can get the driver interface implementation if it needs to.
func openTopic(nc *nats.Conn, topicName string) driver.Topic {
	return &topic{nc, topicName}
}

// SendBatch implements driver.Topic.SendBatch.
func (t *topic) SendBatch(ctx context.Context, msgs []*driver.Message) error {
	if t == nil || t.nc == nil {
		return errNotExist
	}

	// Reuse if possible.
	var em encMsg
	var raw [1024]byte
	b := raw[:0]
	enc := codec.NewEncoderBytes(&b, &mh)
	defer enc.Release()

	for _, m := range msgs {
		if err := ctx.Err(); err != nil {
			return err
		}
		enc.ResetBytes(&b)
		em.Body, em.Metadata = m.Body, m.Metadata
		if err := enc.Encode(em); err != nil {
			return err
		}
		if err := t.nc.Publish(t.subj, b); err != nil {
			return err
		}
	}
	// Per specification this is supposed to only return after
	// a message has been sent. Normally NATS is very efficient
	// at sending messages in batches on its on and also handles
	// disconnected buffering during a reconnect event. We will
	// let NATS handle this for now. If needed we could add a
	// FlushWithContext() call which ensures the connected server
	// has processed all the messages.
	return nil
}

// IsRetryable implements driver.Topic.IsRetryable.
func (*topic) IsRetryable(error) bool { return false }

// As implements driver.Topic.As.
func (t *topic) As(i interface{}) bool {
	c, ok := i.(**nats.Conn)
	if !ok {
		return false
	}
	*c = t.nc
	return true
}

// ErrorAs implements driver.Topic.ErrorAs
func (*topic) ErrorAs(error, interface{}) bool {
	return false
}

// ErrorCode implements driver.Topic.ErrorCode
func (*topic) ErrorCode(error) gcerrors.ErrorCode { return gcerrors.Unknown }

type subscription struct {
	nc   *nats.Conn
	nsub *nats.Subscription
	oerr error
}

// OpenSubscription returns a *pubsub.Subscription representing a NATS subscription.
// TODO(dlc) - Options for queue groups?
func OpenSubscription(nc *nats.Conn, subscriptionName string) *pubsub.Subscription {
	return pubsub.NewSubscription(openSubscription(nc, subscriptionName), nil)
}

func openSubscription(nc *nats.Conn, subscriptionName string) driver.Subscription {
	sub, err := nc.SubscribeSync(subscriptionName)
	return &subscription{nc, sub, err}
}

// ReceiveBatch implements driver.ReceiveBatch.
func (s *subscription) ReceiveBatch(ctx context.Context, maxMessages int) ([]*driver.Message, error) {
	if s == nil || s.nsub == nil {
		return nil, nats.ErrBadSubscription
	}

	// Return right away if the ctx has an error.
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	// Just in case
	if maxMessages == 0 {
		maxMessages = 1
	}

	var ms []*driver.Message

	// We will assume the desired goal is at least one message since the public API only has Receive().
	// We will load all the messages that are already queued up first.
	for {
		msg, err := s.nsub.NextMsg(0)
		if err == nil {
			dm, err := decode(msg)
			if err != nil {
				return nil, err
			}
			ms = append(ms, dm)
			if len(ms) >= maxMessages {
				break
			}
		} else if err == nats.ErrTimeout {
			break
		} else {
			return nil, err
		}
	}

	// If we have anything go ahead and return here.
	if len(ms) > 0 {
		return ms, nil
	}

	// Here we will assume the ctx has a deadline and will allow it to do its thing. We may
	// get a burst of messages but we will only wait once here and let the next call grab them.
	// The reason is so that if deadline is not properly set we will not needlessly wait here
	// for more messages when the user most likely only wants one.
	msg, err := s.nsub.NextMsgWithContext(ctx)
	if err != nil {
		return nil, err
	}
	dm, err := decode(msg)
	if err != nil {
		return nil, err
	}
	ms = append(ms, dm)
	return ms, nil
}

// Convert NATS msgs to *driver.Message.
func decode(msg *nats.Msg) (*driver.Message, error) {
	if msg == nil {
		return nil, nats.ErrInvalidMsg
	}
	var dm driver.Message
	// Everything is in the msg.Data
	dec := codec.NewDecoderBytes(msg.Data, &mh)
	defer dec.Release()

	dec.Decode(&dm)
	dm.AckID = -1 // Not applicable to NATS
	dm.AsFunc = messageAsFunc(msg)

	return &dm, nil
}

func messageAsFunc(msg *nats.Msg) func(interface{}) bool {
	return func(i interface{}) bool {
		p, ok := i.(**nats.Msg)
		if !ok {
			return false
		}
		*p = msg
		return true
	}
}

// SendAcks implements driver.Subscription.SendAcks. NATS does not need Acks since
// it is At-Most-Once QoS.
func (s *subscription) SendAcks(ctx context.Context, ids []driver.AckID) error {
	return nil
}

// IsRetryable implements driver.Subscription.IsRetryable.
func (s *subscription) IsRetryable(error) bool { return false }

// As implements driver.Subscription.As.
func (s *subscription) As(i interface{}) bool {
	c, ok := i.(**nats.Subscription)
	if !ok {
		return false
	}
	*c = s.nsub
	return true
}

// ErrorAs implements driver.Subscription.ErrorAs
func (*subscription) ErrorAs(error, interface{}) bool {
	return false
}

// ErrorCode implements driver.Subscription.ErrorCode
func (*subscription) ErrorCode(error) gcerrors.ErrorCode { return gcerrors.Unknown }
