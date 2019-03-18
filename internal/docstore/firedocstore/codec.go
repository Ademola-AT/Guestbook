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

// Package firedocstore implements the docstore API for Cloud Firestore.
//
// Docstore types not supported by Firestore:
// - complex64/128: encoded as an array of two float32/64s.
//
// Docstore types handled specially by Firestore:
// - time.Time: encoded as a ts.Timestamp [ts = "github.com/golang/protobuf/ptypes/timestamp"]
// - *ts.Timestamp: encoded as itself
// - *latlng.LatLng: encoded as itself [latlng = "google.golang.org/genproto/googleapis/type/latlng"]
//
// Firestore types not supported by Docstore:
// - Document reference (a pointer to another Firestore document)
// TODO(jba): figure out how to support this
package firedocstore

// Encoding and decoding between supported docstore types and Firestore protos.

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/golang/protobuf/ptypes"
	ts "github.com/golang/protobuf/ptypes/timestamp"
	"gocloud.dev/internal/docstore/driver"
	pb "google.golang.org/genproto/googleapis/firestore/v1"
	"google.golang.org/genproto/googleapis/type/latlng"
)

// encodeDoc encodes a driver.Document into Firestore's representation.
// A Firestore document (*pb.Document) is just a Go map from strings to *pb.Values.
func encodeDoc(doc driver.Document) (*pb.Document, error) {
	var e encoder
	if err := doc.Encode(&e); err != nil {
		return nil, err
	}
	return &pb.Document{Fields: e.pv.GetMapValue().Fields}, nil
}

// encodeValue encodes a Go value as a Firestore Value.
// The Firestore proto definition for Value is a oneof of various types,
// including basic types like string as well as lists and maps.
func encodeValue(x interface{}) (*pb.Value, error) {
	var e encoder
	if err := driver.Encode(reflect.ValueOf(x), &e); err != nil {
		return nil, err
	}
	return e.pv, nil
}

// encoder implements driver.Encoder. Its job is to encode a single Firestore value.
type encoder struct {
	pv *pb.Value
}

var nullValue = &pb.Value{ValueType: &pb.Value_NullValue{}}

func (e *encoder) EncodeNil()            { e.pv = nullValue }
func (e *encoder) EncodeBool(x bool)     { e.pv = &pb.Value{ValueType: &pb.Value_BooleanValue{x}} }
func (e *encoder) EncodeInt(x int64)     { e.pv = &pb.Value{ValueType: &pb.Value_IntegerValue{x}} }
func (e *encoder) EncodeUint(x uint64)   { e.pv = &pb.Value{ValueType: &pb.Value_IntegerValue{int64(x)}} }
func (e *encoder) EncodeBytes(x []byte)  { e.pv = &pb.Value{ValueType: &pb.Value_BytesValue{x}} }
func (e *encoder) EncodeFloat(x float64) { e.pv = floatval(x) }
func (e *encoder) EncodeString(x string) { e.pv = &pb.Value{ValueType: &pb.Value_StringValue{x}} }

func (e *encoder) ListIndex(int) { panic("impossible") }
func (e *encoder) MapKey(string) { panic("impossible") }

func (e *encoder) EncodeComplex(x complex128) {
	vals := []*pb.Value{floatval(real(x)), floatval(imag(x))}
	e.pv = &pb.Value{ValueType: &pb.Value_ArrayValue{&pb.ArrayValue{Values: vals}}}
}

func (e *encoder) EncodeList(n int) driver.Encoder {
	s := make([]*pb.Value, n)
	e.pv = &pb.Value{ValueType: &pb.Value_ArrayValue{&pb.ArrayValue{Values: s}}}
	return &listEncoder{s: s}
}

func (e *encoder) EncodeMap(n int) driver.Encoder {
	m := make(map[string]*pb.Value, n)
	e.pv = &pb.Value{ValueType: &pb.Value_MapValue{&pb.MapValue{Fields: m}}}
	return &mapEncoder{m: m}
}

