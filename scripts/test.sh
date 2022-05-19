#!/bin/bash
# This script performs tests against the project, specifically:
#
#   * gofmt         (https://golang.org/cmd/gofmt)
#   * goimports     (https://godoc.org/cmd/goimports)
#   * golint        (https://github.com/golang/lint)
#   * go vet        (https://golang.org/cmd/vet)
#   * test coverage (https://blog.golang.org/cover)
#
# It outputs test and coverage reports in a way that Jenkins can understand,
# with test results in JUnit format and test coverage in Cobertura format.
# The reports are saved to bin/$SUBDIR/{test-reports,coverage-reports}/*.xml
#
set -e
set -o pipefail
export PATH="${GOPATH}/bin:${PATH}"
export GO111MODULE=on

PACKAGES="$(go list ./... | grep -v /vendor/ | grep -v /internal/keto)"
SUBDIRS=$(go list -f {{.Dir}} ./... | grep -v /vendor/)
SOURCE_DIR=$(git rev-parse --show-toplevel)
BUILD_DIR="${SOURCE_DIR}/bin"

function logmsg() {
  echo -e "\n\n*** $1 ***"
}

function _gofmt() {
  logmsg "Running 'gofmt' ..."
  test -z "$(gofmt -l -d ${SUBDIRS} | tee /dev/stderr)"
  echo "PASS"
}

function _goimports() {
  logmsg "Running 'goimports' ..."
  echo "Files (if any) listed below are not properly formatted with goimports:"
  go install golang.org/x/tools/cmd/goimports@latest
  test -z "$(goimports -l ${SUBDIRS} | tee /dev/stderr)"
  echo "(N/A)"
  echo "PASS"
}

function _golint() {
  logmsg "Running 'go lint' ..."
  go install golang.org/x/lint/golint@latest
  for pkg in $PACKAGES; do
    golint -set_exit_status $pkg
  done
  echo "PASS"
}

function _govet() {
  logmsg "Running 'go vet' ..."
  go vet ${PACKAGES}
  echo "PASS"
}

function _unittest_with_coverage() {
  local covermode="atomic"
  logmsg "Running unit tests ..."

  echo "cleaning test cache..."
  go clean -testcache

  go install github.com/jstemmer/go-junit-report@latest
  go install github.com/smartystreets/goconvey@latest
  go install golang.org/x/tools/cmd/cover@latest
  go install github.com/axw/gocov/gocov@latest
  go install github.com/AlekSi/gocov-xml@latest

  # We can't' use the test profile flag with multiple packages. Therefore,
  # run 'go test' for each package, and concatenate the results into
  # 'profile.cov'.
  mkdir -p ${BUILD_DIR}/{test-reports,coverage-reports}
  echo "mode: ${covermode}" >${BUILD_DIR}/coverage-reports/profile.cov

  for import_path in ${PACKAGES}; do
    package=$(basename ${import_path})

    go test -mod=mod -v -race -covermode=$covermode \
      -coverprofile="${BUILD_DIR}/coverage-reports/profile_${package}.cov" \
      $import_path | tee /dev/stderr |
      go-junit-report >"${BUILD_DIR}/test-reports/${package}-report.xml"

  done

  # Concatenate per-package coverage reports into a single file.
  for f in ${BUILD_DIR}/coverage-reports/profile_*.cov; do
    tail -n +2 ${f} >>${BUILD_DIR}/coverage-reports/profile.cov
    rm $f
  done

  go tool cover -func ${BUILD_DIR}/coverage-reports/profile.cov
  gocov convert ${BUILD_DIR}/coverage-reports/profile.cov |
    gocov-xml >"${BUILD_DIR}/coverage-reports/coverage.xml"
  echo "PASS"
}

function main() {
  logmsg "Excluding Keto Tests"
  _gofmt
  _goimports
  # _golint
  _govet
  _unittest_with_coverage
}

main
