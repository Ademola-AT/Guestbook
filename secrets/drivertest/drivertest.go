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

// Package drivertest provides a conformance test for implementations of
// the secrets driver.
package drivertest // import "gocloud.dev/secrets/drivertest"

import (
	"bytes"
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gocloud.dev/secrets"
	"gocloud.dev/secrets/driver"
)

// Harness descibes the functionality test harnesses must provide to run
// conformance tests.
type Harness interface {
	// MakeDriver returns a pair of driver.Keeper, each backed by a different key.
	MakeDriver(ctx context.Context) (driver.Keeper, driver.Keeper, error)

	// Close is called when the test is complete.
	Close()
}

// HarnessMaker describes functions that construct a harness for running tests.
// It is called exactly once per test.
type HarnessMaker func(ctx context.Context, t *testing.T) (Harness, error)

// RunConformanceTests runs conformance tests for provider implementations of secret management.
func RunConformanceTests(t *testing.T, newHarness HarnessMaker) {
	t.Run("TestEncryptDecrypt", func(t *testing.T) {
		testEncryptDecrypt(t, newHarness)
	})
	t.Run("TestMultipleEncryptionsNotEqual", func(t *testing.T) {
		testMultipleEncryptionsNotEqual(t, newHarness)
	})
	t.Run("TestMultipleKeys", func(t *testing.T) {
		testMultipleKeys(t, newHarness)
	})
	t.Run("TestDecryptMalformedError", func(t *testing.T) {
		testDecryptMalformedError(t, newHarness)
	})
}

// testEncryptDecrypt tests the functionality of encryption and decryption
func testEncryptDecrypt(t *testing.T, newHarness HarnessMaker) {
	ctx := context.Background()
	harness, err := newHarness(ctx, t)
	if err != nil {
		t.Fatal(err)
	}
	defer harness.Close()

	drv, _, err := harness.MakeDriver(ctx)
	if err != nil {
		t.Fatal(err)
	}
	keeper := secrets.NewKeeper(drv)

	msg := []byte("I'm a secret message!")
	encryptedMsg, err := keeper.Encrypt(ctx, msg)
	if err != nil {
		t.Fatal(err)
	}
	if cmp.Equal(msg, encryptedMsg) {
		t.Errorf("Got encrypted message %v, want it to differ from original message %v", string(msg), string(encryptedMsg))
	}
	decryptedMsg, err := keeper.Decrypt(ctx, encryptedMsg)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(msg, decryptedMsg) {
		t.Errorf("Got decrypted message %v, want it to match original message %v", string(msg), string(decryptedMsg))
	}

}

// testMultipleEncryptionsNotEqual tests that encrypting a plaintext multiple
// times with the same key works, and that the encrypted bytes are different.
func testMultipleEncryptionsNotEqual(t *testing.T, newHarness HarnessMaker) {
	ctx := context.Background()
	harness, err := newHarness(ctx, t)
	if err != nil {
		t.Fatal(err)
	}
	defer harness.Close()

	drv, _, err := harness.MakeDriver(ctx)
	if err != nil {
		t.Fatal(err)
	}
	keeper := secrets.NewKeeper(drv)

	msg := []byte("I'm a secret message!")
	encryptedMsg1, err := keeper.Encrypt(ctx, msg)
	if err != nil {
		t.Fatal(err)
	}
	encryptedMsg2, err := keeper.Encrypt(ctx, msg)
	if err != nil {
		t.Fatal(err)
	}
	if cmp.Equal(encryptedMsg1, encryptedMsg2) {
		t.Errorf("Got same encrypted messages from multiple encryptions %v, want them to be different", string(encryptedMsg1))
	}
	decryptedMsg, err := keeper.Decrypt(ctx, encryptedMsg1)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(decryptedMsg, msg) {
		t.Errorf("got decrypted %q want %q", string(decryptedMsg), string(msg))
	}
	decryptedMsg, err = keeper.Decrypt(ctx, encryptedMsg2)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(decryptedMsg, msg) {
		t.Errorf("got decrypted %q want %q", string(decryptedMsg), string(msg))
	}
}

// testMultipleKeys tests that encrypting the same text with different
// keys works, and that the encrypted bytes are different.
func testMultipleKeys(t *testing.T, newHarness HarnessMaker) {
	ctx := context.Background()
	harness, err := newHarness(ctx, t)
	if err != nil {
		t.Fatal(err)
	}
	defer harness.Close()

	drv1, drv2, err := harness.MakeDriver(ctx)
	if err != nil {
		t.Fatal(err)
	}
	keeper1 := secrets.NewKeeper(drv1)
	keeper2 := secrets.NewKeeper(drv2)

	msg := []byte("I'm a secret message!")
	encryptedMsg1, err := keeper1.Encrypt(ctx, msg)
	if err != nil {
		t.Fatal(err)
	}
	encryptedMsg2, err := keeper2.Encrypt(ctx, msg)
	if err != nil {
		t.Fatal(err)
	}
	if cmp.Equal(encryptedMsg1, encryptedMsg2) {
		t.Errorf("Got same encrypted messages from multiple encryptions %v, want them to be different", string(encryptedMsg1))
	}

	// We cannot assert that decrypting encryptedMsg1 with keeper2 fails,
	// or that decrypting encryptedMsg2 with keeper1 fails, as Decrypt is allowed
	// to decrypt using a different key than the one given to Keeper.

	decryptedMsg, err := keeper1.Decrypt(ctx, encryptedMsg1)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(decryptedMsg, msg) {
		t.Errorf("got decrypted %q want %q", string(decryptedMsg), string(msg))
	}

	decryptedMsg, err = keeper2.Decrypt(ctx, encryptedMsg2)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(decryptedMsg, msg) {
		t.Errorf("got decrypted %q want %q", string(decryptedMsg), string(msg))
	}
}

// testDecryptMalformedError tests decryption returns an error when the
// ciphertext is malformed.
func testDecryptMalformedError(t *testing.T, newHarness HarnessMaker) {
	ctx := context.Background()
	harness, err := newHarness(ctx, t)
	if err != nil {
		t.Fatal(err)
	}
	defer harness.Close()

	drv, _, err := harness.MakeDriver(ctx)
	if err != nil {
		t.Fatal(err)
	}
	keeper := secrets.NewKeeper(drv)

	if _, err := keeper.Decrypt(ctx, []byte("malformed cipher message")); err == nil {
		t.Error("Got nil, want decrypt error")
	}
}
