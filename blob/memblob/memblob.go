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

// Package memblob provides an in-memory blob implementation.
// Use OpenBucket to construct a *blob.Bucket.
//
// As
//
// memblob does not support any types for As.
package memblob // import "gocloud.dev/blob/memblob"

import (
	"bytes"
	"context"
	"crypto/md5"
	"errors"
	"hash"
	"io"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"gocloud.dev/blob"
	"gocloud.dev/blob/driver"
	"gocloud.dev/gcerrors"
)

const defaultPageSize = 1000

var (
	errNotFound       = errors.New("blob not found")
	errNotImplemented = errors.New("not implemented")
)

func init() {
	blob.DefaultURLMux().RegisterBucket(Scheme, &URLOpener{})
}

// Scheme is the URL scheme memblob registers its URLOpener under on
// blob.DefaultMux.
const Scheme = "mem"

// URLOpener opens URLs like "mem://".
type URLOpener struct{}

// OpenBucketURL returns a new in-memory bucket.
func (*URLOpener) OpenBucketURL(ctx context.Context, u *url.URL) (*blob.Bucket, error) {
	return OpenBucket(nil), nil
}

// Options sets options for constructing a *blob.Bucket backed by memory.
type Options struct{}

type blobEntry struct {
	Content    []byte
	Attributes driver.Attributes
}

type bucket struct {
	mu    sync.Mutex
	blobs map[string]*blobEntry
}

// openBucket creates a driver.Bucket backed by memory.
func openBucket(_ *Options) driver.Bucket {
	return &bucket{
		blobs: map[string]*blobEntry{},
	}
}

// OpenBucket creates a *blob.Bucket backed by memory.
func OpenBucket(opts *Options) *blob.Bucket {
	return blob.NewBucket(openBucket(opts))
}

func (b *bucket) ErrorCode(err error) gcerrors.ErrorCode {
	switch err {
	case errNotFound:
		return gcerrors.NotFound
	case errNotImplemented:
		return gcerrors.Unimplemented
	default:
		return gcerrors.Unknown
	}
}

