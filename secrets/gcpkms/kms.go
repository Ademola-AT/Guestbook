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

// Package gcpkms provides a secrets implementation backed by Google Cloud KMS.
// Use NewKeeper to construct a *secrets.Keeper.
//
// As
//
// gcpkms exposes the following type for As:
//  - Error: *google.golang.org/grpc/status.Status
package gcpkms // import "gocloud.dev/secrets/gcpkms"

import (
	"context"
	"fmt"

	cloudkms "cloud.google.com/go/kms/apiv1"
	"gocloud.dev/gcerrors"
	"gocloud.dev/gcp"
	"gocloud.dev/internal/gcerr"
	"gocloud.dev/internal/useragent"
	"gocloud.dev/secrets"
	"google.golang.org/api/option"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
	"google.golang.org/grpc/status"
)

// endPoint is the address to access Google Cloud KMS API.
const endPoint = "cloudkms.googleapis.com:443"

// Dial returns a client to use with Cloud KMS and a clean-up function to close
// the client after used.
func Dial(ctx context.Context, ts gcp.TokenSource) (*cloudkms.KeyManagementClient, func(), error) {
	c, err := cloudkms.NewKeyManagementClient(ctx, option.WithTokenSource(ts), useragent.ClientOption("secrets"))
	return c, func() { c.Close() }, err
}

// NewKeeper returns a *secrets.Keeper that uses Google Cloud KMS.
// See the package documentation for an example.
func NewKeeper(client *cloudkms.KeyManagementClient, ki *KeyID, opts *KeeperOptions) *secrets.Keeper {
	return secrets.NewKeeper(&keeper{
		keyID:  ki,
		client: client,
	})
}

// KeyID includes related information to construct a key name that is managed
// by Cloud KMS.
// See https://cloud.google.com/kms/docs/object-hierarchy#key for more
// information.
type KeyID struct {
	ProjectID, Location, KeyRing, Key string
}

func (ki *KeyID) String() string {
	return fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
		ki.ProjectID, ki.Location, ki.KeyRing, ki.Key)
}

// keeper contains information to construct the pull path of a key.
type keeper struct {
	keyID  *KeyID
	client *cloudkms.KeyManagementClient
}

// Decrypt decrypts the ciphertext using the key constructed from ki.
func (k *keeper) Decrypt(ctx context.Context, ciphertext []byte) ([]byte, error) {
	req := &kmspb.DecryptRequest{
		Name:       k.keyID.String(),
		Ciphertext: ciphertext,
	}
	resp, err := k.client.Decrypt(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.GetPlaintext(), nil
}

// Encrypt encrypts the plaintext into a ciphertext.
func (k *keeper) Encrypt(ctx context.Context, plaintext []byte) ([]byte, error) {
	req := &kmspb.EncryptRequest{
		Name:      k.keyID.String(),
		Plaintext: plaintext,
	}
	resp, err := k.client.Encrypt(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.GetCiphertext(), nil
}

// ErrorAs implements driver.Keeper.ErrorAs.
func (k *keeper) ErrorAs(err error, i interface{}) bool {
	s, ok := status.FromError(err)
	if !ok {
		return false
	}
	p, ok := i.(**status.Status)
	if !ok {
		return false
	}
	*p = s
	return true
}

// ErrorCode implements driver.ErrorCode.
func (k *keeper) ErrorCode(err error) gcerrors.ErrorCode {
	return gcerr.GRPCCode(err)
}

// KeeperOptions controls Keeper behaviors.
// It is provided for future extensibility.
type KeeperOptions struct{}
