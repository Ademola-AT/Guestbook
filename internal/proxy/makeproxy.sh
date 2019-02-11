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

# This script downloads and arranges packages in a structure for a
# Go Modules proxy.

# https://coderwall.com/p/fkfaqq/safer-bash-scripts-with-set-euxo-pipefail
set -euxo pipefail

if [[ $# -ne 1 ]]; then
  echo "usage: makeproxy.sh <dir>" 1>&2
  exit 64
fi

# Set GOPATH to argument. go requires an absolute path, so append to
# working directory if need be.
export GOPATH="$1"
if ! [[ "$GOPATH" =~ ^/ ]]; then
  GOPATH="$(pwd)/$GOPATH"
fi

# Script is in internal/proxy, so repository root is two levels up.
repo_root="$( cd "$( dirname "${BASH_SOURCE[0]}" )/../.." >/dev/null 2>&1 && pwd )"

# Download modules for each of the modules in our repo.
for path in "." "internal/contributebot" "samples/appengine"; do
  cd "$repo_root/$path"
  go mod download
  # We need to get goveralls because Travis uses it for tests and will look
  # for it in the module proxy.
  go get github.com/mattn/goveralls
done
