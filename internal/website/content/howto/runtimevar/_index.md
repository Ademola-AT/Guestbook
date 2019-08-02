---
title: "Runtime Configuration"
date: 2019-07-11T12:00:00-07:00
lastmod: 2019-07-29T12:00:00-07:00
showInSidenav: true
toc: true
---

The runtimevar package provides an easy and portable way to watch runtime
configuration variables. This guide shows how to work with runtime configuration
variables using the Go CDK.

<!--more-->

## Opening a Variable {#opening}

The first step in watching a variable is to instantiate a
[`*runtimevar.Variable`][].

The easiest way to do so is to use [`runtimevar.OpenVariable`][] and a URL pointing
to the variable, making sure you ["blank import"][] the driver package to link
it in. See [Concepts: URLs][] for more details. If you need fine-grained control
over the connection settings, you can call the constructor function in the
driver package directly (like `etcdvar.OpenVariable`).

When opening the variable, you can provide a [decoder][] parameter (either as a
[query parameter][] for URLs, or explicitly to the constructor) to specify
whether the raw value stored in the variable is interpreted as a `string`, a
`[]byte`, or as JSON. Here's an example of using a JSON encoder:

{{< goexample src="gocloud.dev/runtimevar.Example_jsonDecoder" imports="0" >}}

See the [guide below][] for usage of both forms for each supported provider.

[`*runtimevar.Variable`]: https://godoc.org/gocloud.dev/runtimevar#Variable
[`runtimevar.OpenVariable`]: https://godoc.org/gocloud.dev/runtimevar#OpenVariable
["blank import"]: https://golang.org/doc/effective_go.html#blank_import
[Concepts: URLs]: {{< ref "/concepts/urls.md" >}}
[decoder]: https://godoc.org/gocloud.dev/runtimevar#Decoder
[guide below]: {{< ref "#services" >}}
[query parameter]: https://godoc.org/gocloud.dev/runtimevar#DecoderByName

## Using a Variable {#using}

Once you have opened a `runtimevar.Variable` for the provider you want, you can
use it portably.

### Latest {#latest}

The easiest way to a `Variable` is to use the [`Variable.Latest`][] method. It
returns the latest good [`Snapshot`][] of the variable value, blocking if no
good value has *ever* been received. The dynamic type of `Snapshot.Value`
depends on the decoder you provided when creating the `Variable`.

To avoid blocking, you can pass an already-`Done` context.

{{< goexample src="gocloud.dev/runtimevar.ExampleVariable_Latest" imports="0" >}}

[`Variable.Latest`]: https://godoc.org/gocloud.dev/runtimevar#Variable.Latest
[`Snapshot`]: https://godoc.org/gocloud.dev/runtimevar#Snapshot

### Watch {#watch}

`Variable` also has a [`Watch`][] method for obtaining the value of a variable;
it has different semantics than `Latest` and may be useful in some scenarios. We
recommend starting with `Latest` as it's conceptually simpler to work with.

[`Watch`]: https://godoc.org/gocloud.dev/runtimevar#Variable.Watch

## Supported Services {#services}

### GCP Runtime Configurator {#gcprc}

To open a variable stored in [GCP Runtime Configurator][] via a URL, you can use
the `runtimevar.OpenVariable` function as shown in the example below.

[GCP Runtime Configurator]: https://cloud.google.com/deployment-manager/runtime-configurator/

`runtimevar.OpenVariable` will use Application Default Credentials; if you have
authenticated via [`gcloud auth login`][], it will use those credentials. See
[Application Default Credentials][GCP creds] to learn about authentication
alternatives, including using environment variables.

[GCP creds]: https://cloud.google.com/docs/authentication/production
[`gcloud auth login`]: https://cloud.google.com/sdk/gcloud/reference/auth/login

{{< goexample
"gocloud.dev/runtimevar/gcpruntimeconfig.Example_openVariableFromURL" >}}

#### GCP Constructor {#gcprc-ctor}

