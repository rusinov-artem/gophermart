#!/bin/bash
set -e

cd /app

if [ ! -f test.log ]; then
    echo "test.log не найден, запускаю test.bash..."
    ./test/test.bash
fi

FAILED_TESTS=$(jq -r 'select(.Action=="fail" and .Test) | .Test' test.log | sort -u)

if [ -z "$FAILED_TESTS" ]; then
    echo "Нет упавших тестов"
    exit 0
fi

SELECTED=$(echo "$FAILED_TESTS" | fzf --height=~50% --border)

if [ -z "$SELECTED" ]; then
    echo "Тест не выбран"
    exit 0
fi

echo "Запускаю тест: $SELECTED"
go test -v -count=1 -run "^${SELECTED}$" ./...