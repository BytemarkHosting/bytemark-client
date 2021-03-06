#!/bin/bash

EXPECTED=$(grep 'sha256sum:' HACKING)
EXPECTED=$(echo -e "${EXPECTED##*sha256sum: }" | tr -d '[[:space:]]')

TREE=$(find . -type d \! -path './vendor/*' \! -path './.*' | sort)

ACTUAL=$(echo -e "$TREE" | sha256sum)
ACTUAL=$(echo -e "${ACTUAL%% -}" | tr -d '[[:space:]]')

[ "$EXPECTED" == "$ACTUAL" ]
