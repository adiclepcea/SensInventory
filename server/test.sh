#!/bin/sh

unformatted_files=$(find . -name "*.go" | xargs gofmt -l)

[ -z "$unformatted_files" ] && exit 0

echo "\nUnformated files: \n###############\n$unformatted_files\n##############\nPlease run go fmt before commiting\n"

exit 1
