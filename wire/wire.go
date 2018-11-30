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

// Package wire is the deprecated package for Wire code generation directives.
// It exists only for backward compatibility. Please update your import paths to
// "github.com/google/wire".
//
// For an overview of working with Wire, see the user guide at
// https://github.com/google/wire/blob/master/README.md
//
// Deprecated: Please import "github.com/google/wire" instead of this package.
package wire

import "github.com/google/wire"

// ProviderSet is a marker type that collects a group of providers.
type ProviderSet = wire.ProviderSet

// NewSet creates a new provider set that includes the providers in its
// arguments. Each argument is a function value, a struct (zero) value, a
// provider set, a call to Bind, a call to Value, or a call to InterfaceValue.
//
// Passing a function value to NewSet declares that the function's first
// return value type will be provided by calling the function. The arguments
// to the function will come from the providers for their types. As such, all
// the parameters must be of non-identical types. The function may optionally
// return an error as its last return value and a cleanup function as the
// second return value. A cleanup function must be of type func() and is
// guaranteed to be called before the cleanup function of any of the
// provider's inputs. If any provider returns an error, the injector function
// will call all the appropriate cleanup functions and return the error from
// the injector function.
//
// Passing a struct value of type S to NewSet declares that both S and *S will
// be provided by creating a new value of the appropriate type by filling in
// each field of S using the provider of the field's type.
//
// Passing a ProviderSet to NewSet is the same as if the set's contents
// were passed as arguments to NewSet directly.
//
// The behavior of passing the result of a call to other functions in this
// package are described in their respective doc comments.
func NewSet(args ...interface{}) wire.ProviderSet {
	return wire.NewSet(args...)
}

// Build is placed in the body of an injector function template to declare the
// providers to use. The Wire code generation tool will fill in an
// implementation of the function. The arguments to Build are interpreted the
// same as NewSet: they determine the provider set presented to Wire's
// dependency graph. Build returns an error message that can be sent to a call
// to panic().
//
// The parameters of the injector function are used as inputs in the dependency
// graph.
//
// Similar to provider functions passed into NewSet, the first return value is
// the output of the injector function, the optional second return value is a
// cleanup function, and the optional last return value is an error. If any of
// the provider functions in the injector function's provider set return errors
// or cleanup functions, the corresponding return value must be present in the
// injector function template.
//
// Examples:
//
//	func injector(ctx context.Context) (*sql.DB, error) {
//		wire.Build(otherpkg.FooSet, myProviderFunc)
//		return nil, nil
//	}
//
//	func injector(ctx context.Context) (*sql.DB, error) {
//		panic(wire.Build(otherpkg.FooSet, myProviderFunc))
//	}
func Build(args ...interface{}) string {
	return wire.Build(args...)
}

// A Binding maps an interface to a concrete type.
type Binding = wire.Binding

// Bind declares that a concrete type should be used to satisfy a
// dependency on the type of iface, which must be a pointer to an
// interface type.
//
// Example:
//
//	type Fooer interface {
//		Foo()
//	}
//
//	type MyFoo struct{}
//
//	func (MyFoo) Foo() {}
//
//	var MySet = wire.NewSet(
//		MyFoo{},
//		wire.Bind(new(Fooer), new(MyFoo)))
func Bind(iface, to interface{}) wire.Binding {
	return wire.Bind(iface, to)
}

// A ProvidedValue is an expression that is copied to the generated injector.
type ProvidedValue = wire.ProvidedValue

// Value binds an expression to provide the type of the expression.
// The expression may not be an interface value; use InterfaceValue for that.
//
// Example:
//
//	var MySet = wire.NewSet(wire.Value([]string(nil)))
func Value(v interface{}) wire.ProvidedValue {
	return wire.Value(v)
}

// InterfaceValue binds an expression to provide a specific interface type.
//
// Example:
//
//	var MySet = wire.NewSet(wire.InterfaceValue(new(io.Reader), os.Stdin))
func InterfaceValue(typ interface{}, x interface{}) wire.ProvidedValue {
	return wire.InterfaceValue(typ, x)
}
