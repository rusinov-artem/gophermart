#!/bin/bash

echo "Compaling ..."
go build -o ./test/bintest/app ./cmd/gophermart/main.go
R_VAL=$?

if [[ R_VAL -ne "0" ]] ; then
      echo "!!! compilation failed !!!"
      exit 1
fi

echo "OK! Compilation succeeded"

echo "Runing tests..."

TEST_CMD="go test -v -count=1 -json -coverpkg=./... -coverprofile=coverage.out ./..."

TEST_OUT=$(${TEST_CMD})

R_VAL=$?
if [[ R_VAL -ne "0" ]] ; then
  FAIL="TRUE"
fi

# Filter failed tests
 FAILED_TESTS=$( echo "${TEST_OUT}" | jq -c 'select(.Action=="fail")')
 R_VAL=$?
 if [[ R_VAL -ne "0" ]] ; then
   echo "!!! TESTS FAILED !!!"
   exit 2
 fi

if [[ ${FAIL} ]]; then
    echo "${FAILED_TESTS}"
    echo "!!! TESTS FAILED !!!"
    exit 3 # needed to make build red
else
   echo "!!! TESTS PASSED !!!"
fi

