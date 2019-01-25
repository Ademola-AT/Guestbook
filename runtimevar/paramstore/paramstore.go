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

// Package paramstore provides a runtimevar implementation with variables
// read from AWS Systems Manager Parameter Store
// (https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-paramstore.html)
// Use NewVariable to construct a *runtimevar.Variable.
//
// As
//
// paramstore exposes the following types for As:
//  - Snapshot: *ssm.GetParameterOutput, *ssm.DescribeParametersOutput
//  - Error: awserr.Error
package paramstore // import "gocloud.dev/runtimevar/paramstore"

import (
	"context"
	"fmt"
	"time"

	"gocloud.dev/runtimevar"
	"gocloud.dev/runtimevar/driver"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// Options sets options.
type Options struct {
	// WaitDuration controls the rate at which Parameter Store is polled.
	// Defaults to 30 seconds.
	WaitDuration time.Duration
}

// NewVariable constructs a *runtimevar.Variable backed by the variable name in
// AWS Systems Manager Parameter Store.
// Parameter Store returns raw bytes; provide a decoder to decode the raw bytes
// into the appropriate type for runtimevar.Snapshot.Value.
// See the runtimevar package documentation for examples of decoders.
func NewVariable(sess client.ConfigProvider, name string, decoder *runtimevar.Decoder, opts *Options) (*runtimevar.Variable, error) {
	return runtimevar.New(newWatcher(sess, name, decoder, opts)), nil
}

func newWatcher(sess client.ConfigProvider, name string, decoder *runtimevar.Decoder, opts *Options) *watcher {
	if opts == nil {
		opts = &Options{}
	}
	return &watcher{
		sess:    sess,
		name:    name,
		wait:    driver.WaitDuration(opts.WaitDuration),
		decoder: decoder,
	}
}

// state implements driver.State.
type state struct {
	val        interface{}
	rawGet     *ssm.GetParameterOutput
	rawDesc    *ssm.DescribeParametersOutput
	updateTime time.Time
	version    int64
	err        error
}

// Value implements driver.State.Value.
func (s *state) Value() (interface{}, error) {
	return s.val, s.err
}

// UpdateTime implements driver.State.UpdateTime.
func (s *state) UpdateTime() time.Time {
	return s.updateTime
}

// As implements driver.State.As.
func (s *state) As(i interface{}) bool {
	if s.rawGet == nil {
		return false
	}
	switch p := i.(type) {
	case **ssm.GetParameterOutput:
		*p = s.rawGet
	case **ssm.DescribeParametersOutput:
		*p = s.rawDesc
	default:
		return false
	}
	return true
}

// errorState returns a new State with err, unless prevS also represents
// the same error, in which case it returns nil.
func errorState(err error, prevS driver.State) driver.State {
	s := &state{err: err}
	if prevS == nil {
		return s
	}
	prev := prevS.(*state)
	if prev.err == nil {
		// New error.
		return s
	}
	if equivalentError(err, prev.err) {
		// Same error, return nil to indicate no change.
		return nil
	}
	return s
}

// equivalentError returns true iff err1 and err2 represent an equivalent error;
// i.e., we don't want to return it to the user as a different error.
func equivalentError(err1, err2 error) bool {
	if err1 == err2 || err1.Error() == err2.Error() {
		return true
	}
	var code1, code2 string
	if awsErr, ok := err1.(awserr.Error); ok {
		code1 = awsErr.Code()
	}
	if awsErr, ok := err2.(awserr.Error); ok {
		code2 = awsErr.Code()
	}
	return code1 != "" && code1 == code2
}

type watcher struct {
	// sess is the AWS session to use to talk to AWS.
	sess client.ConfigProvider
	// name is the parameter to retrieve.
	name string
	// wait is the amount of time to wait between querying AWS.
	wait time.Duration
	// decoder is the decoder that unmarshals the value in the param.
	decoder *runtimevar.Decoder
}

// WatchVariable implements driver.WatchVariable.
func (w *watcher) WatchVariable(ctx context.Context, prev driver.State) (driver.State, time.Duration) {
	lastVersion := int64(-1)
	if prev != nil {
		lastVersion = prev.(*state).version
	}
	// GetParameter from S3 to get the current value and version.
	svc := ssm.New(w.sess)
	getResp, err := svc.GetParameter(&ssm.GetParameterInput{Name: aws.String(w.name)})
	if err != nil {
		return errorState(err, prev), w.wait
	}
	if getResp.Parameter == nil {
		return errorState(fmt.Errorf("unable to get %q parameter", w.name), prev), w.wait
	}
	getP := getResp.Parameter
	if *getP.Version == lastVersion {
		// Version hasn't changed, so no change; return nil.
		return nil, w.wait
	}

	// DescribeParameters from S3 to get the LastModified date.
	descResp, err := svc.DescribeParameters(&ssm.DescribeParametersInput{
		Filters: []*ssm.ParametersFilter{
			{Key: aws.String("Name"), Values: []*string{&w.name}},
		},
	})
	if err != nil {
		return errorState(err, prev), w.wait
	}
	if len(descResp.Parameters) != 1 || *descResp.Parameters[0].Name != w.name {
		return errorState(fmt.Errorf("unable to get single %q parameter", w.name), prev), w.wait
	}
	descP := descResp.Parameters[0]

	// New value (or at least, new version). Decode it.
	val, err := w.decoder.Decode([]byte(*getP.Value))
	if err != nil {
		return errorState(err, prev), w.wait
	}
	return &state{val: val, rawGet: getResp, rawDesc: descResp, updateTime: *descP.LastModifiedDate, version: *getP.Version}, w.wait
}

// Close implements driver.Close.
func (w *watcher) Close() error {
	return nil
}

// ErrorAs implements driver.ErrorAs.
func (w *watcher) ErrorAs(err error, i interface{}) bool {
	switch v := err.(type) {
	case awserr.Error:
		if p, ok := i.(*awserr.Error); ok {
			*p = v
			return true
		}
	}
	return false
}

// IsNotExist implements driver.IsNotExist.
func (*watcher) IsNotExist(err error) bool {
	if awsErr, ok := err.(awserr.Error); ok {
		return awsErr.Code() == "ParameterNotFound"
	}
	return false
}
