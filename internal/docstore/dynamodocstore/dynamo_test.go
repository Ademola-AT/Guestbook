// Copyright 2019 The Go Cloud Development Kit Authors
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
	"context"
	"testing"

	"gocloud.dev/internal/docstore"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	dyn "github.com/aws/aws-sdk-go/service/dynamodb"
	"gocloud.dev/internal/docstore/driver"
	"gocloud.dev/internal/docstore/drivertest"
	"gocloud.dev/internal/testing/setup"
)

const (
	region         = "us-east-2"
	collectionName = "docstore-test"
	partKey        = "_kind"
	sortKey        = "_id"
)

type harness struct {
	sess   *session.Session
	closer func()
}

func newHarness(ctx context.Context, t *testing.T) (drivertest.Harness, error) {
	sess, _, done := setup.NewAWSSession(t, region)
	return &harness{sess: sess, closer: done}, nil
}

func (h *harness) Close() {
	h.closer()
}

func (h *harness) MakeCollection(context.Context) (driver.Collection, error) {
	return newCollection(dyn.New(h.sess), collectionName, partKey, sortKey), nil
}

func TestConformance(t *testing.T) {
	if *setup.Record {
		clearTable(t)
	}
	drivertest.MakeUniqueStringDeterministicForTesting(1)
	drivertest.RunConformanceTests(t, newHarness, &codecTester{})
}

func clearTable(t *testing.T) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		t.Fatal(err)
	}
	db := dyn.New(sess)
	in := &dyn.ScanInput{
		TableName:            aws.String(collectionName),
		ProjectionExpression: aws.String("#pk, #sk"),
		ExpressionAttributeNames: map[string]*string{
			"#pk": aws.String(partKey),
			"#sk": aws.String(sortKey),
		},
	}
	for {
		out, err := db.Scan(in)
		if err != nil {
			t.Fatal(err)
		}
		if len(out.Items) > 0 {
			bwin := &dyn.BatchWriteItemInput{
				RequestItems: map[string][]*dyn.WriteRequest{},
			}
			var wrs []*dyn.WriteRequest
			for _, item := range out.Items {
				wrs = append(wrs, &dyn.WriteRequest{
					DeleteRequest: &dyn.DeleteRequest{Key: item},
				})
			}
			bwin.RequestItems[collectionName] = wrs
			if _, err := db.BatchWriteItem(bwin); err != nil {
				t.Fatal(err)
			}
		}
		if out.LastEvaluatedKey == nil {
			break
		}
		in.ExclusiveStartKey = out.LastEvaluatedKey
	}
}

// Dynamodocstore-specific tests.

func TestOpenKeeper(t *testing.T) {
	tests := []struct {
		URL     string
		WantErr bool
	}{
		// OK.
		{"dynamodb://mytable?partition_key=_kind", false},
		// OK.
		{"dynamodb://mytable?partition_key=_kind&sort_key=_id", false},
		// OK, overriding region.
		{"dynamodb://mytable?partition_key=_kind&region=us-west1", false},
		// Unknown parameter.
		{"dynamodb://mytable?partition_key=_kind&param=value", true},
		// With path.
		{"dynamodb://mytable/subcoll?partition_key=_kind", true},
		// Missing partition_key.
		{"dynamodb://mytable?sort_key=_id", true},
	}

	ctx := context.Background()
	for _, test := range tests {
		_, err := docstore.OpenCollection(ctx, test.URL)
		if (err != nil) != test.WantErr {
			t.Errorf("%s: got error %v, want error %v", test.URL, err, test.WantErr)
		}
	}
}
