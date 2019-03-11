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

package dynamodocstore

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	dyn "github.com/aws/aws-sdk-go/service/dynamodb"
	"gocloud.dev/internal/docstore/driver"
)

var nullValue = new(dyn.AttributeValue).SetNULL(true)

type encoder struct {
	av *dyn.AttributeValue
}

func (e *encoder) EncodeNil()            { e.av = nullValue }
func (e *encoder) EncodeBool(x bool)     { e.av = new(dyn.AttributeValue).SetBOOL(x) }
func (e *encoder) EncodeInt(x int64)     { e.av = new(dyn.AttributeValue).SetN(strconv.FormatInt(x, 10)) }
func (e *encoder) EncodeUint(x uint64)   { e.av = new(dyn.AttributeValue).SetN(strconv.FormatUint(x, 10)) }
func (e *encoder) EncodeBytes(x []byte)  { e.av = new(dyn.AttributeValue).SetB(x) }
func (e *encoder) EncodeFloat(x float64) { e.av = encodeFloat(x) }

func (e *encoder) ListIndex(int) { panic("impossible") }
func (e *encoder) MapKey(string) { panic("impossible") }

func (e *encoder) EncodeString(x string) {
	if len(x) == 0 {
		e.av = nullValue
	} else {
		e.av = new(dyn.AttributeValue).SetS(x)
	}
}

func (e *encoder) EncodeComplex(x complex128) {
	e.av = new(dyn.AttributeValue).SetL([]*dyn.AttributeValue{encodeFloat(real(x)), encodeFloat(imag(x))})
}

func (e *encoder) EncodeList(n int) driver.Encoder {
	s := make([]*dyn.AttributeValue, n)
	e.av = new(dyn.AttributeValue).SetL(s)
	return &listEncoder{s: s}
}

func (e *encoder) EncodeMap(n int) driver.Encoder {
	m := make(map[string]*dyn.AttributeValue, n)
	e.av = new(dyn.AttributeValue).SetM(m)
	return &mapEncoder{m: m}
}

func (*encoder) EncodeSpecial(reflect.Value) (bool, error) {
	return false, nil
}

type listEncoder struct {
	s []*dyn.AttributeValue
	encoder
}

func (e *listEncoder) ListIndex(i int) { e.s[i] = e.av }

type mapEncoder struct {
	m map[string]*dyn.AttributeValue
	encoder
}

func (e *mapEncoder) MapKey(k string) { e.m[k] = e.av }

func encodeDoc(doc driver.Document) (map[string]*dyn.AttributeValue, error) {
	var e encoder
	if err := doc.Encode(&e); err != nil {
		return nil, err
	}
	return e.av.M, nil
}

// TODO: improve driver.Encoder to make this efficient.
func encodeDocKeyFields(doc driver.Document, pkey, skey string) (map[string]*dyn.AttributeValue, error) {
	m, err := encodeDoc(doc)
	if err != nil {
		return nil, err
	}
	hasP := false
	hasS := false
	if skey == "" {
		hasS = true
	}
	for f := range m {
		switch f {
		case pkey:
			hasP = true
		case skey:
			hasS = true
		default:
			delete(m, f)
		}
	}
	if !hasP || !hasS {
		return nil, errors.New("missing key field(s)")
	}
	return m, nil
}

func encodeValue(v interface{}) (*dyn.AttributeValue, error) {
	var e encoder
	if err := driver.Encode(reflect.ValueOf(v), &e); err != nil {
		return nil, err
	}
	return e.av, nil
}

func encodeFloat(f float64) *dyn.AttributeValue {
	return new(dyn.AttributeValue).SetN(strconv.FormatFloat(f, 'f', -1, 64))
}

////////////////////////////////////////////////////////////////

func decodeDoc(doc driver.Document, item map[string]*dyn.AttributeValue) error {
	return doc.Decode(decoder{&dyn.AttributeValue{M: item}})
}