The [`gcpruntimeconfig.OpenVariable`][] constructor opens a Runtime Configurator
variable.

{{< goexample "gocloud.dev/runtimevar/gcpruntimeconfig.ExampleOpenVariable" >}}

[`gcpruntimeconfig.OpenVariable`]: https://godoc.org/gocloud.dev/runtimevar/gcpruntimeconfig#OpenVariable

### AWS Parameter Store {#awsps}

To open a variable stored in [AWS Parameter Store][] via a URL, you can use the
`runtimevar.OpenVariable` function as shown in the example below.

[AWS Parameter Store]:
https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-parameter-store.html

`runtimevar.OpenVariable` will create a default AWS Session with the
`SharedConfigEnable` option enabled; if you have authenticated with the AWS CLI,
it will use those credentials. See [AWS Session][] to learn about authentication
alternatives, including using environment variables.

[AWS Session]: https://docs.aws.amazon.com/sdk-for-go/api/aws/session/

{{< goexample
"gocloud.dev/runtimevar/awsparamstore.Example_openVariableFromURL" >}}

#### AWS Constructor {#awsps-ctor}

The [`awsparamstore.OpenVariable`][] constructor opens a Parameter Store
variable.

{{< goexample "gocloud.dev/runtimevar/awsparamstore.ExampleOpenVariable" >}}

[`awsparamstore.OpenVariable`]:
https://godoc.org/gocloud.dev/runtimevar/awsparamstore#OpenVariable

### etcd {#etcd}

To open a variable stored in [etcd][] via URL, you can use the
`runtimevar.OpenVariable` function as follows.

{{< goexample "gocloud.dev/runtimevar/etcdvar.Example_openVariableFromURL" >}}

[etcd]: https://etcd.io/

#### etcd Constructor {#etcd-ctor}

The [`etcdvar.OpenVariable`][] constructor opens an `etcd` variable.

[`etcdvar.OpenVariable`]:
https://godoc.org/gocloud.dev/runtimevar/etcdvar#OpenVariable

{{< goexample "gocloud.dev/runtimevar/etcdvar.ExampleOpenVariable" >}}

### HTTP {#http}

`httpvar` supports watching a variable via an HTTP request. Use
`runtimevar.OpenVariable` with a regular URL starting with `http` or `https`.
`httpvar` will periodically make an HTTP `GET` request to that URL, with the
`decode` URL parameter removed (if present).

{{< goexample "gocloud.dev/runtimevar/httpvar.Example_openVariableFromURL" >}}

#### HTTP Constructor {#http-ctor}

The [`httpvar.OpenVariable`][] constructor opens a variable with a `http.Client`
and a URL.

{{< goexample "gocloud.dev/runtimevar/httpvar.ExampleOpenVariable" >}}

[`httpvar.OpenVariable`]: https://godoc.org/gocloud.dev/runtimevar/httpvar#OpenVariable

### Blob {#blob}

`blobvar` supports watching a variable based on the contents of a
[Go CDK blob][]. Set the environment variable `BLOBVAR_BUCKET_URL` to the URL
of the bucket, and then use `runtimevar.OpenVariable` as shown below.
`blobvar` will periodically re-fetch the contents of the blob.

{{< goexample "gocloud.dev/runtimevar/blobvar.Example_openVariableFromURL" >}}

[Go CDK blob]: https://gocloud.dev/howto/blob/

You can also use [`blobvar.OpenVariable`][].

[`blobvar.OpenVariable`]: https://godoc.org/gocloud.dev/runtimevar/blobvar#OpenVariable

### Local {#local}

You can create an in-memory variable (useful for testing) using `constantvar`:

{{< goexample "gocloud.dev/runtimevar/constantvar.Example_openVariableFromURL" >}}

Alternatively, you can create a variable based on the contents of a file using
`filevar`:

{{< goexample "gocloud.dev/runtimevar/filevar.Example_openVariableFromURL" >}}
