#! /bin/bash

input=$(cat)
package=$(echo "$input" | jq '.Package' | tr -d '"')
test=$(echo "$input" | jq '.Test' | tr -d '"')

echo "package $package"
echo "test $test"

go test -v "$package" -run "$test"
