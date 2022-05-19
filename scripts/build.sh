#!/bin/bash

# Copyright 2016 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

if [ -z "${OS:-}" ]; then
    echo "OS must be set"
    exit 1
fi
if [ -z "${ARCH:-}" ]; then
    echo "ARCH must be set"
    exit 1
fi
if [ -z "${BINARY_NAME:-}" ]; then
    echo "BINARY_NAME must be set"
    exit 1
fi
if [ -z "${BUILD_TARGET:-}" ]; then
    echo "BUILD_TARGET must be set"
    exit 1
fi
if [ -z "${BINARY_OUTPUT_ROOT:-}" ]; then
    echo "BINARY_OUTPUT_ROOT must be set"
    exit 1
fi
if [ -z "${VERSION:-}" ]; then
    echo "VERSION must be set"
    exit 1
fi

export CGO_ENABLED=0
export GOARCH="${ARCH}"
export GOOS="${OS}"
export GO111MODULE=on

mkdir -p "${BINARY_OUTPUT_ROOT}"/"${OS}"_"${ARCH}"

X_LDFLAGS="-X $(go list -m)/cmd.version=${VERSION} -X '$(go list -m)/cmd.goVersion=$(go version)' -X $(go list -m)/cmd.buildStamp=$(date -u '+%Y-%m-%dT%H:%M:%SZ')"

# shellcheck disable=SC2086
go build \
    -trimpath \
    -installsuffix "static" \
    -ldflags "-s" \
    -o "${BINARY_OUTPUT_ROOT}/${OS}_${ARCH}/${BINARY_NAME}" \
    ${BUILD_TARGET}
