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

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/go-cloud/blob"
	"github.com/google/go-cloud/blob/gcsblob"
	"github.com/google/go-cloud/blob/s3blob"
	"github.com/google/go-cloud/gcp"

	azureblob "../../blob/azureblob"
)

// setupBucket creates a connection to a particular cloud provider's blob storage.
func setupBucket(ctx context.Context, cloud, bucket string) (*blob.Bucket, error) {
	switch cloud {
	case "aws":
		return setupAWS(ctx, bucket)
	case "gcp":
		return setupGCP(ctx, bucket)
	case "azure": // uses github.com/Azure/azure-storage-blob-go/2018-03-28/azblob
		return setupAzure(ctx, bucket)
	default:
		return nil, fmt.Errorf("invalid cloud provider: %s", cloud)
	}
}

// setupGCP creates a connection to Google Cloud Storage (GCS).
func setupGCP(ctx context.Context, bucket string) (*blob.Bucket, error) {
	// DefaultCredentials assumes a user has logged in with gcloud.
	// See here for more information:
	// https://cloud.google.com/docs/authentication/getting-started
	creds, err := gcp.DefaultCredentials(ctx)
	if err != nil {
		return nil, err
	}
	c, err := gcp.NewHTTPClient(gcp.DefaultTransport(), gcp.CredentialsTokenSource(creds))
	if err != nil {
		return nil, err
	}
	return gcsblob.OpenBucket(ctx, bucket, c)
}

// setupAWS creates a connection to Simple Cloud Storage Service (S3).
func setupAWS(ctx context.Context, bucket string) (*blob.Bucket, error) {
	c := &aws.Config{
		// Either hard-code the region or use AWS_REGION.
		Region: aws.String("us-east-2"),
		// credentials.NewEnvCredentials assumes two environment variables are
		// present:
		// 1. AWS_ACCESS_KEY_ID, and
		// 2. AWS_SECRET_ACCESS_KEY.
		Credentials: credentials.NewEnvCredentials(),
	}
	s := session.Must(session.NewSession(c))
	return s3blob.OpenBucket(ctx, s, bucket)
}

func setupAzure(ctx context.Context, bucket string) (*blob.Bucket, error) {
	settings := azureblob.Settings{
		AccountName:      os.Getenv("ACCOUNT_NAME"),
		AccountKey:       os.Getenv("ACCOUNT_KEY"),
		PublicAccessType: azureblob.PublicAccessBlob,
		SASToken:         "", // Not used when bootstrapping with AccountName & AccountKey
	}
	return azureblob.OpenBucket(ctx, &settings, bucket)
}

func setupAzureWithSASToken(ctx context.Context, bucket string) (*blob.Bucket, error) {

	// Use SASToken scoped at the Storage Account (full permission)
	// with this sasToken, ContainerExists can be either true or false
	sasToken := azureblob.GenerateSampleSASTokenForAccount(os.Getenv("ACCOUNT_NAME"), os.Getenv("ACCOUNT_KEY"), time.Now().UTC(), time.Now().UTC().Add(48*time.Hour))

	// Use SASToken scoped to a container (limited permissions, cannot create new container)
	// with this sasToken, ContainerExists should be true to avoid AuthenticationFailure exceptions
	//sasToken := azureblob.GenerateSampleSASTokenForContainerBlob(os.Getenv("ACCOUNT_NAME"), os.Getenv("ACCOUNT_KEY"), bucket, "", time.Now().UTC(), time.Now().UTC().Add(48*time.Hour))

	settings := azureblob.Settings{
		AccountName:      os.Getenv("ACCOUNT_NAME"),
		AccountKey:       "", // Not used when bootstrapping with SASToken
		PublicAccessType: azureblob.PublicAccessContainer,
		SASToken:         sasToken,
	}
	return azureblob.OpenBucket(ctx, &settings, bucket)
}
