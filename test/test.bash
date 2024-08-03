#!/bin/bash

echo "Compaling ..."
go build -o ./test/bintest/app ./cmd/gophermart/main.go
R_VAL=$?

if [[ R_VAL -ne "0" ]] ; then
      echo "!!! compilation failed !!!"
      exit 1
fi

echo "OK! Compilation succeeded"


