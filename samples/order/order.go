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

package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"gocloud.dev/blob"
	"gocloud.dev/docstore"
	"gocloud.dev/pubsub"
)

var (
	requestTopicURL  = flag.String("request-topic", "mem://requests", "gocloud.dev/pubsub URL for request topic")
	requestSubURL    = flag.String("request-sub", "mem://requests", "gocloud.dev/pubsub URL for request subscription")
	responseTopicURL = flag.String("response-topic", "mem://responses", "gocloud.dev/pubsub URL for response topic")
	responseSubURL   = flag.String("response-sub", "mem://responses", "gocloud.dev/pubsub URL for response subscription")
	bucketURL        = flag.String("bucket", "", "gocloud.dev/blob URL for image bucket")
	collectionURL    = flag.String("collection", "mem://orders/ID", "gocloud.dev/docstore URL for order collection")
	runProcessor     = flag.Bool("processor", true, "run the image processor")
)

func main() {

	// TODO(jba): add frontend
	flag.Parse()
	_, processor, cleanup, err := setup()
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()
	if err := processor.run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func setup() (_ *frontend, _ *processor, cleanup func(), err error) {
	cleanup = func() {}

	ctx := context.Background()
	reqTopic, err := pubsub.OpenTopic(ctx, *requestTopicURL)
	if err != nil {
		return nil, nil, cleanup, err
	}
	reqSub, err := pubsub.OpenSubscription(ctx, *requestSubURL)
	if err != nil {
		return nil, nil, cleanup, err
	}
	resTopic, err := pubsub.OpenTopic(ctx, *responseTopicURL)
	if err != nil {
		return nil, nil, cleanup, err
	}
	resSub, err := pubsub.OpenSubscription(ctx, *responseSubURL)
	if err != nil {
		return nil, nil, cleanup, err
	}

	burl := *bucketURL
	if burl == "" {
		dir, err := ioutil.TempDir("", "gocdk-order")
		if err != nil {
			return nil, nil, cleanup, err
		}
		burl = "file://" + dir
		cleanup = func() { os.Remove(dir) }
	}
	bucket, err := blob.OpenBucket(ctx, burl)
	if err != nil {
		cleanup()
		return nil, nil, nil, err
	}
	coll, err := docstore.OpenCollection(ctx, *collectionURL)
	if err != nil {
		cleanup()
		return nil, nil, nil, err
	}
	f := &frontend{
		requestTopic: reqTopic,
		responseSub:  resSub,
		bucket:       bucket,
		coll:         coll,
	}
	p := &processor{
		requestSub:    reqSub,
		responseTopic: resTopic,
		bucket:        bucket,
		coll:          coll,
	}
	return f, p, cleanup, nil
}
