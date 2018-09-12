# Go Cloud module proxy

The Travis build for Go Cloud uses a [Go module proxy][] for dependencies, for
efficiency and reliability, and to enforce that all of our dependencies meet our
license requirements.

The proxy is set
[here](https://github.com/google/go-cloud/blob/master/.travis.yml#L22).

[Go module proxy]: https://research.swtch.com/vgo-module

## Updating the GCS bucket

When you add a new dependency, the Travis build will fail. To add the new
dependency to the proxy:

1.  Ensure that the new dependency is covered under one of the `notice`,
    `permissive`, or `unencumbered` licences under the categories defined by
    [Google Open Source](https://opensource.google.com/docs/thirdparty/licenses/).
2.  Add the dependency to the GCS bucket:

```bash
./internal/proxy/makeproxy.sh proxy

# -n: preview only
# -r: recurse directories
# -c: compare checksums not write times
# -d: source directory; value should match what you passed to makeproxy.sh
$ gsutil rsync -n -r -c -d proxy gs://go-cloud-modules
...
# If it looks good, repeat without the -n.
```
