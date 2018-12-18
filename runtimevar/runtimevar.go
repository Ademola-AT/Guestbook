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
// limtations under the License.

// Package runtimevar provides an interface for reading runtime variables and
// ability to detect changes and get updates on those variables.
package runtimevar // import "gocloud.dev/runtimevar"

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"reflect"
	"time"

	"gocloud.dev/runtimevar/driver"
)

// Snapshot contains a variable and metadata about it.
type Snapshot struct {
	// Value is an object containing a runtime variable  The type of
	// this object is set by the driver and it should always be the same type for the same Variable
	// object. A driver implementation can provide the ability to configure the object type and a
	// decoding scheme where variables are stored as bytes in the backend service.  Clients
	// should not mutate this object as it can be accessed by multiple goroutines.
	Value interface{}

	// UpdateTime is the time when the last changed was detected.
	UpdateTime time.Time
}

// Variable provides the ability to read runtime variables with its blocking Watch method.
type Variable struct {
	watcher  driver.Watcher
	nextCall time.Time
	prev     driver.State
}

// New constructs a Variable object given a driver.Watcher implementation.
func New(w driver.Watcher) *Variable {
	return &Variable{watcher: w}
}

// Watch blocks until there are variable changes, the Context's Done channel
// closes or a new error occurs.
//
// If the variable changes, the method returns a Snapshot object containing the
// updated value.
//
// If method returns an error, the returned Snapshot object is a zero value and cannot be used.
//
// The first call to this method should return the current variable unless there are errors in
// retrieving the value.
//
// Users should not call this method from multiple goroutines as implementations may not guarantee
// safety in data access.  It is typical to use only one goroutine calling this method in a loop.
//
// To stop this function from blocking, caller can passed in Context object constructed via
// WithCancel and call the cancel function.
func (c *Variable) Watch(ctx context.Context) (Snapshot, error) {
	for {
		wait := c.nextCall.Sub(time.Now())
		if wait > 0 {
			select {
			case <-ctx.Done():
				return Snapshot{}, ctx.Err()
			case <-time.After(wait):
				// Continue.
			}
		}

		cur, wait := c.watcher.WatchVariable(ctx, c.prev)
		c.nextCall = time.Now().Add(wait)
		if cur == nil {
			// No change.
			continue
		}
		// Something new to return!
		c.prev = cur
		v, err := cur.Value()
		if err != nil {
			return Snapshot{}, wrapError(err)
		}
		return Snapshot{Value: v, UpdateTime: cur.UpdateTime()}, nil
	}
}

// Close cleans up any resources used by the Variable object.
func (c *Variable) Close() error {
	err := c.watcher.Close()
	return wrapError(err)
}

// wrappedError is used to wrap all errors returned by drivers so that users
// are not given access to provider-specific errors.
type wrappedError struct {
	err error
}

func wrapError(err error) error {
	if err == nil {
		return nil
	}
	return &wrappedError{err: err}
}

func (w *wrappedError) Error() string {
	return "runtimevar: " + w.err.Error()
}

// Decode is a function type for unmarshaling/decoding a slice of bytes into
// an arbitrary type. Decode functions are used when creating a Decoder via
// NewDecoder. This package provides common Decode functions including
// GobDecode and JSONDecode.
type Decode func([]byte, interface{}) error

// Decoder decodes a slice of bytes into a particular Go object.
//
// This package provides some common Decoders that you can use directly,
// including StringDecoder and BytesDecoder. You can also NewDecoder to
// construct other Decoders.
type Decoder struct {
	typ reflect.Type
	fn Decode
}

// NewDecoder returns a Decoder that uses fn to decode a slice of bytes into
// an object of type obj.
//
// This package provides some common Decode functions, including JSONDecode
// and GobDecode, which can be passed to this function to create Decoders for
// JSON and gob values.
func NewDecoder(obj interface{}, fn Decode) *Decoder {
	return &Decoder{
		typ: reflect.TypeOf(obj),
		fn:  fn,
	}
}

// Decode decodes b into a new instance of the target type.
func (d *Decoder) Decode(b []byte) (interface{}, error) {
	nv := reflect.New(d.typ).Interface()
	if err := d.fn(b, nv); err != nil {
		return nil, err
	}
	ptr := reflect.ValueOf(nv)
	return ptr.Elem().Interface(), nil
}

var (
	// StringDecoder decodes into strings.
	StringDecoder = NewDecoder("", stringDecode)

	// BytesDecoder copies the slice of bytes.
	BytesDecoder = NewDecoder([]byte{}, bytesDecode)

	// JSONDecode can be passed to NewDecoder when decoding JSON (https://golang.org/pkg/encoding/json/).
	JSONDecode = json.Unmarshal
)

// GobDecode can be passed to NewDecoder when decoding gobs (https://golang.org/pkg/encoding/gob/).
func GobDecode(data []byte, obj interface{}) error {
	return gob.NewDecoder(bytes.NewBuffer(data)).Decode(obj)
}

func stringDecode(b []byte, obj interface{}) error {
	v := obj.(*string)
	*v = string(b)
	return nil
}

func bytesDecode(b []byte, obj interface{}) error {
	v := obj.(*[]byte)
	*v = b[:]
	return nil
}
