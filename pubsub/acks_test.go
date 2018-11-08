package pubsub_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/golang/go/src/pkg/math/rand"
	"github.com/google/go-cloud/pubsub"
	"github.com/google/go-cloud/pubsub/driver"
)

type ackingDriverSub struct {
	q        []*driver.Message
	sendAcks func(context.Context, []driver.AckID) error
}

func (s *ackingDriverSub) ReceiveBatch(ctx context.Context) ([]*driver.Message, error) {
	ms := s.q
	s.q = nil
	return ms, nil
}

func (s *ackingDriverSub) SendAcks(ctx context.Context, ackIDs []driver.AckID) error {
	return s.sendAcks(ctx, ackIDs)
}

func (s *ackingDriverSub) Close() error {
	return nil
}

func TestAckTriggersDriverSendAcksForOneMessage(t *testing.T) {
	ctx := context.Background()
	var sentAcks []driver.AckID
	f := func(ctx context.Context, ackIDs []driver.AckID) error {
		sentAcks = ackIDs
		return nil
	}
	id := rand.Int()
	m := &driver.Message{AckID: id}
	ds := &ackingDriverSub{
		q:        []*driver.Message{m},
		sendAcks: f,
	}
	sub := pubsub.NewSubscription(ctx, ds, pubsub.SubscriptionOptions{})
	m2, err := sub.Receive(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if err := m2.Ack(ctx); err != nil {
		t.Fatal(err)
	}
	if len(sentAcks) != 1 {
		t.Fatalf("len(sentAcks) = %d, want exactly 1", len(sentAcks))
	}
	if sentAcks[0] != id {
		t.Errorf("sentAcks[0] = %d, want %d", sentAcks[0], id)
	}
}

func TestMultipleAcksCanGoIntoASingleBatch(t *testing.T) {
	ctx := context.Background()
	var sentAcks []driver.AckID
	f := func(ctx context.Context, ackIDs []driver.AckID) error {
		sentAcks = ackIDs
		return nil
	}
	ids := []int{rand.Int(), rand.Int()}
	ds := &ackingDriverSub{
		q:        []*driver.Message{{AckID: ids[0]}, {AckID: ids[1]}},
		sendAcks: f,
	}
	sopts := pubsub.SubscriptionOptions{
		AckBatchSize: 2,
	}
	sub := pubsub.NewSubscription(ctx, ds, sopts)

	// Receive and ack the messages concurrently.
	var wg sync.WaitGroup
	recv := func() {
		mr, err := sub.Receive(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if err := mr.Ack(ctx); err != nil {
			t.Fatal(err)
		}
		wg.Done()
	}
	wg.Add(2)
	go recv()
	go recv()
	wg.Wait()

	if len(sentAcks) != 2 {
		t.Fatalf("len(sentAcks) = %d, want exactly 2", len(sentAcks))
	}
	if sentAcks[0] != ids[0] {
		t.Errorf("sentAcks[0] = %d, want %d", sentAcks[0], ids[0])
	}
	if sentAcks[1] != ids[1] {
		t.Errorf("sentAcks[0] = %d, want %d", sentAcks[1], ids[1])
	}
}

func TestTooManyAcksForASingleBatchGoIntoMultipleBatches(t *testing.T) {

}

func TestMsgAckReturnsErrorFromSendAcks(t *testing.T) {
	ctx := context.Background()
	e := fmt.Sprintf("%d", rand.Int())
	f := func(ctx context.Context, ackIDs []driver.AckID) error {
		return errors.New(e)
	}
	m := &driver.Message{}
	ds := &ackingDriverSub{
		q:        []*driver.Message{m},
		sendAcks: f,
	}
	sub := pubsub.NewSubscription(ctx, ds, pubsub.SubscriptionOptions{})
	mr, err := sub.Receive(ctx)
	if err != nil {
		t.Fatal(err)
	}
	err = mr.Ack(ctx)
	if err == nil {
		t.Fatal("got nil, want error")
	}
	if err.Error() != e {
		t.Errorf("got error %q, want %q", err.Error(), e)
	}
}

func TestCancelAck(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	f := func(ctx context.Context, ackIDs []driver.AckID) error {
		// Hang.
		c := make(chan struct{})
		<-c
		return nil
	}
	m := &driver.Message{}
	ds := &ackingDriverSub{
		q:        []*driver.Message{m},
		sendAcks: f,
	}
	sub := pubsub.NewSubscription(ctx, ds, pubsub.SubscriptionOptions{})
	mr, err := sub.Receive(ctx)
	if err != nil {
		t.Fatal(err)
	}
	cancel()
	if err := mr.Ack(ctx); err == nil {
		t.Error("got nil, want cancellation error")
	}
}