type decoder struct {
	av *dyn.AttributeValue
}

func (d decoder) String() string {
	return d.av.String()
}

func (d decoder) AsBool() (bool, bool) {
	if d.av.BOOL == nil {
		return false, false
	}
	return *d.av.BOOL, true
}

func (d decoder) AsNull() bool {
	return d.av.NULL != nil
}

func (d decoder) AsString() (string, bool) {
	// Empty string is represented by NULL.
	if d.av.NULL != nil {
		return "", true
	}
	if d.av.S == nil {
		return "", false
	}
	return *d.av.S, true
}

func (d decoder) AsInt() (int64, bool) {
	if d.av.N == nil {
		return 0, false
	}
	i, err := strconv.ParseInt(*d.av.N, 10, 64)
	if err != nil {
		return 0, false
	}
	return i, true
}

func (d decoder) AsUint() (uint64, bool) {
	if d.av.N == nil {
		return 0, false
	}
	u, err := strconv.ParseUint(*d.av.N, 10, 64)
	if err != nil {
		return 0, false
	}
	return u, true
}

func (d decoder) AsFloat() (float64, bool) {
	if d.av.N == nil {
		return 0, false
	}
	f, err := strconv.ParseFloat(*d.av.N, 64)
	if err != nil {
		return 0, false
	}
	return f, true

}

func (d decoder) AsComplex() (complex128, bool) {
	if d.av.L == nil {
		return 0, false
	}
	if len(d.av.L) != 2 {
		return 0, false
	}
	r, ok := decoder{d.av.L[0]}.AsFloat()
	if !ok {
		return 0, false
	}
	i, ok := decoder{d.av.L[1]}.AsFloat()
	if !ok {
		return 0, false
	}
	return complex(r, i), true
}

func (d decoder) AsBytes() ([]byte, bool) {
	if d.av.B == nil {
		return nil, false
	}
	return d.av.B, true
}

func (d decoder) ListLen() (int, bool) {
	if d.av.L == nil {
		return 0, false
	}
	return len(d.av.L), true
}

func (d decoder) DecodeList(f func(i int, vd driver.Decoder) bool) {
	for i, el := range d.av.L {
		if !f(i, decoder{el}) {
			break
		}
	}
}

func (d decoder) MapLen() (int, bool) {
	if d.av.M == nil {
		return 0, false
	}
	return len(d.av.M), true
}

func (d decoder) DecodeMap(f func(key string, vd driver.Decoder) bool) {
	for k, av := range d.av.M {
		if !f(k, decoder{av}) {
			break
		}
	}
}

func (d decoder) AsInterface() (interface{}, error) {
	return toGoValue(d.av)
}

func toGoValue(av *dyn.AttributeValue) (interface{}, error) {
	switch {
	case av.NULL != nil:
		return nil, nil
	case av.BOOL != nil:
		return *av.BOOL, nil
	case av.N != nil:
		f, err := strconv.ParseFloat(*av.N, 64)
		if err != nil {
			return nil, err
		}
		i := int64(f)
		if float64(i) == f {
			return i, nil
		}
		u := uint64(f)
		if float64(u) == f {
			return u, nil
		}
		return f, nil

	case av.B != nil:
		return av.B, nil
	case av.S != nil:
		return *av.S, nil

	case av.L != nil:
		s := make([]interface{}, len(av.L))
		for i, v := range av.L {
			x, err := toGoValue(v)
			if err != nil {
				return nil, err
			}
			s[i] = x
		}
		return s, nil

	case av.M != nil:
		m := make(map[string]interface{}, len(av.M))
		for k, v := range av.M {
			x, err := toGoValue(v)
			if err != nil {
				return nil, err
			}
			m[k] = x
		}
		return m, nil

	default:
		return nil, fmt.Errorf("dynamodocstore: AttributeValue %s not supported", av)
	}
}

func (decoder) AsSpecial(reflect.Value) (bool, interface{}, error) {
	return false, nil, nil
}
