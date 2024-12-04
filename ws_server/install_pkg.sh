#!/bin/bash
while IFS= read -r line; do
    go get "$line"
done < dependencies.txt
