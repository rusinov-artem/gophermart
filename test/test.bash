#!/bin/bash

# needed to run tests with -race flag
export CGO_ENABLED=1
export GOCOVERDIR=/app/test/bintest/coverdir
rm ${GOCOVERDIR:?}/* -r

APP_BIN=/app/test/bintest/app
if [ -f ${APP_BIN} ]; then
  echo "app found and removed"
  rm ${APP_BIN:?}
fi

echo "Compaling ..."
go build -cover -o ./test/bintest/app ./cmd/gophermart
R_VAL=$?

if [[ R_VAL -ne "0" ]] ; then
      echo "!!! compilation failed !!!"
      exit 1
fi

echo "OK! Compilation succeeded"

export GOCOVERDIR=/app/test/bintest/coverdir/Migration
mkdir ${GOCOVERDIR}
./test/bintest/app migrate

echo "Runing tests..."

TEST_CMD="go test -v -count=1 -json -coverpkg=./... -covermode=set -coverprofile=coverage.out ./..."

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


COVERAGE_DIR_LIST=$(find /app/test/bintest/coverdir/ -maxdepth 1 -type d | sort | tail -n +2 | tr "\n" "," | sed 's/,$/\n/')
if [[ ${COVERAGE_DIR_LIST} ]]; then
  go tool covdata textfmt -i="${COVERAGE_DIR_LIST}"  -o bincoverage.out
  gocov-merger coverage.out bincoverage.out > merge.out
else
  cp coverage.out merge.out
fi

sed -i '/\/gophermart\/test/d' merge.out
go tool cover -html=merge.out -o coverage.html
go-cover-treemap -coverprofile merge.out > coverage.svg
echo -n "coverage "
go tool cover -func ./merge.out | tail -1 2>&1 | tr -s '\t' | tr -s ' '
