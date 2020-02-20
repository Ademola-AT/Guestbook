package mqttpubsub

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"gocloud.dev/gcerrors"
	"gocloud.dev/pubsub/driver"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	defaultQOS byte = 0
	pubID           = "publisher"
	subID           = "subscriber"
)

// TODO init sub and pub

var (
	errNoURLEnv         = errors.New("mqttpubsub: no url env provided ")
	errInvalidMessage   = errors.New("mqttpubsub: invalid or empty message")
	errConnRequired     = errors.New("mqttpubsub: mqtt connection is required")
	errStillConnected   = errors.New("mqttpubsub: still connected. Kill all processes manually")
	errMQTTDisconnected = errors.New("mqttpubsub: disconnected")
)

type (
	MQTTMessenger interface {
		Subscriber
		Publisher
	}

	messenger interface {
		Subscriber
		Publisher
	}

	Subscriber interface {
		Subscribe(topic string, handler mqtt.MessageHandler) error
		UnSubscribe(topic string) error
		Close() error
	}

	Publisher interface {
		Publish(topic string, payload interface{}) error
		Stop() error
	}

	subscriber struct {
		subConnect mqtt.Client
	}

	publisher struct {
		pubConnect mqtt.Client

		isStopped bool
		wg        *sync.WaitGroup
	}
)

func defaultSubClient(url string) (_ Subscriber, err error) {
	if url == "" {
		return nil, errNoURLEnv
	}
	var subConnect mqtt.Client

	subConnect, err = makeConnect(subID, url)
	if err != nil {
		return nil, err
	}
	return &subscriber{
		subConnect: subConnect,
	}, nil
}

func defaultPubClient(url string) (_ Publisher, err error) {
	if url == "" {
		return nil, errNoURLEnv
	}
	var pubConnect mqtt.Client

	pubConnect, err = makeConnect(pubID, url)
	if err != nil {
		return nil, err
	}
	return &publisher{
		pubConnect: pubConnect,
		wg:         new(sync.WaitGroup),
	}, nil
}

func makeConnect(clientID, url string) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions()
	opts = opts.AddBroker(url)
	opts.ClientID = clientID
	mqttClient := mqtt.NewClient(opts)
	token := mqttClient.Connect()
	token.Wait()
	if token.Error() != nil {
		return nil, token.Error()
	}
	if !mqttClient.IsConnectionOpen() {
		return nil, errMQTTDisconnected
	}

	return mqttClient, nil
}

func (c *publisher) Publish(topic string, payload interface{}) error {
	if c.isStopped {
		return nil
	}

	token := c.pubConnect.Publish(topic, defaultQOS, false, payload)
	if token.Wait() && token.Error() != nil {
		fmt.Println(topic, token.Error(),"@@@@@@@@@@@")
		return token.Error()
	}
	return nil
}

func (c *publisher) Stop() error {
	if c.pubConnect == nil {
		return errConnRequired
	}
	c.isStopped = true
	c.wg.Wait()
	c.pubConnect.Disconnect(0)
	if c.pubConnect.IsConnected() {
		return errStillConnected
	}
	fmt.Println("STOPPED")

	return nil
}

func (c *subscriber) Subscribe(topic string, handler mqtt.MessageHandler) error {
	if !c.subConnect.IsConnected() {
		return errMQTTDisconnected
	}

	token := c.subConnect.Subscribe(topic, defaultQOS, handler)

	if token.Wait() && token.Error() != nil {
		fmt.Println(topic, token.Error(),"############")

		return token.Error()
	}
	return nil
}

func (c *subscriber) UnSubscribe(topic string) error {
	if !c.subConnect.IsConnected() {
		return errMQTTDisconnected
	}

	token := c.subConnect.Unsubscribe(topic)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (c *subscriber) Close() error {
	fmt.Println("CLOSING")

	if c.subConnect == nil {
		return errConnRequired
	}

	c.subConnect.Disconnect(0)
	if c.subConnect.IsConnected() {
		return errStillConnected
	}
	fmt.Println("CLOSED")

	return nil
}

// Convert MQTT msgs to *driver.Message.
func decode(msg mqtt.Message) (*driver.Message, error) {
	if msg == nil {
		return nil, errInvalidMessage
	}
	var dm driver.Message
	if err := decodeMessage(msg.Payload(), &dm); err != nil {
		return nil, err
	}
	dm.AckID = msg.MessageID() // uint16
	dm.AsFunc = messageAsFunc(msg)
	return &dm, nil
}

func messageAsFunc(msg mqtt.Message) func(interface{}) bool {
	return func(i interface{}) bool {
		p, ok := i.(*mqtt.Message)
		if !ok {
			return false
		}
		*p = msg
		return true
	}
}

func encodeMessage(dm *driver.Message) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if len(dm.Metadata) == 0 {
		return dm.Body, nil
	}
	if err := enc.Encode(dm.Metadata); err != nil {
		return nil, err
	}
	if err := enc.Encode(dm.Body); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decodeMessage(data []byte, dm *driver.Message) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&dm.Metadata); err != nil {
		// This may indicate a normal NATS message, so just treat as the body.
		dm.Metadata = nil
		dm.Body = data
		return nil
	}
	return dec.Decode(&dm.Body)
}

func whichError(err error) gcerrors.ErrorCode {
	switch err {
	case nil:
		return gcerrors.OK
	case context.Canceled:
		return gcerrors.Canceled
	case errMQTTDisconnected, errConnRequired:
		return gcerrors.NotFound
	case mqtt.ErrInvalidTopicEmptyString, mqtt.ErrInvalidQos, mqtt.ErrInvalidTopicMultilevel, errInvalidMessage:
		return gcerrors.FailedPrecondition
	case errStillConnected:
		return gcerrors.Internal
	}
	return gcerrors.Unknown
}