var (
	typeOfGoTime         = reflect.TypeOf(time.Time{})
	typeOfProtoTimestamp = reflect.TypeOf((*ts.Timestamp)(nil))
	typeOfLatLng         = reflect.TypeOf((*latlng.LatLng)(nil))
)

// Encode time.Time, latlng.LatLng, and ts.Timestamp specially, because the Go Firestore
// client does.
func (e *encoder) EncodeSpecial(v reflect.Value) (bool, error) {
	switch v.Type() {
	case typeOfGoTime:
		ts, err := ptypes.TimestampProto(v.Interface().(time.Time))
		if err != nil {
			return false, err
		}
		e.pv = &pb.Value{ValueType: &pb.Value_TimestampValue{ts}}
		return true, nil
	case typeOfProtoTimestamp:
		if v.IsNil() {
			e.pv = nullValue
		} else {
			e.pv = &pb.Value{ValueType: &pb.Value_TimestampValue{v.Interface().(*ts.Timestamp)}}
		}
		return true, nil
	case typeOfLatLng:
		if v.IsNil() {
			e.pv = nullValue
		} else {
			e.pv = &pb.Value{ValueType: &pb.Value_GeoPointValue{v.Interface().(*latlng.LatLng)}}
		}
		return true, nil
	default:
		return false, nil
	}
}

type listEncoder struct {
	s []*pb.Value
	encoder
}

func (e *listEncoder) ListIndex(i int) { e.s[i] = e.pv }

type mapEncoder struct {
	m map[string]*pb.Value
	encoder
}

func (e *mapEncoder) MapKey(k string) { e.m[k] = e.pv }

func floatval(x float64) *pb.Value { return &pb.Value{ValueType: &pb.Value_DoubleValue{x}} }

////////////////////////////////////////////////////////////////

// decodeDoc decodes a Firestore document into a driver.Document.
func decodeDoc(pdoc *pb.Document, ddoc driver.Document) error {
	mv := &pb.Value{ValueType: &pb.Value_MapValue{&pb.MapValue{Fields: pdoc.Fields}}}
	return ddoc.Decode(decoder{mv})
}

type decoder struct {
	pv *pb.Value
}

func (d decoder) String() string { // for debugging
	return fmt.Sprint(d.pv)
}

func (d decoder) AsNull() bool {
	_, ok := d.pv.ValueType.(*pb.Value_NullValue)
	return ok
}

func (d decoder) AsBool() (bool, bool) {
	if b, ok := d.pv.ValueType.(*pb.Value_BooleanValue); ok {
		return b.BooleanValue, true
	}
	return false, false
}

func (d decoder) AsString() (string, bool) {
	if s, ok := d.pv.ValueType.(*pb.Value_StringValue); ok {
		return s.StringValue, true
	}
	return "", false
}

func (d decoder) AsInt() (int64, bool) {
	if i, ok := d.pv.ValueType.(*pb.Value_IntegerValue); ok {
		return i.IntegerValue, true
	}
	return 0, false
}

func (d decoder) AsUint() (uint64, bool) {
	if i, ok := d.pv.ValueType.(*pb.Value_IntegerValue); ok {
		return uint64(i.IntegerValue), true
	}
	return 0, false
}

func (d decoder) AsFloat() (float64, bool) {
	if f, ok := d.pv.ValueType.(*pb.Value_DoubleValue); ok {
		return f.DoubleValue, true
	}
	return 0, false
}

func (d decoder) AsBytes() ([]byte, bool) {
	if bs, ok := d.pv.ValueType.(*pb.Value_BytesValue); ok {
		return bs.BytesValue, true
	}
	return nil, false
}

func (d decoder) AsComplex() (complex128, bool) {
	a := d.pv.GetArrayValue()
	if a == nil {
		return 0, false
	}
	vs := a.Values
	if len(vs) != 2 {
		return 0, false
	}
	r, okr := vs[0].ValueType.(*pb.Value_DoubleValue)
	i, oki := vs[1].ValueType.(*pb.Value_DoubleValue)
	if !okr || !oki {
		return 0, false
	}
	return complex(r.DoubleValue, i.DoubleValue), true
}

