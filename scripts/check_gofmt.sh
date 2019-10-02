#!/bin/bash
bad_files=$(find . -name '*.go' | xargs gofmt -s -l)
if [[ -n "${bad_files}" ]]; then
  echo "!!! go fmt needs to be run on the following files:"
  echo "${bad_files}"
  exit 1
fi
