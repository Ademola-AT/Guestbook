// Copyright 2019 The Go Cloud Authors
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

// Package errors provides an error type for Go Cloud APIs.
package errors

import "fmt"

// An ErrorCode describes the error's category. Programs should act upon an error's
// code, not its message.
type ErrorCode int

const (
	// Returned by the Code function on a nil error. It is not a valid
	// code for an error.
	OK ErrorCode = 0

	// The error could not be categorized.
	Unknown = 1

	// The resource was not found.
	NotFound = 2

	// The resource exists, but it should not.
	AlreadyExists = 3

	// A value given to a Go Cloud API is incorrect.
	InvalidArgument = 4

	// Something unexpected happened. Internal errors always indicate
	// bugs in Go Cloud (or possibly the underlying provider).
	Internal = 5
)

// To get stringer: go get golang.org/x/tools/cmd/stringer

//go:generate stringer -type=ErrorCode

// An Error described a Go Cloud error.
type Error struct {
	Code ErrorCode
	msg  string
	err  error
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v: %s", e.Code, e.msg)
}

// Unwrap returns the error underlying the receiver.
func (e *Error) Unwrap() error {
	return e.err
}

// Code returns the ErrorCode of err if it is an *Error.
// It returns Unknown if err is a non-nil error of a different type.
// If err is nil, it returns the special code OK.
func Code(err error) ErrorCode {
	if err == nil {
		return OK
	}
	if e, ok := err.(*Error); ok {
		return e.Code
	}
	return Unknown
}

// New returns a new error with the given code, underlying error and message.
func New(c ErrorCode, err error, msg string) *Error {
	return &Error{
		Code: c,
		msg:  msg,
		err:  err,
	}
}

// Newf uses format and args to format a message, then calls New.
func Newf(c ErrorCode, err error, format string, args ...interface{}) *Error {
	return New(c, err, fmt.Sprintf(format, args...))
}