// AsInterface decodes the value in d into the most appropriate Go type.
func (d decoder) AsInterface() (interface{}, error) {
	return decodeValue(d.pv)
}

func decodeValue(v *pb.Value) (interface{}, error) {
	switch v := v.ValueType.(type) {
	case *pb.Value_NullValue:
		return nil, nil
	case *pb.Value_BooleanValue:
		return v.BooleanValue, nil
	case *pb.Value_IntegerValue:
		return v.IntegerValue, nil
	case *pb.Value_DoubleValue:
		return v.DoubleValue, nil
	case *pb.Value_StringValue:
		return v.StringValue, nil
	case *pb.Value_BytesValue:
		return v.BytesValue, nil
	case *pb.Value_TimestampValue:
		// Return TimestampValue as time.Time.
		t, err := ptypes.Timestamp(v.TimestampValue)
		if err != nil {
			return nil, err
		}
		return t, nil
	case *pb.Value_ReferenceValue:
		// TODO(jba): support references
		return nil, errors.New("references are not currently supported")
	case *pb.Value_GeoPointValue:
		// Return GeoPointValue as *latlng.LatLng.
		return v.GeoPointValue, nil
	case *pb.Value_ArrayValue:
		s := make([]interface{}, len(v.ArrayValue.Values))
		for i, pv := range v.ArrayValue.Values {
			e, err := decodeValue(pv)
			if err != nil {
				return nil, err
			}
			s[i] = e
		}
		return s, nil
	case *pb.Value_MapValue:
		m := make(map[string]interface{}, len(v.MapValue.Fields))
		for k, pv := range v.MapValue.Fields {
			e, err := decodeValue(pv)
			if err != nil {
				return nil, err
			}
			m[k] = e
		}
		return m, nil
	}
	return nil, fmt.Errorf("unknown firestore value type %T", v)
}

func (d decoder) ListLen() (int, bool) {
	a := d.pv.GetArrayValue()
	if a == nil {
		return 0, false
	}
	return len(a.Values), true
}

func (d decoder) DecodeList(f func(int, driver.Decoder) bool) {
	for i, e := range d.pv.GetArrayValue().Values {
		if !f(i, decoder{e}) {
			return
		}
	}
}
func (d decoder) MapLen() (int, bool) {
	m := d.pv.GetMapValue()
	if m == nil {
		return 0, false
	}
	return len(m.Fields), true
}
func (d decoder) DecodeMap(f func(string, driver.Decoder) bool) {
	for k, v := range d.pv.GetMapValue().Fields {
		if !f(k, decoder{v}) {
			return
		}
	}
}

func (d decoder) AsSpecial(v reflect.Value) (bool, interface{}, error) {
	switch v.Type() {
	case typeOfGoTime:
		if ts, ok := d.pv.ValueType.(*pb.Value_TimestampValue); ok {
			if ts.TimestampValue == nil {
				return true, time.Time{}, nil
			}
			t, err := ptypes.Timestamp(ts.TimestampValue)
			return true, t, err
		}
		return true, nil, fmt.Errorf("expected TimestampValue for time.Time, got %+v", d.pv.ValueType)
	case typeOfProtoTimestamp:
		if ts, ok := d.pv.ValueType.(*pb.Value_TimestampValue); ok {
			return true, ts.TimestampValue, nil
		}
		return true, nil, fmt.Errorf("expected TimestampValue for *ts.Timestamp, got %+v", d.pv.ValueType)

	case typeOfLatLng:
		if ll, ok := d.pv.ValueType.(*pb.Value_GeoPointValue); ok {
			return true, ll.GeoPointValue, nil
		}
		return true, nil, fmt.Errorf("expected GeoPointValue for *latlng.LatLng, got %+v", d.pv.ValueType)

	default:
		return false, nil, nil
	}
}
