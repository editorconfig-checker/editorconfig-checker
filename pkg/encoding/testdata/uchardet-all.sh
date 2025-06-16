#!/usr/bin/env bash

mapfile -t < <(find . -iname '*.txt' -iname '*.html')

for file in "${MAPFILE[@]}"; do
    printf "%-50s: " "${file}"
    uchardet "${file}"
done
