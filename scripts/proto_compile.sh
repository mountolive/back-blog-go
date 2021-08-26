#! /bin/bash

cd proto

for dir in */; do
  cd "$dir"
	TMP_DIR="$(mktemp -d -p .)"
	protoc -I=. --go-grpc_out="$TMP_DIR" *.proto
	rm -rf $TMP_DIR
  cd ..
done
