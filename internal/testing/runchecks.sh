#!/usr/bin/env bash
# Copyright 2018 The Go Cloud Development Kit Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Runs only tests relevant to the current pull request.
# At the moment, this only gates running the Wire test suite.
# See https://github.com/google/go-cloud/issues/28 for solving the
# general case.

# https://coderwall.com/p/fkfaqq/safer-bash-scripts-with-set-euxo-pipefail
set -euxo pipefail

if [[ $# -gt 0 ]]; then
  echo "usage: runchecks.sh" 1>&2
  exit 64
fi

# The following logic lets us skip the (lengthy) installation process and tests
# in some cases where the PR carries trivial changes that don't affect the code
# (such as documentation-only).
if [[ ! -z "$TRAVIS_BRANCH" ]] && [[ ! -z "$TRAVIS_PULL_REQUEST_SHA" ]]; then
  tmpfile=$(mktemp)
  function cleanup() {
    rm -rf "$tmpfile"
  }
  trap cleanup EXIT

  mergebase="$(git merge-base -- "$TRAVIS_BRANCH" "$TRAVIS_PULL_REQUEST_SHA")"
  git diff --name-only "$mergebase" "$TRAVIS_PULL_REQUEST_SHA" -- > $tmpfile

  # Find lines that don't start with internal/website in the diff log; if no such
  # lines are found, it means that we don't have to run tests. grep returns 1 in
  # this case.
  echo "Looking for files that changed"
  if grep -v ^internal/website $tmpfile; then
    echo "Running tests"
  else
    echo "Diff doesn't affect tests; not running them"
    exit 0
  fi
fi

# Run Go tests for the root. Only do coverage for the Linux build
# because it is slow, and codecov will only save the last one anyway.
result=0
if [[ "$TRAVIS_OS_NAME" == "linux" ]]; then
  go test -mod=readonly -race -coverpkg=./... -coverprofile=coverage.out ./... || result=1
  if [ -f coverage.out ] && [ $result -eq 0 ]; then
    # Filter out test and sample packages.
    grep -v test coverage.out | grep -v samples > coverage2.out
    mv coverage2.out coverage.out
    bash <(curl -s https://codecov.io/bash)
  fi
else
  go test -mod=readonly -race ./... || result=1
  # No need to run wire checks or other module tests on OSs other than linux.
  exit $result
fi

# Ensure that the code has no extra dependencies (including transitive
# dependencies) that we're not already aware of by comparing with
# ./internal/testing/alldeps
#
# Whenever project dependencies change, rerun ./internal/testing/listdeps.sh
./internal/testing/listdeps.sh | diff ./internal/testing/alldeps - || {
  echo "FAIL: dependencies changed; compare listdeps.sh output with alldeps" && result=1
}

# Install wire; Moved here from the "install" step because we don't need to
# install wire if the diff doesn't require testing (see condition above).
go install -mod=readonly github.com/google/wire/cmd/wire

wire check ./... || result=1
# "wire diff" fails with exit code 1 if any diffs are detected.
wire diff ./... || { echo "FAIL: wire diff found diffs!" && result=1; }

# Run Go tests for each additional module, without coverage.
for path in "./internal/contributebot" "./samples/appengine"; do
  ( cd "$path" && exec go test -mod=readonly ./... ) || result=1
  ( cd "$path" && exec wire check ./... ) || result=1
  ( cd "$path" && exec wire diff ./... ) || (echo "FAIL: wire diff found diffs!" && result=1)
done
exit $result
