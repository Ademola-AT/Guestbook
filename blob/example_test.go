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

package blob_test

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/google/go-cloud/blob"
	"github.com/google/go-cloud/blob/fileblob"
)

func ExampleBucket_NewReader() {
	// Connect to a bucket when your program starts up.
	// This example uses the file-based implementation.
	dir, cleanup := newTempDir()
	defer cleanup()
	// Write a file to read using the bucket.
	err := ioutil.WriteFile(filepath.Join(dir, "foo.txt"), []byte("Hello, World!\n"), 0666)
	if err != nil {
		log.Fatal(err)
	}
	// Create the file-based bucket.
	bucket, err := fileblob.OpenBucket(dir)
	if err != nil {
		log.Fatal(err)
	}

	// Open a reader using the blob's key.
	ctx := context.Background()
	r, err := bucket.NewReader(ctx, "foo.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	// The blob reader implements io.Reader, so we can use any function that
	// accepts an io.Reader.
	if _, err := io.Copy(os.Stdout, r); err != nil {
		log.Fatal(err)
	}

	// Output:
	// Hello, World!
}

func ExampleBucket_NewRangeReader() {
	// Connect to a bucket when your program starts up.
	// This example uses the file-based implementation.
	dir, cleanup := newTempDir()
	defer cleanup()
	// Write a file to read using the bucket.
	err := ioutil.WriteFile(filepath.Join(dir, "foo.txt"), []byte("Hello, World!\n"), 0666)
	if err != nil {
		log.Fatal(err)
	}
	// Create the file-based bucket.
	bucket, err := fileblob.OpenBucket(dir)
	if err != nil {
		log.Fatal(err)
	}

	// Open a reader using the blob's key at a specific offset at length.
	ctx := context.Background()
	r, err := bucket.NewRangeReader(ctx, "foo.txt", 1, 4)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	// The blob reader implements io.Reader, so we can use any function that
	// accepts an io.Reader.
	if _, err := io.Copy(os.Stdout, r); err != nil {
		log.Fatal(err)
	}

	// Output:
	// ello
}

func ExampleBucket_NewWriter() {
	// Connect to a bucket when your program starts up.
	// This example uses the file-based implementation.
	dir, cleanup := newTempDir()
	defer cleanup()
	bucket, err := fileblob.OpenBucket(dir)
	if err != nil {
		log.Fatal(err)
	}

	// Open a writer using the key "foo.txt" and the default options.
	ctx := context.Background()
	// fileblob doesn't support custom content-type yet, see
	// https://github.com/google/go-cloud/issues/111.
	w, err := bucket.NewWriter(ctx, "foo.txt", &blob.WriterOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		log.Fatal(err)
	}
	// The blob writer implements io.Writer, so we can use any function that
	// accepts an io.Writer. A writer must always be closed.
	_, printErr := fmt.Fprintln(w, "Hello, World!")
	closeErr := w.Close()
	if printErr != nil {
		log.Fatal(printErr)
	}
	if closeErr != nil {
		log.Fatal(closeErr)
	}
	// Copy the written blob to stdout.
	r, err := bucket.NewReader(ctx, "foo.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	if _, err := io.Copy(os.Stdout, r); err != nil {
		log.Fatal(err)
	}

	// Output:
	// Hello, World!
}

func ExampleBucket_ReadAll() {
	// Connect to a bucket when your program starts up.
	// This example uses the file-based implementation.
	dir, cleanup := newTempDir()
	defer cleanup()

	// Create the file-based bucket.
	bucket, err := fileblob.OpenBucket(dir)
	if err != nil {
		log.Fatal(err)
	}

	// Write a blob using WriteAll.
	ctx := context.Background()
	if err := bucket.WriteAll(ctx, "foo.txt", []byte("Go Cloud"), nil); err != nil {
		log.Fatal(err)
	}

	// Read it back using ReadAll.
	b, err := bucket.ReadAll(ctx, "foo.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))

	// Output:
	// Go Cloud
}

func ExampleBucket_List() {
	// Connect to a bucket when your program starts up.
	// This example uses the file-based implementation.
	dir, cleanup := newTempDir()
	defer cleanup()

	// Create the file-based bucket.
	bucket, err := fileblob.OpenBucket(dir)
	if err != nil {
		log.Fatal(err)
	}

	// Create some blob objects for listing: "foo[0..4].txt".
	ctx := context.Background()
	for i := 0; i < 5; i++ {
		if err := bucket.WriteAll(ctx, fmt.Sprintf("foo%d.txt", i), []byte("Go Cloud"), nil); err != nil {
			log.Fatal(err)
		}
	}

	// Iterate over them.
	// This will list the blobs created above because fileblob is strongly
	// consistent, but is not guaranteed to work on all providers.
	iter, err := bucket.List(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		obj, err := iter.Next(ctx)
		if err != nil {
			log.Fatal(err)
		}
		if obj == nil {
			break
		}
		fmt.Println(obj.Key)
	}

	// Output:
	// foo0.txt
	// foo1.txt
	// foo2.txt
	// foo3.txt
	// foo4.txt
}

func ExampleBucket_List_withDelimiter() {
	// Connect to a bucket when your program starts up.
	// This example uses the file-based implementation.
	dir, cleanup := newTempDir()
	defer cleanup()

	// Create the file-based bucket.
	bucket, err := fileblob.OpenBucket(dir)
	if err != nil {
		log.Fatal(err)
	}

	// Create some blob objects in a hierarchy.
	ctx := context.Background()
	for _, key := range []string{
		"dir1/subdir/a.txt",
		"dir1/subdir/b.txt",
		"dir2/c.txt",
		"d.txt",
	} {
		if err := bucket.WriteAll(ctx, key, []byte("Go Cloud"), nil); err != nil {
			log.Fatal(err)
		}
	}

	// list lists files in b starting with prefix. It uses the delimiter "/",
	// and recurses into "directories", adding 2 spaces to indent each time.
	// It will list the blobs created above because fileblob is strongly
	// consistent, but is not guaranteed to work on all providers.
	var list func(context.Context, *blob.Bucket, string, string)
	list = func(ctx context.Context, b *blob.Bucket, prefix, indent string) {
		iter, err := b.List(ctx, &blob.ListOptions{
			Delimiter: "/",
			Prefix:    prefix,
		})
		if err != nil {
			log.Fatal(err)
		}
		for {
			obj, err := iter.Next(ctx)
			if err != nil {
				log.Fatal(err)
			}
			if obj == nil {
				break
			}
			fmt.Printf("%s%s\n", indent, obj.Key)
			if obj.IsDir {
				list(ctx, b, obj.Key, indent+"  ")
			}
		}
	}
	list(ctx, bucket, "", "")

	// Output:
	// d.txt
	// dir1/
	//   dir1/subdir/
	//     dir1/subdir/a.txt
	//     dir1/subdir/b.txt
	// dir2/
	//   dir2/c.txt
}
func ExampleBucket_As() {
	// Connect to a bucket when your program starts up.
	// This example uses the file-based implementation.
	dir, cleanup := newTempDir()
	defer cleanup()

	// Create the file-based bucket.
	bucket, err := fileblob.OpenBucket(dir)
	if err != nil {
		log.Fatal(err)
	}
	// This example uses As to try to fill in a string variable. As will return
	// false because fileblob doesn't support any types for Bucket.As.
	// See the package documentation for your provider (e.g., gcsblob or s3blob)
	// to see what type(s) it supports.
	var providerSpecific string
	if bucket.As(&providerSpecific) {
		fmt.Println("fileblob supports the `string` type for Bucket.As")
		// Use providerSpecific.
	} else {
		fmt.Println("fileblob does not support the `string` type for Bucket.As")
	}

	// This example sets WriterOptions.BeforeWrite to be called before the
	// provider starts writing. In the callback, it uses asFunc to try to fill in
	// a *string. Again, asFunc will return false because fileblob doesn't support
	// any types for Writer.
	fn := func(asFunc func(i interface{}) bool) error {
		var mutableProviderSpecific *string
		if asFunc(&mutableProviderSpecific) {
			fmt.Println("fileblob supports the `*string` type for WriterOptions.BeforeWrite")
			// Use mutableProviderSpecific.
		} else {
			fmt.Println("fileblob does not support the `*string` type for WriterOptions.BeforeWrite")
		}
		return nil
	}
	ctx := context.Background()
	if err := bucket.WriteAll(ctx, "foo.txt", []byte("Go Cloud"), &blob.WriterOptions{BeforeWrite: fn}); err != nil {
		log.Fatal(err)
	}
	// Output:
	// fileblob does not support the `string` type for Bucket.As
	// fileblob does not support the `*string` type for WriterOptions.BeforeWrite
}

func ExampleOpen() {
	// Connect to a bucket using a URL.
	// This example uses the file-based implementation, which registers for
	// the "file" protocol.
	dir, cleanup := newTempDir()
	defer cleanup()

	ctx := context.Background()
	if _, err := blob.Open(ctx, "file:///nonexistentpath"); err == nil {
		log.Fatal("Expected an error opening nonexistent path")
	}
	fmt.Println("Error opening nonexistentpath")

	if _, err := blob.Open(ctx, "file://" + dir); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Got bucket for valid path")
	}
	// Output:
	// Error opening nonexistentpath
	// Got bucket for valid path
}

func newTempDir() (string, func()) {
	dir, err := ioutil.TempDir("", "go-cloud-blob-example")
	if err != nil {
		panic(err)
	}
	return dir, func() { os.RemoveAll(dir) }
}
