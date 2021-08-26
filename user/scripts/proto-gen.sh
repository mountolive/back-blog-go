#! /bin/bash

BASE_DIR="grpc"

PKG_DIR="$(mkdir -p "$(pwd)/$BASE_DIR")"

cd ../proto/user

for dir in */; do
  cd "$dir"
	protoc -I=. --go-grpc_out="$PKG_DIR" *.proto
  cd ..
done
