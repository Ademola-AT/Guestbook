// Copyright 2018 Google LLC
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

//+build wireinject

package main

import (
	"context"

	"github.com/google/go-x-cloud/blob"
	"github.com/google/go-x-cloud/blob/gcsblob"
	"github.com/google/go-x-cloud/gcp"
	"github.com/google/go-x-cloud/gcp/gcpcloud"
	"github.com/google/go-x-cloud/mysql/cloudmysql"
	"github.com/google/go-x-cloud/runtimevar"
	"github.com/google/go-x-cloud/runtimevar/runtimeconfigurator"
	"github.com/google/go-x-cloud/wire"
)

// This file wires the generic interfaces up to Google Cloud Platform (GCP). It
// won't be directly included in the final binary, since it includes a Wire
// injector template function (setupGCP), but the declarations will be copied
// into wire_gen.go when gowire is run.

// setupGCP is a Wire injector function that sets up the application using GCP.
func setupGCP(ctx context.Context, flags *cliFlags) (*application, func(), error) {
	// This will be filled in by gowire with providers from the provider sets in
	// wire.Build.

	panic(wire.Build(
		gcpcloud.GCP,
		applicationSet,
		gcpBucket,
		gcpMOTDVar,
		gcpSQLParams,
	))
}

// gcpBucket is a Wire provider function that returns the GCS bucket based on
// the command-line flags.
func gcpBucket(ctx context.Context, flags *cliFlags, ts gcp.TokenSource) (*blob.Bucket, error) {
	return gcsblob.New(ctx, flags.bucket, &gcsblob.BucketOptions{
		TokenSource: ts,
	})
}

// gcpSQLParams is a Wire provider function that returns the Cloud SQL
// connection parameters based on the command-line flags. Other providers inside
// gcpcloud.GCP use the parameters to construct a *sql.DB.
func gcpSQLParams(id gcp.ProjectID, flags *cliFlags) *cloudmysql.Params {
	return &cloudmysql.Params{
		ProjectID: string(id),
		Region:    flags.cloudSQLRegion,
		Instance:  flags.dbHost,
		Database:  flags.dbName,
		User:      flags.dbUser,
		Password:  flags.dbPassword,
	}
}

// gcpMOTDVar is a Wire provider function that returns the Message of the Day
// variable from Runtime Configurator.
func gcpMOTDVar(ctx context.Context, client *runtimeconfigurator.Client, project gcp.ProjectID, flags *cliFlags) (*runtimevar.Variable, func(), error) {
	name := runtimeconfigurator.ResourceName{
		ProjectID: string(project),
		Config:    flags.runtimeConfigName,
		Variable:  flags.motdVar,
	}
	v, err := client.NewVariable(ctx, name, runtimevar.StringDecoder, &runtimeconfigurator.WatchOptions{
		WaitTime: flags.motdVarWaitTime,
	})
	if err != nil {
		return nil, nil, err
	}
	return v, func() { v.Close() }, nil
}
