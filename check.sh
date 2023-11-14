#!/bin/bash

set -e

if [ -t 1 ]
then
  YELLOW='\033[0;33m'
  GREEN='\033[0;32m'
  RED='\033[0;31m'
  NC='\033[0m'
fi

yellow() { printf "${YELLOW}%s${NC}" "$*"; }
green() { printf "${GREEN}%s${NC}" "$*"; }
red() { printf "${RED}%s${NC}" "$*"; }

good() {
  echo "$(green "● success:")" "$@"
}

bad() {
  ret=$1
  shift
  echo "$(red "● failed:")" "$@"
  exit "$ret"
}

try() {
  "$@" || bad $? "$@" && good "$@"
}

cmd_exists() {
  type "$1" >/dev/null 2>&1
}

cmd_exists go || bad 1 "uhhh, where's your go install? what are you even doing here!"
cmd_exists gofmt || bad 1 "so you have go installed, but not gofmt?^^"
cmd_exists golangci-lint || bad 1 "cannot find golangci-lint; check PATH or go to https://github.com/golangci/golangci-lint"
cmd_exists gitleaks || bad 1 "cannot find gitleaks 8.8.8 https://github.com/zricethezav/gitleaks"
cmd_exists sqlc || bad 1 "cannot find sqlc https://sqlc.dev/"

# Confirm that unique codes are actually unique
try ./unique_code.py

# Check that sqlc has been run as expected
try sqlc diff

# Ensure we are not leaking secrets by having them checked in
try gitleaks detect -v -c gitleaks.toml
try gitleaks protect -v -c gitleaks.toml

# Check for common vulnerability https://go.dev/blog/govulncheck
# NB you may want to disable this if you are using github's tools
#try govulncheck ./...

# Vendor dependencies
try go mod tidy
try go mod vendor

# If generates are set run them now
try go generate ./...
try sqlc generate

# Run basic fixes
try go fix ./...

# Lint code and autofix where possible
try golangci-lint run --enable "gofmt,stylecheck,misspell,gosec,errchkjson,errname,containedctx,bodyclose,bidichk,unused" ./...

# Start the Go service and then call the integration tests
try pkill necroteuch | true
try go build -o necroteuch && ./necroteuch > necroteuch.log 2>&1 &

# Wait for the service to be ready
timeout=10
while ! nc -z localhost 8080; do
    sleep 1
    ((timeout--))
    if [ $timeout -eq 0 ]; then
        echo "Timeout reached, the Go service is not ready."
        # If the service is not ready within the timeout, kill the process and exit
        pkill necroteuch
        exit 1
    fi
done

# Run integration tests and clean up now
try go test -count=1 --tags=integration ./...
try pkill necroteuch | true
{
 {
   opt='shopt -s extglob nullglob'
   gofmt='gofmt -s -w -l !(vendor)/ *.go'
   notice="    running: ( $opt; $gofmt; )"
   prefix="    $(yellow modified)"
   trap 'echo "$notice"; $opt; $gofmt | sed -e "s#^#$prefix #g"' EXIT
 }

 # comma separate linters (e.g. "gofmt,stylecheck")
 additional_linters="gofmt"
 try golangci-lint run --enable $additional_linters ./...
 trap '' EXIT
}
