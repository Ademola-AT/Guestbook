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

package localsecrets

import (
	"context"
	"errors"
	"log"
	"testing"

	"gocloud.dev/secrets"
	"gocloud.dev/secrets/driver"
	"gocloud.dev/secrets/drivertest"
)

type harness struct{}

func (h *harness) MakeDriver(ctx context.Context) (driver.Keeper, driver.Keeper, error) {
	secret1, err := NewRandomKey()
	if err != nil {
		log.Fatal(err)
	}
	secret2, err := NewRandomKey()
	if err != nil {
		log.Fatal(err)
	}
	return &keeper{secretKey: secret1}, &keeper{secretKey: secret2}, nil
}

func (h *harness) Close() {}

func newHarness(ctx context.Context, t *testing.T) (drivertest.Harness, error) {
	return &harness{}, nil
}

func TestConformance(t *testing.T) {
	drivertest.RunConformanceTests(t, newHarness, []drivertest.AsTest{verifyAs{}})
}

type verifyAs struct{}

func (v verifyAs) Name() string {
	return "verify As function"
}

func (v verifyAs) ErrorCheck(k *secrets.Keeper, err error) error {
	var s string
	if k.ErrorAs(err, &s) {
		return errors.New("Keeper.ErrorAs expected to fail")
	}
	return nil
}

func TestOpenKeeper(t *testing.T) {
	tests := []struct {
		URL     string
		WantErr bool
	}{
		// OK.
		{"stringkey://my-special-very-secret-secretkey", false},
		// Invalid parameter.
		{"stringkey://my-special-very-secret-secretkey?param=value", true},
		// OK.
		{"base64key://smGbjm71Nxd1Ig5FS0wj9SlbzAIrnolCz9bQQ6uAhl4=", false},
		// Invalid parameter.
		{"base64key://smGbjm71Nxd1Ig5FS0wj9SlbzAIrnolCz9bQQ6uAhl4=?param=value", true},
		// Invalid base64 key.
		{"base64key://not-valid-base64", true},
	}

	ctx := context.Background()
	for _, test := range tests {
		keeper, err := secrets.OpenKeeper(ctx, test.URL)
		if (err != nil) != test.WantErr {
			t.Errorf("%s: got error %v, want error %v", test.URL, err, test.WantErr)
		}
		if err == nil {
			if err = keeper.Close(); err != nil {
				t.Errorf("%s: got error during close: %v", test.URL, err)
			}
		}
	}
}
