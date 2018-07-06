# The Go Cloud Project
_Write once, run on any cloud ☁️_

[![Build Status](https://travis-ci.com/google/go-cloud.svg?branch=master)][travis]
[![godoc](https://godoc.org/github.com/google/go-cloud?status.svg)][godoc]

The Go Cloud Project is an initiative that will allow application developers to seamlessly deploy cloud applications on any combination of cloud providers. It does this by providing stable, idiomatic interfaces for common uses like storage and databases. Think `database/sql` for cloud products.

A key part of the project is to also provide a code generator called [Wire](https://github.com/google/go-cloud/blob/master/wire/README.md). It creates human-readable code that only imports the cloud SDKs for providers you use. This allows Go Cloud to grow to support any number of cloud providers, without increasing compile times or binary sizes, and avoiding any side effects from `init()` functions.

Imagine writing this to read from blob storage (like Google Cloud Storage or S3):

```go
blobReader, err := bucket.NewReader(context.Background(), "my-blob")
```

and being able to run that code on any cloud you want, avoiding all the ceremony of cloud-specific authorization, tracing, SDKs and all the other code required to make an application portable across cloud platforms.

## Installation instructions
Installation is easy, but does require `vgo`. `vgo` is not yet stable, and so builds may break with `vgo` changes, but experience has shown this to be rare.

```shell
$ go get -u golang.org/x/vgo
$ git clone https://github.com/google/go-cloud.git
$ cd go-cloud
$ vgo install ./wire/cmd/gowire
```
Go Cloud builds at the latest stable release of Go. Previous Go versions may compile but are not supported.

## Samples
[`samples/tutorial`](https://github.com/google/go-cloud/tree/master/samples/tutorial) shows how to get started with the project by using blob storage.

[`samples/guestbook`](https://github.com/google/go-cloud/tree/master/samples/guestbook) contains an example guestbook application (just like it's 1999!) that can be run locally, on Google Cloud Platform or on Amazon Web Services. The instructions take about 5 minutes to follow if running locally. If you want to see the guestbook app running on cloud resources, it will take about 30 minutes to follow, and uses [Terraform](http://terraform.io) to automatically provision the resources needed.

## Project status
While in alpha, the API is subject to breaking changes so is not yet suitable for production. We encourage you to experiment with Go Cloud and make contributions to help evolve it to meet your needs!

[travis]: https://travis-ci.com/google/go-cloud
[godoc]: http://godoc.org/github.com/google/go-cloud

The GitHub repository at [google/go-cloud](https://github.com/google/go-cloud) currently contains [Google Cloud Platform](http://cloud.google.com) and [Amazon Web Services](http://aws.amazon.com) implementations as examples to prove everything is working. If you create a repository that implements the Go Cloud interfaces for other cloud providers, let us know and we would be happy to link to it here.

## Current features

Go Cloud provides generic APIs for:
-   Unstructured binary (blob) storage
-   Variables that change at runtime (configuration)
-   Connecting to MySQL databases
-   Server startup and diagnostics: request logging, tracing, and health
    checking
