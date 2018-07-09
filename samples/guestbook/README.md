# Guestbook Sample

Guestbook is a sample application that records visitors' messages, displays a
cloud banner, and an administrative message. The main business logic is
written in a cloud-agnostic manner using MySQL, the generic blob API, and the
generic runtimevar API. All platform-specific code is set up by 
[Wire](https://github.com/google/go-cloud/tree/master/wire).

## Prerequisites

You will need to install the following software to run this sample:

- [Go](https://golang.org/doc/install) and
  [vgo](https://go.googlesource.com/vgo)
- [gowire](https://github.com/google/go-cloud/blob/master/wire/README.md#installing)
- [Docker](https://docs.docker.com/install/)
- [Terraform][TF]
- [jq](https://stedolan.github.io/jq/download/)
- [gcloud CLI](https://cloud.google.com/sdk/downloads), if you want to use GCP
- [aws CLI](https://docs.aws.amazon.com/cli/latest/userguide/installing.html),
  if you want to use AWS

## Building

`gowire` is not compatible with `vgo` yet, so you must run `vgo mod -vendor`
first to download all the dependencies in `go.mod`. Running `gowire` generates
the Wire code. Run the following in the `samples/guestbook` directory:

```shell
# First time, for gowire.
$ vgo mod -vendor

# Now build:
$ gowire && vgo build
```

## Running Locally

You will need to run a local MySQL database server, set up a directory that
simulates a bucket, and create a local message of the day. `localdb.sh` is a
script that runs a temporary database using Docker:

```shell
./localdb.sh
```

In another terminal, you can run:

```shell
# Set a local Message of the Day
echo 'Hello, World!' > motd.txt

# Run the server.
./guestbook -env=local -bucket=blobs -motd_var=motd.txt
```

Your server should be running on http://localhost:8080/.

You can stop the MySQL database server with Ctrl-\\. MySQL ignores Ctrl-C
(SIGINT).

## Running on Google Cloud Platform (GCP)

If you want to run this sample on GCP, you need to create a project, download
the gcloud SDK, install `kubectl` and log in.

``` shell
# Install kubectl.
gcloud components install kubectl

# Opens a browser to log you into GCP.
gcloud auth login
```

You can then use [Terraform][TF], a tool for initializing cloud resources, to
set up your project. Finally, this sample provides a script for building the
Guestbook binary and deploying it to the Kubernetes cluster created by
Terraform.

```shell
gcloud auth application-default login
cd gcp
terraform init

# Terraform will prompt you for your GCP project ID, desired region,
# and desired zone.
terraform apply

./deploy.sh
```

The deploy script will display the URL of your running service.

To clean up the created resources, run `terraform destroy` inside the `gcp`
directory using the same variables you entered during `terraform apply`.

## Running on Amazon Web Services (AWS)

If you want to run this sample on AWS, you need to set up an account, download
the AWS command line interface, and log in. You will also need an SSH key. If you
don't already have one, you can follow [this guide from GitHub][GitHub SSH].

Once everything is set up, you can then use Terraform, a tool for initializing
cloud resources, to set up your project. This will create an EC2 instance you can
SSH to and run your binary.

```shell
aws configure
ssh-add
GOOS=linux GOARCH=amd64 vgo build
cd aws
terraform init
terraform apply -var region=us-west-1 -var ssh_public_key="$(cat ~/.ssh/id_rsa.pub)"

# SSH into the EC2 instance.
ssh "admin@$( terraform output instance_host )"
```

When you're connected to the server, run the server binary. Replace the
command-line flag values with values from the output of `terraform apply`.

```
AWS_REGION=us-west-1 ./guestbook -env=aws \
  -bucket=... -db_host=... -motd_var=...
```

You can then visit the server at `http://INSTANCE_HOST:8080/`, where
`INSTANCE_HOST` is the value of `terraform output instance_host` run on your
local machine.

To clean up the created resources, run `terraform destroy` inside the `aws`
directory using the same variables you entered during `terraform apply`.

[GitHub SSH]: https://help.github.com/articles/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent/

## Gophers

The Go gopher was designed by Renee French and used under the [Creative Commons
3.0 Attributions](https://creativecommons.org/licenses/by/3.0/) license.

[TF]: https://www.terraform.io/intro/getting-started/install.html