// ListPaged implements driver.ListPaged.
// The implementation largely mirrors the one in fileblob.
func (b *bucket) ListPaged(ctx context.Context, opts *driver.ListOptions) (*driver.ListPage, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// pageToken is a returned NextPageToken, set below; it's the last key of the
	// previous page.
	var pageToken string
	if len(opts.PageToken) > 0 {
		pageToken = string(opts.PageToken)
	}
	pageSize := opts.PageSize
	if pageSize == 0 {
		pageSize = defaultPageSize
	}

	var keys []string
	for key := range b.blobs {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// If opts.Delimiter != "", lastPrefix contains the last "directory" key we
	// added. It is used to avoid adding it again; all files in this "directory"
	// are collapsed to the single directory entry.
	var lastPrefix string
	var result driver.ListPage
	for _, key := range keys {
		// Skip keys that don't match the Prefix.
		if !strings.HasPrefix(key, opts.Prefix) {
			continue
		}

		entry := b.blobs[key]
		obj := &driver.ListObject{
			Key:     key,
			ModTime: entry.Attributes.ModTime,
			Size:    entry.Attributes.Size,
			MD5:     entry.Attributes.MD5,
		}

		// If using Delimiter, collapse "directories".
		if opts.Delimiter != "" {
			// Strip the prefix, which may contain Delimiter.
			keyWithoutPrefix := key[len(opts.Prefix):]
			// See if the key still contains Delimiter.
			// If no, it's a file and we just include it.
			// If yes, it's a file in a "sub-directory" and we want to collapse
			// all files in that "sub-directory" into a single "directory" result.
			if idx := strings.Index(keyWithoutPrefix, opts.Delimiter); idx != -1 {
				prefix := opts.Prefix + keyWithoutPrefix[0:idx+len(opts.Delimiter)]
				// We've already included this "directory"; don't add it.
				if prefix == lastPrefix {
					continue
				}
				// Update the object to be a "directory".
				obj = &driver.ListObject{
					Key:   prefix,
					IsDir: true,
				}
				lastPrefix = prefix
			}
		}

		// If there's a pageToken, skip anything before it.
		if pageToken != "" && obj.Key <= pageToken {
			continue
		}

		// If we've already got a full page of results, set NextPageToken and return.
		if len(result.Objects) == pageSize {
			result.NextPageToken = []byte(result.Objects[pageSize-1].Key)
			return &result, nil
		}
		result.Objects = append(result.Objects, obj)
	}
	return &result, nil
}

// As implements driver.As.
func (b *bucket) As(i interface{}) bool { return false }

// As implements driver.ErrorAs.
func (b *bucket) ErrorAs(err error, i interface{}) bool { return false }

// Attributes implements driver.Attributes.
func (b *bucket) Attributes(ctx context.Context, key string) (driver.Attributes, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	entry, found := b.blobs[key]
	if !found {
		return driver.Attributes{}, errNotFound
	}
	return entry.Attributes, nil
}

// NewRangeReader implements driver.NewRangeReader.
func (b *bucket) NewRangeReader(ctx context.Context, key string, offset, length int64, opts *driver.ReaderOptions) (driver.Reader, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	entry, found := b.blobs[key]
	if !found {
		return nil, errNotFound
	}

	r := bytes.NewReader(entry.Content)
	if offset > 0 {
		if _, err := r.Seek(offset, io.SeekStart); err != nil {
			return nil, err
		}
	}
	var ior io.Reader = r
	if length >= 0 {
		ior = io.LimitReader(r, length)
	}
	return &reader{
		r: ior,
		attrs: driver.ReaderAttributes{
			ContentType: entry.Attributes.ContentType,
			ModTime:     entry.Attributes.ModTime,
			Size:        entry.Attributes.Size,
		},
	}, nil
}

type reader struct {
	r     io.Reader
	attrs driver.ReaderAttributes
}

func (r *reader) Read(p []byte) (int, error) {
	return r.r.Read(p)
}

func (r *reader) Close() error {
	return nil
}

func (r *reader) Attributes() driver.ReaderAttributes {
	return r.attrs
}

func (r *reader) As(i interface{}) bool { return false }

// NewTypedWriter implements driver.NewTypedWriter.
func (b *bucket) NewTypedWriter(ctx context.Context, key string, contentType string, opts *driver.WriterOptions) (driver.Writer, error) {
	if key == "" {
		return nil, errors.New("invalid key (empty string(")
	}
	b.mu.Lock()
	defer b.mu.Unlock()

	if opts.BeforeWrite != nil {
		if err := opts.BeforeWrite(func(interface{}) bool { return false }); err != nil {
			return nil, err
		}
	}
	md := map[string]string{}
	for k, v := range opts.Metadata {
		md[k] = v
	}
	return &writer{
		ctx:         ctx,
		b:           b,
		key:         key,
		contentType: contentType,
		metadata:    md,
		opts:        opts,
		md5hash:     md5.New(),
	}, nil
}

type writer struct {
	ctx         context.Context
	b           *bucket
	key         string
	contentType string
	metadata    map[string]string
	opts        *driver.WriterOptions
	buf         bytes.Buffer
	// We compute the MD5 hash so that we can store it with the file attributes,
	// not for verification.
	md5hash hash.Hash
}

func (w *writer) Write(p []byte) (n int, err error) {
	if _, err := w.md5hash.Write(p); err != nil {
		return 0, err
	}
	return w.buf.Write(p)
}

func (w *writer) Close() error {
	// Check if the write was cancelled.
	if err := w.ctx.Err(); err != nil {
		return err
	}

	md5sum := w.md5hash.Sum(nil)
	content := w.buf.Bytes()
	entry := &blobEntry{
		Content: content,
		Attributes: driver.Attributes{
			CacheControl:       w.opts.CacheControl,
			ContentDisposition: w.opts.ContentDisposition,
			ContentEncoding:    w.opts.ContentEncoding,
			ContentLanguage:    w.opts.ContentLanguage,
			ContentType:        w.contentType,
			Metadata:           w.metadata,
			Size:               int64(len(content)),
			ModTime:            time.Now(),
			MD5:                md5sum,
		},
	}
	w.b.mu.Lock()
	defer w.b.mu.Unlock()
	w.b.blobs[w.key] = entry
	return nil
}

// Delete implements driver.Delete.
func (b *bucket) Delete(ctx context.Context, key string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.blobs[key] == nil {
		return errNotFound
	}
	delete(b.blobs, key)
	return nil
}

func (b *bucket) SignedURL(ctx context.Context, key string, opts *driver.SignedURLOptions) (string, error) {
	return "", errNotImplemented
}
